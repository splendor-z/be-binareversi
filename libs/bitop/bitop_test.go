package bitop

import (
	"testing"
)

func Test01_ApplyBitOperation_Addition(t *testing.T) {
	row := [8]int{7, 1, 0, 1, 1, 7, 7, 7}
	value := 4
	operator := "+"
	expected := [8]int{7, 1, 1, 1, 1, 7, 7, 7}

	newRow, err := ApplyBitOperation(row, value, operator)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if newRow != expected {
		t.Errorf("Expected %v, got %v", expected, newRow)
	}
}

func Test02_ApplyBitOperation_Multiplication(t *testing.T) {
	row := [8]int{7, 1, 0, 1, 1, 7, 7, 7}
	value := 3
	operator := "*"
	// 1011 * 3 = 33 = 100001
	// 最下位4bit = 1000
	expected := [8]int{7, 1, 0, 0, 0, 7, 7, 7}

	newRow, err := ApplyBitOperation(row, value, operator)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if newRow != expected {
		t.Errorf("Expected %v, got %v", expected, newRow)
	}
}

func Test03_ApplyBitOperation_NoTargets(t *testing.T) {
	row := [8]int{7, 7, 7, 7, 7, 7, 7, 7}
	value := 1
	operator := "+"
	expected := row

	newRow, err := ApplyBitOperation(row, value, operator)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if newRow != expected {
		t.Errorf("Expected no changes, got %v", newRow)
	}
}

func Test04_ApplyBitOperation_Padding(t *testing.T) {
	row := [8]int{7, 1, 0, 7, 7, 7, 7, 7}
	value := 1
	operator := "+"
	// 10 + 1 = 11 -> "11"
	expected := [8]int{7, 1, 1, 7, 7, 7, 7, 7}

	newRow, err := ApplyBitOperation(row, value, operator)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if newRow != expected {
		t.Errorf("Expected %v, got %v", expected, newRow)
	}
}

func Test05_ApplyBitOperation_UnsupportedOperator(t *testing.T) {
	row := [8]int{7, 1, 0, 1, 1, 7, 7, 7}
	value := 1
	operator := "-"
	_, err := ApplyBitOperation(row, value, operator)
	if err == nil {
		t.Error("Expected error for unsupported operator")
	}
}

// func Test06_ApplyBitOperation_LargeResultOverflow(t *testing.T) {
// 	row := [8]int{1, 0, 0, 1}
// 	value := 100
// 	operator := "+"
// 	// 1001 + 100 = 109 = 1101101
// 	// 最上位4bit = 1101
// 	expected := [8]int{1, 1, 0, 1} // <- LSB側に詰める実装なのでこれが正解

// 	newRow, err := ApplyBitOperation(row, value, operator)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// 	if !reflect.DeepEqual(newRow[:4], expected[:4]) {
// 		t.Errorf("Expected %v, got %v", expected[:4], newRow[:4])
// 	}
// }

func Test07_ApplyBitOperation_WithInterleaved7s(t *testing.T) {
	row := [8]int{7, 7, 1, 7, 1, 0, 1, 7}
	value := 1
	operator := "+"
	// 1,1,0,1 → 1101 + 1 = 1110
	// 適用順: index 2,4,5,6 → 1,1,1,0
	expected := [8]int{7, 7, 1, 7, 1, 1, 0, 7}

	newRow, err := ApplyBitOperation(row, value, operator)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if newRow != expected {
		t.Errorf("Expected %v, got %v", expected, newRow)
	}
}

func Test08_ApplyBitOperation_MultiplyByZero(t *testing.T) {
	row := [8]int{1, 0, 1, 1, 7, 7, 7, 7}
	value := 0
	operator := "*"
	// any x 0 = 0 → 0000
	expected := [8]int{0, 0, 0, 0, 7, 7, 7, 7}

	newRow, err := ApplyBitOperation(row, value, operator)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if newRow != expected {
		t.Errorf("Expected %v, got %v", expected, newRow)
	}
}

// func Test09_ApplyBitOperation_MultiplicationOverflow(t *testing.T) {
// 	row := [8]int{1, 1, 0, 1}
// 	value := 9 // 1101 * 9 = 1000001
// 	operator := "*"
// 	// 最上位4bit = 1000
// 	expected := [8]int{1, 0, 0, 0}

// 	newRow, err := ApplyBitOperation(row, value, operator)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// 	if !reflect.DeepEqual(newRow[:4], expected[:4]) {
// 		t.Errorf("Expected %v, got %v", expected[:4], newRow[:4])
// 	}
// }

func Test10_ApplyBitOperation_Multiplication_Interleaved7(t *testing.T) {
	row := [8]int{7, 1, 7, 0, 1, 7, 7, 7}
	value := 3
	operator := "*"
	// bits: 1,0,1 → 101 * 3 = 1111 → 上位3bit = 111
	// index: 1,3,4 → 1,1,1
	expected := [8]int{7, 1, 7, 1, 1, 7, 7, 7}

	newRow, err := ApplyBitOperation(row, value, operator)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if newRow != expected {
		t.Errorf("Expected %v, got %v", expected, newRow)
	}
}
