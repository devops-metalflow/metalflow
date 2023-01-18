package utils

import (
	"reflect"
	"testing"
)

func TestStruct2Json(t *testing.T) {
	tests := []struct {
		name, want string
		obj        any
	}{
		{
			name: "testError",
			want: "",
			obj:  make(chan struct{}),
		},
		{
			name: "testSuccess",
			want: `{"name":"jack","age":18}`,
			obj: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{
				Name: "jack",
				Age:  18,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Struct2Json(tt.obj); got != tt.want {
				t.Errorf(`[%s] wantErr %v, but Struct2Json(%q) = %v`, tt.name, tt.want, tt.obj, got)
			}
		})
	}
}

func TestJson2Struct(t *testing.T) {
	tests := []struct {
		name, str string
		obj, want any
	}{
		{
			name: "testError",
			str:  `{"name":}`,
			obj: &struct {
				Name string `json:"name"`
			}{},
			want: &struct {
				Name string `json:"name"`
			}{},
		},
		{
			name: "testSuccess",
			str:  `{"name":"jack"}`,
			obj: &struct {
				Name string `json:"name"`
			}{},
			want: &struct {
				Name string `json:"name"`
			}{Name: "jack"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Json2Struct(tt.str, tt.obj)
			if !reflect.DeepEqual(tt.obj, tt.want) {
				t.Errorf(`[%s] Json2Struct() = %v, wantErr %v`, tt.name, tt.obj, tt.want)
			}
		})
	}
}

func TestCompareDifferenceStructByJson(t *testing.T) {
	tests := []struct {
		name                 string
		oldStruct, newStruct any
		update, want         *map[string]any
	}{
		{
			name: "testStruct1",
			oldStruct: struct {
				Name string `json:"name"`
			}{
				Name: "jack",
			},
			newStruct: struct {
				Name string `json:"name"`
			}{
				Name: "tom",
			},
			update: &map[string]any{},
			want:   &map[string]any{"name": "tom"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CompareDifferenceStructByJson(tt.oldStruct, tt.newStruct, tt.update)
			if !reflect.DeepEqual(tt.update, tt.want) {
				t.Errorf(`[%s] CompareDifferenceStructByJson() = %v, wantErr %v`, tt.name, tt.update, tt.want)
			}
		})
	}
}
