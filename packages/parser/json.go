package parser

import (
	"bug/m/submodules/schema"
	"encoding/json"
	"log"
)

type JSONParser struct {
	data   *ParserData
	hasher func(text string) string
}

func NewJSON(hsr func(txt string) string) *JSONParser {
	var mp map[string]interface{}
	return &JSONParser{data: &ParserData{
		ParsedContent:   &schema.DatabaseTable{},
		parsedInterface: mp,
	}, hasher: hsr}
}

func (p *JSONParser) extractName() Parser {
	ifc := p.data.GetParsedInterface()
	//log.Print(ifc)
	name := (ifc)["name"].(string)
	p.data.GetParsedContent().Name = name
	return p
}

func (p *JSONParser) extractContentHash() Parser {
	ifc := p.data.GetParsedInterface()
	contentHash := (ifc)["article_body"].(map[string]interface{})["html"].(string)
	p.data.GetParsedContent().ContentHash = p.hasher(contentHash)
	return p
}

func (p *JSONParser) extractDateModified() Parser {
	ifc := p.data.GetParsedInterface()
	dateModified := (ifc)["date_modified"].(string)
	p.data.GetParsedContent().DateModified = dateModified
	return p
}

func (p *JSONParser) extractIdentifier() Parser {
	ifc := p.data.GetParsedInterface()
	identifier := (ifc)["identifier"].(float64)
	p.data.GetParsedContent().Identifier = identifier
	return p
}

func (p *JSONParser) extractURL() Parser {
	ifc := p.data.GetParsedInterface()
	url := (ifc)["url"].(string)
	p.data.GetParsedContent().URL = url
	return p
}

func (p *JSONParser) extractVersionIdentifier() Parser {
	ifc := p.data.GetParsedInterface()
	version := (ifc)["version"].(map[string]interface{})
	identifier := version["identifier"].(float64)
	p.data.GetParsedContent().Identifier = identifier
	return p
}

func (p *JSONParser) Parse(data []byte) Parser {
	err := json.Unmarshal(data, &p.data.parsedInterface)
	if err != nil {
		log.Fatal("Error parsing JSON data: ", err)
	}
	p.extractName().extractContentHash().extractDateModified().extractIdentifier().extractURL().extractVersionIdentifier()
	return p
}

func (p *JSONParser) GetContent() *schema.DatabaseTable {
	return p.data.GetParsedContent()
}
