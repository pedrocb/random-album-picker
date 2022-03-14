package internal

import (
	"testing"
)

func TestBucketBinarySearch(t *testing.T) {
	buckets := []float64{0.1, 0.4, 0.6, 1.0}
	result := BucketBinarySearch(buckets, 0.0)
	expected := 0
	if result != expected {
		t.Errorf("Expected %d got %d", expected, result)
	}

	result = BucketBinarySearch(buckets, 0.256)
	expected = 1
	if result != expected {
		t.Errorf("Expected %d got %d", expected, result)
	}

	result = BucketBinarySearch(buckets, 0.864)
	expected = 3
	if result != expected {
		t.Errorf("Expected %d got %d", expected, result)
	}

	result = BucketBinarySearch(buckets, 2.854)
	expected = -1
	if result != expected {
		t.Errorf("Expected %d got %d", expected, result)
	}
}
