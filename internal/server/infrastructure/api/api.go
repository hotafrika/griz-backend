package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hotafrika/griz-backend/internal/server/app"
	"github.com/hotafrika/griz-backend/internal/server/domain"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/api/resources"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"time"
)

const userIdInCtx = "user_id"

type Rest struct {
	bindAddr    string
	timeout     time.Duration
	scanTimeout time.Duration
	logger      zerolog.Logger
	service     app.CodeService
	router      chi.Router
	server      *http.Server
}

// NewRest creates Rest api
// TODO
func NewRest(bindAddr string, timeout time.Duration, parseTimeout time.Duration, logger zerolog.Logger, service app.CodeService) *Rest {
	r := &Rest{
		bindAddr:    bindAddr,
		timeout:     timeout,
		scanTimeout: parseTimeout,
		logger:      logger,
		service:     service,
		router:      chi.NewRouter(),
	}
	r.configureRouter()
	r.server = &http.Server{
		Addr:    bindAddr,
		Handler: r.router,
	}
	return r
}

// Start starts http listeners
func (rest *Rest) Start() error {
	rest.server = &http.Server{
		Addr:    rest.bindAddr,
		Handler: rest.router,
	}
	return rest.server.ListenAndServe()
}

func (rest *Rest) configureRouter() {
	// content block
	rest.router.Get("/", rest.homepageHandler)
	rest.router.Get("/apps", rest.downloadAppsHandler)

	// /api
	rest.router.Route("/api", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(rest.timeout))

		// api/v1
		r.Route("/v1", func(r chi.Router) {
			// with auth
			r.Group(func(r chi.Router) {
				r.Use(rest.authMiddleware)
				// api/v1/code...
				r.Mount("/code", rest.CodesRouter())
				r.Get("/self", rest.userSelfHandler)
			})
			// api/v1/public/...
			r.Route("/public", func(r chi.Router) {
				r.Post("/url", rest.urlHandler)
				r.Group(func(r chi.Router) {
					//r.Use(middleware.Throttle(10))
					r.Use(middleware.Timeout(rest.scanTimeout))
					r.Post("/scan", rest.scanHandler)
				})
			})
			// api/v1/token
			r.Group(func(r chi.Router) {
				//r.Use(middleware.Throttle(10))
				r.Post("/token", rest.tokenHandler)
			})
		})
	})
}

// HANDLERS
func (rest *Rest) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	rest.writeErrorCode(w, http.StatusNotFound, "page not found")
}

func (rest *Rest) homepageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("main page"))
}

func (rest *Rest) downloadAppsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("download page"))
}

func (rest *Rest) userSelfHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint64)
	if !ok {
		rest.writeErrorCode(w, http.StatusBadRequest, "user is not defined")
		return
	}

	user, err := rest.service.GetUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			rest.writeErrorCode(w, http.StatusNotFound, "user not found")
		}
		rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
		return
	}

	resourceUser := resources.SelfUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	body, err := json.Marshal(resourceUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

func (rest *Rest) tokenHandler(w http.ResponseWriter, r *http.Request) {
	tr := resources.AuthTokenRequest{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "unable to read body")
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(reqBody, &tr)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "unable to deserialize body")
		return
	}

	err = tr.Validate()
	if err != nil {
		rest.writeErrorCode(w, http.StatusUnprocessableEntity, "incomplete credentials")
		return
	}

	user := entities.User{
		Username: tr.Username,
		Password: tr.Password,
	}

	authToken, err := rest.service.CreateAuthToken(r.Context(), user)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			rest.writeErrorCode(w, http.StatusUnauthorized, "wrong credentials")
			return
		}
		rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
		return
	}

	body, err := json.Marshal(resources.AuthTokenResponse{Token: authToken})
	if err != nil {
		rest.writeErrorCode(w, http.StatusInternalServerError, "error during building response")
		return
	}

	w.Write(body)
}

func (rest *Rest) urlHandler(w http.ResponseWriter, r *http.Request) {
	l := resources.LinkHashRequest{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "unable to read body")
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(reqBody, &l)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "unable to deserialize body")
		return
	}

	token, err := l.Parse()
	if err != nil {
		rest.writeErrorCode(w, http.StatusUnprocessableEntity, "link is not compatible")
		return
	}

	link, err := rest.service.FindCodeByHash(r.Context(), token)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			rest.writeErrorCode(w, http.StatusNotFound, "link not found")
			return
		}
		rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
		return
	}

	body, err := json.Marshal(resources.LinkHashResponse{URL: link})
	if err != nil {
		rest.writeErrorCode(w, http.StatusInternalServerError, "error during building response")
		return
	}

	w.Write(body)
}

func (rest *Rest) scanHandler(w http.ResponseWriter, r *http.Request) {
	sl := resources.SocialLinkRequest{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "unable to read body")
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(reqBody, &sl)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "unable to deserialize body")
		return
	}

	err = sl.ValidateIsInsta()
	if err != nil {
		rest.writeErrorCode(w, http.StatusUnprocessableEntity, "link is not compatible")
		return
	}

	link, err := rest.service.FindCodeBySocial(r.Context(), sl.URL)
	if err != nil {
		rest.writeErrorCode(w, http.StatusUnprocessableEntity, "unable to process link")
		return
	}

	body, err := json.Marshal(resources.SocialLinkResponse{URL: link})
	if err != nil {
		rest.writeErrorCode(w, http.StatusInternalServerError, "error during building response")
		return
	}

	w.Write(body)
}

// MIDDLEWARE
func (rest *Rest) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			rest.writeErrorCode(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID, err := rest.service.GetUserIDByAuthToken(r.Context(), authToken)
		if err != nil {
			if errors.Is(err, domain.ErrCacheNotExist) {
				rest.writeErrorCode(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), userIdInCtx, userID))

		next.ServeHTTP(w, r)
	})
}

func (rest *Rest) writeErrorCode(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	b, _ := json.Marshal(
		struct {
			Message string `json:"message"`
		}{
			Message: message,
		})
	w.Write(b)
}
