package helper

import (
	"fmt"
	"log"
	"strconv"
)

func ConvertNumberSliceToString(i interface{}) ([]string, error) {
	stringSlice := []string{}

	switch i.(type) {
	case []interface{}:
		for _, v := range i.([]interface{}) {
			stringSlice = append(stringSlice, strconv.Itoa(int(v.(int))))
		}
	case []int64:
		for _, v := range i.([]int64) {
			stringSlice = append(stringSlice, strconv.Itoa(int(v)))
		}
	case []int:
		for _, v := range i.([]int) {
			stringSlice = append(stringSlice, strconv.Itoa(int(v)))
		}
	case []float64:
		for _, v := range i.([]float64) {
			stringSlice = append(stringSlice, strconv.Itoa(int(v)))
		}
	default:
		msg := "failed to handle %+v (%T)\n"
		log.Printf(msg, i, i)
		return nil, fmt.Errorf(msg, i, i)
		// fmt.Printf("failed to handle %+v (%T)\n", v, v)
	}

	return stringSlice, nil
}
