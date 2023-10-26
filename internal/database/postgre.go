package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	*ConnectionInfo
	TableName string
}

func (m *Postgres) Select(queries map[string]interface{}) *sqlx.Rows {
	fmt.Printf("select from Postgres, table: %+v\n", m.TableName)
	return nil
}
func (m *Postgres) Save(records Records) *sqlx.Rows {
	fmt.Printf("save from Postgres, table: %+v\n", m.TableName)
	return nil
}

// func (m *Postgres) Update() {
// 	fmt.Printf("update from Postgres, table: %+v\n", m.TableName)
// }
func (m *Postgres) Delete(ids *[]int64) error {
	fmt.Printf("delete from Postgres, table: %+v\n", m.TableName)
	return nil
}
