package slices

import "fmt"

func RotateRight(s string) string {
	slice1 := []rune(s)
	size := len(slice1)
	slice2 := make([]rune, size)
	last := size - 1
	for i := 0; i < size; i++ {
		if i == 0 {
			slice2[0] = slice1[last]
		} else {
			slice2[i] = slice1[i-1]
		}
	}
	return string(slice2)
}

func RotateLeft(s string) string {
	slice1 := []rune(s)
	size := len(slice1)
	slice2 := make([]rune, size)
	last := size - 1
	for i := 0; i < size; i++ {
		if i == last {
			slice2[last] = slice1[0]
		} else {
			slice2[i] = slice1[i+1]
		}
	}
	return string(slice2)
}

func main() {
	var slice1 string

	slice1 = "1234"

	fmt.Println(slice1)
	resultado := RotateRight(slice1)
	fmt.Println(resultado)
	resultado = RotateLeft(slice1)
	fmt.Println(resultado)

}
