package model

import (
	"controlplane/kafka"
	"controlplane/util"
	"encoding/json"
	"errors"
)

type Api struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	TeamId            uint      `json:"team_id"`
	RateLimit         uint      `json:"rate_limit"`
	ConnectionTimeout uint      `json:"connection_timeout"`
	Protocols         string    `json:"protocol" default:"http"`
	Retries           uint      `json:"retries"`
	Balance           string    `json:"balance"`
	Type              string    `json:"type" default:"REST"`
	Backends          []Backend `json:"backends" gorm:"-"`
}

func (api *Api) Create() (string, error) {
	util.GetDB().Create(api)

	if api.ID <= 0 {
		return "", errors.New("failed to create api, connection error")
	}

	for i := 0; i < len(api.Backends); i++ {
		backend := &api.Backends[i]
		backend.ApiId = api.ID
	}
	err := util.GetDB().Create(api.Backends).Error
	if err != nil {
		return "", errors.New("failed to create backend, connection error")
	}

	latestApi, _ := api.Get()
	apiData, _ := json.Marshal(latestApi)
	b, _ := json.Marshal(
		kafka.Data{
			Event: "API",
			Data:  string(apiData),
		})
	kafka.Publish(b)

	return "api created successfully", nil
}

func (api *Api) Get() (*Api, error) {
	result := &Api{}
	err := util.GetDB().Table("api").Where("id = ?", api.ID).First(result).Error
	if err != nil {
		return &Api{}, errors.New("api not found")
	}

	result.Backends, _ = api.GetBackends()
	return result, nil
}

func (api *Api) GetAll() []Api {
	var result []Api
	err := util.GetDB().Find(&result).Error
	if err != nil {
		return make([]Api, 0)
	}
	return result
}

func (api *Api) Update() (string, error) {
	existingApi, err := api.Get()
	if err != nil {
		return "", err
	}
	var name string
	var rateLimit uint
	var connectionTimeout uint
	if api.Name != "" {
		name = api.Name
	} else {
		name = existingApi.Name
	}
	if api.RateLimit > 0 {
		rateLimit = api.RateLimit
	} else {
		rateLimit = existingApi.RateLimit
	}
	if api.ConnectionTimeout > 0 {
		connectionTimeout = api.ConnectionTimeout
	} else {
		connectionTimeout = existingApi.ConnectionTimeout
	}
	err = util.GetDB().Model(Api{}).Where("id = ?", api.ID).Updates(
		Api{Name: name, RateLimit: rateLimit, ConnectionTimeout: connectionTimeout},
	).Error
	if err != nil {
		return "", errors.New("failed to update api, connection error")
	}

	latestApi, _ := api.Get()
	apiData, _ := json.Marshal(latestApi)
	b, _ := json.Marshal(
		kafka.Data{
			Event: "API",
			Data:  string(apiData),
		})
	kafka.Publish(b)

	return "api updated successfully", nil
}

func (api *Api) Delete() (string, error) {
	err := util.GetDB().Where("id = ?", api.ID).Delete(&Api{}).Error
	if err != nil {
		return "", errors.New("api not found")
	}
	return "api deleted successfully", nil
}

func (api *Api) GetBackends() ([]Backend, error) {
	var result []Backend
	err := util.GetDB().Raw("select * from backend where api_id = ?", api.ID).Scan(&result).Error
	if err != nil {
		return make([]Backend, 0), errors.New("error while fetching backend")
	}
	return result, nil
}
