package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var redisClient *redis.Client = nil
var ctx = context.Background()

type redisObject struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type randomObj struct {
	Timestamp int64 `json:"timestamp"`
	RandomNum int64 `json:"random_num"`
}

func main() {
	app := fiber.New()
	app.Get("/", health)
	app.Post("/redis", put)
	app.Get("/redis", get)
	app.Get("/postman-get", postmanGet)
	app.Get("/random", random)
	app.Listen(":8090")
}

func health(c *fiber.Ctx) error {
	_, err := c.WriteString("Healthy")
	return err
}

func put(c *fiber.Ctx) error {
	input := new(redisObject)
	if err := c.BodyParser(input); err != nil {
		return err
	}
	_, err := getRedisClient().Set(ctx, input.Key, input.Value, 0).Result()
	if err != nil {
		c.Status(500).SendString(err.Error())
	}
	return c.JSON(input)
}

func get(c *fiber.Ctx) error {
	key := c.Query("key")
	if key == "" {
		return c.Status(400).SendString("Missing redis key")
	}
	value, err := getRedisClient().Get(ctx, key).Result()
	if err != nil {
		return c.Status(404).SendString("Not found")
	}
	return c.JSON(redisObject{
		Key:   key,
		Value: value,
	})
}

func postmanGet(c *fiber.Ctx) error {
	agent := fiber.Get("https://postman-echo.com/get?test=123")
	status, body, errs := agent.Bytes()
	if len(errs) != 0 {
		return c.Status(500).SendString(errs[0].Error())
	}
	c.Response().Header.Set("Content-Type", "application/json")
	return c.Status(status).Send(body)
}

func random(c *fiber.Ctx) error {
	timestamp := time.Now().Unix()
	rand.Seed(time.Now().UnixMilli())
	return c.JSON(randomObj{
		Timestamp: timestamp,
		RandomNum: rand.Int63(),
	})
}

func getRedisClient() *redis.Client {
	if redisClient == nil {
		opts := redis.Options{
			Addr:     "redis-master.default.svc.cluster.local:6379",
			Password: "",
			DB:       0,
		}
		redisClient = redis.NewClient(&opts)
	}
	return redisClient
}
