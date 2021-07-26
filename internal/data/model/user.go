package model

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type User struct {
	Id        string    `db:"id" json:"id"`
	CreatedOn time.Time `db:"created_on" json:"created_on"`
	Credentials
}

type Credentials struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"pass"`
}

func (o *Credentials) UnmarshalJSON(data []byte) error {
	type AliasCredentials Credentials
	credentialsOverwritten := &struct {
		*AliasCredentials
	}{
		AliasCredentials: (*AliasCredentials)(o),
	}

	if err := json.Unmarshal(data, &credentialsOverwritten); err != nil {
		return err
	}

	if o.Username == "" || strings.Contains(o.Username, " ") ||
		o.Password == "" || strings.Contains(o.Password, " ") {
		return errors.New("Please provide valid credentials")
	}

	return nil
}
