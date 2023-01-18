package utils

import (
	"crypto/md5" //nolint:gosec
	"encoding/base64"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// RandString 生成指定长度随机字符串
// nolint:gosec
func RandString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	data := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, data[r.Intn(len(data))])
	}
	return string(result)
}

// StrIsEmpty 是否空字符串
func StrIsEmpty(str string) bool {
	return str == "null" || strings.TrimSpace(str) == ""
}

// Str2UintArr 字符串转uint数组, 默认逗号分割
func Str2UintArr(str string) (ids []uint) {
	idArr := strings.Split(str, ",")
	for _, v := range idArr {
		ids = append(ids, Str2Uint(v))
	}
	return
}

// Str2Uint 字符串转uint
func Str2Uint(str string) uint {
	num, err := strconv.ParseUint(str, 10, 32) //nolint:gomnd
	if err != nil || num == math.MaxUint {
		return 0
	}
	return uint(num)
}

// Str2Float64 字符串转float64
func Str2Float64(str string) float64 {
	num, err := strconv.ParseFloat(str, 64) //nolint:gomnd
	if err != nil || math.IsNaN(num) {
		return 0
	}
	return num
}

var (
	camelRe = regexp.MustCompile("(_)([a-zA-Z]+)")
	snakeRe = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// CamelCase 字符串转为驼峰
func CamelCase(str string) string {
	camel := camelRe.ReplaceAllString(str, " $2")
	ca := cases.Title(language.Und, cases.NoLower)
	camel = ca.String(camel)
	camel = strings.Replace(camel, " ", "", -1)
	return camel
}

// CamelCaseLowerFirst 字符串转为驼峰(首字母小写)
func CamelCaseLowerFirst(str string) string {
	camel := CamelCase(str)
	for i, v := range camel {
		return string(unicode.ToLower(v)) + camel[i+1:]
	}
	return camel
}

// SnakeCase 驼峰式写法转为下划线蛇形写法
func SnakeCase(str string) string {
	snake := snakeRe.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

// EncodeStr2Base64 加密base64字符串
func EncodeStr2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// DecodeStrFromBase64 解密base64字符串
func DecodeStrFromBase64(str string) string {
	decodeBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes)
}

// Md5V2 对字符串进行md5
func Md5V2(str string) string {
	data := []byte(str)
	has := md5.Sum(data) //nolint:gosec
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
