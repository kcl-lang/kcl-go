package inline

type TypeMeta struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
}

type ObjectMeta struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type AppInline struct {
	TypeMeta   `json:",inline,omitempty"`
	ObjectMeta `json:"metadata,omitempty"`
}

type AppNoInline struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata"`
}
