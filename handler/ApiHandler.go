package handler

import (
	"controlplane/model"
	"controlplane/transferobject"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

func CreateApi(w http.ResponseWriter, r *http.Request) {
	var api transferobject.Api
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&api); err != nil {
		respondWithJSON(w, http.StatusBadRequest, transferobject.Response{
			Message: "failed to parse input json",
		})
		return
	}
	defer r.Body.Close()

	backends := make([]model.Backend, 0)
	for _, backend := range api.Backends {
		backends = append(backends, model.Backend{Url: backend.Url, Weight: backend.Weight})
	}

	dbApi := model.Api{
		Name:              api.Name,
		TeamId:            api.TeamId,
		RateLimit:         api.RateLimit,
		ConnectionTimeout: api.ConnectionTimeout,
		Backends:          backends,
	}
	response, err := dbApi.Create()

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func UploadSpec(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	version, err := strconv.ParseUint(vars["version"], 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, transferobject.Response{
			Message: "invalid id",
		})
		return
	}
	defer r.Body.Close()

	spec, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, transferobject.Response{
			Message: "unable to read body",
		})
		return
	}

	dbSpec := model.Spec{ApiId: uint(id), Version: uint(version), Spec: string(spec)}
	response, err := dbSpec.Create()

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func GetApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, transferobject.Response{
			Message: "invalid id",
		})
		return
	}

	dbApi := model.Api{ID: uint(id)}
	response, err := dbApi.Get()

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func GetAllApi(w http.ResponseWriter, r *http.Request) {
	dbApi := model.Api{}
	respondWithJSON(w, http.StatusOK, dbApi.GetAll())
}

func UpdateApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, transferobject.Response{
			Message: "invalid id",
		})
		return
	}

	var api transferobject.Api
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&api); err != nil {
		respondWithJSON(w, http.StatusBadRequest, transferobject.Response{
			Message: "failed to parse input json",
		})
		return
	}
	defer r.Body.Close()

	dbApi := model.Api{
		ID:                uint(id),
		Name:              api.Name,
		TeamId:            api.TeamId,
		RateLimit:         api.RateLimit,
		ConnectionTimeout: api.ConnectionTimeout,
	}
	response, err := dbApi.Update()

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, transferobject.Response{
			Message: err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, transferobject.Response{
		Message: response,
	})
}

func DeleteApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, transferobject.Response{
			Message: "invalid id",
		})
		return
	}

	dbApi := model.Api{ID: uint(id)}
	response, err := dbApi.Delete()

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, transferobject.Response{
			Message: err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, transferobject.Response{
		Message: response,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
