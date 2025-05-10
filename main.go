package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rnikrozoft/pramool.in.th-backend/handler"
	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
	"github.com/rnikrozoft/pramool.in.th-backend/repository"
	"github.com/rnikrozoft/pramool.in.th-backend/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var gormDB *gorm.DB

func init() {
	dsn := "host=db user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	gormDB = db

	gormDB.AutoMigrate(&entity.User{})
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	validate := validator.New()

	registerRepository := repository.NewRegisterRepository(gormDB)
	registerService := service.NewRegisterService(registerRepository)
	registerHandler := handler.NewRegisterHandler(validate, registerService)

	app.Post("/register", registerHandler.Register)

	app.Listen(":3001")
}
