package qrcode

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/karmdip-mi/go-fitz"

	"github.com/disintegration/imaging"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

type Service struct {
	ctx *fiber.Ctx
}

func NewService() *Service {
	return &Service{nil}
}

func PdfToImg(fileBytes []byte, filename string) (string, error) {
	imagesPath := ""
	doc, err := fitz.NewFromMemory(fileBytes)
	if err != nil {
		fmt.Printf("fitz.NewFromMemory failed, err: %+v\n", err.Error())
		return "", err
	}

	// Extract pages as images
	for n := 0; n < doc.NumPage(); n++ {
		if n > 0 {
			break
		}
		img, err := doc.Image(n)
		if err != nil {
			fmt.Printf("failed to get image from pdf, err: %+v\n", err.Error())
			return "", err
		}
		width := img.Bounds().Dx()
		// height := img.Bounds().Dy()
		topImg := imaging.Fill(img, 500, 500, imaging.Top, imaging.CatmullRom)
		// resizedImg := imaging.Fit(topRightImg, width/2, height/2, imaging.Lanczos)
		resizedImg := imaging.Resize(topImg, width/2, 0, imaging.CatmullRom)

		err = os.MkdirAll("qrcodes/", 0755)
		if err != nil {
			fmt.Printf("failed create directory qrcodes, err: %+v\n", err.Error())
			return "", err
		}

		// imagesLocation = filepath.Join("qrcodes/", fmt.Sprintf("image-%05d.jpg", n))
		imagesPath = filepath.Join("qrcodes/", fmt.Sprintf("%s.jpg", filename))
		f, err := os.Create(imagesPath)
		if err != nil {
			fmt.Printf("failed create image, err: %+v\n", err.Error())
			return "", err
		}

		err = jpeg.Encode(f, resizedImg, nil)
		if err != nil {
			fmt.Printf("failed encode resizedImg jpeg, err: %+v\n", err.Error())
			return "", err
		}

		f.Close()
	}

	return imagesPath, nil
}

func GetContentFromImg(path string) (*string,error) {
	// open and decode image file
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("os.Open failed, err: %+v\n", err.Error())
		return nil, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("image.Decode failed, err: %+v\n", err.Error())
		return nil, err
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		fmt.Printf("gozxing.NewBinaryBitmapFromImage failed, err: %+v\n", err.Error())
		return nil, err
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		fmt.Printf("qrReader.Decode failed, err: %+v\n", err.Error())
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("no qrcode detect in given image")
	}
	content := result.String()
	return &content,nil
}

func (s *Service) GetQrcodeContentFromPdf() error {
	fmt.Printf("qrcode service GetQrcodeContentFromPdf\n")
	return nil
}
