package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {

	var (
		fileName     = flag.String("file-name", "", "File name that you want to parse")
		datePosition = flag.Int("date-position", -1, "indicate the position where the year is set in your csv")
		year         = flag.String("year-to-keep", "2019", "year you want to keep")
	)

	flag.Parse()
	if *fileName == "" {
		fmt.Println("You have to give the file-name as argument : ./excutable --file-name [fileName] --")
		log.Fatalln("No file name")
	}
	if *datePosition < 0 {
		fmt.Println("you should give a positive number to indicate the postion of date in your csv\n for example if your csv is 12,1643,2019-06-06: ./excutable --date-postion 3")
		log.Fatalln("position to find not set")
	}

	data, err := os.Open(*fileName)
	if err != nil {
		log.Fatalln("couldn't read file")
	}
	realFileName := strings.Split(*fileName, ".")

	writeData, err := os.Create(realFileName[0] + "_"+*year+".csv")
	defer writeData.Close()

	w := csv.NewWriter(writeData)
	r := csv.NewReader(data)

	// dont change first line
	words, err := r.Read()
	if err := w.Write(words); err != nil {
		log.Fatalln("cound not write csv:", err)
	}

	// check that postion is not creater than index in csv
	if *datePosition > len(words) {
		fmt.Printf("You should give a position smaller than number of fields in your csv %v\n", len(words))
		log.Fatalln("Position can be greater that size of csv fields")
	}
	// make position match array index
	*datePosition -= 1

	yearAlreadyFound := false
	lastLineFound := false
	// read all lines containing data
	for {
		words, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if match, _ := regexp.MatchString(".*"+*year+".*", words[*datePosition]); match {
			if err := w.Write(words); err != nil {
				log.Fatalln("cound not write csv:", err)
			}
			yearAlreadyFound = true
			lastLineFound = true
		} else {
			lastLineFound = false
		}

		if yearAlreadyFound && !lastLineFound {
			break
		}

	}

	if err := w.Error(); err != nil {
		log.Fatalln("error writting csv", err)
	}

	w.Flush()
	fmt.Println("convertion Done")

}
