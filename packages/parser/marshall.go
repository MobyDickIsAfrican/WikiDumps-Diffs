package parser

import (
	"bug/m/submodules/schema"
	"log"
)

type DataParser struct {
	data       *ParserData
	hasher     func(text string) string
	marshaller func(data []byte, v any) error
}

func NewParser(hsr func(txt string) string, marshaller func(data []byte, v any) error) *DataParser {
	var mp map[string]interface{}
	return &DataParser{data: &ParserData{
		ParsedContent:   &schema.DatabaseTable{},
		parsedInterface: mp,
	}, hasher: hsr, marshaller: marshaller}
}

func (p *DataParser) extractName() Parser {
	ifc := p.data.GetParsedInterface()
	name := (ifc)["name"].(string)
	p.data.GetParsedContent().Name = name
	return p
}

func (p *DataParser) extractContentHash() Parser {
	ifc := p.data.GetParsedInterface()
	contentHash := (ifc)["article_body"].(map[string]interface{})["html"].(string)
	p.data.GetParsedContent().ContentHash = p.hasher(contentHash)
	return p
}

func (p *DataParser) extractDateModified() Parser {
	ifc := p.data.GetParsedInterface()
	dateModified := (ifc)["date_modified"].(string)
	p.data.GetParsedContent().DateModified = dateModified
	return p
}

func (p *DataParser) extractIdentifier() Parser {
	ifc := p.data.GetParsedInterface()
	identifier := (ifc)["identifier"].(float64)
	p.data.GetParsedContent().Identifier = identifier
	return p
}

func (p *DataParser) extractURL() Parser {
	ifc := p.data.GetParsedInterface()
	url := (ifc)["url"].(string)
	p.data.GetParsedContent().URL = url
	return p
}

func (p *DataParser) extractVersionIdentifier() Parser {
	ifc := p.data.GetParsedInterface()
	version := (ifc)["version"].(map[string]interface{})
	identifier := version["identifier"].(float64)
	p.data.GetParsedContent().VersionIdentifier = identifier
	return p
}

func (p *DataParser) Parse(data []byte) Parser {
	err := p.marshaller(data, &p.data.parsedInterface)
	if err != nil {
		log.Fatal("Error parsing JSON data: ", err)
	}
	p.extractName().extractContentHash().extractDateModified().extractIdentifier().extractURL().extractVersionIdentifier()
	return p
}

func (p *DataParser) GetContent() *schema.DatabaseTable {
	return p.data.GetParsedContent()
}
