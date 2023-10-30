package main

import (
	"context"
	"fmt"
	"github.com/Dimoonevs/orders-api.git/application"
	"os"
	"os/signal"
)

func main() {
	app := application.New()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	err := app.Start(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
