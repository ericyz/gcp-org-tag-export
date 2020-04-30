/**
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
 
package main
import (
	"bufio"
	"encoding/csv"
	"os"
	"io"
	"log"
	"strconv"
	"strings"
)

type Tag struct {
	Name string
	Value string
	Resource string
}

type ResourceOutputItem struct {
	Resource string
	Identifier string
	Count string
}

func main() {
	
	tags := parseFromFile("inventory-reports/scorecard.csv")
	// Count Keys
	keyCount := make(map[string]int)
	valueCount := make(map[string]int)
	keyValueCount := make(map[string]int)
	resourceMap := make(map[string][]Tag)
	for _, tag:= range tags {
		keyCount[tag.Name] += 1
		valueCount[tag.Value] += 1
		keyValueCount[tag.Name + ":" + tag.Value] += 1
		resourceMap[tag.Resource] = append(resourceMap[tag.Resource], tag)
	}

	// Write to CSV
	outputKeyValueToFile("inventory-reports/key-counts.csv", keyCount)
	outputKeyValueToFile("inventory-reports/value-counts.csv", valueCount)
	outputKeyValueToFile("inventory-reports/keyvalue-counts.csv", keyValueCount)

	var resourceKeyItems []ResourceOutputItem
	var resourceValueItems []ResourceOutputItem
	var resourceKeyValueItems []ResourceOutputItem

	for resource, tags := range resourceMap {
		keyCount := make(map[string]int)
		valueCount := make(map[string]int)
		keyValueCount := make(map[string]int)	
		
		for _, tag:= range tags {
			keyCount[tag.Name] += 1
			valueCount[tag.Value] += 1
			keyValueCount[tag.Name + ":" + tag.Value] += 1
		}

		for tag, count := range keyCount {
			resourceKeyItems = append(resourceKeyItems, ResourceOutputItem {
				Resource : resource,
				Identifier : tag,
				Count: strconv.Itoa(count),
			})
		}

		for tag, count := range valueCount {
			resourceValueItems = append(resourceValueItems, ResourceOutputItem {
				Resource : resource,
				Identifier : tag,
				Count: strconv.Itoa(count),
			})
		}

		for tag, count := range keyValueCount {
			resourceKeyValueItems = append(resourceKeyValueItems, ResourceOutputItem {
				Resource : resource,
				Identifier : tag,
				Count: strconv.Itoa(count),
			})
		}
	}


	
	outputResourceKeyValueToFile("inventory-reports/key-counts-by-resource.csv", resourceKeyItems)
	outputResourceKeyValueToFile("inventory-reports/value-counts-by-resource.csv", resourceValueItems)
	outputResourceKeyValueToFile("inventory-reports/keyvalue-counts-by-resource.csv", resourceKeyValueItems)
}

func parseFromFile(filePath string) []Tag {
	csvFile, _ := os.Open(filePath)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var tags []Tag
	count := 0

	for {
		line, error := reader.Read()
		
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		if count > 0 && line[5] != "" {
			resource := strings.Split(line[3], "/")
			tags = append(tags, Tag{
				Resource: resource[len(resource)-1],
				Name: line[4],
				Value: line[5],
			})
		}
		count += 1
	}
	return tags
}

func outputKeyValueToFile (filePath string, countMap map[string]int) {
	outputFile, _ := os.Create(filePath)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()
	
	writer.Write([]string {"Identifier", "Count"})

	for key, value := range countMap {
		err := writer.Write([]string {key, strconv.Itoa(value)})
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}
}

func outputResourceKeyValueToFile(filePath string, resourceResultItems []ResourceOutputItem) {
	outputFile, _ := os.Create(filePath)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()
	
	writer.Write([]string {"Resource", "Identifier", "Count"})
	
	for _, item := range resourceResultItems {
		err := writer.Write([]string {item.Resource, item.Identifier, item.Count} )
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}
}
