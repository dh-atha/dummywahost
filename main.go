package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LoginHost(c echo.Context) error {
	auth := c.Request().Header["Authorization"]
	if auth == nil {
		return c.JSON(401, "missing basic auth")
	}
	expired := time.Now().Add(120 * time.Hour).Format(time.RFC3339)
	expired = strings.Join(strings.Split(expired, "T"), " ")
	return c.JSON(200, map[string]interface{}{
		"users": []map[string]string{
			{
				"token":         "eyJhbGciOiAiSFMyNTYiLCAidHlwIjogIkpXVCJ9.eyJ1c2VyIjoiT0NBXzFfd2EiLCJpYXQiOjE2NjUzNjk2NDksImV4cCI6MTY2NTk3NDQ0OSwid2E6cmFuZCI6IjA0YmVkNzc4NjVjYTlhOWU3Y2E2NGNiMzcwNGM1ZTc3In0.NgFf8Vg-PitIY_9lfWzFuUCwS15x4aZzfaQ257IOy80",
				"expires_after": expired,
			},
		},
		"meta": map[string]string{
			"version":    "v2.41.2",
			"api_status": "stable",
		},
	})
}

func LoginBRIHost(c echo.Context) error {
	username, password, ok := c.Request().BasicAuth()
	if !ok {
		return c.JSON(401, "error getting basic auth header")
	}

	if username != "admin" || password != "Jatis865*" {
		return c.JSON(401, "invalid credential")
	}

	expired := time.Now().Add(120 * time.Hour).Format(time.RFC3339)
	expired = strings.Join(strings.Split(expired, "T"), " ")
	return c.JSON(200, map[string]interface{}{
		"users": []map[string]string{
			{
				"token":         "eyJhbGciOiAiSFMyNTYiLCAidHlwIjogIkpXVCJ9.eyJ1c2VyIjoiT0NBXzFfd2EiLCJpYXQiOjE2NjUzNjk2NDksImV4cCI6MTY2NTk3NDQ0OSwid2E6cmFuZCI6IjA0YmVkNzc4NjVjYTlhOWU3Y2E2NGNiMzcwNGM1ZTc3In0.NgFf8Vg-PitIY_9lfWzFuUCwS15x4aZzfaQ257IOy80",
				"expires_after": expired,
			},
		},
		"meta": map[string]string{
			"version":    "v2.41.2",
			"api_status": "stable",
		},
	})
}

func MessageSuccess(c echo.Context) error {
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
			"version":    "v2.41.2",
		},
	})
}

func MessageError(c echo.Context) error {
	auth := c.Request().Header["Authorization"]
	if auth == nil || strings.ToUpper(strings.Split(auth[0], " ")[0]) != "BEARER" {
		return c.JSON(401, "missing bearer auth")
	}
	return c.JSON(200, map[string]interface{}{
		"errors": []map[string]interface{}{
			{
				"code":   1000,
				"title":  "error send message",
				"detail": "message error",
			},
		},
		"meta": map[string]string{
			"api_status": "stable",
			"version":    "v2.41.2",
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
			"version":    "v2.41.2",
		},
		"errors": []map[string]interface{}{},
	})
}

func ContactError(c echo.Context) error {
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
			"version":    "v2.41.2",
		},
		"errors": []map[string]interface{}{
			{
				"code":   1000,
				"title":  "error check contact",
				"detail": "contact error",
			},
		},
	})
}

var i int
var start time.Time
var lock sync.Mutex

func Counter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			lock.Lock()
			if start.Equal(time.Time{}) {
				start = time.Now()
			}
			i++
			lock.Unlock()
			if i == 150 {
				log.Println(i)
				fmt.Println(time.Since(start))
			}
			return next(c)
		}
	}
}

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	// e.Use(middleware.Logger())

	// e.POST("/v1/users/login", LoginHost)
	e.POST("/v1/users/login", LoginBRIHost)
	// e.POST("/v1/messages", Message, Counter())
	e.POST("/v1/messages", MessageSuccess, Counter())
	e.POST("/v1/contacts", Contact)
	// e.POST("/v1/messages", MessageError)

	e.Logger.Fatal(e.Start(":3000"))
}

type ContactRequest struct {
	Blocking   string   `json:"blocking"`
	Contacts   []string `json:"contacts"`
	ForceCheck bool     `json:"force_check"`
}

func (cr *ContactRequest) ToResponse() (contactResponse []ContactResponse) {
	for _, contact := range cr.Contacts {
		if len(contact) < 3 {
			continue
		}
		contactResponse = append(contactResponse, ContactResponse{
			WAID:   cr.Normalize(contact),
			Input:  contact,
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
	WAID   string `json:"wa_id"`
	Input  string `json:"input"`
	Status string `json:"status"`
}
