package response

import "gorm.io/datatypes"

type SecureImages struct {
	Dockers datatypes.JSON `json:"dockers"`
}

type SecureImageReport struct {
	Sbom datatypes.JSON `json:"sbom,omitempty"`
	Vul  datatypes.JSON `json:"vul,omitempty"`
}
