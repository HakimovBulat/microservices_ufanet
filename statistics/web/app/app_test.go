package app

import (
	"context"
	"net/http"
	. "statistics/models"
	"testing"
	"time"
)

func BenchmarkIndex(b *testing.B) {
	req, _ := http.NewRequest("GET", "/", nil)
	app := Setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response, err := app.Test(req)
		if err != nil {
			b.Fatal()
		}
		if response.StatusCode != 200 {
			b.Fatal("Expected status code 200, got: ", response.StatusCode)
		}
	}
}

func BenchmarkSaleList(b *testing.B) {
	req, _ := http.NewRequest("GET", "/api/sale", nil)
	app := Setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response, err := app.Test(req)
		if err != nil {
			b.Fatal()
		}
		if response.StatusCode != 200 {
			b.Fatal("Expected status code 200, got: ", response.StatusCode)
		}
	}
}

func BenchmarkCategoryList(b *testing.B) {
	req, _ := http.NewRequest("GET", "/api/category", nil)
	app := Setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response, err := app.Test(req)
		if err != nil {
			b.Fatal()
		}
		if response.StatusCode != 200 {
			b.Fatal("Expected status code 200, got: ", response.StatusCode)
		}
	}
}

func BenchmarkSaleOne(b *testing.B) {
	req, _ := http.NewRequest("GET", "/api/sale/1", nil)
	app := Setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response, err := app.Test(req)
		if err != nil {
			b.Fatal()
		}
		if response.StatusCode != 200 {
			b.Fatal("Expected status code 200, got: ", response.StatusCode)
		}
	}
}

func BenchmarkCategoryOne(b *testing.B) {
	req, _ := http.NewRequest("GET", "/api/category/1", nil)
	app := Setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response, err := app.Test(req)
		if err != nil {
			b.Fatal()
		}
		if response.StatusCode != 200 {
			b.Fatal("Expected status code 200, got: ", response.StatusCode)
		}
	}
}
func BenchmarkCategoryStatistics(b *testing.B) {
	req, _ := http.NewRequest("GET", "/category", nil)
	app := Setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response, err := app.Test(req)
		if err != nil {
			b.Fatal()
		}
		if response.StatusCode != 200 {
			b.Fatal("Expected status code 200, got: ", response.StatusCode)
		}
	}
}

func BenchmarkSaleStatistics(b *testing.B) {
	req, _ := http.NewRequest("GET", "/sale", nil)
	app := Setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response, err := app.Test(req)
		if err != nil {
			b.Fatal()
		}
		if response.StatusCode != 200 {
			b.Fatal("Expected status code 200, got: ", response.StatusCode)
		}
	}
}

func BenchmarkCategoryTopic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(ctxBackground, 50*time.Millisecond)
		defer cancel()
		getMessagesTopic[Category]("wal_listener.public_billboard_category", ctx)
	}
}
func BenchmarkSaleTopic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(ctxBackground, 50*time.Millisecond)
		defer cancel()
		getMessagesTopic[Sale]("wal_listener.public_billboard_sale", ctx)
	}
}
