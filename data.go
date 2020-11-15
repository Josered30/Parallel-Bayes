package main

import (
	"sort"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Data ...
type Data struct {
	word string
	spamCount float64
	hamCount float64	
}

func newData(word string, spamCount float64, hamCount float64) Data {
	return Data{
		word: word,
		spamCount: spamCount,
		hamCount: hamCount,
	}
}

func (data Data) copy() Data {
	return Data {
		word: data.word,
		spamCount: data.spamCount,
		hamCount: data.hamCount,
	}
}

func csvReader(filename string) ([][]string, error) {
	recordFile, err := os.Open(filename)
	defer recordFile.Close()

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	reader := csv.NewReader(recordFile)
	records, _ := reader.ReadAll()
	return records, err

}

func mapToList(data map[string]*Data, limit int) []*Data {
	values := make([]*Data, len(data))
	i := 0
	for _, value := range data {
		values[i] = value
		i++
	}

	if limit == len(values) {
		return values
	}

	sort.Slice(values, func(i, j int) bool {
		return (values[i].spamCount + values[i].hamCount) > (values[j].spamCount + values[j].hamCount)
	})
	return values[:limit]
}

func loadData(filename string) (map[string]*Data, error) {
	records, err := csvReader(filename)
	if err != nil {
		return nil, err
	}
	records = records[1:]
	dataMap := make(map[string]*Data)

	for _ , record := range records {
		spamCount, _ := strconv.ParseFloat(record[1], 64)
		hamCount, _ := strconv.ParseFloat(record[2], 64)

		data := newData(record[0], spamCount, hamCount)
		dataMap[record[0]] = &data
	}
	return dataMap , err
}

func loadDataList(filename string) ([]*Data, error) {
	records, err := csvReader(filename)
	if err != nil {
		return nil, err
	}
	records = records[1:]
	dataList := make([]*Data, len(records))

	for index, record := range records {
		spamCount, _ := strconv.ParseFloat(record[1], 64)
		hamCount, _ := strconv.ParseFloat(record[2], 64)
		data := newData(record[0], spamCount, hamCount)
		dataList[index] = &data
	}
	return dataList, err
}

func saveData(data []*Data, filename string) {
	recordFile, err := os.Create(filename)
	defer recordFile.Close()

	if err != nil {
		fmt.Println("An error encountered ::", err)
	}

	writer := csv.NewWriter(recordFile)
	defer writer.Flush()

	writer.Write([]string{"word", "spam", "ham"})

	for _, value := range data {
		spamCount := fmt.Sprintf("%f",value.spamCount)
		hamCount := fmt.Sprintf("%f",value.hamCount)
		writer.Write([]string{value.word, spamCount, hamCount})
	}
}



