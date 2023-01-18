package request

type Image struct {
	Repo string `json:"repo,omitempty" form:"repo"`
	Tag  string `json:"tag,omitempty" form:"tag"`
}

type SecureImage struct {
	Category string   `json:"category,omitempty" form:"category"`
	Images   []*Image `json:"images,omitempty" form:"images"`
}

type SecureBare struct {
	Category string   `json:"category,omitempty"`
	Paths    []string `json:"paths,omitempty"`
}

type SecureFix struct {
	CveId string `json:"cveId,omitempty"`
}
