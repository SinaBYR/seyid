package types

import (
	"net/http"
	"regexp"
	"time"
)

type Category struct {
	Id int64
	Title string
	Icon string
}

type UserAccount struct {
	Id int64
	Name string
	Avatar string
}

type Receipt struct {
	Id							int64
	Description 		string
	Amount					int64
	Datetime				time.Time
	Nickname				string
	Avatar					string
	CategoryTitle		string
	CategoryIcon		string
}

type Film struct {
	// Id int64
	Title string
	Director string
	ReleasedAt time.Time
}

type Route struct {
	Method string
	Regex *regexp.Regexp
	Handler http.HandlerFunc
}

