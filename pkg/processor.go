package pkg

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func ProcessTestRecords(records []JUnitTestSuites) error {
	var (
		timeId        = fmt.Sprintf("%s", time.Now().Format("2006-01-02-15:04:05"))
		passCount     int32
		failCount     int32
		skippedCount  int32
		totalExecTime float64
		status        int32
	)
	for _, set := range records {
		for _, suite := range set.Suites {
			if len(suite.TestCases) == 0 {
				log.Printf("Invalid testsuite. Please check: %v, %v", suite)
				continue
			}
			// this is our class name logic, you should update it based on your org's naming conventions
			classNameArr := strings.Split(suite.TestCases[0].Classname, ".") //Classname = {product}.{sub-product}.{service/server}.{api/feature}
			for _, test := range suite.TestCases {
				testName := test.Name
				execTime, _ := strconv.ParseFloat(test.Time, 64)

				product := classNameArr[0]
				subProductLine := classNameArr[1]
				service := classNameArr[2]
				api := classNameArr[len(classNameArr)-1]

				totalExecTime += execTime
				if test.Failure != nil {
					status = StatusFail
					failCount++
				} else if test.SkipMessage != nil {
					status = StatusSkipped
					skippedCount++
				} else {
					status = StatusPass
					passCount++
				}
				err := WriteToInflux(AllRecordsMeasurement, map[string]string{
					"run_id": timeId,
				}, map[string]interface{}{
					"product":          product,
					"sub_product_line": subProductLine,
					"service":          service,
					"api":              api,
					"case":             testName,
					"exec_time":        execTime,
					"status":           status,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
