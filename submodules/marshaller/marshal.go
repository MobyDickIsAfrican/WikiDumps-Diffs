package marshaller

import (
	"encoding/json"
)

func JSONMarshaller(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	return err
}
