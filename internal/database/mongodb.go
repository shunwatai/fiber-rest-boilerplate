package database

import (
	"context"
	"fmt"
	"golang-api-starter/internal/helper"
	"log"
	"math"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongodb struct {
	*ConnectionInfo
	TableName string
	db        *mongo.Client
	ctx       *context.Context
}

type MongoRows struct {
	cur *mongo.Cursor
	ctx context.Context
}

func (mr *MongoRows) StructScan(result interface{}) error {
	if err := mr.cur.Decode(result); err != nil {
		fmt.Printf("mongo decode err: %+v", err.Error())
		return err
	}
	return nil
}
func (mr *MongoRows) Next() bool {
	return mr.cur.Next(mr.ctx)
}
func (mr *MongoRows) Close() error {
	return mr.cur.Close(mr.ctx)
}

func (m *Mongodb) Connect() *mongo.Client {
	fmt.Printf("connecting to Mongodb... \n")
	// fmt.Printf("Table: %+v\n", m.TableName)
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin&sslmode=disable", *m.User, *m.Pass, *m.Host, *m.Port, *m.Database)
	fmt.Printf("ConnString: %+v\n", connectionString)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// Useless for mongodb, but needed for implement the IDatabase...
func (m *Mongodb) constructSelectStmtFromQuerystring(
	queries map[string]interface{},
) (string, *helper.Pagination, map[string]interface{}) {
	return "", &helper.Pagination{}, map[string]interface{}{}
}

// return select statment and *pagination by the req querystring
func (m *Mongodb) getConditionsFromQuerystring(
	queries map[string]interface{},
	countFunc func(interface{}) (int64, error),
	// ) (string, *helper.Pagination, map[string]interface{}) {
) (bson.D, *options.FindOptions, *helper.Pagination) {
	exactMatchCols := map[string]bool{"id": true, "_id": true} // default id(PK) & _id(mongo) have to be exact match
	// fmt.Printf("mongo query: %+v\n\n",queries)
	if queries["exactMatch"] != nil {
		for k := range queries["exactMatch"].(map[string]bool) {
			exactMatchCols[k] = true
		}
	}

	// bindvarMap := map[string]interface{}{}
	if queries["columns"] ==nil{
		fmt.Printf("error: queries[\"columns\"] is nil")
	}
	cols := queries["columns"]
	pagination := helper.GetPagination(queries)
	dateRangeStmt := getDateRangeBson(queries)
	// fmt.Printf("dateRangeStmt: %+v, len: %+v\n", dateRangeStmt, len(dateRangeStmt))
	helper.SanitiseQuerystring(cols.([]string), queries)

	// countAllStmt := fmt.Sprintf("SELECT COUNT(*) FROM %s", m.TableName)
	// selectStmt := fmt.Sprintf(`SELECT * FROM %s`, m.TableName)
	// countAllStmt := bson.D{}
	selectStmt := bson.D{}

	fmt.Printf("queries: %+v, len: %+v\n", queries, len(queries))
	if len(queries) != 0 || len(dateRangeStmt) != 0 { // add where clause
		// whereClauses := []string{}
		whereClauses := bson.D{}
		for k, v := range queries {
			fmt.Printf("%+v: %+v(%T)\n", k, v, v)
			switch v.(type) {
			case []string:
				// placeholders := []string{}
				if exactMatchCols[k] {
					oids := []primitive.ObjectID{}
					if strings.Contains(k, "_id") {
						for _, value := range v.([]string) {
							oid, _ := primitive.ObjectIDFromHex(value)
							oids = append(oids, oid)
						}
						whereClauses = append(whereClauses, bson.E{k, bson.D{{"$in", oids}}})
					} else {
						whereClauses = append(whereClauses, bson.E{k, bson.D{{"$in", v.([]string)}}})
					}
					break
				}

				// mongo $or ref: https://stackoverflow.com/a/58359960
				orWildcard := bson.A{}
				for _, value := range v.([]string) {
					orWildcard = append(orWildcard, bson.D{{k, primitive.Regex{Pattern: fmt.Sprintf(".*%s.*", value), Options: "i"}}})
				}
				whereClauses = append(whereClauses, bson.E{"$or", orWildcard})
			default:
				if exactMatchCols[k] {
					if k == "id" {
						oid, _ := primitive.ObjectIDFromHex(v.(string))
						whereClauses = append(whereClauses, bson.E{"_id", oid})
					} else if strings.Contains(k, "_id") {
						oid, _ := primitive.ObjectIDFromHex(v.(string))
						whereClauses = append(whereClauses, bson.E{k, oid})
					} else {
						whereClauses = append(whereClauses, bson.E{k, v})
					}
					break
				}

				whereClauses = append(whereClauses, bson.E{k, primitive.Regex{Pattern: fmt.Sprintf(".*%s.*", v), Options: "i"}})
			}
		}

		if len(dateRangeStmt) > 0 {
			whereClauses = append(whereClauses, dateRangeStmt...)
		}
		selectStmt = append(selectStmt, whereClauses...)
	}
	// fmt.Printf("countAllStmt: %+v, bindvarmap: %+v\n", countAllStmt, bindvarMap)

	// if totalRow, err := m.db.NamedQuery(countAllStmt, bindvarMap); err != nil {
	// 	log.Printf("Queryx Count(*) err: %+v\n", err.Error())
	// } else if totalRow.Next() {
	// 	defer totalRow.Close()
	// 	totalRow.Scan(&pagination.Count)
	// }
	if count, err := countFunc(selectStmt); err != nil {
		fmt.Printf("count error: %+v\n", err)
	} else {
		fmt.Printf("count: %+v\n", count)
		pagination.Count = count
		if pagination.Items > 0 {
			pagination.TotalPages = int64(math.Ceil(float64(pagination.Count) / float64(pagination.Items)))
		}
		fmt.Printf("pagination: %+v\n", pagination)
	}

	var limit int64
	var offset int64 = (pagination.Page - 1) * pagination.Items
	if pagination.Items == 0 {
		limit = pagination.Count
	} else {
		limit = pagination.Items
	}
	options := options.Find()
	options.SetSkip(offset)
	options.SetLimit(limit)
	options.SetSort(
		bson.M{
			pagination.OrderBy["key"]: func() int {
				if pagination.OrderBy["by"] == "desc" {
					return -1
				}
				return 1
			}(),
		},
	)

	// selectStmt = fmt.Sprintf(`%s
	// 		ORDER BY %s %s
	// 		LIMIT %s OFFSET %s
	// 	`,
	// 	selectStmt,
	// 	pagination.OrderBy["key"], pagination.OrderBy["by"],
	// 	limit, offset,
	// )

	// return selectStmt, pagination, bindvarMap
	return selectStmt, options, pagination
}

// Get all columns []string by m.TableName
func (m *Mongodb) GetColumns() []string {
	return []string{}
}

func (m *Mongodb) Select(queries map[string]interface{}) (Rows, *helper.Pagination) {
	fmt.Printf("select from Mongodb, table: %+v\n", m.TableName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	m.db = m.Connect()
	defer m.db.Disconnect(ctx)
	collection := m.db.Database(fmt.Sprintf("%s", *m.Database)).Collection(fmt.Sprintf("%s", m.TableName))

	fmt.Printf("queries: %+v\n", queries)
	var (
		cur *mongo.Cursor
		err error
	)

	// fmt.Printf("len(queries):%+v, %+v\n", queries, len(queries))
	// if queries["_id"] != nil { // add where clause
	// 	// if len(queries) != 0 { // add where clause
	// 	cur, err = collection.Find(ctx, bson.M{
	// 		"_id": bson.M{"$in": queries["_id"]},
	// 	})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// } else {
	// 	cur, err = collection.Find(ctx, bson.D{})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	var countFunc = func(filter interface{}) (int64, error) {
		count, err := collection.CountDocuments(ctx, filter)
		return count, err
	}

	conditions, findOptions, pagination := m.getConditionsFromQuerystring(queries, countFunc)
	// conditions:=bson.D{{"task", "/take passport/"}}
	fmt.Printf("m conditions: %+v\n", conditions)
	cur, err = collection.Find(ctx, conditions, findOptions)
	// cur, err = collection.Find(ctx, bson.D{{"created_at", bson.D{{"$lt", primitive.NewDateTimeFromTime(time.Now())}}}}, findOptions)

	// oid1, _ := primitive.ObjectIDFromHex("6551ee5f53a746ae0824c3ee")
	// oid2, _ := primitive.ObjectIDFromHex("65519d29973632f67580045d")
	// cur, err = collection.Find(ctx, bson.D{bson.E{"_id", bson.D{{"$in", []primitive.ObjectID{oid1, oid2}}}}})
	if err != nil {
		log.Fatal(err)
	}

	// selectStmt, pagination, bindvarMap := m.constructSelectStmtFromQuerystring(queries)
	// fmt.Printf("bindvarMap: %+v\n", bindvarMap)
	// fmt.Printf("selectStmt: %+v\n", selectStmt)
	//
	// rows, err := m.db.NamedQuery(selectStmt, bindvarMap)
	// if err != nil {
	// 	log.Printf("Queryx err: %+v\n", err.Error())
	// }
	//
	// if rows.Err() != nil {
	// 	log.Printf("rows.Err(): %+v\n", err.Error())
	// }

	return &MongoRows{cur, ctx}, pagination
}

func (m *Mongodb) Save(records Records) (Rows, error) {
	fmt.Printf("save from Mongodb, table: %+v\n", m.TableName)
	// fmt.Printf("records: %+v\n", records)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m.db = m.Connect()
	defer m.db.Disconnect(ctx)
	collection := m.db.Database(fmt.Sprintf("%s", *m.Database)).Collection(fmt.Sprintf("%s", m.TableName))

	opts := options.Update().SetUpsert(true)

	upsertedIds := []string{}
	recordsMap := records.StructToMap()
	for _, record := range recordsMap {
		filter := bson.D{}
		fmt.Printf("record: %+v\n", record)
		if record["_id"] != nil {
			id, _ := primitive.ObjectIDFromHex(record["_id"].(string))
			filter = bson.D{{Key: "_id", Value: id}}
			record["updated_at"] = time.Now()
			if record["created_at"] != nil {
				createdAt, _ := time.Parse(time.RFC3339, record["created_at"].(string))
				record["created_at"] = createdAt
			} else {
				delete(record, "created_at")
			}
			delete(record, "_id")
			upsertedIds = append(upsertedIds, id.Hex())
		} else {
			filter = bson.D{{Key: "_id", Value: primitive.NewObjectID()}}
			record["updated_at"] = time.Now()
			record["created_at"] = time.Now()
		}

		res, err := collection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: record},
		}, opts)
		if err != nil {
			log.Fatal(err)
		}

		/* only new created records has res.UpsertedID, existing's Ids appended in the if condition above */
		if res.UpsertedID == nil {
			continue
		}
		upsertedIds = append(upsertedIds, res.UpsertedID.(primitive.ObjectID).Hex())
	}
	fmt.Printf("upsertedIds: %+v\n", upsertedIds)

	// cols := m.GetColumns()
	//
	// // fmt.Printf("cols: %+v\n", cols)
	// var colWithColon, colUpdateSet []string
	// for _, col := range cols {
	// 	// use in SQL's VALUES()
	// 	if col == "id" {
	// 		colWithColon = append(colWithColon, fmt.Sprintf("COALESCE(:%s, nextval('%s_id_seq'))", col, m.TableName))
	// 	} else if strings.Contains(col, "_at") {
	// 		colWithColon = append(colWithColon, fmt.Sprintf("COALESCE(:%s, CURRENT_TIMESTAMP)", col))
	// 	} else {
	// 		colWithColon = append(colWithColon, fmt.Sprintf(":%s", col))
	// 	}
	//
	// 	// use in SQL's ON DUPLICATE KEY UPDATE
	// 	if strings.Contains(col, "_at") {
	// 		colUpdateSet = append(colUpdateSet, fmt.Sprintf("%s=COALESCE(EXCLUDED.%s, %s.%s)", col, col, m.TableName, col))
	// 		continue
	// 	}
	// 	colUpdateSet = append(colUpdateSet, fmt.Sprintf("%s=COALESCE(EXCLUDED.%s, %s.%s)", col, col, m.TableName, col))
	// }
	//
	// insertStmt := fmt.Sprintf(
	// 	`INSERT INTO %s (%s) VALUES (%s)
	// 	ON CONFLICT (id) DO UPDATE SET
	//    %s
	// 	RETURNING id;`,
	// 	m.TableName,
	// 	fmt.Sprintf(strings.Join(cols[:], ",")),
	// 	fmt.Sprintf(strings.Join(colWithColon[:], ",")),
	// 	fmt.Sprintf(strings.Join(colUpdateSet[:], ",\n")),
	// )
	// fmt.Printf("%+v \n", insertStmt)
	//
	// insertedIds := []string{}
	// sqlResult, err := m.db.NamedQuery(insertStmt, records)
	// if err != nil {
	// 	log.Printf("insert error: %+v\n", err)
	// }
	// // fmt.Printf("sqlResult: %+v\n", sqlResult)
	//
	// for sqlResult.Next() {
	// 	var id string
	// 	err := sqlResult.Scan(&id)
	// 	if err != nil {
	// 		log.Fatalf("Scan: %v", err)
	// 	}
	// 	insertedIds = append(insertedIds, id)
	// }
	//
	// fmt.Printf("insertedIds: %+v\n", insertedIds)
	// rows, _ := m.Select(map[string]interface{}{"id": insertedIds})

	// return &sqlx.Rows{}

	rows, _ := m.Select(map[string]interface{}{
		"_id":     upsertedIds,
		"columns": records.GetTags("bson"),
	})
	return rows, nil
}

// func (m *Mongodb) Update() {
// 	fmt.Printf("update from Mongodb, table: %+v\n", m.TableName)
// }
func (m *Mongodb) Delete(ids []string) error {
	fmt.Printf("delete ids: %+v from Mongodb, table: %+v\n", ids, m.TableName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m.db = m.Connect()
	defer m.db.Disconnect(ctx)
	collection := m.db.Database(fmt.Sprintf("%s", *m.Database)).Collection(fmt.Sprintf("%s", m.TableName))
	objectIds := []primitive.ObjectID{}
	for _, id := range ids {
		if oid, err := primitive.ObjectIDFromHex(id); err != nil {
			fmt.Printf("ObjectIDFromHex err: %+v\n", err)
		} else {
			objectIds = append(objectIds, oid)
		}
	}
	filter := bson.D{{"_id", bson.D{{"$in", objectIds}}}}
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		fmt.Printf("DeleteMany err: %+v\n", err)
	}
	fmt.Printf("result: %+v\n", result)

	// m.db = m.Connect()
	// defer m.db.Close()
	//
	// deleteStmt, args, err := sqlx.In(
	// 	fmt.Sprintf("DELETE FROM %s WHERE id IN (?);", m.TableName),
	// 	*ids,
	// )
	// if err != nil {
	// 	log.Printf("sqlx.In err: %+v\n", err.Error())
	// 	return err
	// }
	// deleteStmt = m.db.Rebind(deleteStmt)
	// fmt.Printf("stmt: %+v, args: %+v\n", deleteStmt, args)
	//
	// _, err = m.db.Exec(deleteStmt, args...)
	// if err != nil {
	// 	log.Printf("Delete Query err: %+v\n", err.Error())
	// 	return err
	// }

	return nil
}

// useless for mongo, it implemented by sqlite, postgres, mariadb
func (m *Mongodb) RawQuery(sql string) *sqlx.Rows {
	// fmt.Printf("raw query from Mongodb\n")
	return &sqlx.Rows{}
}

// mongo version of RawQuery
func (m *Mongodb) runCommands(cmds []bson.D) error {
	for _, cmd := range cmds {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		m.db = m.Connect()
		defer m.db.Disconnect(ctx)
		db := m.db.Database(fmt.Sprintf("%s", *m.Database)).Collection(fmt.Sprintf("%s", m.TableName))

		err := db.Database().RunCommand(ctx, cmd).Err()
		if err != nil {
			log.Printf("mongo cmd failed: %+v\n", err)
			return err
		}
	}
	return nil
}
