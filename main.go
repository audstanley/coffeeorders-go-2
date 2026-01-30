package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoffeeOrder struct {
	ID           primitive.ObjectID `json:"_id"`
	Coffee       string             `json:"coffee,omitempty"`
	EmailAddress string             `json:"emailAddress,omitempty"`
	Flavor       string             `json:"flavor,omitempty"`
	Strength     uint8              `json:"strength,omitempty"`
}

var db *leveldb.DB

func main() {
	var err error
	db, err = leveldb.OpenFile("coffeeorders.db", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()

	app.Get("/", root)
	app.Get("/coffeeorders", listAll)
	app.Get("/coffeeorders/:email", listByEmail)
	app.Post("/coffeeorders", createOrder)
	app.Delete("/coffeeorders/:email", deleteNewest)
	app.Delete("/coffeeorders", deleteAll)
	app.Put("/coffeeorders/:email", replaceNewest)

	go func() {
		for {
			now := time.Now()

			// If it's the 1st of the month at 00:00–00:59
			if now.Day() == 1 && now.Hour() == 0 {
				fmt.Println("Monthly cleanup: deleting all coffee orders")

				iter := db.NewIterator(nil, nil)
				for iter.Next() {
					key := string(iter.Key())
					if strings.HasPrefix(key, "coffeeorder:") {
						db.Delete(iter.Key(), nil)
					}
				}
				iter.Release()

				// Sleep for an hour so we don't run multiple times in the same hour
				time.Sleep(time.Hour)
			}

			// Check once per minute
			time.Sleep(time.Minute)
		}
	}()
	log.Fatal(app.Listen(":3300"))
}

func root(c *fiber.Ctx) error {
	return c.SendString("Try /coffeeorders")
}

func timeUntilDeletion() string {
	now := time.Now()
	year, month, _ := now.Date()
	loc := now.Location()

	// First day of next month at midnight
	nextMonth := month + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}

	next := time.Date(nextYear, nextMonth, 1, 0, 0, 0, 0, loc)
	return next.Sub(now).String()
}

func listAll(c *fiber.Ctx) error {
	orders := readAll()
	return c.JSON(fiber.Map{
		"data":              orders,
		"timeUntilDeletion": timeUntilDeletion(),
	})
}

func listByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	orders := readByEmail(email)
	return c.JSON(fiber.Map{
		"data":              orders,
		"timeUntilDeletion": timeUntilDeletion(),
	})
}

func createOrder(c *fiber.Ctx) error {
	var co CoffeeOrder
	if err := c.BodyParser(&co); err != nil {
		return c.Status(400).JSON(fiber.Map{"err": "Malformed JSON"})
	}

	co.ID = primitive.NewObjectID()
	key := fmt.Sprintf("coffeeorder:%s", co.ID.Hex())

	data, _ := json.Marshal(co)
	db.Put([]byte(key), data, nil)

	return c.JSON(readAll())
}

func deleteNewest(c *fiber.Ctx) error {
	email := c.Params("email")
	key := newestKey(email)
	if key != "" {
		db.Delete([]byte(key), nil)
	}
	return c.JSON(readAll())
}

func deleteAll(c *fiber.Ctx) error {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := string(iter.Key())
		if strings.HasPrefix(key, "coffeeorder:") {
			db.Delete(iter.Key(), nil)
		}
	}
	iter.Release()
	return c.JSON([]string{})
}

func replaceNewest(c *fiber.Ctx) error {
	email := c.Params("email")

	var co CoffeeOrder
	if err := c.BodyParser(&co); err != nil {
		return c.Status(400).JSON(fiber.Map{"err": "Malformed JSON"})
	}

	// delete newest
	key := newestKey(email)
	if key != "" {
		db.Delete([]byte(key), nil)
	}

	// insert new
	co.ID = primitive.NewObjectID()
	newKey := fmt.Sprintf("coffeeorder:%s", co.ID.Hex())
	data, _ := json.Marshal(co)
	db.Put([]byte(newKey), data, nil)

	return c.JSON(readAll())
}

func readAll() []CoffeeOrder {
	iter := db.NewIterator(nil, nil)
	var orders []CoffeeOrder

	for iter.Next() {
		key := string(iter.Key())
		if strings.HasPrefix(key, "coffeeorder:") {
			var co CoffeeOrder
			json.Unmarshal(iter.Value(), &co)
			orders = append(orders, co)
		}
	}
	iter.Release()
	return orders
}

func readByEmail(email string) []CoffeeOrder {
	iter := db.NewIterator(nil, nil)
	var orders []CoffeeOrder

	for iter.Next() {
		key := string(iter.Key())
		if strings.HasPrefix(key, "coffeeorder:") {
			var co CoffeeOrder
			json.Unmarshal(iter.Value(), &co)
			if co.EmailAddress == email {
				orders = append(orders, co)
			}
		}
	}
	iter.Release()
	return orders
}

func newestKey(email string) string {
	iter := db.NewIterator(nil, nil)
	var newest string

	for iter.Next() {
		key := string(iter.Key())
		if strings.HasPrefix(key, "coffeeorder:") {
			var co CoffeeOrder
			json.Unmarshal(iter.Value(), &co)
			if co.EmailAddress == email {
				newest = key // last one wins
			}
		}
	}
	iter.Release()
	return newest
}
