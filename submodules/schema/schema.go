package schema

type DatabaseTable struct {
	Name              string  `json:"name"`
	Identifier        float64 `json:"identifier"`
	VersionIdentifier float64 `json:"version_identifier"`
	URL               string  `json:"url"`
	DateModified      string  `json:"date_modified"`
	ContentHash       string  `json:"content_hash"`
}
