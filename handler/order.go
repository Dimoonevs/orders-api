package handler

import (
	"fmt"
	"net/http"
)

type Order struct {
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Order")
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List Order")
}

func (o *Order) GetById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetById Order")
}

func (o *Order) UpdateById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateById Order")
}

func (o *Order) DeleteById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteById Order")
}
