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

func PdfToImg(fileBytes []byte, filename string) string {
	imagesPath := ""
	doc, err := fitz.NewFromMemory(fileBytes)
	if err != nil {
		panic(err)
	}

	// Extract pages as images
	for n := 0; n < doc.NumPage(); n++ {
		if n > 0 {
			break
		}
		img, err := doc.Image(n)
		if err != nil {
			panic(err)
		}
		width := img.Bounds().Dx()
		// height := img.Bounds().Dy()
		topImg := imaging.Fill(img, 500, 500, imaging.Top, imaging.CatmullRom)
		// resizedImg := imaging.Fit(topRightImg, width/2, height/2, imaging.Lanczos)
		resizedImg := imaging.Resize(topImg, width/2, 0, imaging.CatmullRom)

		err = os.MkdirAll("qrcodes/", 0755)
		if err != nil {
			panic(err)
		}

		// imagesLocation = filepath.Join("qrcodes/", fmt.Sprintf("image-%05d.jpg", n))
		imagesPath = filepath.Join("qrcodes/", fmt.Sprintf("%s.jpg", filename))
		f, err := os.Create(imagesPath)
		if err != nil {
			panic(err)
		}

		err = jpeg.Encode(f, resizedImg, nil)
		if err != nil {
			panic(err)
		}

		f.Close()
	}

	return imagesPath
}

func GetContentFromImg(path string) *string {
	// open and decode image file
	file, _ := os.Open(path)
	img, _, _ := image.Decode(file)

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, _ := qrReader.Decode(bmp, nil)

	if result == nil {
		return nil
	}
	content := result.String()
	return &content
}

func (s *Service) GetQrcodeContentFromPdf() error {
	fmt.Printf("qrcode service GetQrcodeContentFromPdf\n")
	return nil
}
