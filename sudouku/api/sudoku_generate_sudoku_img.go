package api

import (
	"awesomeProject/sudouku/util"
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang/freetype"
)

func addLabel(img *image.RGBA, x, y int, label string) {
	fontBytes, err := ioutil.ReadFile("./HonyaJi-Re.ttf")
	if err != nil {
		util.LogUnexpectedErr(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}
	c := freetype.NewContext()
	c.SetDst(img)
	c.SetFont(f)
	c.SetFontSize(90.0)
	c.SetSrc(image.Black)
	c.SetClip(img.Bounds())
	pt := freetype.Pt(x, y)
	_, err = c.DrawString(label, pt)
	if err != nil {
		util.LogUnexpectedErr(err)
	}
}

func hLine(img *image.RGBA, x1, y, x2 int, col color.Color) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

func vLine(img *image.RGBA, x, y1, y2 int, col color.Color) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

func fillRect(img *image.RGBA, col color.Color) {
	rect := img.Rect
	for h := rect.Min.Y; h < rect.Max.Y; h++ {
		for v := rect.Min.X; v < rect.Max.X; v++ {
			img.Set(v, h, col)
		}
	}
}

func SudokuGenerateImg(problem [9][9]uint8) (imgBytes []byte) {
	x := 0
	y := 0
	size := 20
	length := 47 * size
	outSideSpace := 200
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}
	img := image.NewRGBA(image.Rect(x-outSideSpace, y-outSideSpace, length+outSideSpace, length+outSideSpace))
	fillRect(img, white)
	j := 0
	borderWidth := 3
	for i := size; i < length; i += (length - (size)*2) / 9 {
		for k := 0; k < borderWidth; k++ {
			hLine(img, size, i+k, (length - (size)), black)
			vLine(img, i+k, size, (length - (size)), black)
			hLine(img, size, i-k, (length - (size)), black)
			vLine(img, i-k, size, (length - (size)), black)
		}

		for k := 0; k < 9; k++ {
			if j < 9 {
				if problem[j][k] > 0 {
					num := fmt.Sprint(problem[j][k])
					addLabel(img, (size*2)+(100*k), i+size*4, num)
				}
			}
		}
		if j%3 == 0 {
			for k := borderWidth; k < borderWidth*2; k++ {
				hLine(img, size, i+k, (length - (size)), black)
				vLine(img, i+k, size, (length - (size)), black)
				hLine(img, size, i-k, (length - (size)), black)
				vLine(img, i-k, size, (length - (size)), black)
			}
		}
		j++
	}
	b := bytes.Buffer{}
	writer := bufio.NewWriter(&b)
	err := jpeg.Encode(writer, img, &jpeg.Options{Quality: 100})
	if err != nil {
		util.LogUnexpectedErr(err)
	}
	return b.Bytes()
}

func SudokuGenerateImgAPI(c *gin.Context) {
	sudokuSolveRequest := SudokuSolveRequest{}
	err := c.BindJSON(&sudokuSolveRequest)
	if err != nil {
		util.LogUnexpectedErr(err)
		return
	}
	imgBytes := SudokuGenerateImg(sudokuSolveRequest.Problem)
	util.JPEGStatusOK(c, imgBytes)
	return
}
