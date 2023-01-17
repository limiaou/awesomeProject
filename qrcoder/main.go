//func main() {
//err := qrcode.WriteFile("你好，我是玲娜贝儿，我正在森林探险。现在给我转49.9，待我今天肯德基疯狂星期四购买老北京鸡肉卷补充体力后，今晚，我来你家探险。\n📞贝儿向您发起了探险邀请\n\n🔘不接受     丨     🔘不想接受", qrcode.Medium, 256, "qr.jpg")
//_, err = qrcode.Encode("૮꒰ ˶• ༝ •˶꒱ა ", qrcode.Medium, 256)
//err = qrcode.WriteColorFile("૮꒰ ˶• ༝ •˶꒱ა ", qrcode.Medium, 256, color.White, color.Black, "qr.jpg")
//CreateQrCodeWithLogo("你好，我是玲娜贝儿，我正在森林探险。现在给我转49.9，待我今天肯德基疯狂星期四购买老北京鸡肉卷补充体力后，今晚，我来你家探险。\n📞贝儿向您发起了探险邀请\n\n🔘不接受     丨     🔘不想接受", "C:\\Users\\Helen.Wang\\Desktop\\bell.jpg", "qr.jpg", qrcode.Medium, 256)
//if err != nil {
//	fmt.Println("write error")
//}
//}
//
//func read() {
//	//tuotoo/qrcode识别二维码
//	fi, err := os.Open("qr2.png")
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	defer fi.Close()
//	qrmatrix, err := qrcode.Encode(fi)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	fmt.Println(qrmatrix)
//}

package main

import (
	"flag"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/LyricTian/logger"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
)

var (
	text    string
	logo    string
	percent int
	size    int
	out     string
)

func init() {
	flag.StringVar(&text, "t", "你好，我是玲娜贝儿，我正在森林探险。现在给我转49.9，待我今天肯德基疯狂星期四购买老北京鸡肉卷补充体力后，今晚，我来你家探险。\\n📞贝儿向您发起了探险邀请\\n\\n🔘接受     丨     🔘想接受", "二维码内容")
	flag.StringVar(&logo, "l", "C:\\Users\\Helen.Wang\\Desktop\\opacity50Percent.png", "二维码Logo(png)")
	flag.IntVar(&percent, "p", 100, "二维码Logo的显示比例(默认15%)")
	flag.IntVar(&size, "s", 128, "二维码的大小(默认256)")
	flag.StringVar(&out, "o", "test1.jpg", "输出文件")
}

func main1() {
	flag.Parse()

	if text == "" {
		logger.Fatalf("请指定二维码的生成内容")
	}

	if out == "" {
		logger.Fatalf("请指定输出文件")
	}

	if exists, err := checkFile(out); err != nil {
		logger.Fatalf("检查输出文件发生错误：%s", err.Error())
	} else if exists {
		logger.Fatalf("输出文件已经存在，请重新指定")
	}

	code, err := qrcode.New(text, qrcode.Highest)
	if err != nil {
		logger.Fatalf("创建二维码发生错误：%s", err.Error())
	}

	srcImage := code.Image(size)
	if logo != "" {
		logoSize := float64(size) * float64(percent) / 100

		srcImage, err = addLogo(srcImage, logo, int(logoSize))
		if err != nil {
			logger.Fatalf("增加Logo发生错误：%s", err.Error())
		}
	}

	outAbs, err := filepath.Abs(out)
	if err != nil {
		logger.Fatalf("获取输出文件绝对路径发生错误：%s", err.Error())
	}

	os.MkdirAll(filepath.Dir(outAbs), 0777)
	outFile, err := os.Create(outAbs)
	if err != nil {
		logger.Fatalf("创建输出文件发生错误：%s", err.Error())
	}
	defer outFile.Close()

	jpeg.Encode(outFile, srcImage, &jpeg.Options{Quality: 100})

	logger.Infof("二维码生成成功，文件路径：%s", outAbs)
}

func checkFile(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resizeLogo(logo string, size uint) (image.Image, error) {
	file, err := os.Open(logo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	img = resize.Resize(size, size, img, resize.Lanczos3)
	return img, nil
}

func addLogo(srcImage image.Image, logo string, size int) (image.Image, error) {
	logoImage, err := resizeLogo(logo, uint(size))
	if err != nil {
		return nil, err
	}

	offset := image.Pt((srcImage.Bounds().Dx()-logoImage.Bounds().Dx())/2, (srcImage.Bounds().Dy()-logoImage.Bounds().Dy())/2)
	b := srcImage.Bounds()
	m := image.NewNRGBA(b)
	draw.Draw(m, b, srcImage, image.ZP, draw.Src)
	draw.Draw(m, logoImage.Bounds().Add(offset), logoImage, image.ZP, draw.Over)

	return m, nil
}

//CreateQrCodeWithLogo 带logo的二维码图片生成 content-二维码内容   level-容错级别,Low,Medium,High,Highest   size-像素单位  outPath-输出路径  logoPath-logo文件路径
func CreateQrCodeWithLogo(content, logoPath, outPath string, level qrcode.RecoveryLevel, size int) interface{} {
	code, err := qrcode.New(content, level)
	if err != nil {
		return err.Error()
	}
	//设置文件大小并创建画板
	qrcodeImg := code.Image(size)
	outImg := image.NewRGBA(qrcodeImg.Bounds())

	//读取logo文件
	logoFile, err := os.Open(logoPath)
	if err != nil {
		panic(err)
	}
	//logoImg, err := png.Decode(logoFile)
	logoImg, err := jpeg.Decode(logoFile)
	logoImg = resize.Resize(uint(size/6), uint(size/6), logoImg, resize.Lanczos3)

	//logo和二维码拼接
	draw.Draw(outImg, outImg.Bounds(), qrcodeImg, image.Pt(0, 0), draw.Over)
	//offset := image.Pt(0, 0)
	offset := image.Pt((outImg.Bounds().Max.X-logoImg.Bounds().Max.X)/2, (outImg.Bounds().Max.Y-logoImg.Bounds().Max.Y)/2)
	draw.Draw(outImg, outImg.Bounds().Add(offset), logoImg, image.Pt(0, 0), draw.Over)

	//draw.DrawMask(outImg, outImg.Bounds().Add(offset), logoImg, image.Pt(0, 0), outImg, image.Pt(0, 0), draw.Over)

	f, err := os.Create(outPath)
	if err != nil {
		return err.Error()
	}
	png.Encode(f, outImg)
	return nil
}
