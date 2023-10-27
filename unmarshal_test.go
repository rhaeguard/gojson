package gojson

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	var inputJson = `{
    "value": [
        1239,
        123.45
    ],
    "name": "renault",
    "token": true,
    "hello": null
}`
	t.Run("check json", func(t *testing.T) {
		json, err := parseJson(inputJson)
		if err != nil {
			t.FailNow()
		}
		expected := JsonValue{
			ValueType: OBJECT,
			Value: map[string]JsonValue{
				"value": {
					ValueType: ARRAY,
					Value: []JsonValue{
						{ValueType: NUMBER, Value: float64(1239)},
						{ValueType: NUMBER, Value: float64(123.45)},
					},
				},
				"name": {
					ValueType: STRING,
					Value:     "renault",
				},
				"token": {
					ValueType: BOOL,
					Value:     true,
				},
				"hello": {
					ValueType: NULL,
					Value:     nil,
				},
			},
		}

		if !reflect.DeepEqual(json, expected) {
			t.Fail()
		}
	})
}

func TestNumber(t *testing.T) {
	var numberCandidates = map[string]JsonValue{
		"12345":            {ValueType: NUMBER, Value: float64(12345)},
		"2500":             {ValueType: NUMBER, Value: float64(2500)},
		"0":                {ValueType: NUMBER, Value: float64(0)},
		"-123":             {ValueType: NUMBER, Value: float64(-123)},
		"10":               {ValueType: NUMBER, Value: float64(10)},
		"1234567890123456": {ValueType: NUMBER, Value: float64(1234567890123456)},
		"3.14159":          {ValueType: NUMBER, Value: 3.14159},
		"-0.005":           {ValueType: NUMBER, Value: -0.005},
		"-0.000123":        {ValueType: NUMBER, Value: -0.000123},
		"1.23456789012345": {ValueType: NUMBER, Value: 1.23456789012345},
		"0.3":              {ValueType: NUMBER, Value: 0.3},
		"3.14":             {ValueType: NUMBER, Value: 3.14},
	}

	for inputJson, expected := range numberCandidates {
		name := fmt.Sprintf("numbers(%s)", inputJson)
		t.Run(name, func(t *testing.T) {
			json, err := parseJson(inputJson)
			if err != nil {
				t.FailNow()
			}
			fmt.Printf("%s\n", json)

			if !reflect.DeepEqual(json, expected) {
				t.FailNow()
			}
		})
	}

}

func TestErrorHandling(t *testing.T) {
	var inputJson = `{
    "value": [
        1239,
        12345
    ],
    "name": "renault",
    "token": true,
    "hello": nill
}`
	t.Run("error check", func(t *testing.T) {
		if _, err := parseJson(inputJson); err == nil {
			t.Logf("error value was required")
			t.Fail()
		}
	})
}

func TestErrorHandling2(t *testing.T) {
	var inputJson = `{
    "value": [
        1239,
        12345,
    "name": "renault",
    "token": true,
    "hello": null
}`
	t.Run("error check", func(t *testing.T) {
		if _, err := parseJson(inputJson); err == nil {
			t.Logf("error value was required")
			t.Fail()
		} else {
			fmt.Printf("%s\n", err.Error())
		}
	})
}
