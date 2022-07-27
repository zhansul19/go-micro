package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Config) HandleSubmission(wr http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := c.readJson(wr, r, &requestPayload)
	if err != nil {
		c.errorJson(wr, err)
		return
	}
	switch requestPayload.Action {
	case "auth":
		c.authenticate(wr, requestPayload.Auth)
	default:
		c.errorJson(wr, errors.New("action undefined"))
	}
}

func (c *Config) authenticate(wr http.ResponseWriter, auth AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(auth, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	log.Println(request)

	if err != nil {
		c.errorJson(wr, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	log.Println(response)

	if err != nil {
		c.errorJson(wr, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		c.errorJson(wr, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		c.errorJson(wr, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		c.errorJson(wr, err)
		return
	}

	if jsonFromService.Error {
		c.errorJson(wr, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	c.writeJson(wr, http.StatusAccepted, payload)
}
