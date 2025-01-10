package helper

import (
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"math"
)

// GetStandardDeviation calculates the standard deviation of a slice of float64 numbers.
func GetStandardDeviation(nums []float64) (float64, error) {
	if len(nums) == 0 {
		return 0, logger.Errorf("len(nums) must be greater than 0")
	}

	var sum, mean, sd float64

	// Calculate sum
	for _, num := range nums {
		sum += num
	}
	mean = sum / float64(len(nums))

	// Calculate standard deviation
	for _, num := range nums {
		sd += math.Pow(num-mean, 2)
	}
	sd = math.Sqrt(sd / float64(len(nums)))

	return sd, nil
}
