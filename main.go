package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/luohuahuang/junit-influx-client/pkg"
	"io/ioutil"
	"os"
)

// part of the credit to https://github.com/jstemmer/go-junit-report/blob/master/formatter/formatter.go
var (
	file = flag.String("suites", "", `provide xml file name`)
)

func main() {
	flag.Parse()
	records := ExtractSuites(file)
	err := pkg.ProcessTestRecords(records)
	if err != nil {
		fmt.Println("Error when processing test records: ", err)
	}
}

func ExtractSuites(file *string) []pkg.JUnitTestSuites {
	final := []pkg.JUnitTestSuites{}
	data, err := os.Open(*file)
	if err != nil {
		fmt.Println("Error opening file ", err)
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fileName := scanner.Text()
		if fileName == "" {
			continue
		}
		suites := pkg.JUnitTestSuites{}
		content, err := os.Open(fileName)
		if err != nil {
			fmt.Println("Error opening XML file: ", err, fileName)
		}
		xmlByte, err := ioutil.ReadAll(content)
		if err != nil {
			fmt.Println("Error reading XML file: ", err, fileName)
		}
		err = xml.Unmarshal(xmlByte, &suites)
		if err != nil {
			fmt.Println("Error unmarshalling XML file for file: ", err, fileName)
		}
		final = append(final, suites)
	}
	return final
}
