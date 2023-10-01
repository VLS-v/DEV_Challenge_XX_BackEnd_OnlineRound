package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"spreadsheets/controllers"
	"spreadsheets/utils/saves"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const SPREADSHEETS_FILE_NAME = "saves.txt"

func main() {
	e := echo.New()

	saves := saves.Saves{}

	err := saves.Open(SPREADSHEETS_FILE_NAME)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer saves.SavesFile.Close()

	err = saves.Load()
	if err != nil {
		e.Logger.Fatal(err)
	}

	c := controllers.New(saves)
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/api/v1/:sheet_id/:cell_id", c.SetCellValue)
	e.GET("/api/v1/:sheet_id/", c.GetSheet)
	e.GET("/api/v1/:sheet_id/:cell_id", c.GetCell)

	/*
	 running the server in a separate routine allows us
	 to catch the server shutdown signal later and shut down the server correctly.
	*/
	go func() {
		e.Logger.Info(e.Start(":8080"))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("Server stopped")
}
