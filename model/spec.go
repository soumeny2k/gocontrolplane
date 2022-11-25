package model

import (
	"controlplane/util"
	"errors"
)

type Spec struct {
	ApiId   uint   `json:"api_id"`
	Spec    string `json:"spec"`
	Version uint   `json:"version"`
}

func (spec *Spec) Create() (string, error) {
	var dbSpec *Spec
	err := util.GetDB().Raw("SELECT * FROM spec WHERE api_id = ? AND version = ?", spec.ApiId, spec.Version).Scan(&dbSpec).Error
	if err != nil {
		err = util.GetDB().Create(spec).Error
		if err != nil {
			return "", errors.New("failed to create spec")
		}

	} else {
		err := util.GetDB().Model(Spec{}).Where("api_id = ? AND version = ?", spec.ApiId, spec.Version).Updates(
			Spec{Spec: spec.Spec},
		).Error
		if err != nil {
			return "", errors.New("failed to update spec, connection error")
		}
	}

	return "spec created successfully", nil
}
