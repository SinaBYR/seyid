package main

import (
	"net/http"

	_ "github.com/lib/pq"

	"sinabyr/seyid/cmd"
)

func main() {
	// serve "./public" directory contents under "/static" url path
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/extra/users", cmd.GetUsers)
	http.HandleFunc("/extra/categories", cmd.GetCategories)
	http.HandleFunc("/categories", cmd.CategoriesPageHandler)
	http.HandleFunc("/createReceipt", cmd.CreateReceiptHandler)
	http.HandleFunc("/", cmd.HomePageHandler)
	http.ListenAndServe(":8000", nil)
}

