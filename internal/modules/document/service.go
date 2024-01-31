package document

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"golang-api-starter/internal/helper"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Service struct {
	repo *Repository
	ctx  *fiber.Ctx
}

func NewService(r *Repository) *Service {
	return &Service{r, nil}
}

func (s *Service) GetIdMap(documents Documents) map[string]*Document {
	documentMap := map[string]*Document{}
	for _, document := range documents {
		documentMap[document.GetId()] = document
	}
	return documentMap
}

func (s *Service) Get(queries map[string]interface{}) ([]*Document, *helper.Pagination) {
	fmt.Printf("document service get\n")
	return s.repo.Get(queries)
}

func (s *Service) GetById(queries map[string]interface{}) ([]*Document, error) {
	fmt.Printf("document service getById\n")

	records, _ := s.repo.Get(queries)
	if len(records) == 0 {
		return nil, fmt.Errorf("%s with id: %s not found", tableName, queries["id"])
	}
	return records, nil
}

func (s *Service) Create(form *multipart.Form) ([]*Document, *helper.HttpErr) {
	fmt.Printf("document service create\n")
	timer := helper.Timer(time.Now())
	defer timer()

	/* create upload folder if not exists */
	baseUploadDir := "./uploads"
	if _, err := os.Stat(baseUploadDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(baseUploadDir, os.ModePerm)
		if err != nil {
			log.Println("upload path create failed ", err)
		}
	}

	documents := []*Document{}
	documentsMap := map[string]string{} // for keep track on duplicated same file in form.File["file"]

	/* extract files from the form-data and copy them into ./uploads */
	if form.File["file"] == nil {
		return nil, &helper.HttpErr{fiber.StatusUnprocessableEntity, fmt.Errorf("key \"file\" missing")}
	}

	claims := s.ctx.Locals("claims").(jwt.MapClaims)
	for _, fh := range form.File["file"] {
		file, err := fh.Open()
		if err != nil {
			log.Println("failed to open file", err)
			return nil, &helper.HttpErr{fiber.StatusUnprocessableEntity, err}
		}
		defer file.Close()

		t := time.Now()
		filename := fmt.Sprintf("%s-%s", t.Format("20060102150405"), fh.Filename)
		uploadPath := fmt.Sprintf("%s/%s", baseUploadDir, filename)
		out, err := os.OpenFile(uploadPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("failed to create file", err)
			return nil, &helper.HttpErr{fiber.StatusInternalServerError, err}
		}
		defer out.Close()
		// fmt.Printf("file?: %T\n", file)

		hash := sha1.New()
		f := io.TeeReader(file, hash)

		_, copyError := io.Copy(out, f)
		if copyError != nil {
			log.Println("failed to copy file", copyError)
			return nil, &helper.HttpErr{fiber.StatusInternalServerError, copyError}
		}
		// fmt.Println("uploaded to ", uploadPath)

		sha1Sum := hex.EncodeToString(hash.Sum(nil))
		// fmt.Println("file hash: ", sha1Sum, hash.Sum(nil))

		document := &Document{
			Name:     fh.Filename,
			FilePath: uploadPath,
			FileType: strings.Split(fh.Header["Content-Type"][0], "/")[1],
			FileSize: fh.Size,
			Hash:     sha1Sum,
			Public:   true,
		}

		prevUploadPath, exists := documentsMap[sha1Sum]
		if exists {
			os.Remove(uploadPath)
			document.FilePath = prevUploadPath
			documentsMap[sha1Sum] = prevUploadPath
		} else {
			documentsMap[sha1Sum] = uploadPath
		}
		recordsWithSameHash, _ := s.repo.Get(map[string]interface{}{"hash": sha1Sum})
		// fmt.Println("sameRecord",recordsWithSameHash,len(recordsWithSameHash) )

		if document.UserId == nil {
			document.UserId = claims["userId"]
		}
		if validErr := helper.ValidateStruct(*document); validErr != nil {
			return nil, &helper.HttpErr{fiber.StatusUnprocessableEntity, validErr}
		}

		/* use same file, remove newly uploaded same file */
		if len(recordsWithSameHash) > 0 {
			os.Remove(document.FilePath)
			document.FilePath = recordsWithSameHash[0].FilePath
		}

		documents = append(documents, document)
	}

	results, err := s.repo.Create(documents)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Update(documents []*Document) ([]*Document, *helper.HttpErr) {
	fmt.Printf("document service update\n")
	results, err := s.repo.Update(documents)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Delete(ids []string) ([]*Document, error) {
	fmt.Printf("document service delete\n")
	var (
		records    = []*Document{}
		conditions = map[string]interface{}{}
	)

	cfg.LoadEnvVariables()
	if cfg.DbConf.Driver == "mongodb" {
		conditions["_id"] = ids
	} else {
		conditions["id"] = ids
	}

	records, _ = s.repo.Get(conditions)
	fmt.Printf("records: %+v\n", records)
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to delete, %s with id: %+v not found", tableName, ids)
	}

	return records, s.repo.Delete(ids)
}

func (s *Service) GetDocument(queries map[string]interface{}) ([]byte, string, string, error) {
	fmt.Printf("GetDocument service\n")
	var size int64 = 0
	if queries["size"] != nil {
		size, _ = strconv.ParseInt(queries["size"].(string), 10, 64)
	}
	repoData, _ := s.repo.Get(queries)

	if len(repoData) == 0 {
		return nil, "", "", fmt.Errorf("not found")
	}

	// fmt.Printf("filePath: %+v\n", *repoData.DocumentRequest.FilePath)
	f, err := os.Open(repoData[0].FilePath)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to open file, %+v", err.Error())
	}
	defer f.Close()

	fileType, err := GetFileContentType(f)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get file type, %+v", err.Error())
	}
	fmt.Println("fileType: ", fileType)
	fileBytes, fileErr := os.ReadFile(repoData[0].FilePath)
	if fileErr != nil {
		return nil, "", "", fmt.Errorf("failed to get file type, %+v", fileErr.Error())
	}

	fileName := repoData[0].Name
	/* use imaging to resize the image by given size for thumbnail */
	jpgPngRegex := regexp.MustCompile(`png|jpg|jpeg|jpe`)
	if size != 0 && strings.Contains(fileType, "image") && jpgPngRegex.MatchString(fileType) {
		buf := new(bytes.Buffer)
		img, _, err := image.Decode(bytes.NewReader(fileBytes))
		if err != nil {
			log.Fatalln("image.Decode err: ", err)
		}
		resizedImg := imaging.Resize(img, int(size), 0, imaging.Lanczos)
		err = jpeg.Encode(buf, resizedImg, nil)
		if err != nil {
			log.Fatalln("jpeg.Encode err: ", err)
		}

		return buf.Bytes(), "image/jpeg", fileName, nil
	}

	return fileBytes, fileType, fileName, nil
}

func GetFileContentType(ouput *os.File) (string, error) {
	// to sniff the content type only the first 512 bytes are used.
	buf := make([]byte, 512)

	_, err := ouput.Read(buf)

	if err != nil {
		return "", err
	}

	/* get mime type */
	contentType := http.DetectContentType(buf)

	/* reset the file point to beginning for further actions like decode img etc. */
	ouput.Seek(0, 0)

	return contentType, nil
}
