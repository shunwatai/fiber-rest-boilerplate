package database

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MariaDb struct {
	*ConnectionInfo
	TableName string
}

func (m *MariaDb) Connect() *sqlx.DB {
	fmt.Printf("connecting to MariaDb... \n")
	fmt.Printf("Table: %+v\n", m.TableName)
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", *m.User, *m.Pass, *m.Host, *m.Port, *m.Database)
	fmt.Printf("ConnString: %+v\n", connectionString)
	// os.Remove(dbFile)

	db, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	return db
}

func (m *MariaDb) Select(queries map[string]interface{}) *sqlx.Rows {
	fmt.Printf("select from MariaDB, table: %+v\n", m.TableName)
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

func (m *MariaDb) Save(records Records) *sqlx.Rows {
	fmt.Printf("save from MariaDB, table: %+v\n", m.TableName)
	// fmt.Printf("records: %+v\n", records)
	db := m.Connect()
	defer db.Close()
	selectStmt := fmt.Sprintf("select * from %s limit 1;", m.TableName)

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
	var colWithColon, colUpdateSet []string
	for _, col := range cols {
		// use in SQL's VALUES()
		colWithColon = append(colWithColon, fmt.Sprintf(":%s", col))

		// use in SQL's ON DUPLICATE KEY UPDATE
		if strings.Contains(col, "_at") {
			colUpdateSet = append(colUpdateSet, fmt.Sprintf("%s=IFNULL(VALUES(%s), CURRENT_TIMESTAMP)", col, col))
			continue
		}
		colUpdateSet = append(colUpdateSet, fmt.Sprintf("%s=VALUES(%s)", col, col))
	}

	insertStmt := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s) 
		ON DUPLICATE KEY UPDATE
    %s
		RETURNING id;`,
		m.TableName,
		fmt.Sprintf(strings.Join(cols[:], ",")),
		fmt.Sprintf(strings.Join(colWithColon[:], ",")),
		fmt.Sprintf(strings.Join(colUpdateSet[:], ",\n")),
	)
	fmt.Printf("%+v \n", insertStmt)

	insertedIds := []string{}
	sqlResult, err := db.NamedQuery(insertStmt, records)
	if err != nil {
		log.Printf("insert error: %+v\n", err)
	}
	// fmt.Printf("sqlResult: %+v\n", sqlResult)

	for sqlResult.Next() {
		var id string
		err := sqlResult.Scan(&id)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		insertedIds = append(insertedIds, id)
	}

	fmt.Printf("insertedIds: %+v\n", insertedIds)
	return m.Select(map[string]interface{}{"id": insertedIds})
}

// func (m *MariaDb) Update() {
// 	fmt.Printf("update from MariaDB, table: %+v\n", m.TableName)
// }
func (m *MariaDb) Delete(ids *[]int64) error {
	fmt.Printf("delete from MariaDB, table: %+v\n", m.TableName)
	db := m.Connect()

	deleteStmt, args, err := sqlx.In(
		fmt.Sprintf("DELETE FROM %s WHERE id IN (?);", m.TableName),
		*ids,
	)
	if err != nil {
		log.Printf("sqlx.In err: %+v\n", err.Error())
		return err
	}
	deleteStmt = db.Rebind(deleteStmt)
	fmt.Printf("stmt: %+v, args: %+v\n", deleteStmt, args)

	_, err = db.Exec(deleteStmt, args...)
	if err != nil {
		log.Printf("Delete Query err: %+v\n", err.Error())
		return err
	}

	return nil
}
