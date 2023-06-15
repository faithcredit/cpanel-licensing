package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func FormatDateTime(dateString string) string {

	date, error := time.Parse("2006-01-02", dateString)

	if error != nil {
		fmt.Println(error)
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(" %2dth %s %2d", date.Day(), date.Month().String(), date.Year()))
	return buffer.String()
}

func StrToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return num
}
func Str2int(str string) int {
	ret, _ := strconv.Atoi(str)
	return ret
}

func Strpos(haystack string, needle string) int {
	return strings.Index(haystack, needle)
}

func Str_replace(search string, replace string, subject string) string {
	return strings.ReplaceAll(subject, search, replace)
}
func Strtolower(str string) string {
	return strings.ToLower(str)
}

func Map2str(mapData map[string]string) string {
	jsonData, err := json.Marshal(mapData)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	jsonStr := string(jsonData)

	fmt.Println(jsonStr)

	return jsonStr
}

func MapToByteSlice(dataMap map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, err
	}
	return []byte(jsonData), nil
}

func MapToString(dataMap map[interface{}]interface{}) string {
	str := fmt.Sprintf("%v", dataMap)
	return str
}
func Interface2Str(data interface{}) string {
	return fmt.Sprintf("%v", data)
}
