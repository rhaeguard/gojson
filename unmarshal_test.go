package gojson

import (
	"fmt"
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

func TestErrorHandling34(t *testing.T) {
	t.Run("unmarshalling test", func(t *testing.T) {
		var refString string
		var refBool bool
		var refInt int
		var refFloat64 float64

		// string
		json, _ := parseJson(`"hello world"`)
		if synErr := json.Unmarshal(&refString); synErr != nil {
			t.Fatalf("%s", synErr.Error())
		}

		// bool
		json, _ = parseJson(`true`)
		if synErr := json.Unmarshal(&refBool); synErr != nil {
			t.Fatalf("%s", synErr.Error())
		}

		// int
		json, _ = parseJson(`-1243`)
		if synErr := json.Unmarshal(&refInt); synErr != nil {
			t.Fatalf("%s", synErr.Error())
		}

		// float64
		json, _ = parseJson(`-0.9912`)
		if synErr := json.Unmarshal(&refFloat64); synErr != nil {
			t.Fatalf("%s", synErr.Error())
		}

		// slice
		var nums []float64
		json, _ = parseJson(`[1, 2, 3]`)
		if synErr := json.Unmarshal(&nums); synErr != nil {
			t.Fatalf("%s", synErr.Error())
		}

		// object
		var person Person
		json, _ = parseJson(`{"Name": "John", "Age": 25}`)
		if synErr := json.Unmarshal(&person); synErr != nil {
			t.Fatalf("%s", synErr.Error())
		}

		// array of objects
		var persons []Person
		json, _ = parseJson(`[{"Name": "John", "Age": 25},{"Name": "Jane", "Age": 23}]`)
		if synErr := json.Unmarshal(&persons); synErr != nil {
			t.Fatalf("%s", synErr.Error())
		}

		// complex
		var cPerson ComplexPerson
		json, _ = parseJson(`{"Person": {"Name": "John", "Age": 25}, "Job": "Plumber", "LuckyNumbers": [-1, 0, 1, 1022]}`)
		if synErr := json.Unmarshal(&cPerson); synErr != nil {
			t.Fatalf("%s", synErr.Error())
		}

		fmt.Printf("Value is: '%s'\n", refString)
		fmt.Printf("Value is: '%t'\n", refBool)
		fmt.Printf("Value is: '%d'\n", refInt)
		fmt.Printf("Value is: '%f'\n", refFloat64)
		fmt.Printf("Value is: '%-v'\n", nums)
		fmt.Printf("Person is: '%-v'\n", person)
		fmt.Printf("Persons is: '%-v'\n", persons)
		fmt.Printf("ComplexPerson is: '%-v'\n", cPerson)
	})
}
