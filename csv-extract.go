package main
import (
	"bufio"
	"encoding/csv"
	"os"
	"io"
	"log"
	"fmt"
	"strconv"
)

type Tag struct {
	Name string
	Value string
}

func main() {
	
	tags := parseFromFile("inventory-reports/scorecard.csv")
	fmt.Print(tags)

	// Count Keys
	keyCount := make(map[string]int)
	valueCount := make(map[string]int)
	keyValueCount := make(map[string]int)
	for _, tag:= range tags {
		keyCount[tag.Name] += 1
		valueCount[tag.Value] += 1
		keyValueCount[tag.Name + ":" + tag.Value] += 1
	}
	fmt.Print(keyCount)
	fmt.Print(valueCount)
	fmt.Print(keyValueCount)
	
	// Write to CSV
	outputToFile("inventory-reports/key.csv", keyCount)
	outputToFile("inventory-reports/values.csv", valueCount)
	outputToFile("inventory-reports/key-value-pair.csv", keyValueCount)
}

func parseFromFile(filePath string) []Tag {
	csvFile, _ := os.Open(filePath)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var tags []Tag
	count := 0

	for {
		line, error := reader.Read()
		count += 1

		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		if line[5] != "" {
			tags = append(tags, Tag{
				Name: line[4],
				Value: line[5],
			})
		}
	}
	return tags
}

func outputToFile (filePath string, countMap map[string]int) {
	outputFile, _ := os.Create(filePath)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()
	
	writer.Write([]string {"Label", "Count"})

	for key, value := range countMap {
		err := writer.Write([]string {key, strconv.Itoa(value)})
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}
}