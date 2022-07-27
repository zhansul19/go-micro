package main

import (
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
	response := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("user %s logged in", request.Email),
		Data:    user,
	}
	c.writeJson(wr, http.StatusAccepted, response)
}
