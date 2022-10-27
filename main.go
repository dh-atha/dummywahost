package main

import (
	"strings"
	"time"

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
			{
				"token": "eyJhbGciOiAiSFMyNTYiLCAidHlwIjogIkpXVCJ9.eyJ1c2VyIjoiT0NBXzFfd2EiLCJpYXQiOjE2NjUzNjk2NDksImV4cCI6MTY2NTk3NDQ0OSwid2E6cmFuZCI6IjA0YmVkNzc4NjVjYTlhOWU3Y2E2NGNiMzcwNGM1ZTc3In0.NgFf8Vg-PitIY_9lfWzFuUCwS15x4aZzfaQ257IOy80",
				"expires_after": time.Now().Add(120 * time.Hour).Format(time.RFC3339),
			},
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
		"contacts": request.ToResponse(),
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

func (cr *ContactRequest) ToResponse() (contactResponse []ContactResponse) {
	for _, contact := range cr.Contacts {
		if len(contact) < 3 {
			continue
		}
		contactResponse = append(contactResponse, ContactResponse{
			WAID: cr.Normalize(contact),
			Input: contact,
			Status: "valid",
		})
	}
	return 
}

func (cr *ContactRequest) Normalize(input string) string {
	if input[:3] == "+62" {
		return input[1:]
	} else if string(input[0]) == "0" {
		return "62" + input[1:]
	} else if string(input[0]) == "8" {
		return "62" + input
	}
	return input
}

type ContactResponse struct {
	WAID string `json:"wa_id"`
	Input string `json:"input"`
	Status string `json:"status"`
}