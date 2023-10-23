package database

import "fmt"

type MariaDb struct {
	Host   string
	Port   string
	User   string
	Pass   string
	TableName string
}
func (m *MariaDb) Select() {
	fmt.Printf("select from MariaDB, table: %+v\n", m.TableName)
}
func (m *MariaDb) Save() {
	fmt.Printf("save from MariaDB, table: %+v\n", m.TableName)
}
func (m *MariaDb) Update() {
	fmt.Printf("update from MariaDB, table: %+v\n", m.TableName)
}
func (m *MariaDb) Delete() {
	fmt.Printf("delete from MariaDB, table: %+v\n", m.TableName)
}
