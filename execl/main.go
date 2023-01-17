package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func create() {
	f := excelize.NewFile()
	// 创建一个工作表
	f.NewSheet("Sheet2")

	//index = f.GetSheetIndex("Sheet1")
	//fmt.Println(index)
	//index := f.NewSheet("Sheet2")
	// 设置单元格的值
	//f.SetCellValue("Sheet2", "A2", "Hello world.")
	//f.SetCellValue("Sheet1", "B2", 100)
	// 设置工作簿的默认工作表
	//f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs("Test2.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func open(filepath string) {
	file, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Printf("Open file failed", "%v\n", err)
	}
	defer file.Close()
	//defer func() {
	//	if err := f.Close(); err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//value, err := file.GetCellValue("Sheet1", "B2")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(value)
	// 获取 Sheet1 上所有单元格
	rows, err := file.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range rows {
		for _, c := range v {
			fmt.Println(c)
		}
	}
}

//
// categories
//  @Description: 插入图表
//
func categories() {
	file, err := excelize.OpenFile("Test1.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	if file.GetSheetIndex("Sheet2") != -1 {
		file.NewSheet("Sheet2")
	}
	categories := map[string]string{
		"A2": "Small", "A3": "Normal", "A4": "Large", "B1": "Apple", "C1": "Orange", "D1": "Pear"}
	values := map[string]int{
		"B2": 2, "C2": 3, "D2": 3, "B3": 5, "C3": 2, "D3": 4, "B4": 6, "C4": 7, "D4": 8}
	f := excelize.NewFile()
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	for k, v := range values {
		f.SetCellValue("Sheet1", k, v)
	}
	if err := f.AddChart("Sheet2", "E1", `{
        "type": "col3DClustered",
        "series": [
        {
            "name": "Sheet1!$A$2",
            "categories": "Sheet1!$B$1:$D$1",
            "values": "Sheet1!$B$2:$D$2"
        },
        {
            "name": "Sheet1!$A$3",
            "categories": "Sheet1!$B$1:$D$1",
            "values": "Sheet1!$B$3:$D$3"
        },
        {
            "name": "Sheet1!$A$4",
            "categories": "Sheet1!$B$1:$D$1",
            "values": "Sheet1!$B$4:$D$4"
        }],
        "title":
        {
            "name": "Fruit 3D Clustered Column Chart"
        }
    }`); err != nil {
		fmt.Println(err)
		return
	}
	// 根据指定路径保存文件
	if err := f.Save(); err != nil {
		//if err := f.SaveAs("Test1.xlsx"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success")
	}
}

func image() {
	file, err := excelize.OpenFile("Test1.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	// 插入图片
	if err := file.AddPicture("Sheet1", "A32", "C:\\Users\\Helen.Wang\\Desktop\\mmm.jpg", ""); err != nil {
		fmt.Println(err)
	}
	// 在工作表中插入图片，并设置图片的缩放比例
	/*	if err := file.AddPicture("Sheet1", "A2", "C:\\Users\\Helen.Wang\\Desktop\\mmm.jpg",
		`{
				"x_scale": 0.5
				"y_scale": 0.5
				}`); err != nil {
		fmt.Println(err)
	}*/
	// 在工作表中插入图片，并设置图片的打印属性
	/*if err := file.AddPicture("Sheet1", "H2", "image.gif", `{
	        "x_offset": 15,
	        "y_offset": 10,
	        "print_obj": true,
	        "lock_aspect_ratio": false,
	        "locked": false
	    }`); err != nil {
			fmt.Println(err)
		}*/
	// 保存文件
	if err = file.Save(); err != nil {
		fmt.Println(err)
	}
}

func Test() {
	file, err := excelize.OpenFile("C:\\Users\\Helen.Wang\\Desktop\\saas-uat.xlsx")
	if err != nil {
		fmt.Printf("Open file failed", "%v\n", err)
	}
	defer file.Close()
	rows, err := file.GetRows("Sheet1")
	for _, v := range rows {
		for _, c := range v {
			fmt.Println(c)
		}
	}
	cart := make(map[int][]string)
	for i, c := range rows {
		if _, ok := cart[i]; !ok {
			cart[i] = c
			fmt.Println(cart[i])
		}
	}
}

func main() {
	create()
	//time.Sleep(10 * time.Second)
	categories()
}
