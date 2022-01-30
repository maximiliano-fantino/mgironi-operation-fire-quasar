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
	Name     string   `json:"name"`
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

type TopSecretRequest struct {
	Satellites []SatelliteInfoRequest `json:"satellites"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
