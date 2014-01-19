package handlers

import (
	"github.com/felixge/quantastic/backend/version"
)

func NewGetVersion(version version.Version) GetVersion {
	return getVersionHandler{getVersionResponse{version}}
}

type getVersionHandler struct {
	response getVersionResponse
}

func (h getVersionHandler) GetVersion(req GetVersionRequest) GetVersionResponse {
	return h.response
}

type getVersionResponse struct {
	version version.Version
}

func (res getVersionResponse) Version() version.Version {
	return res.version
}

type GetVersion interface {
	GetVersion(GetVersionRequest) GetVersionResponse
}

type GetVersionRequest interface{}

type GetVersionResponse interface {
	Version() version.Version
}
