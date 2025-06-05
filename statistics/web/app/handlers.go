package app

import (
	"context"
	"encoding/json"
	"log"
	"slices"
	"sort"
	. "statistics/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/segmentio/kafka-go"
)

var ctxBackground = context.Background()

func GetCategory(c fiber.Ctx) error {
	var category Category
	db.Table("billboard_category").First(&category, c.Params("id"))
	return c.JSON(category)
}

func GetCategoryList(c fiber.Ctx) error {
	var categories []Category
	db.Table("billboard_category").Find(&categories)
	return c.JSON(categories)
}

func GetSale(c fiber.Ctx) error {
	var sale Sale
	db.Table("billboard_sale").First(&sale, c.Params("id"))
	return c.JSON(sale)
}

func DeleteSale(c fiber.Ctx) error {
	var sale Sale
	err := db.Table("billboard_sale").Delete(&sale, c.Params("id")).Error
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON("sale was deleted")
}

func DeleteCategory(c fiber.Ctx) error {
	var category Category
	err := db.Table("billboard_category").Delete(&category, c.Params("id")).Error
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON("category was deleted")
}

func PostSale(c fiber.Ctx) error {
	var sale Sale
	if err := json.Unmarshal(c.Body(), &sale); err != nil {
		return c.JSON(err)
	}
	if sale.Photo == "" {
		sale.Photo = "static/media/category-new_MROs2uS.png"
	}
	db.Table("billboard_sale").Create(&sale)
	return c.JSON("ok")
}

func CheckCategory(data []byte) (Category, error) {
	var category Category
	if err := json.Unmarshal(data, &category); err != nil {
		return Category{}, err
	}
	if category.Photo == "" {
		category.Photo = "static/media/category-new_MROs2uS.png"
	}
	return category, nil
}

func PostCategory(c fiber.Ctx) error {
	category, err := CheckCategory(c.Body())
	if err != nil {
		return c.JSON(err)
	}
	db.Table("billboard_category").Create(&category)
	return c.JSON("ok")
}

func GetSaleList(c fiber.Ctx) error {
	var sales []Sale
	db.Table("billboard_sale").Find(&sales)
	return c.JSON(sales)
}

func IndexHandler(c fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func SaleStatistics(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(ctxBackground, 50*time.Millisecond)
	defer cancel()

	saleMessages, popularitySalesStruct := getMessagesTopic[Sale]("wal_listener.public_billboard_sale", ctx)

	return c.Render("sale", fiber.Map{
		"saleMessages":          saleMessages,
		"popularitySalesStruct": popularitySalesStruct,
	})
}

func PostCategoryStatistics(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(ctxBackground, 50*time.Millisecond)
	defer cancel()

	categoryMessages, popularityCategoriesStruct := getMessagesTopic[Category]("wal_listener.public_billboard_category", ctx)

	mapResponse := make(map[string]string)
	c.Bind().Body(&mapResponse)
	search := strings.ToLower(mapResponse["category_search"])

	var searchCategories []Category
	db.Table("billboard_category").Where("LOWER(title) LIKE ?", "%"+search+"%").Find(&searchCategories)

	return c.Render("category", fiber.Map{
		"searchWord":                 search,
		"searchCategories":           searchCategories,
		"categoryMessages":           categoryMessages,
		"popularityCategoriesStruct": popularityCategoriesStruct,
	})
}

func PostSaleStatistics(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(ctxBackground, 50*time.Millisecond)
	defer cancel()

	saleMessages, popularitySalesStruct := getMessagesTopic[Sale]("wal_listener.public_billboard_sale", ctx)

	mapResponse := make(map[string]string)
	c.Bind().Body(&mapResponse)
	search := strings.ToLower(mapResponse["sale_search"])

	var searchSales []Sale
	db.Table("billboard_sale").Where(
		"LOWER(title) LIKE ?", "%"+search+"%",
	).Or(
		"LOWER(description) LIKE ?", "%"+search+"%",
	).Or(
		"LOWER(about_partner) LIKE ?", "%"+search+"%",
	).Or(
		"LOWER(subtitle) LIKE ?", "%"+search+"%",
	).Find(&searchSales)

	return c.Render("sale", fiber.Map{
		"searchWord":            search,
		"searchSales":           searchSales,
		"saleMessages":          saleMessages,
		"popularitySalesStruct": popularitySalesStruct,
	})
}

func CategoryStatistics(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(ctxBackground, 50*time.Millisecond)
	defer cancel()

	categoryMessages, popularityCategoriesStruct := getMessagesTopic[Category]("wal_listener.public_billboard_category", ctx)
	return c.Render("category", fiber.Map{
		"categoryMessages":           categoryMessages,
		"popularityCategoriesStruct": popularityCategoriesStruct,
	})
}

func getMessagesTopic[T SaleCategory](topic string, ctx context.Context) ([]Message[T], []Popularity[T]) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:29092"},
<<<<<<< HEAD
		MaxWait: 50 * time.Millisecond,
=======
		MaxWait: 1 * time.Second,
>>>>>>> 045509d916bc1ae38e993185a1e6f4a27a8faabe
		Topic:   topic,
	})
	defer reader.Close()

	var messages []Message[T]
	popularity := make(map[T]int)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if err.Error() != "fetching message: context deadline exceeded" {
				log.Println(err)
			}
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
		popularityStruct = append(popularityStruct, Popularity[T]{Struct: key, View: value})
	}

	sort.Slice(popularityStruct, func(i, j int) bool {
		return popularityStruct[i].View < popularityStruct[j].View
	})

	slices.Reverse(popularityStruct)
	return messages, popularityStruct
}
