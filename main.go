package main

import (
	"context"
	"fmt"
	"github.com/Dimoonevs/orders-api.git/application"
)

func main() {
	app := application.New()
	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
}
