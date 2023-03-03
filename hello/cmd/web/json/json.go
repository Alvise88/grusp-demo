package json

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Envelope map[string]any

func Write(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)

	if err != nil {
		return fmt.Errorf("unable to convert to json: %w", err)
	}

	return nil
}
