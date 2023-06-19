package schema

import (
	"encoding/json"
)

type DatabaseTable struct {
	Name              string      `json:"name"`
	Identifier        json.Number `json:"identifier"`
	VersionIdentifier json.Number `json:"version_identifier"`
	URL               string      `json:"url"`
	DateModified      string      `json:"date_modified"`
	ContentHash       string      `json:"content_hash"`
}
