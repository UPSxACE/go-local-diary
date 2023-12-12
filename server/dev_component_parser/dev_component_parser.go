package dev_component_parser

import (
	"encoding/json"
	"io"
	"os"
)

type DevComponentParser struct {
	Data []Category
}

func Init() *DevComponentParser {
	devCompPars := DevComponentParser{}
	devCompPars.ParseJsonConfigFile()
	return &devCompPars
}

func (devComponentParser *DevComponentParser) ParseJsonConfigFile() *DevComponentParser {
	jsonFile, err := os.Open("./server/dev-components.json")
	
	if err != nil {
		devComponentParser.Data = []Category{};
		return devComponentParser;
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	
	var parsedData []Category

	err = json.Unmarshal([]byte(byteValue), &parsedData)
	if err != nil {
		devComponentParser.Data = []Category{};
		return devComponentParser;
	}

	devComponentParser.Data = parsedData
	return devComponentParser;
}



