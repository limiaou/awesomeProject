package api

import (
	"awesomeProject/sudouku/util"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type SudokuGenerateProblemResponse struct {
	Status  string      `json:"status"`
	Problem [9][9]uint8 `json:"problem"`
}

func rundomValue(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func SudokuGenerate() (answer [9][9]uint8) {
	problem := [9][9]uint8{
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
	oldProblem := problem
	for {
		// exchange
		// 0: 0and1
		// 1: 0and2
		// 2: 1and2
		for i := 0; i < 50; i++ {
			boxNum := uint8(rundomValue(3))
			exchangeSeed := rundomValue(3)
			var a uint8
			var b uint8
			if exchangeSeed == 0 {
				a = 0 + 3*boxNum
				b = 1 + 3*boxNum
			} else if exchangeSeed == 1 {
				a = 0 + 3*boxNum
				b = 2 + 3*boxNum
			} else {
				a = 1 + 3*boxNum
				b = 2 + 3*boxNum
			}
			var c uint8
			var d uint8
			if rundomValue(2) == 0 {
				c = a
				d = b
			} else {
				c = b
				d = a
			}
			for j := 0; j < 9; j++ {
				tempVal := problem[c][j]
				problem[c][j] = problem[d][j]
				problem[d][j] = tempVal
			}
		}
		for i := 0; i < 65; i++ {
			j := rundomValue(9)
			k := rundomValue(9)
			problem[j][k] = 0
		}
		_, _, _, err := SudokuSolve(problem)
		if err == nil {
			return problem
		}
		problem = oldProblem
	}
}

func SudokuGenerateProblemAPI(c *gin.Context) {
	problem := SudokuGenerate()
	resType := c.Query("type")
	if resType == "img" {
		imgBytes := SudokuGenerateImg(problem)
		util.JPEGStatusOK(c, imgBytes)
		return
	}
	sudokuGenerateProblemResponse := SudokuGenerateProblemResponse{
		Problem: problem,
		Status:  "ok",
	}
	util.JSONStatusOK(c, sudokuGenerateProblemResponse)
	return
}
