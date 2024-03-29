package main

import (
	"context"
	"fmt"
	"github.com/alvarezcarlos/payment/app/application"
	"github.com/alvarezcarlos/payment/app/config"
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/domain/repository"
	"github.com/alvarezcarlos/payment/app/infrastructure/postgres/connection"
	"github.com/alvarezcarlos/payment/app/interface/rest"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

func main() {
	config.Environment()

	file := setLogger()
	defer file.Close()

	//DBConnection
	db := connection.NewPostgresConnection(&gorm.Config{}, slog.Default())
	conn := db.GetConnection()
	migrator := connection.NewMigrate(conn, slog.Default())
	initDBMigrations(migrator)
	//Repositories
	merchantRepo := repository.NewMerchantRepository(conn)
	paymentRepo := repository.NewPaymentRepository(conn)
	//UseCases
	merchantUseCase := application.NewMerchantUseCase(merchantRepo, slog.Default())
	paymentUseCase := application.NewPaymentUseCase(paymentRepo, slog.Default())
	e := echo.New()

	//Controllers
	rest.NewMerchantController(e, merchantUseCase)
	rest.NewPaymentController(e, paymentUseCase)

	go startServer(e)
	gracefulShutdown(e)
}

func setLogger() *os.File {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	w := io.MultiWriter(file, os.Stderr)

	replace := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == "source" {
			src := a.Value.Any().(*slog.Source)
			return slog.String("source", fmt.Sprintf("method: %s, line: %d", src.Function, src.Line))
		}
		if a.Key == "time" {
			t := a.Value.Time()
			return slog.String("time", t.Format(time.DateTime))
		}
		return a
	}

	var level slog.Level
	if config.Config().Environment == "local" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	slogOptions := slog.HandlerOptions{
		AddSource:   true,
		Level:       level,
		ReplaceAttr: replace,
	}

	logger := slog.New(slog.NewJSONHandler(w, &slogOptions))
	slog.SetDefault(logger)
	return file
}

func startServer(e *echo.Echo) {
	if err := e.Start(fmt.Sprintf(":%s", config.Config().Port)); err != nil {
		e.Logger.Info("shutting down the server")
	}
}

func gracefulShutdown(e *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func initDBMigrations(migrator connection.MigrateInterface) {
	tables := []interface{}{
		&entity.Customer{},
		&entity.Merchant{},
		&entity.Payment{},
		&entity.Card{},
		//&entity.PaymentState{},
	}
	migrator.AutoMigrateAll(tables...)
}
