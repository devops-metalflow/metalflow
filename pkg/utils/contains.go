package utils

import "github.com/thoas/go-funk"

// Contains Determine if the array arr contains the item element
func Contains(arr, item any) bool {
	//nolint:gocritic
	switch arr.(type) {
	case []uint:
		// funk does not have a strong type []uint
		if val, ok := item.(uint); ok {
			return ContainsUint(arr.([]uint), val)
		}
	case []int:
		if val, ok := item.(int); ok {
			return funk.ContainsInt(arr.([]int), val)
		}
	case []string:
		if val, ok := item.(string); ok {
			return funk.ContainsString(arr.([]string), val)
		}
	case []int32:
		if val, ok := item.(int32); ok {
			return funk.ContainsInt32(arr.([]int32), val)
		}
	case []int64:
		if val, ok := item.(int64); ok {
			return funk.ContainsInt64(arr.([]int64), val)
		}
	case []float32:
		if val, ok := item.(float32); ok {
			return funk.ContainsFloat32(arr.([]float32), val)
		}
	case []float64:
		if val, ok := item.(float64); ok {
			return funk.ContainsFloat64(arr.([]float64), val)
		}
	}
	// funk uses reflection by default, and the performance is not as good as strong typing
	return funk.Contains(arr, item)
}

// ContainsUint check if uint array contains item element
func ContainsUint(arr []uint, item uint) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

// ContainsUintIndex determine whether the uint array contains item elements, return index
func ContainsUintIndex(arr []uint, item uint) int {
	for i, v := range arr {
		if v == item {
			return i
		}
	}
	return -1
}
