package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	ReSetCsv("zbx_problems_export (6).csv")
}

func ReSetCsv(filename string) {
	create, err := os.Create("temp_" + filename)
	if err != nil {
		panic(err)
	}
	open, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(open)
	writer := bufio.NewWriter(create)
	_, err = writer.WriteString("\xEF\xBB\xBF")
	if err != nil {
		panic(err)
	}
	var exec string
	for {
		exec, err = reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		_, err = writer.WriteString(exec)
		if err != nil {
			panic(err)
		}
		err = writer.Flush()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("转码完成")
	err = open.Close()
	if err != nil {
		panic(err)
	}
	err = create.Close()
	if err != nil {
		panic(err)
	}
	err = os.Remove(filename)
	if err != nil {
		panic(err)
	}
	err = os.Rename("temp_"+filename, filename)
	if err != nil {
		panic(err)
	}
}
