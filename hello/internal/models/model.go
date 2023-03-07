package models

type Info struct {
	Namespace string

	Container string
	Pod       string
	Node      string
	Message   string
}

type TemplateData struct {
	Info
	CurrentYear int
}
