/*
Copyright © 2025 rnikrozoft rnikrozoft.dev@gmail.com
*/
package cmd

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	_ "github.com/rnikrozoft/pramool.in.th-backend/docs"
	"github.com/rnikrozoft/pramool.in.th-backend/handler"
	"github.com/rnikrozoft/pramool.in.th-backend/repository"
	"github.com/rnikrozoft/pramool.in.th-backend/service"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// @title           Pramool Backend API
// @version         1.0
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3001

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func serve() {
	app := fiber.New()
	app.Use(cors.New())

	validate := validator.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger")
	})
	app.Get("/swagger/*", swagger.HandlerDefault)

	userRepository := repository.NewUserRepository(conn)
	authenticationService := service.NewAuthenticationService(appConfigs, userRepository)

	registerRepository := repository.NewRegisterRepository(conn)
	registerService := service.NewRegisterService(registerRepository, authenticationService)
	registerHandler := handler.NewRegisterHandler(validate, registerService)

	app.Post("/register", registerHandler.Register)

	app.Listen(":3001")
}
