package model

import (
	"controlplane/kafka"
	"controlplane/util"
	"encoding/json"
	"errors"
)

type Route struct {
	ID                uint          `json:"id"`
	ApiId             uint          `json:"api_id"`
	Name              string        `json:"name"`
	Path              string        `json:"path"`
	PathPattern       string        `json:"path_pattern"`
	Method            string        `json:"method"`
	Retries           uint          `json:"retries"`
	RateLimit         uint          `json:"rate_limit"`
	ConnectionTimeout uint          `json:"connection_timeout"`
	CacheEnabled      bool          `json:"cache_enabled"`
	Headers           []RouteHeader `json:"headers" gorm:"-"`
}

func (route *Route) Create() (string, error) {
	util.GetDB().Create(route)

	if route.ID <= 0 {
		return "", errors.New("failed to create route, connection error")
	}

	for i := 0; i < len(route.Headers); i++ {
		header := &route.Headers[i]
		header.RouteId = route.ID
	}
	err := util.GetDB().Create(route.Headers).Error
	if err != nil {
		return "", errors.New("failed to create header, connection error")
	}

	latestRoute, _ := route.Get()
	routeData, _ := json.Marshal(latestRoute)
	b, _ := json.Marshal(
		kafka.Data{
			Event: "ROUTE",
			Data:  string(routeData),
		})
	kafka.Publish(b)

	return "route created successfully", nil
}

func (route *Route) Update() (string, error) {
	err := util.GetDB().Model(Route{}).Where("id = ?", route.ID).Updates(
		Route{Name: route.Name, RateLimit: route.RateLimit},
	).Error
	if err != nil {
		return "", errors.New("failed to update route, connection error")
	}

	latestRoute, _ := route.Get()
	routeData, _ := json.Marshal(latestRoute)
	b, _ := json.Marshal(
		kafka.Data{
			Event: "ROUTE",
			Data:  string(routeData),
		})
	kafka.Publish(b)

	return "route updated successfully", nil
}

func (route *Route) Get() (*Route, error) {
	result := &Route{}
	err := util.GetDB().Table("route").Where("id = ?", route.ID).First(result).Error
	if err != nil {
		return &Route{}, errors.New("route not found")
	}
	result.Headers, _ = route.GetHeaders()
	return result, nil
}

func (route *Route) GetHeaders() ([]RouteHeader, error) {
	var result []RouteHeader
	err := util.GetDB().Raw("select * from route_header where route_id = ?", route.ID).Scan(&result).Error
	if err != nil {
		return make([]RouteHeader, 0), errors.New("error while fetching header")
	}
	return result, nil
}
