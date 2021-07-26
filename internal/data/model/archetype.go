package model

import (
	"database/sql"
	"time"
)

type Archetype struct {
	Id        int          `db:"id" json:"id"`
	Title     string       `db:"title" json:"title"`
	UrlImage  string       `db:"url_image" json:"url_image"`
	CreatedOn time.Time    `db:"created_on" json:"created_on"`
	UpdatedOn sql.NullTime `db:"updated_on" json:"updated_on"`
}
