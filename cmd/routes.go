package cmd

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"sinabyr/seyid/lib"
	"sinabyr/seyid/types"
	"strconv"
	"time"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := lib.InitDatabase()
	defer db.Close()

	rows, err := db.Query(`
		SELECT *
		FROM users;
	`)
	if err != nil {
		panic(err)
	}

	var users []types.UserAccount
	for rows.Next() {
		var id int64
		var nickname string
		var avatar sql.NullString

		err := rows.Scan(&id, &nickname, &avatar)
		if err != nil {
			panic(err)
		}

		var avatarString string
		if avatar.Valid {
			avatarString = avatar.String
		} else {
			avatarString = ""
		}

		users = append(users, types.UserAccount{
			Id: id,
			Nickname: nickname,
			Avatar: avatarString,
		})
	}

	// TODO use Avatar
	htmlStr := ""
	for _, c := range users {
		htmlStr += fmt.Sprintf(
			`<option value="%d">%s</option>%s`, c.Id, c.Nickname, "\n",
		)
	}

	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	db := lib.InitDatabase()
	defer db.Close()

	rows, err := db.Query(`
		SELECT *
		FROM categories;
	`)
	if err != nil {
		panic(err)
	}

	var categories []types.Category
	for rows.Next() {
		var id int64
		var title string
		var icon sql.NullString

		err := rows.Scan(&id, &title, &icon)
		if err != nil {
			panic(err)
		}

		var iconString string
		if icon.Valid {
			iconString = icon.String
		} else {
			iconString = ""
		}

		categories = append(categories, types.Category{
			Id: id,
			Title: title,
			Icon: iconString,
		})
	}

	// TODO use Icon
	htmlStr := ""
	for _, c := range categories {
		htmlStr += fmt.Sprintf(
			`<option value="%d">%s</option>%s`, c.Id, c.Title, "\n",
		)
	}

	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func CategoriesPageHandler(w http.ResponseWriter, r *http.Request) {
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
		var icon sql.NullString

		err := rows.Scan(&id, &title, &icon)
		if err != nil {
			panic(err)
		}

		var iconString string
		if icon.Valid {
			iconString = icon.String
		} else {
			iconString = ""
		}

		categories = append(categories, types.Category{
			Id: id,
			Title: title,
			Icon: iconString,
		})
	}

	data := map[string][]types.Category{ "categories": {{Id: 1, Title: "hello", Icon: "world"}} }

	tmpl := template.Must(template.ParseFiles("templates/categories.html"))
	tmpl.Execute(w, data)
}

func CreateReceiptHandler(w http.ResponseWriter, r *http.Request) {
	description := r.PostFormValue("description")
	amount, _ := strconv.ParseInt(r.PostFormValue("amount"), 10, 64) // convert to int64
	datetimeEpoch, _ := strconv.ParseInt(r.PostFormValue("datetimeEpoch"), 10, 64) // convert to int64
	userId, _ := strconv.ParseInt(r.PostFormValue("userId"), 10, 64) // convert to int64
	categoryId, _ := strconv.ParseInt(r.PostFormValue("categoryId"), 10, 64) // convert to int64
	fmt.Printf("%d\n", datetimeEpoch)

	db := lib.InitDatabase()
	defer db.Close()

	row := db.QueryRow(`
	WITH inserted_receipt as (
		INSERT INTO receipts (description, amount, datetime, user_id, category_id)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING *
	) SELECT r.id, r.description, r.amount, r.datetime, u.avatar, c.title, c.icon
		FROM inserted_receipt r, users u, categories c
		WHERE r.user_id = u.id AND r.category_id = c.id
	`, description, amount, time.Unix(datetimeEpoch, 0), userId, categoryId)
	fmt.Printf("row: %v\n", row)

	var rowId int64
	var rowDescription string
	var rowAmount int64
	var rowDatetime time.Time
	// var rowNickname string
	var rowAvatar sql.NullString
	var rowCategoryTitle string
	var rowCategoryIcon sql.NullString

	err := row.Scan(
		&rowId,
		&rowDescription,
		&rowAmount,
		&rowDatetime,
		// &rowNickname,
		&rowAvatar,
		&rowCategoryTitle,
		&rowCategoryIcon,
	)
	if err != nil {
		panic(err)
	}

	var rowCategoryIconString string
	var rowAvatarString string

	if rowCategoryIcon.Valid {
		rowCategoryIconString = rowCategoryIcon.String
	} else {
		rowCategoryIconString = ""
	}

	if rowAvatar.Valid {
		rowAvatarString = rowAvatar.String
	} else {
		rowAvatarString = ""
	}

	htmlStr := fmt.Sprintf(`
		<tr>
			<td class="p-4 text-white text-left">%s</td>
			<td class="p-4 text-white text-left">%d</td>
			<td class="p-4 text-white text-left">%s</td>
			<td class="p-4 text-white text-left">%s</td>
		</tr>
	`, rowDescription, rowAmount, rowCategoryIconString, rowAvatarString)

	tmpl, err := template.New("t").Parse(htmlStr)
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, nil)
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
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
		var avatar sql.NullString
		var categoryTitle string
		var categoryIcon sql.NullString

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

		var avatarString string
		var categoryIconString string

		if avatar.Valid {
			avatarString = avatar.String
		} else {
			avatarString = ""
		}

		if categoryIcon.Valid {
			categoryIconString = categoryIcon.String
		} else {
			categoryIconString = ""
		}

		receipts = append(receipts, types.Receipt{
			Id: id,
			Description: description,
			Amount: amount,
			Datetime: datetime,
			Nickname: nickname,
			Avatar: avatarString,
			CategoryTitle: categoryTitle,
			CategoryIcon: categoryIconString,
		})
	}

	data := map[string][]types.Receipt{ "receipts": receipts }

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

