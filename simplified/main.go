package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	dirname := "C:\\Users\\Helen.Wang\\Desktop\\change"
	dir, err := ioutil.ReadDir(dirname)
	//getwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error reading directory failed: ", err)
	}
	fileName := make([]string, 0)
	//rs := make([]string, 0)
	//rs = append(rs, rS)
	file, err := os.Create("test.csv")
	file.WriteString("\xEF\xBB\xBF")
	//err = csv.NewWriter(file).Write(rs)
	//if err != nil {
	//	fmt.Println("Error writing file", err)
	//}
	for i, f := range dir {
		if !f.IsDir() {
			fileName = append(fileName, f.Name())
			fmt.Println(fileName[i])
			open, err := os.Open(fmt.Sprint(dirname + "\\" + fileName[i]))
			if err != nil {
				fmt.Println("Error opening file: ", err)
			}
			//for {
			reader := csv.NewReader(open)
			reader.FieldsPerRecord = -1
			all, err := reader.ReadAll()
			if err != nil {
				panic(err)
			}
			r := make([]string, 0)
			for _, v := range all {
				for _, k := range v {
					r = append(r, k)
					file.WriteString(fmt.Sprint(k + ","))
				}
				file.WriteString("\n")
				fmt.Println("r[%d]:%s", i, r)
			}

			//rd := bufio.NewReader(open)
			//rS, err := rd.ReadString('\n')
			//if err != nil || io.EOF == err {
			//	break
			//}
			//_, err = file.WriteString(rS)
			//if err != nil {
			//	fmt.Println("Error writing file", err)
			//}
			//}
		}
	}
}
