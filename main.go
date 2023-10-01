package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"sinabyr/seyid/lib"
	"sinabyr/seyid/types"
)

func homeHandlerFilms(w http.ResponseWriter, r *http.Request) {
	db := lib.InitDatabase()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM films;")
	if err != nil {
		panic(err)
	}

	var Films []types.Film
	for rows.Next() {
		var id int64
		var title string
		var director string
		var released_at time.Time

		err := rows.Scan(&id, &title, &director, &released_at)
		if err != nil {
			panic(err)
		}

		Films = append(Films, types.Film{Title: title, Director: director, ReleasedAt: released_at})
	}

	data := map[string][]types.Film{ "films": Films } 

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.InitDatabase()
	defer db.Close()

	rows, err := db.Query(`
		SELECT r.id, r.description, r.amount, r.datetime, u.nickname, u.avatar, c.title, c.icon
		FROM receipts r, categories c, users u
		WHERE r.category_id = c.id AND r.user_id = u.id
	`)
	if err != nil {
		panic(err)
	}

	var receipts []types.Receipt
	for rows.Next() {
		var id int64
		var description string
		var amount int64
		var datetime time.Time
		var nickname string
		var avatar string
		var categoryTitle string
		var categoryIcon string

		err := rows.Scan(
			&id,
			&description,
			&amount,
			&datetime,
			&nickname,
			&avatar,
			&categoryTitle,
			&categoryIcon,
		)
		if err != nil {
			panic(err)
		}

		receipts = append(receipts, types.Receipt{
			Id: id,
			Description: description,
			Amount: amount,
			Datetime: datetime,
			Nickname: nickname,
			Avatar: avatar,
			CategoryTitle: categoryTitle,
			CategoryIcon: categoryIcon,
		})
	}

	data := map[string][]types.Receipt{ "receipts": receipts }

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.InitDatabase()
	defer db.Close()

	rows, err := db.Query(`
		SELECT *
		FROM categories
	`)
	if err != nil {
		panic(err)
	}

	var categories []types.Category
	for rows.Next() {
		var id int64
		var title string
		var icon string

		err := rows.Scan(&id, &title, &icon)
		if err != nil {
			panic(err)
		}

		categories = append(categories, types.Category{
			Id: id,
			Title: title,
			Icon: icon,
		})
	}
	fmt.Println(categories)

	data := map[string][]types.Category{ "categories": {{Id: 1, Title: "hello", Icon: "world"}} }

	tmpl := template.Must(template.ParseFiles("templates/categories.html"))
	tmpl.Execute(w, data)
}

func main() {
	// serve "./public" directory contents under "/static" url path
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
	// http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8000", nil)
}

