/*
Copyright © 2025 rnikrozoft rnikrozoft.dev@gmail.com
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	_ "github.com/rnikrozoft/pramool-core/docs"
	"github.com/rnikrozoft/pramool-core/handler"
	"github.com/rnikrozoft/pramool-core/internal/telemetry"
	"github.com/rnikrozoft/pramool-core/middleware"
	"github.com/rnikrozoft/pramool-core/repository"
	"github.com/rnikrozoft/pramool-core/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
	logger, err := telemetry.NewZapLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	shutdownTelemetry, err := telemetry.Init("pramool-core")
	if err != nil {
		logger.Fatal("otel init", zap.Error(err))
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTelemetry(ctx); err != nil {
			logger.Warn("otel shutdown", zap.Error(err))
		}
	}()

	app := fiber.New(fiber.Config{
		AppName: "pramool-core",
	})
	app.Use(otelfiber.Middleware())
	app.Use(telemetry.AccessLogWithZap(logger))
	telemetry.MountHealth(app, "pramool-core")

	corsOrigins := strings.TrimSpace(appConfigs.CorsAllowOrigins)
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowOriginsFunc: corsAllowDevLAN,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Cookie",
	}))

	validate := validator.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger")
	})
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Static("/uploads", "./uploads")

	userRepository := repository.NewUserRepository(conn)
	bankRepository := repository.NewBankRepository(conn)
	userService := service.NewUserService(userRepository)
	bankService := service.NewBankService(bankRepository)

	authenticationService := service.NewAuthenticationService(appConfigs, userService)
	accessH := appConfigs.Jwt.ExpireTime
	if accessH <= 0 {
		accessH = 1
	}
	accessCookieSec := accessH * 3600
	refreshH := appConfigs.Jwt.RefreshExpireTime
	if refreshH <= 0 {
		refreshH = 168
	}
	refreshCookieSec := refreshH * 3600

	authenticationHandler := handler.NewAuthenticationHandler(validate, authenticationService, accessCookieSec, refreshCookieSec)

	registerRepository := repository.NewRegisterRepository(conn)
	registerService := service.NewRegisterService(registerRepository, userService)
	registerHandler := handler.NewRegisterHandler(validate, authenticationService, registerService, userService, accessCookieSec, refreshCookieSec)

	otpHTTP := telemetry.DefaultHTTPClient(45 * time.Second)
	otpService := service.NewOTPService(
		logger,
		appConfigs.ThaiBulkSMS.AddressRequest,
		appConfigs.ThaiBulkSMS.AddressVerify,
		appConfigs.ThaiBulkSMS.APIKey,
		appConfigs.ThaiBulkSMS.APISecret,
		otpHTTP,
	)
	otpHandler := handler.NewOTPHandler(validate, otpService, registerService, userService)

	userhandler := handler.NewUserHandler(validate, userService)
	bankHandler := handler.NewBankHandler(bankService)
	m := middleware.Middleware{JWTSecret: appConfigs.Jwt.Secret}

	user := app.Group("/users")
	user.Get("/", m.JWTMiddleware, userhandler.GetMyInformation)
	user.Get("/profile", m.JWTMiddleware, userhandler.GetMyInformation)
	user.Get("/onboarding-status", m.JWTMiddleware, userhandler.GetOnboardingStatus)
	user.Put("/profile", m.JWTMiddleware, userhandler.UpdateProfile)
	user.Get("/:tel", userhandler.IsTelAlreadyUsed)
	user.Post("/", registerHandler.Register)
	app.Get("/banks", bankHandler.List)

	app.Post("/otp/request", otpHandler.RequestOTP)
	app.Post("/otp/verify", otpHandler.VerifyOTP)
	app.Post("/otp/timeout", otpHandler.RecordTimeout)

	app.Post("/login/tel", authenticationHandler.LoginByTel)
	app.Post("/auth/signup", registerHandler.Signup)
	app.Post("/auth/refresh", authenticationHandler.Refresh)
	app.Post("/logout", authenticationHandler.Logout)

	const listenAddr = ":3001"
	go func() {
		if err := app.Listen(listenAddr); err != nil {
			logger.Error("listen stopped", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error("fiber shutdown", zap.Error(err))
	}
}
