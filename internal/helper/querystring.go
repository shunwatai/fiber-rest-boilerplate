package helper

import (
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

func parseOrderBy(rawOrderBy string) map[string]string {
	orderBy := make(map[string]string)
	splitedOrderBy := strings.Split(rawOrderBy, ".")
	// fmt.Printf("splitedOrderBy:  %+v\n", splitedOrderBy)

	orderBy["by"] = "desc"
	if len(rawOrderBy) == 0 {
		orderBy["key"] = "id"
	} else if len(splitedOrderBy) != 2 {
		orderBy["key"] = strcase.ToSnake(splitedOrderBy[0])
	} else {
		orderBy["key"] = strcase.ToSnake(splitedOrderBy[0])
		orderBy["by"] = strings.ToLower(splitedOrderBy[1])
	}

	return orderBy
}

type Pagination struct {
	Page       int64             `json:"page"`       // current page
	Items      int64             `json:"items"`      // records per page
	Count      int64             `json:"count"`      // total records
	OrderBy    map[string]string `json:"orderBy"`    // orderBy
	TotalPages int64             `json:"totalPages"` // ceil(count / items)
}

func getDefaultPagination() *Pagination {
	orderBy := map[string]string{
		"key": "id",
		"by":  "desc",
	}

	return &Pagination{
		Page:       1,
		Items:      0,
		Count:      0,
		OrderBy:    orderBy,
		TotalPages: 1,
	}
}

func GetPagination(cols []string, queries map[string]interface{}) *Pagination{
	pagination := getDefaultPagination()

	if queries["page"] != nil && queries["items"] != nil {
		pagination.Page, _ = strconv.ParseInt(queries["page"].(string), 10, 64)
		pagination.Items, _ = strconv.ParseInt(queries["items"].(string), 10, 64)
	}

	if queries["order_by"] != nil {
		pagination.OrderBy = parseOrderBy(queries["order_by"].(string))
	}

	tmpColsMap := map[string]struct{}{}
	for _, col := range cols {
		tmpColsMap[col] = struct{}{}
	}
	for k := range queries {
		_, ok := tmpColsMap[k]
		if !ok {
			delete(queries, k)
		}
	}

	return pagination
}
