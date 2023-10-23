package database

import "fmt"

type Postgre struct {
	Host   string
	Port   string
	User   string
	Pass   string
	TableName string
}
func (m *Postgre) Select() {
	fmt.Printf("select from Postgre, table: %+v\n", m.TableName)
}
func (m *Postgre) Save() {
	fmt.Printf("save from Postgre, table: %+v\n", m.TableName)
}
func (m *Postgre) Update() {
	fmt.Printf("update from Postgre, table: %+v\n", m.TableName)
}
func (m *Postgre) Delete() {
	fmt.Printf("delete from Postgre, table: %+v\n", m.TableName)
}
