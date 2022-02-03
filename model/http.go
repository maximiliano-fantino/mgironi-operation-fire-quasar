package model

type CoordinatesResponse struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type TopSecretResponse struct {
	Position CoordinatesResponse `json:"position"`
	Message  string              `json:"message"`
}

type SatelliteInfoRequest struct {
	Name     string   `json:"name" example:"kenobi" redis:"name"`
	Distance float32  `json:"distance" example:"100.23" redis:"distance"`
	Message  []string `json:"message" example:",is,a,,message" redis:"message"`
}

type Dataset struct {
	Satellites []SatelliteInfoRequest
	Key        string
	Operation  string
}

type TopSecretRequest struct {
	Satellites []SatelliteInfoRequest `json:"satellites"`
}

type TopSecretSplitRequest struct {
	*SatelliteInfoRequest
}

type TopSecretSplitPOSTResponse struct {
	Operation string `json:"operation"`
}

type ErrorResponse struct {
	Message string `json:"error" example:"this is an error message description"`
}
