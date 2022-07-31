package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zhansul19/go-micro/broker/event"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
	case "log":
		c.logItemViaRabbit(wr, requestPayload.Log)
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
func (c *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.errorJson(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		c.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		c.errorJson(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	c.writeJson(w, http.StatusAccepted, payload)

}
func (c *Config) logItemViaRabbit(wr http.ResponseWriter, l LogPayload) {
	err := c.pushToQueue(l.Name, l.Data)
	if err != nil {
		c.errorJson(wr, err)
	}
	response := jsonResponse{
		Error:   false,
		Message: "logged via rabbit",
	}

	c.writeJson(wr, http.StatusAccepted, response)
}
func (c *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(c.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
