package main

import (
	"fmt"
	"time"
	"slices"
	"os"

	"github.com/joho/godotenv"
	
	"github.com/lukinkratas/t212-to-digrin-cli-go/internal/t212"
	"github.com/lukinkratas/t212-to-digrin-cli-go/internal/dataframe"
	"github.com/lukinkratas/t212-to-digrin-cli-go/internal/utils"
)


func GetInputDt() *string {

	// var currentDt time.Time = 
	// var previousMonthDt time.Time = 
	var previousMonthDtStr string = time.Now().AddDate(0, -1, 0).Format("2006-01")

	var inputDtStr *string = new(string)
	fmt.Println("Reporting Year Month in \"YYYY-mm\" format: ")
	fmt.Printf("Or confirm default \"%v\" by ENTER.\n", previousMonthDtStr)
	fmt.Scanln(inputDtStr)

	if *inputDtStr == "" {
		inputDtStr = &previousMonthDtStr
	}

	return inputDtStr

}


func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var bucketName string = os.Getenv("BUCKET_NAME")

	t212Client := t212.APIClient{os.Getenv("T212_API_KEY")}

	var inputDtStr *string = GetInputDt()

	// var inputDt time.Time
	// inputDt, err = time.Parse("2006-01", *inputDtStr)
	// if err != nil {
	// 	panic(err)
	// }

	// var fromDt time.Time = utils.GetFirstDayOfMonth(&inputDt)
	// var toDt time.Time = utils.GetFirstDayOfNextMonth(&inputDt)

	// var createdReportId *uint 

	// for {
		
	// 	createdReportId = t212Client.CreateReport(&fromDt, &toDt)

	// 	if *createdReportId != 0 {
	// 		break
	// 	}

	// 	time.Sleep(10 * time.Second)

	// }
	
	// createdReportId Mock Up
	var createdReportId *uint = new(uint)
	*createdReportId = 1695548

	fmt.Printf("  *createdReportId: %v\n", *createdReportId)

	// optimized wait time for report creation
    time.Sleep(10 * time.Second)

	var createdReport t212.Report

	for {

		var reportsList []t212.Report = t212Client.ListReports()

		// report list is empty
		if len(reportsList) == 0 {
			time.Sleep(10 * time.Second)
			continue
		}
		
		var startTime time.Time = time.Now()

		// reverse order for loop, cause latest export is expected to be at the end
		slices.Reverse(reportsList)

		for _, report := range reportsList {

			if report.Id == *createdReportId {
		        createdReport = report
				break
			}

		}

		fmt.Printf("  Took %v\n", time.Since(startTime))


		if createdReport.Status == "Finished" {
			break
		}

	}

	// // downloadLink Mock Up
	// downloadLink = "https://tzswiy3zk5dms05cfeo.s3.eu-central-1.amazonaws.com/from_2025-03-01_to_2025-04-01_MTc0MzU4MDY0MDE0Mw.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20250402T075723Z&X-Amz-SignedHeaders=host&X-Amz-Expires=604799&X-Amz-Credential=AKIARJCCZCDEKCUWYOXG%2F20250402%2Feu-central-1%2Fs3%2Faws4_request&X-Amz-Signature=857a3b30cb532fdc0d52137a8af7602cbdfd84f597de0c74f61727403c71be3c"

	fmt.Printf("  createdReport.downloadLink: %v\n", createdReport.DownloadLink)

	var t212CsvEncoded []byte = createdReport.Download()

	var fileName string = fmt.Sprintf("%s.csv", *inputDtStr)

	var keyName string = fmt.Sprintf("t212/%s", fileName)
	utils.S3PutObject(t212CsvEncoded, &bucketName, &keyName)

	var t212DataFrame []dataframe.Schema = dataframe.DecodeCSV(t212CsvEncoded)
	t212DataFrame = dataframe.Transform(t212DataFrame)

	var digrinCsvEncoded []byte = dataframe.Encode(t212DataFrame)
	keyName = fmt.Sprintf("digrin/%s", fileName)
	utils.S3PutObject(digrinCsvEncoded, &bucketName, &keyName)
	var digrinCsvUrl *string = utils.S3GetPresignedURL(&bucketName, &keyName)
	fmt.Printf("  *digrinCsvUrl: %v\n", *digrinCsvUrl)

}
