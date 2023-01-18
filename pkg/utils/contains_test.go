package utils

import (
	"testing"
)

// nolint: funlen
func TestContains(t *testing.T) {
	containsTest := []struct {
		name      string
		arr, item any
		want      bool
	}{
		{
			name: "uint1",
			arr:  []uint{1, 2, 3},
			item: uint(1),
			want: true,
		},
		{
			name: "uint2",
			arr:  []uint{3, 4, 5},
			item: uint(8),
			want: false,
		},
		{
			name: "int1",
			arr:  []int{1, -1, 3},
			item: -1,
			want: true,
		},
		{
			name: "int2",
			arr:  []int{2, 4, -3},
			item: -1,
			want: false,
		},
		{
			name: "string1",
			arr:  []string{"1", "b", "e232323"},
			item: "b",
			want: true,
		},
		{
			name: "string2",
			arr:  []string{"dsdf", "333", "sd222"},
			item: "2323s",
			want: false,
		},
		{
			name: "int32_1",
			arr:  []int32{100, 333, 0xa},
			item: int32(10),
			want: true,
		},
		{
			name: "int32_2",
			arr:  []int32{100, 333, 0xe},
			item: int32(10),
			want: false,
		},
		{
			name: "int64_1",
			arr:  []int64{100, 333, 0xa},
			item: int64(10),
			want: true,
		},
		{
			name: "int64_2",
			arr:  []int64{100, 333, 0xe},
			item: int64(10),
			want: false,
		},
		{
			name: "float32_1",
			arr:  []float32{1.2, 3.3, 33.0},
			item: float32(33.0),
			want: true,
		},
		{
			name: "float32_2",
			arr:  []float32{12.3, 33.4, 22.0},
			item: float32(33.0),
			want: false,
		},
		{
			name: "float64_1",
			arr:  []float64{1.33, 4.3, 5.555},
			item: 4.3,
			want: true,
		},
		{
			name: "float64_2",
			arr:  []float64{1.33, 4.3, 5.333},
			item: 5.33,
			want: false,
		},
		{
			name: "other",
			arr:  []rune{555, 666, 777, 'a'},
			item: 'a',
			want: true,
		},
	}
	for _, tt := range containsTest {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.arr, tt.item); got != tt.want {
				t.Errorf(`[%s] wantErr %v, but Contains(%q, %q) = %v`, tt.name, tt.want, tt.arr, tt.item, got)
			}
		})
	}
}

func TestContainsUint(t *testing.T) {
	tests := []struct {
		name string
		arr  []uint
		item uint
		want bool
	}{
		{
			name: "indexTrue",
			arr:  []uint{11, 22, 33},
			item: 11,
			want: true,
		},
		{
			name: "indexFalse",
			arr:  []uint{},
			item: 0,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsUint(tt.arr, tt.item); got != tt.want {
				t.Errorf(`[%s] wantErr %v, but ContainsUint(%q, %q) = %v`, tt.name, tt.want, tt.arr, tt.item, got)
			}
		})
	}
}

func TestContainsUintIndex(t *testing.T) {
	tests := []struct {
		name string
		arr  []uint
		item uint
		want int
	}{
		{
			name: "indexTrue",
			arr:  []uint{11, 22, 33},
			item: 11,
			want: 0,
		},
		{
			name: "indexFalse",
			arr:  []uint{},
			item: 0,
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsUintIndex(tt.arr, tt.item); got != tt.want {
				t.Errorf(`[%s] wantErr %v, but ContainsUintIndex(%q, %q) = %v`, tt.name, tt.want, tt.arr, tt.item, got)
			}
		})
	}
}
