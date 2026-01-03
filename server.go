package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/miceremwirigi/journey-family-sermons/m/cmd/config"
	routes "github.com/miceremwirigi/journey-family-sermons/m/pkg"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/databases"
)

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1:3000, http://127.0.0.1:5500, https://journey-family-sermons.onrender.com, https://journey-family-sermons-miceremwirigi268-ccj540ay.leapcell.dev",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, OPTIONS",
	}))

	app.Static("/static", "./templates")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	conf := config.LoadConfig()

	db := databases.StartDatabase(conf.Environment)

	routes.RegisterRoutes(app, db)

	log.Fatal(app.Listen(":3000"))
}
