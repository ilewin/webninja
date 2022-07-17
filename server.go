package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"webp.ninja/handlers"
	"webp.ninja/utils"
)

func main() {
	config := utils.GetConfig()
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // this is the default limit of 100MB
	})
	api := app.Group("/api/v1")
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.Cors_Origin,
		AllowHeaders: config.Cors_Header,
		AllowMethods: config.Cors_Methods,
	}))
	api.Post("/convert", handlers.ConvertHandler)

	app.Post("/submit", handlers.ContactHandler)

	app.Static("/", "./public", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 31536000 * time.Second,
		MaxAge:        31536000,
	})

	log.Fatal(app.Listen(":" + config.App_Port))
}
