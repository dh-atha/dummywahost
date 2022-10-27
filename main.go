package main

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LoginHost(c echo.Context) error {
	auth := c.Request().Header["Authorization"]
	if auth == nil {
		return c.JSON(401, "missing basic auth")
	}
	return c.JSON(200, map[string]interface{}{
		"users": []map[string]string{
			{"token": "eyJhbGciOiAiSFMyNTYiLCAidHlwIjogIkpXVCJ9.eyJ1c2VyIjoiT0NBXzFfd2EiLCJpYXQiOjE2NjUzNjk2NDksImV4cCI6MTY2NTk3NDQ0OSwid2E6cmFuZCI6IjA0YmVkNzc4NjVjYTlhOWU3Y2E2NGNiMzcwNGM1ZTc3In0.NgFf8Vg-PitIY_9lfWzFuUCwS15x4aZzfaQ257IOy80",
			"expires_after": "2022-10-31 02:40:49+00:00"},
		},
		"meta": map[string]string{
			"version": "v2.41.2",
			"api_status": "stable",
		},
	})
} 

func Message(c echo.Context) error {
	auth := c.Request().Header["Authorization"]
	if auth == nil || strings.ToUpper(strings.Split(auth[0], " ")[0]) != "BEARER" {
		return c.JSON(401, "missing bearer auth")
	}
	return c.JSON(200, map[string]interface{}{
		"messages": []map[string]string{
			{"id": "gBGGFlBxZ0M_AgkSjW7mD4HxMUw"},
		},
		"meta": map[string]string{
			"api_status": "stable",
			"version": "v2.41.2",
		},
	})
}

func Contact(c echo.Context) error {
	auth := c.Request().Header["Authorization"]
	if auth == nil || strings.ToUpper(strings.Split(auth[0], " ")[0]) != "BEARER" {
		return c.JSON(401, "missing bearer auth")
	}
	var request ContactRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(500, "error parsing request body")
	}

	return c.JSON(200, map[string]interface{}{
		"contacts": []map[string]string{
			{
				"wa_id": "6282213604539",
				"input": request.Contacts[0],
				"status": "valid",
			  },
		},
		"meta": map[string]string{
			"api_status": "stable",
			"version": "v2.41.2",
		},
	})
}

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())

	e.POST("/v1/users/login", LoginHost)
	e.POST("/v1/messages", Message)
	e.POST("/v1/contacts", Contact)

	e.Logger.Fatal(e.Start(":3000"))
}

type ContactRequest struct {
	Blocking string `json:"blocking"`
	Contacts []string `json:"contacts"`
	ForceCheck bool `json:"force_check"`
}