package testhelper

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func printFailed(){
	_, file, no, ok := runtime.Caller(2)
    if ok {
        fmt.Printf("    Test failed on %s #%d\n", file, no)
    }
}

func ExpectDifferent(t *testing.T, outputToTest any, expectedOutput any){
	if reflect.DeepEqual(outputToTest, expectedOutput) {
		printFailed()
		t.Fatalf("Expected %#v to be different", outputToTest)
	}
}

func ExpectEqual(t *testing.T, outputToTest any, expectedOutput any){
	if !reflect.DeepEqual(outputToTest, expectedOutput) {
		printFailed()
		t.Fatalf("Expected %#v, got %#v", expectedOutput, outputToTest)
	}
}

func ExpectError(t *testing.T, errorToTest error) {
	if errorToTest == nil {
		printFailed()
		t.Fatal("Expected an error")
	}
}

func ExpectNoError(t *testing.T, errorToTest error) {
	if errorToTest != nil {
		printFailed()
		t.Fatal(errorToTest)
	}
}

func ExpectPanic(t *testing.T, callback func()) (didPanic bool, recoverValue any){
	didPanic = false
	recoverValue = nil;
	func() {
		defer func() {
			r := recover()
			if r != nil {
				didPanic = true;
				recoverValue = r;
			}
			if r == nil {
				printFailed()
				t.Fatal("Expected it to panic")
			}
		}()
		callback()
	}()
	return didPanic, recoverValue
}

func ExpectType[T any](t *testing.T, thingToTest any) (convertedThing T){
	converted, ok := thingToTest.(T);
	var zeroVal T;
	if(!ok){
		printFailed()
		t.Fatalf("Expected object of type %T, got object of type %T", zeroVal, thingToTest)
	}
	return converted;
}

