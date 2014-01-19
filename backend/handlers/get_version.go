package handlers

import (
	"github.com/felixge/quantastic/backend/version"
)

func NewGetVersionHandler(version version.Version) GetVersionHandler {
	return GetVersionHandler{GetVersionResponse{version}}
}

type GetVersionHandler struct {
	response GetVersionResponse
}

func (h GetVersionHandler) GetVersion(req GetVersionRequest) GetVersionResponse {
	return h.response
}

type GetVersionRequest interface{}

type GetVersionResponse struct {
	version version.Version
}

func (res GetVersionResponse) Version() version.Version {
	return res.version
}
