package qrcode

import (
	"fmt"
	"golang-api-starter/internal/config"
	"io"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) Controller {
	return Controller{s}
}

var cfg = config.Cfg
var respCode = fiber.StatusInternalServerError

func (c *Controller) GetQrcodeContentFromPdf(ctx *fiber.Ctx) error {
	fmt.Printf("qrcode ctrl GetQrcodeContentFromPdf\n")
	form, err := ctx.MultipartForm()
	if err != nil { /* handle error */
		fmt.Printf("failed to get multipartForm, err: %+v\n", err.Error())
		return err
	}

	result := make(chan struct {
		filename  string
		logNumber *string
	})

	var wg sync.WaitGroup
	start := time.Now()
	for formFieldName, fileHeaders := range form.File {
		for _, header := range fileHeaders {
			wg.Add(1)
			// process uploaded file here
			fmt.Printf("fieldName: %+v, fileName: %+v, fileType: %+v, fileSize: %+v\n", formFieldName, header.Filename, header.Header["Content-Type"][0], header.Size)

			file, err := header.Open()
			defer file.Close()
			if err != nil {
				fmt.Printf("failed to open file, err: %+v\n", err.Error())
				return err
			}

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				fmt.Printf("failed to read file, err: %+v\n", err.Error())
				return err
			}

			go func(fb []byte, filename string, wg *sync.WaitGroup) {
				defer wg.Done()
				imagesLocation := PdfToImg(fb, filename)
				logNumber := GetContentFromImg(imagesLocation)
				result <- struct {
					filename  string
					logNumber *string
				}{filename, logNumber}
			}(fileBytes, header.Filename, &wg)
		}
	}

	go func() {
		wg.Wait()
		close(result)
		fmt.Printf("duration: %+v\n", time.Since(start))
	}()

	// get the results from chan
	for r := range result {
		if r.logNumber == nil {
			fmt.Printf("result: %+v --> %+v\n", r.filename, nil)
			continue
		}
		fmt.Printf("result: %+v --> %+v\n", r.filename, *r.logNumber)
	}

	c.service.GetQrcodeContentFromPdf()
	return nil
}
