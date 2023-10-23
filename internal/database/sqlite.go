package database

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	*ConnectionInfo
	TableName string
}

func (m *Sqlite) Connect() *sqlx.DB {
	fmt.Printf("connecting to Sqlite... \n")
	fmt.Printf("Table: %+v\n", m.TableName)
	dbFile := fmt.Sprintf("%s.db", m.Database)
	connectionString := fmt.Sprintf("./%s?_auth&_auth_user=%s&_auth_pass=%s&_auth_crypt=sha1&parseTime=true", dbFile, m.User, m.Pass)
	fmt.Printf("ConnString: %+v\n", connectionString)
	// os.Remove(dbFile)

	db, err := sqlx.Open("sqlite3", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	return db
}

func (m *Sqlite) GetConnectionInfo() ConnectionInfo {
	return *m.ConnectionInfo
}

func (m *Sqlite) Select(queries map[string]interface{}) *sqlx.Rows {
	fmt.Printf("select from Sqlite, table: %+v\n", m.TableName)
	db := m.Connect()
	defer db.Close()

	selectStmt := fmt.Sprintf(
		"SELECT * FROM %s",
		m.TableName,
	)

	fmt.Printf("queries: %+v, len: %+v\n", queries, len(queries))
	if len(queries) != 0 { // add where clause
		whereClauses := []string{}
		for k, v := range queries {
			fmt.Printf("%+v: %+v(%T)\n", k, v, v)
			switch v.(type) {
			case []string:
				whereClauses = append(whereClauses, fmt.Sprintf("%s IN ('%s')", k, strings.Join(v.([]string), "','")))
			default:
				whereClauses = append(whereClauses, fmt.Sprintf("%s='%s'", k, v))
			}
		}

		selectStmt = fmt.Sprintf("%s WHERE %s", selectStmt, strings.Join(whereClauses, " AND "))
	}

	fmt.Printf("selectStmt: %+v\n", selectStmt)
	rows, err := db.Queryx(selectStmt)
	if err != nil {
		log.Printf("Queryx err: %+v\n", err.Error())
	}
	err = rows.Err()
	if err != nil {
		log.Printf("rows.Err(): %+v\n", err.Error())
	}

	return rows
}

func (m *Sqlite) Save(records Records) *sqlx.Rows {
	fmt.Printf("save from Sqlite, table: %+v\n", m.TableName)
	// fmt.Printf("records: %+v\n", records)
	db := m.Connect()
	defer db.Close()
	selectStmt := fmt.Sprintf("select * from %s", m.TableName)

	rows, err := db.Queryx(selectStmt)
	defer rows.Close()
	if err != nil {
		log.Printf("%+v\n", err)
	}

	cols, err := rows.Columns()
	if err != nil {
		log.Printf("%+v\n", err)
	}

	// fmt.Printf("cols: %+v\n", cols)
	colWithColon := []string{}
	for _, col := range cols {
		colWithColon = append(colWithColon, fmt.Sprintf(":%s", col))
	}

	insertStmt := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		m.TableName,
		fmt.Sprintf(strings.Join(cols[:], ",")),
		fmt.Sprintf(strings.Join(colWithColon[:], ",")),
	)
	fmt.Printf("%+v \n", insertStmt)

	insertedIds := []string{}
	mapsResults := records.StructToMap()
	for _, record := range mapsResults {
		sqlResult, err := db.NamedExec(insertStmt, record)
		if err != nil {
			log.Printf("insert error: %+v\n", err)
		}
		id, _ := sqlResult.LastInsertId()
		insertedIds = append(insertedIds, strconv.Itoa(int(id)))
	}

	fmt.Printf("insertedIds: %+v\n", insertedIds)
	return m.Select(map[string]interface{}{"id": insertedIds})
}

func (m *Sqlite) Update() {
	fmt.Printf("update from Sqlite, table: %+v\n", m.TableName)
}

func (m *Sqlite) Delete() {
	fmt.Printf("delete from Sqlite, table: %+v\n", m.TableName)
}
