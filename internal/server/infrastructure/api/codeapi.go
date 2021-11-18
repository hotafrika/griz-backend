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
		rest.writeErrorCode(w, http.StatusBadRequest, "user is not defined")
		return
	}

	cr := resources.CodeCreateRequest{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "unable to read body")
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(reqBody, &cr)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "unable to deserialize body")
		return
	}

	err = cr.Validate()
	if err != nil {
		rest.writeErrorCode(w, http.StatusUnprocessableEntity, "wrong data")
		return
	}

	code := entities.Code{
		UserID: userID,
		SrcURL: cr.URL,
	}

	id, err := rest.service.CreateCode(r.Context(), code)
	if err != nil {
		rest.logger.Error().Err(err).Send()
		rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
		return
	}

	body, err := json.Marshal(resources.CodeCreateResponse{ID: id})
	if err != nil {
		rest.logger.Error().Err(err).Send()
		rest.writeErrorCode(w, http.StatusInternalServerError, "error during building response")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (rest *Rest) listCodes(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		rest.writeErrorCode(w, http.StatusBadRequest, "user is not defined")
		return
	}

	codes, err := rest.service.GetCodes(r.Context(), userID)
	if err != nil {
		if !errors.Is(err, domain.ErrCodeNotFound) {
			rest.logger.Error().Err(err).Send()
			rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
			return
		}
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
		rest.logger.Error().Msg(err.Error())
		rest.writeErrorCode(w, http.StatusInternalServerError, "error during building response")
		return
	}

	w.Write(body)
}

func (rest *Rest) getCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		rest.writeErrorCode(w, http.StatusBadRequest, "user is not defined")
		return
	}

	codeIDString := chi.URLParam(r, "codeID")
	codeID, err := strconv.ParseUint(codeIDString, 10, 64)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "wrong code id")
		return
	}

	code, err := rest.service.GetCode(r.Context(), codeID)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			rest.writeErrorCode(w, http.StatusNotFound, "not found")
			return
		}
		rest.logger.Error().Err(err).Send()
		rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
		return
	}

	if code.UserID != userID {
		rest.writeErrorCode(w, http.StatusForbidden, "unauthorized")
		return
	}

	body, err := json.Marshal(resources.GetCodeResponse{
		ID:  code.ID,
		URL: code.SrcURL,
	})
	if err != nil {
		rest.logger.Error().Err(err).Send()
		rest.writeErrorCode(w, http.StatusInternalServerError, "error during building response")
		return
	}

	w.Write(body)
}

func (rest *Rest) downloadCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		rest.writeErrorCode(w, http.StatusBadRequest, "user is not defined")
		return
	}

	codeIDString := chi.URLParam(r, "codeID")
	codeID, err := strconv.ParseUint(codeIDString, 10, 64)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "wrong code id")
		return
	}

	code, err := rest.service.GetCode(r.Context(), codeID)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			rest.writeErrorCode(w, http.StatusNotFound, "not found")
			return
		}
		rest.logger.Error().Err(err).Send()
		rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
		return
	}

	if code.UserID != userID {
		rest.writeErrorCode(w, http.StatusForbidden, "unauthorized")
		return
	}

	encodedQR, err := rest.service.DownloadCodeByHash(r.Context(), code.Hash)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			rest.writeErrorCode(w, http.StatusNotFound, "not found")
			return
		}
		rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
		return
	}

	body, err := json.Marshal(resources.DownloadCodeResponse{Code: encodedQR})
	if err != nil {
		rest.logger.Error().Err(err).Send()
		rest.writeErrorCode(w, http.StatusInternalServerError, "error during building response")
		return
	}

	w.Write(body)
}

func (rest *Rest) updateCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		rest.writeErrorCode(w, http.StatusBadRequest, "user is not defined")
		return
	}

	codeIDString := chi.URLParam(r, "codeID")
	codeID, err := strconv.ParseUint(codeIDString, 10, 64)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "wrong code id")
		return
	}

	cr := resources.CodeCreateRequest{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "unable to read body")
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(reqBody, &cr)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "unable to deserialize body")
		return
	}

	err = cr.Validate()
	if err != nil {
		rest.writeErrorCode(w, http.StatusUnprocessableEntity, "wrong data")
		return
	}

	code, err := rest.service.GetCode(r.Context(), codeID)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			rest.writeErrorCode(w, http.StatusNotFound, "not found")
			return
		}
		rest.logger.Error().Err(err).Send()
		rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
		return
	}

	if code.UserID != userID {
		rest.writeErrorCode(w, http.StatusForbidden, "unauthorized")
		return
	}

	code.SrcURL = cr.URL

	err = rest.service.UpdateCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			rest.writeErrorCode(w, http.StatusNotFound, "not found")
			return
		}
		rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
		return
	}

	body, err := json.Marshal(resources.GetCodeResponse{ID: code.ID, URL: code.SrcURL})
	if err != nil {
		rest.logger.Error().Err(err).Send()
		rest.writeErrorCode(w, http.StatusInternalServerError, "error during building response")
		return
	}

	w.Write(body)
}

func (rest *Rest) deleteCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdInCtx).(uint64)
	if !ok {
		rest.writeErrorCode(w, http.StatusBadRequest, "user is not defined")
		return
	}

	codeIDString := chi.URLParam(r, "codeID")
	codeID, err := strconv.ParseUint(codeIDString, 10, 64)
	if err != nil {
		rest.writeErrorCode(w, http.StatusBadRequest, "wrong code id")
		return
	}

	code, err := rest.service.GetCode(r.Context(), codeID)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			rest.writeErrorCode(w, http.StatusNotFound, "not found")
			return
		}
		rest.logger.Error().Err(err).Send()
		rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
		return
	}

	if code.UserID != userID {
		rest.writeErrorCode(w, http.StatusForbidden, "unauthorized")
		return
	}

	err = rest.service.DeleteCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, domain.ErrCodeNotFound) {
			rest.writeErrorCode(w, http.StatusNotFound, "not found")
			return
		}
		rest.writeErrorCode(w, http.StatusInternalServerError, "internal error")
		return
	}

	body, err := json.Marshal(resources.DeleteCodeResponse{Status: "ok"})
	if err != nil {
		rest.logger.Error().Err(err).Send()
		rest.writeErrorCode(w, http.StatusInternalServerError, "error during building response")
		return
	}

	w.Write(body)
}
