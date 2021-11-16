package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/hotafrika/griz-backend/internal/server/domain"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/api/resources"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strconv"
)

// CodesRouter returns router for
func (rest *Rest) CodesRouter() chi.Router {
	router := chi.NewRouter()

	router.Post("/", rest.createCode)
	router.Get("/", rest.listCodes)

	router.Route("/{codeID}", func(r chi.Router) {
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
		w.Write([]byte("unable to read body"))
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

	codes, err := rest.service.GetCodes(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	newCodes := make([]resources.GetCodeResponse, 0, len(codes))
	for _, code := range codes {
		newCodes = append(newCodes, resources.GetCodeResponse{
			ID:  code.ID,
			URL: code.SrcURL,
		})
	}
	newCodesR := resources.GetCodesResponse{
		Codes: newCodes,
	}

	body, err := json.Marshal(newCodesR)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error during building response"))
		return
	}

	w.Write(body)
}

func (rest *Rest) getCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	codeIDString := chi.URLParam(r, "codeID")
	codeID, err := strconv.ParseUint(codeIDString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong code id"))
		return
	}

	code, err := rest.service.GetCode(r.Context(), codeID)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	if code.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("unauthorized"))
		return
	}

	body, err := json.Marshal(resources.GetCodeResponse{
		ID:  code.ID,
		URL: code.SrcURL,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error during building response"))
		return
	}

	w.Write(body)
}

func (rest *Rest) downloadCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	codeIDString := chi.URLParam(r, "codeID")
	codeID, err := strconv.ParseUint(codeIDString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong code id"))
		return
	}

	code, err := rest.service.GetCode(r.Context(), codeID)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	if code.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("unauthorized"))
		return
	}

	encodedQR, err := rest.service.DownloadCodeByHash(r.Context(), code.Hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	body, err := json.Marshal(resources.DownloadCodeResponse{Code: encodedQR})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error during building response"))
		return
	}

	w.Write(body)
}

func (rest *Rest) updateCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	codeIDString := chi.URLParam(r, "codeID")
	codeID, err := strconv.ParseUint(codeIDString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong code id"))
		return
	}

	cr := resources.CodeCreateRequest{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to read body"))
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

	code, err := rest.service.GetCode(r.Context(), codeID)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	if code.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("unauthorized"))
		return
	}

	code.SrcURL = cr.URL

	rest.service.UpdateCode(r.Context(), code)

	body, err := json.Marshal(resources.GetCodeResponse{ID: code.ID, URL: code.SrcURL})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error during building response"))
		return
	}

	w.Write(body)
}

func (rest *Rest) deleteCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not defined"))
		return
	}

	codeIDString := chi.URLParam(r, "codeID")
	codeID, err := strconv.ParseUint(codeIDString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong code id"))
		return
	}

	code, err := rest.service.GetCode(r.Context(), codeID)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	if code.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("unauthorized"))
		return
	}

	err = rest.service.DeleteCode(r.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	body, err := json.Marshal(resources.DeleteCodeResponse{Status: "ok"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error during building response"))
		return
	}

	w.Write(body)
}
