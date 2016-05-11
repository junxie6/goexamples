package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func ReadInput() (string, error) {
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

func WriteFile(fileName string, data []byte) error {
	return ioutil.WriteFile(fileName, data, 0644)
}

func ReadFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func ReadWebContent(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode: %s", resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func DaysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func Now() string {
	// Mon Jan 2 15:04:05 -0700 MST 2006
	return time.Now().Format("2006-01-02 15:04:05")
}

func ConvStrToTime(str string) (time.Time, error) {
	layout := "2006-01-02" // Mon Jan 2 15:04:05 -0700 MST 2006
	return time.Parse(layout, str)
}

// WaitForServer attempts to contact the server of a URL.
// It tries for one minute using exponential backoff.
// It reports an error if all attempts fail.
func WaitForServer(url string) error {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			return nil // success
		}
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential backoff
	}
	return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}

func HashPassword(plaintextPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
}

func ValidatePassword(hashed string, plaintextPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintextPassword))
}

// this function can be used for rows.Scan() for setting the value for database fields from SQL query.
func StrutToSliceOfFieldAddress(theStruct interface{}) []interface{} {
	fieldArr := reflect.ValueOf(theStruct).Elem()

	fieldAddrArr := make([]interface{}, fieldArr.NumField())

	for i := 0; i < fieldArr.NumField(); i++ {
		f := fieldArr.Field(i)
		fieldAddrArr[i] = f.Addr().Interface()
	}

	return fieldAddrArr
}

// Fill a slice with values.
func SliceFill(num int, str string) []string {
	slice := make([]string, num)

	for k, _ := range slice {
		slice[k] = str
	}

	return slice
}

// Generate the placeholders for SQL query.
func Placeholder(num int) string {
	return strings.Join(SliceFill(num, "?"), ",")
}

func RandomNumInSlice(slice []int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return slice[rand.Intn(len(slice))]
}

func PrintStructJSON(s interface{}) {
	if strJSON, err := json.MarshalIndent(s, "", " "); err != nil {
		log.Printf("JSON marshaling failed: %s\n", err)
	} else {
		fmt.Printf("%s\n", strJSON)
	}
}

func PrintJSON(rowArr []interface{}) {
	// produces neatly indented output
	if data, err := json.MarshalIndent(rowArr, "", " "); err != nil {
		log.Printf("JSON marshaling failed: %s\n", err)
	} else {
		fmt.Printf("%s\n", data)
	}
}

func PrintErrJSON(rowArr []error) {
	b := make([]interface{}, len(rowArr))
	for i := range rowArr {
		b[i] = rowArr[i].Error()
	}
	PrintJSON(b)
}

func ConvErrArrToJSON(errArr []error) string {
	strArr := ConvErrArrToStringArr(errArr)

	outMap := map[string]interface{}{
		"Status": false,
		"ErrArr": strArr,
	}

	var byteJSON []byte
	var err error

	if byteJSON, err = json.Marshal(outMap); err != nil {
		return `{"Status":false,"ErrArr":["` + err.Error() + `"]}`
	}

	return string(byteJSON)
}

func ConvSliceToInterface(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		log.Printf("ConvSliceToInterface() given a non-slice type")
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func ConvErrArrToStringArr(errArr []error) []string {
	strArr := make([]string, len(errArr))
	for i := range errArr {
		strArr[i] = errArr[i].Error()
	}
	return strArr
}

func FormValueArr(r *http.Request) map[string]string {
	mapArr := map[string]string{}

	if r.Form == nil {
		r.ParseMultipartForm(32 << 20) // 32 MB
	}

	for k, vs := range r.Form {
		if len(vs) > 0 {
			mapArr[k] = vs[0]
		} else {
			mapArr[k] = ""
		}
	}

	return mapArr
}

func Atoi(num string) int {
	i, _ := strconv.ParseInt(num, 10, 0)
	return int(i)
}

func Atoi64(num string) int64 {
	i, _ := strconv.ParseInt(num, 10, 64)
	return i
}

func StructFieldNameArr(s interface{}) []string {
	sFields := reflect.TypeOf(s)
	fieldNameArr := make([]string, sFields.NumField())

	for i := 0; i < sFields.NumField(); i++ {
		fieldNameArr[i] = sFields.Field(i).Name
	}

	return fieldNameArr
}
