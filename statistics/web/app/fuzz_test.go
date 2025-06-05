package app

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func FuzzApp(f *testing.F) {
	example_json := map[string]string{
		"category_search": "сертификаты",
	}
	body, err := json.Marshal(&example_json)
	if err != nil {
		log.Println(err)
	}
	f.Add([]byte(body))

	f.Fuzz(func(t *testing.T, data []byte) {
		app := Setup()
		req := httptest.NewRequest("POST", "/category", bytes.NewBuffer(data))

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}

		if resp.StatusCode == fiber.StatusBadRequest && json.Valid(data) {
			t.Errorf("Expected valid JSON but got Bad Request for input: %s", data)
		}

	})
}
