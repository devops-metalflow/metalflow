//go:build go1.18

package utils

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestStr2UintArr(t *testing.T) {
	tests := []struct {
		name, str string
		want      []uint
	}{
		{
			name: "test1",
			str:  "1,2,3,4,5",
			want: []uint{1, 2, 3, 4, 5},
		},
		{
			name: "test2",
			str:  "0,1,-2,3,-4,5",
			want: []uint{0, 1, 0, 3, 0, 5},
		},
		{
			name: "test3",
			str:  "1,s,3,44,jack,4",
			want: []uint{1, 0, 3, 44, 0, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Str2UintArr(tt.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`[%s] wantErr %v, but Str2UintArr(%q) = %v`, tt.name, tt.want, tt.str, got)
			}
		})
	}
}

func FuzzStr2Uint(f *testing.F) {
	testCases := []string{"11", "22", "33", "-1", "ss", "sdf"}
	for _, s := range testCases {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, i string) {
		var n uint
		num, err := strconv.ParseUint(i, 10, 32) //nolint:gomnd
		if err != nil || num == math.MaxUint {
			n = 0
		} else {
			n = uint(num)
		}
		if got := Str2Uint(i); got != n {
			t.Errorf(` wantErr %v, but Str2Uint(%q) = %v`, n, i, got)
		}
	})
}

func FuzzStr2Float64(f *testing.F) {
	testCases := []string{"11.11", "22.22", "33.33", "-1.1", "ss", "sdf"}
	for _, s := range testCases {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, i string) {
		var n float64
		num, err := strconv.ParseFloat(i, 64) //nolint:gomnd
		if err != nil || math.IsNaN(num) {
			n = 0
		} else {
			n = num
		}
		if got := Str2Float64(i); got != n {
			t.Errorf(`wantErr %v, but Str2Float64(%q) = %v`, n, i, got)
		}
	})
}

func FuzzCamelCase(f *testing.F) {
	camelRe := regexp.MustCompile("(_)([a-zA-Z]+)")
	testCases := []string{"CamelCase", "testSome", "joker", "AAA", "Time322", "Hss"}
	for _, s := range testCases {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, i string) {
		camel := camelRe.ReplaceAllString(i, " $2")
		ca := cases.Title(language.Und, cases.NoLower)
		camel = ca.String(camel)
		camel = strings.Replace(camel, " ", "", -1)
		assert.Equal(t, camel, CamelCase(i))
	})
}

func FuzzSnakeCase(f *testing.F) {
	testCases := []string{"CamelCase", "testSome", "joker", "AAA", "Time322", "Hss", "hello_a"}
	snakeRe := regexp.MustCompile("([a-z0-9])([A-Z])")
	for _, s := range testCases {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, i string) {
		snake := snakeRe.ReplaceAllString(i, "${1}_${2}")
		snake = strings.ToLower(snake)
		assert.Equal(t, snake, SnakeCase(i))
	})
}
