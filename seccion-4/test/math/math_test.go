package math

import (
	"fmt"
	"testing"
)

type factTest struct {
	arg, expected int
}

var factTests = []factTest{
	factTest{2, 2},
	factTest{4, 24},
	factTest{5, 120},
	factTest{0, 1},
}

func TestAdd(t *testing.T) {

	for _, test := range factTests {
		if output := Factorial(test.arg); output != test.expected {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}
}

func ExampleFact() {
	fmt.Println(Factorial(4))
	// Output: 24
}

/*
// arg1 means argument 1 and arg2 means argument 2, and the expected stands for the 'result we expect'
type addTest struct {
	arg1, arg2, expected int
}

var addTests = []addTest{
	addTest{2, 3, 5},
	addTest{4, 8, 12},
	addTest{6, 9, 15},
	addTest{3, 10, 13},
}

func TestAdd(t *testing.T) {

	for _, test := range addTests {
		if output := Add(test.arg1, test.arg2); output != test.expected {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(4, 6)
	}
}

func ExampleAdd() {
	fmt.Println(Add(4, 6))
	// Output: 10
}
*/
