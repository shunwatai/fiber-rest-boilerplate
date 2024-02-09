package database

import (
	"fmt"
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/helper"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rows interface {
	StructScan(interface{}) error
	Next() bool
	Close() error
}

type IDatabase interface {
	/* Select by raw sql */
	RawQuery(string) *sqlx.Rows

	/* Select by req querystring with pagination */
	Select(map[string]interface{}) (Rows, *helper.Pagination)

	// Insert new records, support upsert when id is present.
	// And also support batch insert/upsert
	Save(Records) (Rows, error)

	/* Delete records by ids(support batch delete) */
	Delete([]string) error
	// GetConnectionInfo() ConnectionInfo
	constructSelectStmtFromQuerystring(queries map[string]interface{}) (string, *helper.Pagination, map[string]interface{})
}

type ConnectionInfo struct {
	Driver   string
	Host     *string
	Port     *string
	User     *string
	Pass     *string
	Database *string
}

type Records interface {
	StructToMap() []map[string]interface{}
	GetTags(string) []string
}

var cfg = config.Cfg

func GetDatabase(tableName string) IDatabase {
	// log.Printf("engine: %+v\n", cfg.DbConf.Driver)

	if cfg.DbConf.Driver == "sqlite" {
		connection := cfg.DbConf.SqliteConf
		return &Sqlite{
			ConnectionInfo: &ConnectionInfo{
				Driver:   cfg.DbConf.Driver,
				Host:     connection.Host,
				Port:     connection.Port,
				User:     connection.User,
				Pass:     connection.Pass,
				Database: connection.Database,
			},
			TableName: tableName,
		}
	}

	if cfg.DbConf.Driver == "mariadb" {
		connection := cfg.DbConf.MariadbConf
		return &MariaDb{
			ConnectionInfo: &ConnectionInfo{
				Driver:   cfg.DbConf.Driver,
				Host:     connection.Host,
				Port:     connection.Port,
				User:     connection.User,
				Pass:     connection.Pass,
				Database: connection.Database,
			},
			TableName: tableName,
		}
	}

	if cfg.DbConf.Driver == "postgres" {
		connection := cfg.DbConf.PostgresConf
		return &Postgres{
			ConnectionInfo: &ConnectionInfo{
				Driver:   cfg.DbConf.Driver,
				Host:     connection.Host,
				Port:     connection.Port,
				User:     connection.User,
				Pass:     connection.Pass,
				Database: connection.Database,
			},
			TableName: tableName,
		}
	}

	if cfg.DbConf.Driver == "mongodb" {
		connection := cfg.DbConf.MongodbConf
		return &Mongodb{
			ConnectionInfo: &ConnectionInfo{
				Driver:   cfg.DbConf.Driver,
				Host:     connection.Host,
				Port:     connection.Port,
				User:     connection.User,
				Pass:     connection.Pass,
				Database: connection.Database,
			},
			TableName: tableName,
		}
	}

	return nil
}

// Generate partial sql(bindvar) for date range filter by the conditions in parameter "queries",
// the date's key will be removed in "queries" after processed.
//
// NOTE: when getting the querystring in controller, GetQueryString() will look for keys end with either "_at" or "date" in order to add the flag "withDateFilter" in queries.
//
// There are 3 accepted input, all indicated by .(DOT):
// 1. records from specified date up to now --> 2023-10-1.
//    e.g. queries --> map[string]interface{}{"created_at":"2023-10-01."} --> sql: created_at >= 2023-10-01
// 2. records before specified date --> .2023-06-30
//    e.g. queries --> map[string]interface{}{"created_at":".2023-06-30"} --> sql: created_at <= 2023-06-30
// 3. records in between 2 dates --> 2023-06-30.2023-10-1
//    e.g. queries --> map[string]interface{}{"created_at":"2023-06-30.2023-10-01"} --> sql: created_at >= 2023-06-30 AND created_at <= 2023-10-01
func getDateRangeStmt(queries, bindvarMap map[string]interface{}) string {
	// fmt.Printf("dd query: %+v\n", queries)
	if queries["withDateFilter"] == nil {
		return ""
	}
	dateRangeConditions := []string{}
	for k, v := range queries {
		if len(k) < 3 || (!strings.Contains(k[len(k)-4:], "date") && !strings.Contains(k[len(k)-3:], "_at")) {
			// fmt.Printf("not date: %+v\n", k)
			continue
		}
		splitedDates := strings.Split(v.(string), ".")
		// fmt.Printf("splitedDates? %+v, len: %+v\n", splitedDates, len(splitedDates))
		if len(splitedDates) == 2 {
			from, to := splitedDates[0], splitedDates[1]
			if from != "" {
				fromKey := fmt.Sprintf("%sFrom", k)
				dateRangeConditions = append(dateRangeConditions, fmt.Sprintf("%s >= :%s", k, fromKey))
				bindvarMap[fromKey] = from
			}
			if to != "" {
				toKey := fmt.Sprintf("%sTo", k)
				dateRangeConditions = append(dateRangeConditions, fmt.Sprintf("%s <= :%s", k, toKey))
				bindvarMap[toKey] = to
			}
		}
		delete(queries, k)
	}

	// fmt.Printf("dateConditions: %+v\n",dateConditions)
	return strings.Join(dateRangeConditions, " AND ")
}

// Generate bson for mongo find's date filtering
func getDateRangeBson(queries map[string]interface{}) bson.D {
	// fmt.Printf("dd query: %+v\n", queries)
	if queries["withDateFilter"] == nil {
		return bson.D{}
	}

	const dateFormat = "2006-01-02"
	dateRangeConditions := bson.D{}
	for k, v := range queries {
		if len(k) < 3 || (!strings.Contains(k[len(k)-4:], "date") && !strings.Contains(k[len(k)-3:], "_at")) {
			// fmt.Printf("not date: %+v\n", k)
			continue
		}
		splitedDates := strings.Split(v.(string), ".")
		// fmt.Printf("splitedDates? %+v, len: %+v\n", splitedDates, len(splitedDates))
		if len(splitedDates) == 2 {
			from, to := splitedDates[0], splitedDates[1]
			if from != "" {
				t, _ := time.Parse(dateFormat, from)
				dateRangeConditions = append(dateRangeConditions, bson.D{{
					k, bson.D{{
						"$gte", primitive.NewDateTimeFromTime(t),
					}},
				}}...)
			}
			if to != "" {
				t, _ := time.Parse(dateFormat, to)
				dateRangeConditions = append(dateRangeConditions, bson.D{{
					k, bson.D{{
						"$lte", primitive.NewDateTimeFromTime(t.AddDate(0, 0, 1)),
					}},
				}}...)
			}
		}
		delete(queries, k)
	}

	// fmt.Printf("dateConditions: %+v\n",dateConditions)
	// return strings.Join(dateRangeConditions, " AND ")
	return dateRangeConditions
}
