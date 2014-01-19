package handlers

import (
	"github.com/felixge/quantastic/backend/version"
)

type GetVersion interface {
	GetVersion(GetVersionRequest) GetVersionResponse
}

type GetVersionRequest interface{}

type GetVersionResponse interface {
	Version() version.Version
}

func NewGetVersion(version version.Version) GetVersion {
	return getVersion{getVersionResponse{version}}
}

type getVersion struct {
	response getVersionResponse
}

func (h getVersion) GetVersion(req GetVersionRequest) GetVersionResponse {
	return h.response
}

type getVersionResponse struct {
	version version.Version
}

func (res getVersionResponse) Version() version.Version {
	return res.version
}
