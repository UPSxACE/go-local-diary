package dev_component_parser

type Category struct {
	Name       string       `json:"categoryName" form:"categoryName" query:"categoryName"`
	Components []Components `json:"components" form:"components" query:"components"`
}

type Components struct {
	Name         string     `json:"componentName" form:"componentName" query:"componentName"`
	Description  string     `json:"componentDescription" form:"componentDescription" query:"componentDescription"`
	TemplateName string     `json:"templateName" form:"templateName" query:"templateName"`
	Examples     []Examples `json:"examples" form:"examples" query:"examples"`
}

type Examples struct {
	Title       string      `json:"exampleTitle" form:"exampleTitle" query:"exampleTitle"`
	Description string      `json:"exampleDescription" form:"exampleDescription" query:"exampleDescription"`
	Data        interface{} `json:"data" form:"data" query:"data"`
}