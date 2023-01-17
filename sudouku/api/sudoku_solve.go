package api

import (
	"awesomeProject/sudouku/util"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SudokuSolveRequest struct {
	Problem [9][9]uint8 `json:"problem"`
}

type SudokuSolveOKResponse struct {
	Status string      `json:"status"`
	Answer [9][9]uint8 `json:"answer"`
	Time   float64     `json:"time"`
}

type CannotSolveSudokuResponse struct {
	Status    string   `json:"status"`
	Reason    string   `json:"reason"`
	Message   string   `json:"message"`
	NGIndexes [2]uint8 `json:"ngIndexes"`
	NGValue   uint8    `json:"ngValue"`
}

type XStruct struct {
	KeyName string
	Values  [2]uint8
}

type XArray []XStruct

type YStruct struct {
	R uint8
	C uint8
	N uint8
}

func sudokuValidCheck(problem [9][9]uint8, isValidCheckAPI bool) (solution []YStruct, reason string, cannotSolveSudokuResponse CannotSolveSudokuResponse, err error) {
	N := uint8(9)
	X := [324]XStruct{}
	keyNames := [4]string{"rc", "rn", "cn", "bn"}
	k := uint16(0)
	tate := [9]map[uint8]bool{}
	box1 := map[uint8]bool{}
	box2 := map[uint8]bool{}
	box3 := map[uint8]bool{}
	for i := uint8(0); i < uint8(N); i++ {
		tate[i] = map[uint8]bool{}
	}
	for i := uint8(0); i < uint8(N); i++ {
		yoko := map[uint8]bool{}
		if i%3 == 0 {
			box1 = map[uint8]bool{}
			box2 = map[uint8]bool{}
			box3 = map[uint8]bool{}
		}
		for j := uint8(0); j < uint8(N); j++ {
			if problem[i][j] > 9 {
				reason = fmt.Sprintf("value: %d in ([row, col] = [%d, %d]) is above 9.\nplease input appropriate value to ([row, col] = [%d, %d]", problem[i][j], i, j, i, j)
				err = errors.New("can't solve sudoku")
				if isValidCheckAPI == false {
					return nil, reason, cannotSolveSudokuResponse, err
				}
				cannotSolveSudokuResponse = CannotSolveSudokuResponse{
					Status:    "ng",
					Reason:    reason,
					Message:   err.Error(),
					NGIndexes: [2]uint8{i, j},
					NGValue:   problem[i][j],
				}
				return nil, reason, cannotSolveSudokuResponse, err
			}
			if problem[i][j] > 0 {
				if _, isDuplicate := yoko[problem[i][j]]; isDuplicate == false {
					yoko[problem[i][j]] = true
				} else {
					reason = fmt.Sprintf("value: %d in ([row, col] = [%d, %d]) is horizontally duplicated .\nplease input appropriate value to ([row, col] = [%d, %d]", problem[i][j], i, j, i, j)
					err = errors.New("can't solve sudoku")
					if isValidCheckAPI == false {
						return nil, reason, cannotSolveSudokuResponse, err
					}
					cannotSolveSudokuResponse = CannotSolveSudokuResponse{
						Status:    "ng",
						Reason:    reason,
						Message:   err.Error(),
						NGIndexes: [2]uint8{i, j},
						NGValue:   problem[i][j],
					}
					return nil, reason, cannotSolveSudokuResponse, err
				}
				if _, isDuplicate := tate[j][problem[i][j]]; isDuplicate == false {
					tate[j][problem[i][j]] = true
				} else {
					reason = fmt.Sprintf("value: %d in ([row, col] = [%d, %d]) is vertically duplicated .\nplease input appropriate value to ([row, col] = [%d, %d]", problem[i][j], i, j, i, j)
					err = errors.New("can't solve sudoku")
					if isValidCheckAPI == false {
						return nil, reason, cannotSolveSudokuResponse, err
					}
					cannotSolveSudokuResponse = CannotSolveSudokuResponse{
						Status:    "ng",
						Reason:    reason,
						Message:   err.Error(),
						NGIndexes: [2]uint8{i, j},
						NGValue:   problem[i][j],
					}
					return nil, reason, cannotSolveSudokuResponse, err
				}
				box := map[uint8]bool{}
				if j < 3 {
					box = box1
				} else if j < 6 {
					box = box2
				} else {
					box = box3
				}
				if _, isDuplicate := box[problem[i][j]]; isDuplicate == false {
					box[problem[i][j]] = true
				} else {
					reason = fmt.Sprintf("value: %d in ([row, col] = [%d, %d]) is duplicated in this box.\nplease input appropriate value to ([row, col] = [%d, %d]", problem[i][j], i, j, i, j)
					err = errors.New("can't solve sudoku")
					if isValidCheckAPI == false {
						return nil, reason, cannotSolveSudokuResponse, err
					}
					cannotSolveSudokuResponse = CannotSolveSudokuResponse{
						Status:    "ng",
						Reason:    reason,
						Message:   err.Error(),
						NGIndexes: [2]uint8{i, j},
						NGValue:   problem[i][j],
					}
					return nil, reason, cannotSolveSudokuResponse, err
				}
			}
			X[k] = XStruct{KeyName: keyNames[0], Values: [2]uint8{i, j}}
			k++
		}
	}
	for h := uint8(0); h < 3; h++ {
		for i := uint8(0); i < uint8(N); i++ {
			for j := uint8(1); j < uint8(N+1); j++ {
				X[k] = XStruct{KeyName: keyNames[h+1], Values: [2]uint8{i, j}}
				k++
			}
		}
	}
	Y := map[YStruct][4]XStruct{}
	for h := uint8(0); h < uint8(N); h++ {
		for i := uint8(0); i < uint8(N); i++ {
			for j := uint8(1); j < uint8(N+1); j++ {
				R, C := uint8(3), uint8(3)
				b := (h/R)*R + (i / C)
				yStruct := YStruct{h, i, j}
				v := [2]uint8{h, i}
				rc := XStruct{KeyName: keyNames[0], Values: v}
				v = [2]uint8{h, j}
				rn := XStruct{KeyName: keyNames[1], Values: v}
				v = [2]uint8{i, j}
				cn := XStruct{KeyName: keyNames[2], Values: v}
				v = [2]uint8{b, j}
				bn := XStruct{KeyName: keyNames[3], Values: v}
				Y[yStruct] = [4]XStruct{rc, rn, cn, bn}
			}
		}
	}
	Z := map[XStruct][9]YStruct{}
	exactCover(X, Y, Z)
	zeroCount := uint8(0)
	for i := uint8(0); i < uint8(len(problem)); i++ {
		for j := uint8(0); j < uint8(len(problem[i])); j++ {
			if problem[i][j] == 0 {
				zeroCount++
			} else {
				choice(Z, Y, YStruct{i, j, problem[i][j]})
			}
		}
	}
	solution = make([]YStruct, zeroCount)
	reason, cannotSolveSudokuResponse, err = solve(Z, Y, solution, isValidCheckAPI)
	if err != nil {
		return nil, reason, cannotSolveSudokuResponse, err
	}
	return solution, "", CannotSolveSudokuResponse{}, nil
}

func exactCover(X [324]XStruct, Y map[YStruct][4]XStruct, Z map[XStruct][9]YStruct) {
	for j := uint16(0); j < uint16(len(X)); j++ {
		XJ := &X[j]
		yStruct := [9]YStruct{}
		i := uint8(0)
		for key, vv := range Y {
			for k := uint8(0); k < uint8(len(vv)); k++ {
				if *XJ == vv[k] {
					yStruct[i] = key
					if i == 8 {
						Z[*XJ] = yStruct
						i = 0
					} else {
						i++
					}
				}
			}
		}
	}
}

func choice(Z map[XStruct][9]YStruct, Y map[YStruct][4]XStruct, r YStruct) {
	isNullYStruct := YStruct{0, 0, 0}
	for i := uint8(0); i < uint8(len(Y[r])); i++ {
		for _, vv := range Z[Y[r][i]] {
			if vv == isNullYStruct {
				continue
			}
			for _, vvv := range Y[vv] {
				if vvv != Y[r][i] {
					ZV := Z[vvv]
					for j := uint8(0); j < uint8(len(ZV)); j++ {
						if ZV[j] == vv {
							ZV[j] = YStruct{0, 0, 0}
							Z[vvv] = ZV
							break
						}
					}
				}
			}
		}
		delete(Z, Y[r][i])
	}
}

func solve(Z map[XStruct][9]YStruct, Y map[YStruct][4]XStruct, solution []YStruct, isValidCheckAPI bool) (reason string, cannotSolveSudokuResponse CannotSolveSudokuResponse, err error) {
	i := uint8(0)
	for i < uint8(len(solution)) {
		isExistSolution := false
		isNullYStruct := YStruct{0, 0, 0}
		for key, v := range Z {
			counter := uint8(0)
			for k := uint8(0); k < uint8(len(v)); k++ {
				if v[k] != isNullYStruct {
					counter++
				}
			}
			min := uint8(10)
			if counter < min {
				min = counter
				if min == 1 {
					for j := uint8(0); j < uint8(len(Z[key])); j++ {
						if Z[key][j] != isNullYStruct {
							solution[i] = Z[key][j]
							choice(Z, Y, solution[i])
							i++
							isExistSolution = true
							break
						}
					}
					break
				}
			}
		}
		if isExistSolution == false {
			reason = "this is not appropriate sudoku because more than one answer exists.\nso please input the appropriate sudoku."
			err = errors.New("can't solve sudoku")
			if isValidCheckAPI == false {
				return reason, cannotSolveSudokuResponse, err
			}
			cannotSolveSudokuResponse = CannotSolveSudokuResponse{
				Status:    "ng",
				Reason:    reason,
				Message:   err.Error(),
				NGIndexes: [2]uint8{10, 10},
				NGValue:   10,
			}
			return reason, cannotSolveSudokuResponse, err
		}
	}
	return "", CannotSolveSudokuResponse{}, nil
}

func SudokuSolve(problem [9][9]uint8) (answer [9][9]uint8, reason string, cannotSolveSudokuResponse CannotSolveSudokuResponse, err error) {
	solution, reason, cannotSolveSudokuResponse, err := sudokuValidCheck(problem, false)
	answer = problem
	for i := uint8(0); i < uint8(len(solution)); i++ {
		answer[solution[i].R][solution[i].C] = solution[i].N
	}
	return answer, reason, cannotSolveSudokuResponse, err
}

func SudokuSolveAPI(c *gin.Context) {
	sudokuSolveRequest := SudokuSolveRequest{}
	err := c.BindJSON(&sudokuSolveRequest)
	if err != nil {
		util.LogUnexpectedErr(err)
		return
	}
	startTime := time.Now()
	answer, reason, _, err := SudokuSolve(sudokuSolveRequest.Problem)
	time := time.Now().Sub(startTime).Seconds()
	if err != nil {
		errStruct := util.ErrStruct{
			Errors: []util.OneErrorStruct{
				util.OneErrorStruct{
					Domain:  "localhost",
					Reason:  reason,
					Message: err.Error(),
				},
			},
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		util.APIErr(c, errStruct, sudokuSolveRequest)
		return
	}
	resType := c.Query("type")
	if resType == "img" {
		imgBytes := SudokuGenerateImg(answer)
		util.JPEGStatusOK(c, imgBytes)
		return
	}
	sudokuSolveOKResponse := SudokuSolveOKResponse{
		Status: "ok",
		Answer: answer,
		Time:   time,
	}
	util.JSONStatusOK(c, sudokuSolveOKResponse)
	return
}
