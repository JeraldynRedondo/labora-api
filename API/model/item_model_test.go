package model

import (
	"testing"
	"time"
)

var itemTest = []Item{
	{1, "John Doe", time.Now(), "Televisor", 2, 150, " ", 0, 0},
	{2, "Jane Smith", time.Now(), "Celular", 3, 200, " ", 0, 0},
	{3, "Michael Smith", time.Now(), "Estufa", 1, 55, " ", 0, 0},
	{4, "Harry Styles", time.Now(), "Ventilador", 3, 70, " ", 0, 0},
	{5, "Zayn Malik", time.Now(), "Mouse", 3, 45, " ", 0, 0},
	{6, "Bob Johnson", time.Now(), "Lavadora", 1, 49, " ", 0, 0},
}

var expectedTotalValue = []int{300, 600, 55, 210, 135, 49}

func TestTotalPriceWorks(t *testing.T) {

	for i, test := range itemTest {

		generatedTotalValue := test.CalculatedTotalPrice()

		if expectedTotalValue[i] != generatedTotalValue {
			t.Errorf("Output %q not equal to expected %q", generatedTotalValue, expectedTotalValue)
		}
	}
}
