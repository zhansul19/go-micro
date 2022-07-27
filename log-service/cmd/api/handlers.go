package main

import (
	"net/http"

	"github.com/zhansul19/log-service/database"
)

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Config) Writelog(wr http.ResponseWriter, r *http.Request) {
	var logPayload LogPayload

	err := c.readJson(wr, r, &logPayload)
	if err != nil {
		c.errorJson(wr, err)
		return
	}

	event:=database.LogEntry{
		Name: logPayload.Name,
		Data: logPayload.Data,
	}

	err=c.Models.LogEntry.Insert(event)
	if err != nil {
		c.errorJson(wr,err)
		return
	}
	resp:=jsonResponse{
		Error: false,
		Message: "logged",
	}

	c.writeJson(wr,http.StatusOK,resp)

}
