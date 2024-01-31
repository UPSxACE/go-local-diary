package dev_component_parser

import (
	"reflect"
	"testing"
)

func init() {
	/* setup for tests */
	jsonPath = "./dev-components-example.json"
}

func TestParseJsonConfigFile(t *testing.T) {
	parser := &DevComponentParser{}
	output := parser.ParseJsonConfigFile().Data
	var expectedOutput []Category = []Category{
		{
			Name: "Category1",
			Components: []Components{
				{
					Name:         "Component1",
					Description:  "Description1",
					TemplateName: "Template1",
					Examples: []Examples{
						{
							Title:       "Example1",
							Description: "ExDescription1",
							Data: map[string]any{
								"key":  "val",
								"int":  float64(123),
								"bool": true,
								"nested": map[string]any{
									"key": "val",
								},
							},
						},
						{
							Title:       "Example2",
							Description: "ExDescription2",
							Data: nil,
						},
					},
				},
				{
					Name:         "Component2",
					Description:  "Description2",
					TemplateName: "Template2",
				},
			},
		},
		{
			Name:       "Category2",
			Components: []Components{},
		},
	}

	// Check if the output matches the expected result
	if !reflect.DeepEqual(expectedOutput, output) {
		t.Fatalf("Expected %#v, got %#v", expectedOutput, output)
	}
}

func TestInit(t *testing.T){
	// Test if Data field is automatically added
	// since ParseJsonConfigFile shall be executed
	// as soon as the parser is initialized
	parser := Init()
	var expectedData []Category = []Category{
		{
			Name: "Category1",
			Components: []Components{
				{
					Name:         "Component1",
					Description:  "Description1",
					TemplateName: "Template1",
					Examples: []Examples{
						{
							Title:       "Example1",
							Description: "ExDescription1",
							Data: map[string]any{
								"key":  "val",
								"int":  float64(123),
								"bool": true,
								"nested": map[string]any{
									"key": "val",
								},
							},
						},
						{
							Title:       "Example2",
							Description: "ExDescription2",
							Data: nil,
						},
					},
				},
				{
					Name:         "Component2",
					Description:  "Description2",
					TemplateName: "Template2",
				},
			},
		},
		{
			Name:       "Category2",
			Components: []Components{},
		},
	}

	// Check if the output matches the expected result
	if !reflect.DeepEqual(expectedData, parser.Data) {
		t.Fatalf("Expected %#v, got %#v", expectedData, parser.Data)
	}
}