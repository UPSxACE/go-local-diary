package app

import (
	"html/template"
	"reflect"
	"testing"

	"github.com/UPSxACE/go-local-diary/utils/testhelper"
)

func TestDefaultFuncList(t *testing.T) {
	type Input = []any
	type Output = []any

	th := testhelper.CreateTestHelper[Input, Output]()

	th.AddTestcase([]any{"item1", "item2", "item3"},nil)
	th.AddTestcase([]any{1, 2, 3, 4}, nil)
	th.AddTestcase([]any{true, false, true, false},nil)
	th.AddTestcase([]interface{}{map[string]string{"key1": "val1", "key2": "val2"}}, nil)

	for _, input := range th.GetInputs() {
		output := list(input...)
		testhelper.ExpectEqual(t, input, output)
	}
}

func TestDefaultFuncObj(t *testing.T) {
	th := testhelper.CreateTestHelper[string,any]()

	th.AddTestcase("key1:val1,,key2:val2", map[string]string{"key1": "val1", "key2": "val2"})
	th.AddTestcase("ss:this is a sentence,,com:sentence with, comma!", map[string]string{"ss": "this is a sentence", "com": "sentence with, comma!"})

	for index, input := range th.GetInputs() {
		output := obj(input)
		testhelper.ExpectEqual(t, output, th.GetOutput(index))
	}
}

func TestDefaultFuncObjError(t *testing.T) {
	type Input = string;
	var ExpectedErrorCode = 1

	th := testhelper.CreateTestHelper[Input, any]()
	th.AddTestcase("key1::val1,,key2:val2", nil)
	th.AddTestcase("ss:this is a sentence,com:sentence with, comma!", nil)
	th.AddTestcase("ac", nil)

	for _, input := range th.GetInputs() {
		didPanic, r := testhelper.ExpectPanic(t, func(){
			obj(input)
		})
		if(didPanic){
			correctType := reflect.TypeOf(r) == reflect.TypeOf(&DefMapInvalidArgs{})
			if (!correctType) {
				t.Fatalf("Not the expected error type")
			} 
			if (correctType){
				rConverted, ok := r.(*DefMapInvalidArgs)

				if !ok {
					t.Fatalf("Conversion to the correct error type in defer function failed")
				}
				
				if ok && rConverted.code != ExpectedErrorCode {
					t.Fatalf("Panicked with the wrong code")
				}
			}
		}
	}
}

func TestDefaultFuncSum(t *testing.T) {
	th := testhelper.CreateTestHelper[[2]int, int]()

	th.AddTestcase([2]int{1, 2},3)
	th.AddTestcase([2]int{4, 5},9)
	th.AddTestcase([2]int{-1, 2},1)
	th.AddTestcase([2]int{-1, -5},-6)
	th.AddTestcase([2]int{0, 0},0)

	for index, input := range th.GetInputs() {
		output := sum(input[0], input[1])

		testhelper.ExpectEqual(t, output, th.GetOutput(index))
	}
}

func TestDefaultFuncSumStr(t *testing.T) {
	th := testhelper.CreateTestHelper[[2]string, int]()

	th.AddTestcase([2]string{"1", "2"},3)
	th.AddTestcase([2]string{"4", "5"},9)
	th.AddTestcase([2]string{"-1", "2"},1)
	th.AddTestcase([2]string{"-1", "-5"},-6)
	th.AddTestcase([2]string{"0", "0"},0)

	for index, input := range th.GetInputs() {
		output := sumStr(input[0], input[1])

		testhelper.ExpectEqual(t, output, th.GetOutput(index))
	}
}

func TestDefaultFuncSumStrError(t *testing.T) {
	type Input = [2]string;
	var ExpectedErrorCode = 2

	th := testhelper.CreateTestHelper[Input, any]()
	th.AddTestcase([2]string{"1.2", "2"}, nil)
	th.AddTestcase([2]string{"4.9", "5"}, nil)
	th.AddTestcase([2]string{"-1", "2,"}, nil)
	th.AddTestcase([2]string{"abc", "-5"}, nil)
	th.AddTestcase([2]string{"0,2", "0"}, nil)

	for _, input := range th.GetInputs() {
		didPanic, r := testhelper.ExpectPanic(t, func(){
			sumStr(input[0], input[1])
		})
		if(didPanic){
			correctType := reflect.TypeOf(r) == reflect.TypeOf(&DefMapInvalidArgs{})
			if (!correctType) {
				t.Fatalf("Not the expected error type")
			} 
			if (correctType){
				rConverted, ok := r.(*DefMapInvalidArgs)

				if !ok {
					t.Fatalf("Conversion to the correct error type in defer function failed")
				}
				
				if ok && rConverted.code != ExpectedErrorCode {
					t.Fatalf("Panicked with the wrong code")
				}
			}
		}
	}
}

func TestDefaultFuncHtmlBreaks(t *testing.T) {
	type Input = string
	type Output = template.HTML
	
	th := testhelper.CreateTestHelper[Input, Output]()

	th.AddTestcase("This is a sentece\nAnother line starts", "This is a sentece<br>Another line starts")
	th.AddTestcase("This is a sentece\nAnother line starts\nAnd another line", "This is a sentece<br>Another line starts<br>And another line")
	th.AddTestcase("\nThis is a sentece\nAnother line starts\nAnd another line\\n", "<br>This is a sentece<br>Another line starts<br>And another line\\n")

	for index, input := range th.GetInputs() {
		output := htmlBreaks(input)
		testhelper.ExpectEqual(t, output, th.GetOutput(index))
	}
}
