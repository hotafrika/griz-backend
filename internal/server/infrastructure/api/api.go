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

// ServeHTTP
//func (rest *Rest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	rest.router.ServeHTTP(w, r)
//}

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
				r.Get("/token", rest.tokenHandler)
			})
		})
	})
}

// HANDLERS
func (rest *Rest) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found page"))
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	user, err := rest.service.GetUser(r.Context(), userID)
	if err != nil {
		// TODO check errors type
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to real body"))
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(reqBody, &tr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to deserialize body"))
		return
	}

	err = tr.Validate()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("incomplete credentials"))
		return
	}

	user := entities.User{
		Username: tr.Username,
		Password: tr.Password,
	}

	authToken, err := rest.service.CreateAuthToken(r.Context(), user)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("wrong credentials"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	body, err := json.Marshal(resources.AuthTokenResponse{Token: authToken})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error during building response"))
		return
	}

	w.Write(body)
}

func (rest *Rest) urlHandler(w http.ResponseWriter, r *http.Request) {
	l := resources.LinkHashRequest{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to real body"))
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(reqBody, &l)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to deserialize body"))
		return
	}

	token, err := l.Parse()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("link is not compatible"))
		return
	}

	link, err := rest.service.FindCodeByHash(r.Context(), token)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("link not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	body, err := json.Marshal(resources.LinkHashResponse{URL: link})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error during building response"))
		return
	}

	w.Write(body)
}

func (rest *Rest) scanHandler(w http.ResponseWriter, r *http.Request) {
	sl := resources.SocialLinkRequest{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to real body"))
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(reqBody, &sl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to deserialize body"))
		return
	}

	err = sl.ValidateIsInsta()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("link is not compatible"))
		return
	}

	link, err := rest.service.FindCodeBySocial(r.Context(), sl.URL)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("unable to process link"))
		return
	}

	body, err := json.Marshal(resources.SocialLinkResponse{URL: link})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error during building response"))
		return
	}

	w.Write(body)
}

// MIDDLEWARE
func (rest *Rest) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		userID, err := rest.service.GetUserIDByAuthToken(r.Context(), authToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), userIdInCtx, userID))

		next.ServeHTTP(w, r)
	})
}
