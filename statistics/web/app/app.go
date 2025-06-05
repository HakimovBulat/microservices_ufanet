package app

import (
	"log"
	"runtime"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup() *fiber.App {
	pc, _, _, _ := runtime.Caller(1)
	engine := html.New("./web/templates", ".html")
	if runtime.FuncForPC(pc).Name() != "main.main" {
		engine = html.New("templates", ".html")
	}

	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", IndexHandler)
	app.Get("/category", CategoryStatistics)
	app.Post("/sale", PostSaleStatistics)
	app.Post("/category", PostCategoryStatistics)
	app.Get("/sale", SaleStatistics)

	app.Post("/api/sale", PostSale)
	app.Post("/api/category", PostCategory)

	app.Get("/api/category/:id", GetCategory)
	app.Get("/api/sale/:id", GetSale)

	app.Get("/api/sale", GetSaleList)
	app.Get("/api/category", GetCategoryList)

	app.Delete("/api/sale/:id", DeleteSale)
	app.Delete("/api/category/:id", DeleteCategory)
	return app
}
