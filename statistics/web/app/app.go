package main

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type SaleCategory interface {
	Sale | Category
}

type Message[T SaleCategory] struct {
	Data       T         `json:"data"`
	Id         uuid.UUID `json:"id"`
	Schema     string    `json:"schema"`
	Action     string    `json:"action"`
	DataOld    T         `json:"dataOld"`
	CommitTime string    `json:"commitTime"`
	Table      string    `json:"table"`
}

var ctxBackground = context.Background()

type Sale struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	Photo        string `json:"photo"`
	AboutPartner string `json:"about_partner"`
	Promocode    string `json:"promocode"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Url          string `json:"url"`
}

type Category struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Photo string `json:"photo"`
}

func main() {
	engine := html.New("./web/templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", indexHandler)
	app.Get("/category", categoryHandler)
	app.Get("/sale", saleHandler)
	// app.Get("/sale/:id", getSaleHandler)
	log.Fatal(app.Listen(":3000"))
}

// func getSaleHandler(c fiber.Ctx) error {
// 	c.Params("id")
// 	return c.Render()
// }

func indexHandler(c fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

type Popularity[T SaleCategory] struct {
	Struct T
	View   int
}

func getMessagesTopic[T SaleCategory](topic string, ctx context.Context) ([]Message[T], []Popularity[T]) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   topic,
	})
	defer reader.Close()

	var messages []Message[T]
	popularity := make(map[T]int)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Println(err)
			break
		}

		var message Message[T]
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			log.Println(err)
		}

		if message.Action == "SELECT" {
			popularity[message.Data]++
		}

		messages = append(messages, message)
	}

	var popularityStruct []Popularity[T]
	for key, value := range popularity {
		popularityStruct = append(popularityStruct, Popularity[T]{key, value})
	}

	sort.Slice(popularityStruct, func(i, j int) bool {
		return popularityStruct[i].View < popularityStruct[j].View
	})

	return messages, popularityStruct
}

func saleHandler(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(ctxBackground, 5*time.Second)
	defer cancel()

	saleMessages, popularitySalesStruct := getMessagesTopic[Sale]("wal_listener.public_billboard_sale", ctx)

	return c.Render("sale", fiber.Map{
		"saleMessages":          saleMessages,
		"popularitySalesStruct": popularitySalesStruct,
	})
}

func categoryHandler(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(ctxBackground, 5*time.Second)
	defer cancel()

	categoryMessages, popularityCategoriesStruct := getMessagesTopic[Sale]("wal_listener.public_billboard_sale", ctx)

	return c.Render("category", fiber.Map{
		"categoryMessages":           categoryMessages,
		"popularityCategoriesStruct": popularityCategoriesStruct,
	})
}
