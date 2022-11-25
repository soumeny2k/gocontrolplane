package handler

import (
	"controlplane/model"
	"controlplane/transferobject"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func CreateRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiId, err := strconv.ParseUint(vars["api_id"], 10, 64)
	var route transferobject.Route
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&route); err != nil {
		fmt.Println(err)
		respondWithJSON(w, http.StatusBadRequest, transferobject.Response{
			Message: "failed to parse input json",
		})
		return
	}
	defer r.Body.Close()

	dbRoute := model.Route{
		Name:  route.Name,
		ApiId: uint(apiId),
		Path:  route.Path,
	}

	headers := make([]model.RouteHeader, 0)
	for _, header := range route.Headers {
		headers = append(headers, model.RouteHeader{Name: header.Name, Value: header.Value})
	}

	dbRoute.Headers = headers
	response, err := dbRoute.Create()

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func UpdateRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiId, err := strconv.ParseUint(vars["api_id"], 10, 64)
	routeId, err := strconv.ParseUint(vars["route_id"], 10, 64)

	var route transferobject.Route
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&route); err != nil {
		fmt.Println(err)
		respondWithJSON(w, http.StatusBadRequest, transferobject.Response{
			Message: "failed to parse input json",
		})
		return
	}
	defer r.Body.Close()

	dbRoute := model.Route{
		ID:    uint(routeId),
		Name:  route.Name,
		ApiId: uint(apiId),
		Path:  route.Path,
	}
	response, err := dbRoute.Update()

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func GetRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiId, err := strconv.ParseUint(vars["api_id"], 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, transferobject.Response{
			Message: "invalid api id",
		})
		return
	}
	routeId, err := strconv.ParseUint(vars["route_id"], 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, transferobject.Response{
			Message: "invalid route id",
		})
		return
	}

	route := model.Route{ID: uint(routeId), ApiId: uint(apiId)}
	response, err := route.Get()

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}
