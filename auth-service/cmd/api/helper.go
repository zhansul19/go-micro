package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (c *Config) readJson(wr http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := int64(1048576) //1MB

	r.Body = http.MaxBytesReader(wr, r.Body, maxBytes)

	decode := json.NewDecoder(r.Body)
	err := decode.Decode(data)
	if err != nil {
		return err
	}

	err = decode.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body cannot have more than 1 JSON value")
	}
	return nil
}
func (c *Config) writeJson(wr http.ResponseWriter, code int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for i, v := range headers[0] {
			wr.Header()[i] = v
		}
	}
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(code)
	_, err = wr.Write(out)
	if err != nil {
		return err
	}

	return nil
}
func (c *Config) errorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return c.writeJson(w, statusCode, payload)
}
