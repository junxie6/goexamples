package gobgzipv2

import (
	"fmt"
	"testing"
)

// Person ...
type Person struct {
	Name   string
	Age    int
	Money  float64
	Houses []House
}

// House ...
type House struct {
	StreetNum     int
	StreetName    string
	PropertyValue float64
}

func GetPerson() Person {
	return Person{
		Name:  "AAA",
		Age:   123,
		Money: 123456.789,
		Houses: []House{
			House{
				StreetNum:     111,
				StreetName:    "Hello 111",
				PropertyValue: 123456.789,
			},
			House{
				StreetNum:     222,
				StreetName:    "Hello 222",
				PropertyValue: 987654.321,
			},
		},
	}
}

func BenchmarkGobGzipV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		person := GetPerson()
		person2 := Person{}

		if b, err := EncodeGobThenGzip(&person); err != nil {
			fmt.Printf("err: %v\n\n", err)
		} else if err := UngzipThenDecodeGob(b, &person2); err != nil {
			fmt.Printf("err: %v\n\n", err)
		}
	}
}
