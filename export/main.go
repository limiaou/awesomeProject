package export

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var fileName string
	fileNames := make([]string, 0)
	for i := 0; i < len(fileNames); i++ {
		fileNames = append(fileNames, fmt.Sprint(fileName, "(", i, ")", ".csv"))
	}
	err := ZipFiles(fmt.Sprint(fileName+".zip"), fileNames)
	if err != nil {
		log.Println(err)
	}
}

// createZip 压缩文件
func createZip(zipName string, originFileName string) (*os.File, error) {
	fileNameSplitList := strings.Split(zipName, ".xlsx")
	fileName := fmt.Sprintf("%s.zip", fileNameSplitList[0])
	zipFile, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer zipFile.Close()
	writer := zip.NewWriter(zipFile)
	defer writer.Close()
	csvFile, err := os.Open(originFileName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer csvFile.Close()
	info, err := csvFile.Stat()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	head, err := zip.FileInfoHeader(info)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	head.Method = zip.Deflate
	fw, err := writer.CreateHeader(head)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = io.Copy(fw, csvFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return zipFile, nil
}

func ZipFiles(filename string, files []string) error {
	fmt.Println("start zip file......")
	// 创建输出文件目录
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()
	// 创建空的zip档案，可以理解为打开zip文件，准备写入
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	// 打开要压缩的文件
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()
	// 获取文件的描述
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	// FileInfoHeader返回一个根据fi填写了部分字段的Header，可以理解成是将fileinfo转换成zip格式的文件信息
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filename

	// 预定义压缩算法。
	// archive/zip包中预定义的有两种压缩方式。一个是仅把文件写入到zip中。不做压缩。一种是压缩文件然后写入到zip中。默认的Store模式。就是只保存不压缩的模式。
	// Store   unit16 = 0  //仅存储文件
	// Deflate unit16 = 8  //压缩文件

	header.Method = zip.Deflate
	// 创建压缩包头部信息
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	// 将源复制到目标，将fileToZip 写入writer   是按默认的缓冲区32k循环操作的，不会将内容一次性全写入内存中,这样就能解决大文件的问题
	_, err = io.Copy(writer, fileToZip)
	return err
}

func DeleteFile(fileNames []string) error {
	var delFile string
	for _, delFile = range fileNames {
		err := os.Remove(delFile)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	err := os.Remove(fmt.Sprint(delFile, ".zip"))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
//	Unzip 将解压缩 zip 存档，将 zip 文件（参数 1）中的所有文件和文件夹移动到输出目录（参数 2）。
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File 创建目录
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
