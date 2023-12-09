package helper

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type CustomDatetime struct {
	*time.Time
	Format *string
}

func ParseInputDatetime(datetime string) (*time.Time, error) {
	var err error
	newTime := time.Time{}

	dateFormats := []string{
		time.RFC3339,
		time.UnixDate,
		time.RFC822Z,
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
	}

	for _, format := range dateFormats {
		if newTime, err = time.Parse(fmt.Sprint(format), datetime); err == nil {
			return &newTime, err
		}
	}

	return nil, fmt.Errorf("failed to parse given datetime: %s\n", datetime)
}

func (t *CustomDatetime) UnmarshalJSON(input []byte) error {
	strInput := strings.Trim(string(input), `"`)
	fmt.Printf("strInput: %+v\n", strInput)

	parsedTime, err := ParseInputDatetime(strInput)
	fmt.Printf("parsedTime: %+v\n", parsedTime)
	if err == nil {
		t.Time = parsedTime
	}
	return err
}

func (t CustomDatetime) MarshalJSON() ([]byte, error) {
	var jsonDatetime string
	// fmt.Printf("date::? %+v\n",t)
	if t.Format == nil {
		jsonDatetime = fmt.Sprintf("\"%s\"", t.Time.Format(time.RFC3339))
		return []byte(jsonDatetime), nil
	}
	jsonDatetime = fmt.Sprintf("\"%s\"", t.Time.Format(*t.Format))
	return []byte(jsonDatetime), nil
}

func (t *CustomDatetime) UnmarshalBSONValue(bt bsontype.Type, value []byte) error {
	// fmt.Printf("UnmarshalBSONValue type:%+v,  value:(%+v)\n",value, bt)
	if bt != bsontype.DateTime {
		return fmt.Errorf("invalid bson value type '%s'", t.String())
	}

	parsedTime, _, ok := bsoncore.ReadTime(value)
	if !ok {
		return fmt.Errorf("invalid bson datetime value")
	}

	// fmt.Printf("parsedTime: %+v\n", parsedTime)
	t.Time = &parsedTime
	return nil
}

// ref: https://stackoverflow.com/a/54921922
func (t CustomDatetime) Value() (driver.Value, error) {
	return time.Time(*t.Time), nil
}

func (t *CustomDatetime) Scan(src interface{}) error {
	if val, ok := src.(time.Time); ok {
		t.Time = &val
	} else {
		return errors.New("time Scanner passed a non-time object")
	}

	return nil
}
