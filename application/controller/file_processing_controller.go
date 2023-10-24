package controller

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/javierBros/backend-balance/application"
	"github.com/javierBros/backend-balance/application/model"
	"github.com/javierBros/backend-balance/application/services"
	"github.com/javierBros/backend-balance/application/utils"
	"io"
	"log"
	"strconv"
	"strings"
)

type FileProcessingController struct {
	summaryProcessingService services.ISummaryProcessingService
	envVariables             *application.EnvironmentVariables
	event                    events.S3Event
	session                  *session.Session
}

func processCSV(data []byte) ([]model.Transaction, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))
	transactions := []model.Transaction{}
	countIteration := 0

	for {
		record, err := reader.Read()
		countIteration++
		if countIteration == 1 {
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf(`err line csv: %v\n`, err.Error())
			return nil, err
		}

		date, err := utils.ConvertMMddStringDateToDate(record[1])
		if err != nil {
			fmt.Printf(`error line csv: Error reading date field.  record[1]->%v, err->%v\n`, record[1], err.Error())
			return nil, errors.New(`error line csv: Error reading date field`)
		}

		amount, err := strconv.ParseFloat(strings.TrimPrefix(record[2], "+"), 64)
		if err != nil {
			fmt.Printf(`error line csv: Error reading amount field:  record[2]->%v, err->%v\n`, record[2], err.Error())
			return nil, errors.New(`error line csv: Error reading amount field`)
		}

		transaction := model.Transaction{
			Date:     date,
			Amount:   amount,
			IsCredit: strings.HasPrefix(record[2], "+"),
		}
		fmt.Printf(`transaction: %v\n`, record)

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func downloadFile(sess *session.Session, bucket, key string) ([]byte, error) {
	s3Client := s3.New(sess)
	output, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Println("Error downloading file from S3:", err)
		return nil, err
	}

	if *output.ContentLength > int64(1*1024*1024) {
		log.Println("File exceeds the maximum allowed size.")
		return nil, fmt.Errorf("File size exceeds the maximum allowed size")
	}

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *FileProcessingController) ProcessBalance() error {

	s3Record := d.event.Records[0].S3
	bucket := s3Record.Bucket.Name
	key := s3Record.Object.Key

	fmt.Printf(`Bucket: %v, Key: %v\n`, bucket, key)

	data, err := downloadFile(d.session, bucket, key)
	if err != nil {
		log.Println("Error downloading file from S3:", err)
		return err
	}

	transactions, err := processCSV(data)
	if err != nil {
		log.Println("Error processiyng CSV:", err)
		return err
	}

	err = d.summaryProcessingService.ProcessSummary(transactions)
	if err != nil {
		return err
	}

	return nil
}

func NewFileProcessingController(summaryProcessingService services.ISummaryProcessingService,
	envVariables *application.EnvironmentVariables,
	event events.S3Event,
	session *session.Session) *FileProcessingController {
	return &FileProcessingController{
		summaryProcessingService: summaryProcessingService,
		envVariables:             envVariables,
		event:                    event,
		session:                  session,
	}
}
