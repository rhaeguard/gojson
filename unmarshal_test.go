package gojson

import (
	"fmt"
	"reflect"
	"testing"
)

type Person struct {
	Name string
	Age  uint8
}

type ComplexPerson struct {
	Person       Person
	Job          string
	LuckyNumbers []int
}

func TestUnmarshall(t *testing.T) {
	doTest := func(refObj interface{}, input string, expected interface{}) {
		name := fmt.Sprintf("unmarshalling-%s", input)
		t.Run(name, func(t *testing.T) {
			json, synErr := ParseJson(input)
			if synErr != nil {
				t.Fatalf("%s", synErr.Error())
			}
			if mError := json.Unmarshal(refObj); mError != nil {
				t.Fatalf("%s", mError.Error())
			} else {
				v := reflect.ValueOf(refObj).Elem().Interface()
				if !reflect.DeepEqual(v, expected) {
					t.Fatalf("expected: '%-v'", expected)
				}
			}
		})
	}

	// string
	var refString string
	doTest(&refString, `"hello world"`, "hello world")

	// bool
	var refBool bool
	doTest(&refBool, `true`, true)

	// int
	var refInt int
	doTest(&refInt, `-1243`, -1243)

	// float64
	var refFloat64 float64
	doTest(&refFloat64, `-0.9912`, -0.9912)

	// slice
	var nums []float64
	doTest(&nums, `[1, 2, 3]`, []float64{1, 2, 3})

	// object
	var person Person
	doTest(&person, `{"Name": "John", "Age": 25}`, Person{Name: "John", Age: 25})

	// array of objects
	var persons []Person
	doTest(&persons, `[{"Name": "John", "Age": 25},{"Name": "Jane", "Age": 23}]`, []Person{
		{Name: "John", Age: 25},
		{Name: "Jane", Age: 23},
	})

	// complex object
	var cPerson ComplexPerson
	doTest(&cPerson, `{"Person": {"Name": "John", "Age": 25}, "Job": "Plumber", "LuckyNumbers": [-1, 0, 1, 1022]}`, ComplexPerson{
		Person:       Person{Name: "John", Age: 25},
		Job:          "Plumber",
		LuckyNumbers: []int{-1, 0, 1, 1022},
	})
}
