package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/api/resources"
	"io"
	"net/http"
)

// CodesRouter returns router for
func (rest *Rest) CodesRouter() chi.Router {
	router := chi.NewRouter()

	router.Post("/", rest.createCode)
	router.Get("/", rest.listCodes)

	router.Route("/{id}", func(r chi.Router) {
		r.Get("/", rest.getCode)
		r.Get("/download", rest.downloadCode)
		r.Put("/", rest.updateCode)
		r.Delete("/", rest.deleteCode)
	})

	return router
}

func (rest *Rest) createCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	cr := resources.CodeCreateRequest{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to real body"))
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(reqBody, &cr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to deserialize body"))
		return
	}

	err = cr.Validate()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("wrong data"))
		return
	}

	code := entities.Code{
		UserID: userID,
		SrcURL: cr.URL,
	}

	id, err := rest.service.CreateCode(r.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	body, err := json.Marshal(resources.CodeCreateResponse{ID: id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error during building response"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (rest *Rest) listCodes(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	// TODO get codes by UserID
	w.Write([]byte(fmt.Sprintf("list codes for %d", userID)))
}

func (rest *Rest) getCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	// TODO validate request

	// TODO get code

	// TODO compare with UserID
	w.Write([]byte(fmt.Sprintf("get code for %d", userID)))
}

func (rest *Rest) downloadCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	// TODO validate request

	// TODO get code

	// TODO compare with UserID

	// TODO generate base64 for image
	w.Write([]byte(fmt.Sprintf("download code for %d", userID)))
}

func (rest *Rest) updateCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	// TODO validate request

	// TODO get code

	// TODO compare with UserID

	// TODO update in DB

	// TODO update in CACHE
	w.Write([]byte(fmt.Sprintf("update code for %d", userID)))
}

func (rest *Rest) deleteCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	// TODO validate request

	// TODO get code

	// TODO compare with UserID

	// TODO delete from CACHE

	// TODO delete from DB
	w.Write([]byte(fmt.Sprintf("delete code for %d", userID)))
}
