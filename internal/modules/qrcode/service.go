package qrcode

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
	"time"

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

func GetContentFromImg(path string) (*string, error) {
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
	return &content, nil
}

func (s *Service) GetQrcodeContentFromPdf(form *multipart.Form) (map[string]interface{}, error) {
	fmt.Printf("qrcode service GetQrcodeContentFromPdf\n")

	result := map[string]interface{}{}
	resultChan := make(chan struct {
		filename  string
		logNumber *string
		err       error
	})

	var wg sync.WaitGroup
	start := time.Now()
	for formFileName, fileHeaders := range form.File {
		for _, header := range fileHeaders {
			wg.Add(1)
			// process uploaded file here
			go func(head *multipart.FileHeader, wg *sync.WaitGroup) {
				fmt.Printf("fieldName: %+v, fileName: %+v, fileType: %+v, fileSize: %+v\n", formFileName, head.Filename, head.Header["Content-Type"][0], head.Size)

				file, err := head.Open()
				defer file.Close()
				if err != nil {
					fmt.Printf("failed to open file, err: %+v\n", err.Error())
					// return result, err
				}

				fileBytes, err := io.ReadAll(file)
				if err != nil {
					fmt.Printf("failed to read file, err: %+v\n", err.Error())
					// return result, err
				}

				// go func(fb []byte, filename string, wg *sync.WaitGroup) {
				defer wg.Done()
				imagesLocation, err := PdfToImg(fileBytes, head.Filename)
				if err != nil {
					fmt.Printf("PdfToImg failed, err: %+v\n", err)
				}

				logNumber, err := GetContentFromImg(imagesLocation)
				if err != nil {
					fmt.Printf("GetContentFromImg failed, err: %+v\n", err)
				}

				resultChan <- struct {
					filename  string
					logNumber *string
					err       error
				}{head.Filename, logNumber, err}
			}(header, &wg)
		}
	}

	go func() {
		wg.Wait()
		close(resultChan)
		fmt.Printf("duration: %+v\n", time.Since(start))
	}()

	// get the results from chan
	for r := range resultChan {
		if r.err != nil {
			fmt.Printf("result err: %+v --> %+v\n", r.filename, r.err.Error())
			continue
		}
		if r.logNumber == nil {
			fmt.Printf("result: %+v --> %+v\n", r.filename, nil)
			continue
		}
		result[r.filename] = r.logNumber
		fmt.Printf("result: %+v --> %+v\n", r.filename, *r.logNumber)
	}

	return result, nil
}
