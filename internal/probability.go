package internal

import (
	"math/rand"
	"time"
)

func GetSampleFromCDF(cdf []float64) int {
	rand.Seed(time.Now().UTC().UnixNano())
	randomFloat := rand.Float64()

	sample := BucketBinarySearch(cdf, randomFloat)

	return sample
}

func BucketBinarySearch(buckets []float64, value float64) int {
	// Finds bucket with value
	// BucketBinarysearch([0.1, 0.4, 0.6, 1.0], 0,5) = 2
	low := 0
	high := len(buckets) - 1
	for low <= high {
		median := (low + high) / 2
		if median == 0 {
			return 0
		}
		if value > buckets[median - 1] && value < buckets[median] {
			return median
		}


		if value > buckets[median] {
			low = median + 1
		} else {
			high = median - 1
		}
	}
	return -1
}
