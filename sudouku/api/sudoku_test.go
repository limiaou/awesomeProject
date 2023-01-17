package api

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSudokuSolve(t *testing.T) {
	grid := [9][9]uint8{
		{0, 6, 1, 0, 0, 7, 0, 0, 3},
		{0, 9, 2, 0, 0, 3, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 8, 5, 3, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 5, 0, 4},
		{5, 0, 0, 0, 0, 8, 0, 0, 0},
		{0, 4, 0, 0, 0, 0, 0, 0, 1},
		{0, 0, 0, 1, 6, 0, 8, 0, 0},
		{6, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	answer, _, _, _ := SudokuSolve(grid)
	correct := [9][9]uint8{
		{4, 6, 1, 9, 8, 7, 2, 5, 3},
		{7, 9, 2, 4, 5, 3, 1, 6, 8},
		{3, 8, 5, 2, 1, 6, 4, 7, 9},
		{1, 2, 8, 5, 3, 4, 7, 9, 6},
		{9, 3, 6, 7, 2, 1, 5, 8, 4},
		{5, 7, 4, 6, 9, 8, 3, 1, 2},
		{8, 4, 9, 3, 7, 5, 6, 2, 1},
		{2, 5, 3, 1, 6, 9, 8, 4, 7},
		{6, 1, 7, 8, 4, 2, 9, 3, 5},
	}
	if reflect.DeepEqual(answer, correct) == false {
		t.Error("数独の答えがあっていません｡")
	}
}

func TestSudokuCheck(t *testing.T) {
	nullCannotSolveSudokuResponse := CannotSolveSudokuResponse{}
	noSol := [9][9]uint8{
		{0, 6, 1, 0, 0, 7, 0, 0, 3},
		{0, 9, 2, 0, 0, 3, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 8, 5, 3, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 5, 0, 4},
		{5, 0, 0, 0, 0, 8, 0, 0, 0},
		{0, 4, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 6, 0, 8, 0, 0},
		{6, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	answer, reason, cannotSolveSudokuResponse, err := SudokuSolve(noSol)
	if answer == [9][9]uint8{} || reason == "" || cannotSolveSudokuResponse != nullCannotSolveSudokuResponse || err == nil {
		t.Error("数独が解無しの場合を正しくチェックできていません｡")
	}
	duplicatedRow := [9][9]uint8{
		{0, 6, 1, 0, 0, 7, 0, 6, 3},
		{0, 9, 2, 0, 0, 3, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 8, 5, 3, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 5, 0, 4},
		{5, 0, 0, 0, 0, 8, 0, 0, 0},
		{0, 4, 0, 0, 0, 0, 0, 0, 1},
		{0, 0, 0, 1, 6, 0, 8, 0, 0},
		{6, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	answer, reason, cannotSolveSudokuResponse, err = SudokuSolve(duplicatedRow)
	if answer == [9][9]uint8{} || reason == "" || cannotSolveSudokuResponse != nullCannotSolveSudokuResponse || err == nil {
		t.Error("数独が行方向に重複している場合を正しくチェックできていません｡")
	}
	duplicatedCol := [9][9]uint8{
		{0, 6, 1, 0, 0, 7, 0, 0, 3},
		{0, 9, 2, 0, 0, 3, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 8, 5, 3, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 5, 0, 4},
		{5, 0, 0, 0, 0, 8, 0, 0, 0},
		{0, 4, 0, 0, 0, 0, 0, 0, 1},
		{0, 4, 0, 1, 6, 0, 8, 0, 0},
		{6, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	answer, reason, cannotSolveSudokuResponse, err = SudokuSolve(duplicatedCol)
	if answer == [9][9]uint8{} || reason == "" || cannotSolveSudokuResponse != nullCannotSolveSudokuResponse || err == nil {
		t.Error("数独が列方向に重複している場合を正しくチェックできていません｡")
	}
	duplicatedInBox := [9][9]uint8{
		{0, 6, 1, 0, 0, 7, 0, 0, 3},
		{0, 9, 2, 0, 0, 3, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 8, 5, 3, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 5, 0, 4},
		{5, 0, 0, 0, 0, 8, 0, 0, 0},
		{0, 4, 0, 0, 0, 0, 0, 0, 1},
		{0, 0, 0, 1, 6, 0, 8, 0, 0},
		{6, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	answer, reason, cannotSolveSudokuResponse, err = SudokuSolve(duplicatedInBox)
	if answer == [9][9]uint8{} || reason == "" || cannotSolveSudokuResponse != nullCannotSolveSudokuResponse || err == nil {
		t.Error("数独が箱内が重複する場合を正しくチェックできていません｡")
	}
}

func BenchmarkSudokuSolve(b *testing.B) {
	grid := [9][9]uint8{
		{0, 6, 1, 0, 0, 7, 0, 0, 3},
		{0, 9, 2, 0, 0, 3, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 8, 5, 3, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 5, 0, 4},
		{5, 0, 0, 0, 0, 8, 0, 0, 0},
		{0, 4, 0, 0, 0, 0, 0, 0, 1},
		{0, 0, 0, 1, 6, 0, 8, 0, 0},
		{6, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	SudokuSolve(grid)
}

func TestSudokuGenerateImg(t *testing.T) {
	generate := SudokuGenerate()
	//img := SudokuGenerateImg(generate)
	fmt.Printf("%v \n", generate)
}

func TestSudokuGenerateProblemAPI(t *testing.T) {
	//api := SudokuGenerateProblemAPI(*gin.Context)

}
