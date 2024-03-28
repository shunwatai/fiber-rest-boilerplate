package user

import (
	"golang-api-starter/internal/helper"
	"golang-api-starter/internal/helper/utils"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetId(t *testing.T) {
	user := &User{
		MongoId: utils.ToPtr("xxxx-xxxx-xxxx-xxxx"),
		Id:      utils.ToPtr(helper.FlexInt(2)),
	}

	tests := []struct {
		name     string
		dbDriver string
		input    *User
		want     string
	}{
		{name: "test Id", dbDriver: "postgres-or-mariadb-or-sqlite", input: user, want: "2"},
		{name: "test MongoId", dbDriver: "mongodb", input: user, want: "xxxx-xxxx-xxxx-xxxx"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cfg.LoadEnvVariables()
			cfg.Vpr.Set("database.engine", testCase.dbDriver)
			if err := cfg.Vpr.Unmarshal(cfg); err != nil {
				t.Logf("failed loading conf, err: %+v\n", err.Error())
			}

			got := user.GetId()
			eq := reflect.DeepEqual(testCase.want, got)

			if !eq {
				t.Errorf("got %q want %q", got, testCase.want)
			}
		})
	}
}

func TestStructToMap(t *testing.T) {
	var id int64 = 2
	now := time.Now()
	customDatetime := &helper.CustomDatetime{&now, utils.ToPtr(time.RFC3339)}
	timeStr, _ := customDatetime.MarshalJSON()
	timeJson := strings.Replace(string(timeStr), "\"", "", -1)
	users := Users{
		&User{
			MongoId:   utils.ToPtr("xxxx-xxxx-xxxx-xxxx"),
			Id:        utils.ToPtr(helper.FlexInt(id)),
			CreatedAt: customDatetime,
			UpdatedAt: customDatetime,
		},
	}

	tests := []struct {
		name  string
		input Users
		want  []map[string]interface{}
	}{
		{name: "test StructToMap", input: users, want: []map[string]interface{}{
			{"_id": "xxxx-xxxx-xxxx-xxxx", "id": float64(2), "created_at": timeJson, "updated_at": timeJson, "first_name": nil, "last_name": nil, "disabled": false, "name": "", "is_oauth": false, "provider": nil},
		}},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := users.StructToMap()
			eq := reflect.DeepEqual(testCase.want, got)

			if !eq {
				t.Errorf("got %q want %q", got, testCase.want)
			}
		})
	}
}

func TestGetTags(t *testing.T) {
	users := Users{
		&User{},
	}

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "test get db tags", input: "db", want: []string{"id", "name", "password", "email", "first_name", "last_name", "disabled", "is_oauth", "provider", "created_at", "updated_at"}},
		{name: "test get bson tags", input: "bson", want: []string{"_id", "id", "name", "password", "email", "first_name", "last_name", "disabled", "is_oauth", "provider", "created_at", "updated_at"}},
		{name: "test get json tags", input: "json", want: []string{"_id", "id", "name", "password", "email", "firstName", "lastName", "disabled", "isOauth", "provider", "createdAt", "updatedAt"}},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := users.GetTags(testCase.input)
			eq := reflect.DeepEqual(testCase.want, got)

			if !eq {
				t.Errorf("got %q want %q", got, testCase.want)
			}
		})
	}
}