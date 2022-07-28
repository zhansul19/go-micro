package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (c *Config) Authenticate(wr http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := c.readJson(wr, r, &request)
	if err != nil {
		c.errorJson(wr, err, http.StatusBadRequest)
	}
	user, err := c.Models.User.GetByEmail(request.Email)
	if err != nil {
		c.errorJson(wr, errors.New("invalid data"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(request.Password)
	if err != nil || !valid {
		c.errorJson(wr, errors.New("invalid data"), http.StatusBadRequest)
		return
	}

	//logging
	err = c.logRequest("authenteication", fmt.Sprintf("user { %s } logged in"))
	if err != nil {
		c.errorJson(wr, errors.New("couldn't log"), http.StatusBadRequest)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("user %s logged in", request.Email),
		Data:    user,
	}
	c.writeJson(wr, http.StatusAccepted, response)
}
func (c *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}

	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
