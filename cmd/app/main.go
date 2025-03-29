package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/IgorAleksandroff/delivery/cmd"
	httpin "github.com/IgorAleksandroff/delivery/internal/adapters/in/http"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres/courierrepo"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres/orderrepo"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
	servers "github.com/IgorAleksandroff/delivery/pkg/servers"
)

func main() {
	httpPort := goDotEnvVariable("HTTP_PORT")
	dbHost := goDotEnvVariable("DB_HOST")
	dbPort := goDotEnvVariable("DB_PORT")
	dbUser := goDotEnvVariable("DB_USER")
	dbPassword := goDotEnvVariable("DB_PASSWORD")
	dbDbName := goDotEnvVariable("DB_DBNAME")
	dbSslMode := goDotEnvVariable("DB_SSLMODE")
	geoServiceGrpcHost := goDotEnvVariable("GEO_SERVICE_GRPC_HOST")

	connectionString, err := makeConnectionString(dbHost, dbPort, dbUser, dbPassword, dbDbName, dbSslMode)
	if err != nil {
		log.Fatal(err.Error())
	}

	crateDbIfNotExists(dbHost, dbPort, dbUser, dbPassword, dbDbName, dbSslMode)
	gormDb := mustGormOpen(connectionString)
	mustAutoMigrate(gormDb)

	compositionRoot := cmd.NewCompositionRoot(gormDb, geoServiceGrpcHost)

	startCron(compositionRoot)
	startWebServer(compositionRoot, httpPort)
}

func startCron(compositionRoot cmd.CompositionRoot) {
	c := cron.New()
	_, err := c.AddFunc("@every 1s", compositionRoot.Jobs.AssignOrdersJob.Run)
	if err != nil {
		log.Fatalf("ошибка при добавлении задачи: %v", err)
	}
	_, err = c.AddFunc("@every 2s", compositionRoot.Jobs.MoveCouriersJob.Run)
	if err != nil {
		log.Fatalf("ошибка при добавлении задачи: %v", err)
	}
	c.Start()
}

func startWebServer(compositionRoot cmd.CompositionRoot, port string) {
	handlers, err := httpin.NewServer(
		compositionRoot.CommandHandlers.CreateOrderCommandHandler,
		compositionRoot.QueryHandlers.GetAllCouriersQueryHandler,
		compositionRoot.QueryHandlers.GetNotCompletedOrdersQueryHandler,
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации HTTP Server: %v", err)
	}

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	e.Pre(middleware.RemoveTrailingSlash())
	registerSwaggerOpenApi(e)
	registerSwaggerUi(e)
	servers.RegisterHandlers(e, handlers)
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func registerSwaggerOpenApi(e *echo.Echo) {
	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := servers.GetSwagger()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to load swagger: "+err.Error())
		}

		data, err := swagger.MarshalJSON()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to marshal swagger: "+err.Error())
		}

		return c.Blob(http.StatusOK, "application/json", data)
	})
}

func registerSwaggerUi(e *echo.Echo) {
	e.GET("/docs", func(c echo.Context) error {
		html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		  <meta charset="UTF-8">
		  <title>Swagger UI</title>
		  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
		</head>
		<body>
		  <div id="swagger-ui"></div>
		  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
		  <script>
			window.onload = () => {
			  SwaggerUIBundle({
				url: "/openapi.json",
				dom_id: "#swagger-ui",
			  });
			};
		  </script>
		</body>
		</html>`
		return c.HTML(http.StatusOK, html)
	})
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func makeConnectionString(host string, port string, user string,
	password string, dbName string, sslMode string) (string, error) {
	if host == "" {
		return "", errs.NewValueIsRequiredError(host)
	}
	if port == "" {
		return "", errs.NewValueIsRequiredError(port)
	}
	if user == "" {
		return "", errs.NewValueIsRequiredError(user)
	}
	if password == "" {
		return "", errs.NewValueIsRequiredError(password)
	}
	if dbName == "" {
		return "", errs.NewValueIsRequiredError(dbName)
	}
	if sslMode == "" {
		return "", errs.NewValueIsRequiredError(sslMode)
	}
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		host,
		port,
		user,
		password,
		dbName,
		sslMode), nil
}

func mustGormOpen(connectionString string) *gorm.DB {
	pgGorm, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN:                  connectionString,
			PreferSimpleProtocol: true,
		},
	), &gorm.Config{})
	if err != nil {
		log.Fatalf("connection to postgres through gorm\n: %s", err)
	}
	return pgGorm
}

func mustAutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&courierrepo.TransportDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	err = db.AutoMigrate(&orderrepo.OrderDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	err = db.AutoMigrate(&courierrepo.CourierDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
}

func crateDbIfNotExists(host string, port string, user string,
	password string, dbName string, sslMode string) {
	dsn, err := makeConnectionString(host, port, user, password, "postgres", sslMode)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer db.Close()

	// Создаём базу данных, если её нет
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		log.Printf("Ошибка создания БД (возможно, уже существует): %v", err)
	}
}
