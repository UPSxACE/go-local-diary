package app_config

import (
	"html/template"
	"reflect"
	"testing"
)

func TestDefaultFuncList(t *testing.T) {
	type TestCase struct {
		input  []interface{}
		output []interface{}
	}

	// Defining test cases
	testCases := []TestCase{
		{
			input:  []interface{}{"item1", "item2", "item3"},
			output: []interface{}{"item1", "item2", "item3"},
		},
		{
			input:  []interface{}{1, 2, 3, 4},
			output: []interface{}{1, 2, 3, 4},
		},
		{
			input:  []interface{}{true, false, true, false},
			output: []interface{}{true, false, true, false},
		},
		{
			input:  []interface{}{map[string]string{"key1": "val1", "key2": "val2"}},
			output: []interface{}{map[string]string{"key1": "val1", "key2": "val2"}},
		},
	}

	// Running the test cases
	for _, tc := range testCases {
		output := list(tc.input...)

		// Check if the output matches the expected result
		if !reflect.DeepEqual(tc.output, output) {
			t.Errorf("Expected %v, got %v", tc.output, output)
		}
	}
}

func TestDefaultFuncObj(t *testing.T) {
	type TestCase struct {
		input  string
		output map[string]string
	}

	// Defining test cases
	testCases := []TestCase{
		{
			input:  "key1:val1,,key2:val2",
			output: map[string]string{"key1": "val1", "key2": "val2"},
		},
		{
			input:  "ss:this is a sentence,,com:sentence with, comma!",
			output: map[string]string{"ss": "this is a sentence", "com": "sentence with, comma!"},
		},
	}

	// Running the test cases
	for _, tc := range testCases {
		output := obj(tc.input)

		// Check if the output matches the expected result
		if !reflect.DeepEqual(tc.output, output) {
			t.Errorf("Expected %v, got %v", tc.output, output)
		}
	}
}

func TestDefaultFuncObjError(t *testing.T) {
	type TestCase struct {
		input   string
		errCode int
	}

	// Defining test cases
	testCases := []TestCase{
		{
			input:   "key1::val1,,key2:val2",
			errCode: 1,
		},
		{
			input:   "ss:this is a sentence,com:sentence with, comma!",
			errCode: 1,
		},
		{
			input:   "ac",
			errCode: 1,
		},
	}

	// Running the test cases
	for _, tc := range testCases {
		func() {
			defer func(code int) {
				r := recover()
				if r == nil {
					t.Errorf("The code did not panic")
				}

				if reflect.TypeOf(r) != reflect.TypeOf(&DefMapInvalidArgs{}) {
					t.Errorf("Not the expected error type")
				} else {
					rConverted, ok := r.(*DefMapInvalidArgs)

					if !ok {
						t.Errorf("Conversion to the correct error type in defer function failed")
					}

					if ok && rConverted.code != code {
						t.Errorf("Panicked with the wrong code")
					}
				}
			}(tc.errCode)
			obj(tc.input)
		}()
	}
}

func TestDefaultFuncSum(t *testing.T) {
	type TestCase struct {
		input  [2]int
		output int
	}

	// Defining test cases
	testCases := []TestCase{
		{
			input:  [2]int{1, 2},
			output: 3,
		},
		{
			input:  [2]int{4,5},
			output: 9,
		},
		{
			input:  [2]int{-1, 2},
			output: 1,
		},
		{
			input:  [2]int{-1, -5},
			output: -6,
		},
		{
			input:  [2]int{0,0},
			output: 0,
		},
	}

	// Running the test cases
	for _, tc := range testCases {
		output := sum(tc.input[0], tc.input[1])

		// Check if the output matches the expected result
		if !reflect.DeepEqual(tc.output, output) {
			t.Errorf("Expected %v, got %v", tc.output, output)
		}
	}
}

func TestDefaultFuncSumStr(t *testing.T) {
	type TestCase struct {
		input  [2]string
		output int
	}

	// Defining test cases
	testCases := []TestCase{
		{
			input:  [2]string{"1", "2"},
			output: 3,
		},
		{
			input:  [2]string{"4","5"},
			output: 9,
		},
		{
			input:  [2]string{"-1", "2"},
			output: 1,
		},
		{
			input:  [2]string{"-1", "-5"},
			output: -6,
		},
		{
			input:  [2]string{"0","0"},
			output: 0,
		},
	}

	// Running the test cases
	for _, tc := range testCases {
		output := sumStr(tc.input[0], tc.input[1])

		// Check if the output matches the expected result
		if !reflect.DeepEqual(tc.output, output) {
			t.Errorf("Expected %v, got %v", tc.output, output)
		}
	}
}

func TestDefaultFuncSumStrError(t *testing.T) {
	type TestCase struct {
		input  [2]string
		errCode int
	}

	// Defining test cases
	testCases := []TestCase{
		{
			input:  [2]string{"1.2", "2"},
			errCode: 2,
		},
		{
			input:  [2]string{"4.9","5"},
			errCode: 2,
		},
		{
			input:  [2]string{"-1", "2,"},
			errCode: 2,
		},
		{
			input:  [2]string{"abc", "-5"},
			errCode: 2,
		},
		{
			input:  [2]string{"0,2","0"},
			errCode: 2,
		},
	}

	// Running the test cases
	for _, tc := range testCases {
		func() {
			defer func(code int) {
				r := recover()
				if r == nil {
					t.Errorf("The code did not panic")
				}

				if reflect.TypeOf(r) != reflect.TypeOf(&DefMapInvalidArgs{}) {
					t.Errorf("Not the expected error type")
				} else {
					rConverted, ok := r.(*DefMapInvalidArgs)

					if !ok {
						t.Errorf("Conversion to the correct error type in defer function failed")
					}

					if ok && rConverted.code != code {
						t.Errorf("Panicked with the wrong code")
					}
				}
			}(tc.errCode)
			sumStr(tc.input[0], tc.input[1])
		}()
	}
}

func TestDefaultFuncHtmlBreaks(t *testing.T){
	type TestCase struct {
		input  string
		output template.HTML
	}

	// Defining test cases
	testCases := []TestCase{
		{
			input:  "This is a sentece\nAnother line starts",
			output: "This is a sentece<br>Another line starts",
		},
		{
			input:  "This is a sentece\nAnother line starts\nAnd another line",
			output: "This is a sentece<br>Another line starts<br>And another line",
		},
		{
			input:  "\nThis is a sentece\nAnother line starts\nAnd another line\\n",
			output: "<br>This is a sentece<br>Another line starts<br>And another line\\n",
		},
	}

	// Running the test cases
	for _, tc := range testCases {
		output := htmlBreaks(tc.input)
		// Check if the output matches the expected result
		if !reflect.DeepEqual(tc.output, output) {
			t.Errorf("Expected %v, got %v", tc.output, output)
		}
	}
}