package slices

import (
	"testing"
)

type rotateTest struct {
	arg, expected string
}

var rightRotateTests = []rotateTest{
	rotateTest{"1234", "4123"},
	rotateTest{"hola", "ahol"},
	rotateTest{"amor", "ramo"},
	rotateTest{"arroz", "zarro"},
}
var leftRotateTests = []rotateTest{
	rotateTest{"1234", "2341"},
	rotateTest{"hola", "olah"},
	rotateTest{"amor", "mora"},
	rotateTest{"arroz", "rroza"},
}

func RotateRightWorks(t *testing.T) {

	for _, test := range rightRotateTests {
		if output := RotateRight(test.arg); output != test.expected {
			t.Errorf("Output %q in RotateRight not equal to expected %q", output, test.expected)
		}
	}
}

func RotateLeftWorks(t *testing.T) {

	for _, test := range leftRotateTests {
		if output := RotateRight(test.arg); output != test.expected {
			t.Errorf("Output %q in RotateLeft not equal to expected %q", output, test.expected)
		}
	}
}
