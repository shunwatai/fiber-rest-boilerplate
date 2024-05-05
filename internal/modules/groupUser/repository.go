package groupUser

import (
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	logger "golang-api-starter/internal/helper/logger/zap_log"

	//"golang-api-starter/internal/modules/user"
	"golang.org/x/exp/maps"
)

type Repository struct {
	db database.IDatabase
}

func NewRepository(db database.IDatabase) *Repository {
	return &Repository{db}
}


func (r *Repository) Get(queries map[string]interface{}) ([]*GroupUser, *helper.Pagination) {
	logger.Debugf("groupUser repo get")
	defaultExactMatch := map[string]bool{
		"id":  true,
		"_id": true,
		//"done": true, // bool match needs exact match, param can be 0(false) & 1(true)
	}
	if queries["exactMatch"] != nil {
		maps.Copy(queries["exactMatch"].(map[string]bool), defaultExactMatch)
	} else {
		queries["exactMatch"] = defaultExactMatch
	}

	queries["columns"] = GroupUser{}.getTags()
	rows, pagination := r.db.Select(queries)

	var records GroupUsers
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	// records.printValue()

	return records, pagination
}

func (r *Repository) Create(groupUsers []*GroupUser) ([]*GroupUser, error) {
	for _, groupUser := range groupUsers {
		logger.Debugf("groupUser repo add: %+v", groupUser)
	}
	database.SetIgnoredCols("search")
	defer database.SetIgnoredCols()
	rows, err := r.db.Save(GroupUsers(groupUsers))

	var records GroupUsers
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records, err
}

func (r *Repository) Update(groupUsers []*GroupUser) ([]*GroupUser, error) {
	logger.Debugf("groupUser repo update")
	rows, err := r.db.Save(GroupUsers(groupUsers))

	var records GroupUsers
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records, err
}

func (r *Repository) Delete(ids []string) error {
	logger.Debugf("groupUser repo delete")
	err := r.db.Delete(ids)
	if err != nil {
		return err
	}

	return nil
}
