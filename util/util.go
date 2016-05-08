package util

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	return strings.Join(sliceFill(num, "?"), ",")
}
