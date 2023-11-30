package user

import (
	"encoding/json"
	"fmt"
	"golang-api-starter/internal/helper"
	"log"

	"github.com/golang-jwt/jwt/v4"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
)

type UserClaims struct {
	UserId    int64  `json:"userId"`
	Username  string `json:"username"`
	TokenType string `json:"tokenType"`
	jwt.StandardClaims
}

type User struct {
	Id        *int64                 `json:"id"   db:"id" example:"2"`
	Name      string                 `json:"name" db:"name" example:"emma"`
	Password  *string                `json:"password,omitempty" db:"password" example:"password"`
	FirstName *string                `json:"firstName" db:"first_name" example:"Emma"`
	LastName  *string                `json:"lastName" db:"last_name" example:"Watson"`
	Disabled  bool                   `json:"disabled" db:"disabled" example:"false"`
	CreatedAt *helper.CustomDatetime `db:"created_at" json:"createdAt"`
	UpdatedAt *helper.CustomDatetime `db:"updated_at" json:"updatedAt"`
}

type Users []*User

func (users Users) StructToMap() []map[string]interface{} {
	mapsResults := []map[string]interface{}{}
	for _, user := range users {
		tmp := map[string]interface{}{}
		result := map[string]interface{}{}
		data, _ := json.Marshal(user)
		json.Unmarshal(data, &tmp)
		for k, v := range tmp {
			result[strcase.ToSnake(k)] = v
		}
		mapsResults = append(mapsResults, result)
	}

	return mapsResults
}

func (users Users) rowsToStruct(rows *sqlx.Rows) []*User {
	defer rows.Close()

	records := make([]*User, 0)
	for rows.Next() {
		var user User
		err := rows.StructScan(&user)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		records = append(records, &user)
	}

	return records
}

func (users *Users) printValue() {
	for _, v := range *users {
		if v.Id != nil {
			fmt.Printf("existing --> id: %+v, v: %+v\n", *v.Id, *v)
		}
		fmt.Printf("new --> v: %+v\n", *v)
	}
}
