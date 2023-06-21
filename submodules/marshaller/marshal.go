package marshaller

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
)

func JSONMarshaller(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	return err
}

func XMLMarshaller(data []byte, v any) error {
	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	// Iterate over the XML tags.
	for {
		// Read the next token.
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Check if the token is a StartElement.
		if se, ok := t.(xml.StartElement); ok {
			// Get the name of the tag.
			name := se.Name.Local
			log.Printf("Token: %#v\n", t)
			log.Println(name)
		}
	}
	return nil
}
