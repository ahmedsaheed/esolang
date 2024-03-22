package object

import "testing"

func TestStringHashKey(t *testing.T) {
	helloFirst := &String{Value: "Hello World"}
	helloSecond := &String{Value: "Hello World"}
	diffFirst := &String{Value: "Goodbye World"}
	diffSecond := &String{Value: "Goodbye World"}

	if helloFirst.HashKey() != helloSecond.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diffFirst.HashKey() != diffSecond.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if helloFirst.HashKey() == diffFirst.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}
