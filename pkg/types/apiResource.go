package types

type APIResource struct {
	Name         string   `json:"name"`
	SingularName string   `json:"singularName"`
	Kind         string   `json:"kind"`
	Verbs        []string `json:"verbs"`
	ShortNames   []string `json:"shortNames,omitempty"`
}

type APIResourceList struct {
	Kind         string        `json:"kind"`
	APIVersion   string        `json:"apiVersion"`
	GroupVersion string        `json:"groupVersion"`
	Resources    []APIResource `json:"resources"`
}
