package main

import (
	"html/template"
	"net/http"
	_ "github.com/lib/pq"

	"sinabyr/seyid/lib"
)

type Film struct {
	// Id int64
	Title string
	Director string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.InitDatabase()

	rows, err := db.Query("SELECT * FROM films;")
	if err != nil {
		panic(err)
	}

	var Films []Film
	for rows.Next() {
		var id int64
		var title string
		var director string

		err := rows.Scan(&id, &title, &director)
		if err != nil {
			panic(err)
		}

		Films = append(Films, Film{Title: title, Director: director})
	}

	data := map[string][]Film{ "films": Films } 

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, data)
}

func createReceipt(w http.ResponseWriter, r *http.Request) {

}

func main() {
	// serve "./public" directory contents under "/static" url path
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/", homeHandler)

	http.ListenAndServe(":8000", nil)
}

