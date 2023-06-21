// Package parser parses data into a desired format.
package parser

import (
	"bug/m/submodules/hash"
	"bug/m/submodules/marshaller"
	"bug/m/submodules/schema"
)

type ParserData struct {
	parsedInterface map[string]interface{}
	ParsedContent   *schema.DatabaseTable
}

func (p *ParserData) GetParsedContent() *schema.DatabaseTable {
	return p.ParsedContent
}

func (p *ParserData) SetParsedContent(parsedContent *schema.DatabaseTable) {
	p.ParsedContent = parsedContent
}

func (p *ParserData) GetParsedInterface() map[string]interface{} {
	return p.parsedInterface
}

type Parser interface {
	extractName() Parser
	extractContentHash() Parser
	extractDateModified() Parser
	extractIdentifier() Parser
	extractURL() Parser
	extractVersionIdentifier() Parser
	Parse(data []byte) Parser
	GetContent() *schema.DatabaseTable
}

func ParseJSON(data []byte) Parser {
	var jsn Parser = NewParser(hash.Hash, marshaller.JSONMarshaller)
	prs := jsn.Parse(data)
	return prs
}

func ParseXML(data []byte) Parser {
	var jsn Parser = NewParser(hash.Hash, marshaller.JSONMarshaller)
	prs := jsn.Parse(data)
	return prs
}
