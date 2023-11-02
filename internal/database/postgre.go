package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type Postgres struct {
	*ConnectionInfo
	TableName string
	db        *sqlx.DB
}

func (m *Postgres) Connect() *sqlx.DB {
	fmt.Printf("connecting to Postgres... \n")
	// fmt.Printf("Table: %+v\n", m.TableName)
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", *m.User, *m.Pass, *m.Host, *m.Port, *m.Database)
	fmt.Printf("ConnString: %+v\n", connectionString)

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	return db
}

func (m *Postgres) GetColumns(selectStmt string) []string {
	rows, err := m.db.Queryx(selectStmt)
	defer rows.Close()
	if err != nil {
		log.Printf("%+v\n", err)
	}

	cols, err := rows.Columns()
	if err != nil {
		log.Printf("%+v\n", err)
	}

	return cols
}

func (m *Postgres) Select(queries map[string]interface{}) *sqlx.Rows {
	fmt.Printf("select from Postgres, table: %+v\n", m.TableName)
	m.db = m.Connect()
	defer m.db.Close()

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
	rows, err := m.db.Queryx(selectStmt)
	if err != nil {
		log.Printf("Queryx err: %+v\n", err.Error())
	}

	if rows.Err() != nil {
		log.Printf("rows.Err(): %+v\n", err.Error())
	}

	return rows
}

func (m *Postgres) Save(records Records) *sqlx.Rows {
	fmt.Printf("save from Postgres, table: %+v\n", m.TableName)
	// fmt.Printf("records: %+v\n", records)
	m.db = m.Connect()
	defer m.db.Close()

	selectStmt := fmt.Sprintf("select * from %s limit 1;", m.TableName)
	cols := m.GetColumns(selectStmt)

	// fmt.Printf("cols: %+v\n", cols)
	var colWithColon, colUpdateSet []string
	for _, col := range cols {
		// use in SQL's VALUES()
		if col == "id" {
			colWithColon = append(colWithColon, fmt.Sprintf("COALESCE(:%s, nextval('%s_id_seq'))", col, m.TableName))
		} else if strings.Contains(col, "_at") {
			colWithColon = append(colWithColon, fmt.Sprintf("COALESCE(:%s, CURRENT_TIMESTAMP)", col))
		} else {
			colWithColon = append(colWithColon, fmt.Sprintf(":%s", col))
		}

		// use in SQL's ON DUPLICATE KEY UPDATE
		if strings.Contains(col, "_at") {
			colUpdateSet = append(colUpdateSet, fmt.Sprintf("%s=COALESCE(EXCLUDED.%s, %s.%s)", col, col, m.TableName, col))
			continue
		}
		colUpdateSet = append(colUpdateSet, fmt.Sprintf("%s=COALESCE(EXCLUDED.%s, %s.%s)", col, col, m.TableName, col))
	}

	insertStmt := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s) 
		ON CONFLICT (id) DO UPDATE SET
    %s
		RETURNING id;`,
		m.TableName,
		fmt.Sprintf(strings.Join(cols[:], ",")),
		fmt.Sprintf(strings.Join(colWithColon[:], ",")),
		fmt.Sprintf(strings.Join(colUpdateSet[:], ",\n")),
	)
	fmt.Printf("%+v \n", insertStmt)

	insertedIds := []string{}
	sqlResult, err := m.db.NamedQuery(insertStmt, records)
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

// func (m *Postgres) Update() {
// 	fmt.Printf("update from Postgres, table: %+v\n", m.TableName)
// }
func (m *Postgres) Delete(ids *[]int64) error {
	fmt.Printf("delete from Postgres, table: %+v\n", m.TableName)
	m.db = m.Connect()
	defer m.db.Close()

	deleteStmt, args, err := sqlx.In(
		fmt.Sprintf("DELETE FROM %s WHERE id IN (?);", m.TableName),
		*ids,
	)
	if err != nil {
		log.Printf("sqlx.In err: %+v\n", err.Error())
		return err
	}
	deleteStmt = m.db.Rebind(deleteStmt)
	fmt.Printf("stmt: %+v, args: %+v\n", deleteStmt, args)

	_, err = m.db.Exec(deleteStmt, args...)
	if err != nil {
		log.Printf("Delete Query err: %+v\n", err.Error())
		return err
	}

	return nil
}
