package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (rest *Rest) CodesRouter() chi.Router {
	router := chi.NewRouter()

	router.Get("/", rest.listCodes)
	router.Post("/", rest.createCode)

	router.Route("/{id}", func(r chi.Router) {
		r.Get("/", rest.getCode)
		r.Get("/download", rest.downloadCode)
		r.Put("/", rest.updateCode)
		r.Delete("/", rest.deleteCode)
	})

	return router
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

func (rest *Rest) createCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	// TODO validate request

	// TODO create code in DB

	// TODO add to CACHE
	w.Write([]byte(fmt.Sprintf("create code for %d", userID)))
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
