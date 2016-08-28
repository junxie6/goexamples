package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
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

func DecodeJSONStreamStruct(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}

func DecodeJSONStreamMap(r *http.Request) (map[string]interface{}, error) {
	var data interface{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.(map[string]interface{}), nil
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

func InArrayV1(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func InArrayV2(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}

func InArrayInt(v int, vArr []int) bool {
	for _, vv := range vArr {
		if v == vv {
			return true
		}
	}

	return false
}

func InArrayStr(v string, vArr []string) bool {
	for _, vv := range vArr {
		if v == vv {
			return true
		}
	}

	return false
}

var FGColor = struct {
	White, Red, Green, Yellow string
}{
	White:  "1;37",
	Red:    "0;31",
	Green:  "0;32",
	Yellow: "1;33",
}

func EchoColor(msg string, color string) string {
	return "\033[" + color + "m" + msg + "\033[0m"
}

// Close is used for defer statement. Example: defer Close(VarResource)
func Close(c io.Closer) {
	// Note: do we need to add recover() here?
	if err := c.Close(); err != nil {
		log.Printf(err.Error())
	}
}

// JSONDeepEqual ...
func JSONDeepEqual(s1 string, s2 string) (bool, error) {
	var m1, m2 map[string]interface{}

	if err := json.Unmarshal([]byte(s1), &m1); err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(s2), &m2); err != nil {
		return false, err
	}

	return reflect.DeepEqual(m1, m2), nil
}

func HelpGenTLSKeys() {
	str := `
To generate the private key and the self-signed certificate:

Use this method if you want to use HTTPS (HTTP over TLS) to secure your Apache HTTP or Nginx web server, and you want to use a Certificate Authority (CA) to issue the SSL certificate. The CSR that is generated can be sent to a CA to request the issuance of a CA-signed SSL certificate. If your CA supports SHA-2, add the -sha256 option to sign the CSR with SHA-2.

# openssl req -newkey rsa:2048 -nodes -subj "/C=CA/ST=British Columbia/L=Vancouver/O=My Company Name/CN=erp.local" -keyout erp.local.key -out erp.local.csr

Note: The -newkey rsa:2048 option specifies that the key should be 2048-bit, generated using the RSA algorithm.
Note: The -nodes option specifies that the private key should not be encrypted with a pass phrase.
Note: The -new option, which is not included here but implied, indicates that a CSR is being generated.

Generate a Self-Signed Certificate:

Use this method if you want to use HTTPS (HTTP over TLS) to secure your Apache HTTP or Nginx web server, and you do not require that your certificate is signed by a CA.

This command creates a 2048-bit private key (domain.key) and a self-signed certificate (domain.crt) from scratch:

# openssl req -newkey rsa:2048 -nodes -subj "/C=CA/ST=British Columbia/L=Vancouver/O=My Company Name/CN=erp.local" -keyout erp.local.key -x509 -days 365 -out erp.local.crt
`
	log.Printf("%v\n", str)
}
