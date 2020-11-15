package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"github.com/bbalet/stopwords"
)

func tokenize(message string) []string {
	cleanMessage := stopwords.CleanString(message, "en", false)

	var re = regexp.MustCompile(`\$[\d]`)
	cleanMessage = re.ReplaceAllString(cleanMessage, "price")

	re = regexp.MustCompile(`\%[\d]`)
	cleanMessage = re.ReplaceAllString(cleanMessage, "percentage")

	re = regexp.MustCompile(`http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`)
	cleanMessage = re.ReplaceAllString(cleanMessage, "url")

	re = regexp.MustCompile(`www.(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`)
	cleanMessage = re.ReplaceAllString(cleanMessage, "url")

	re = regexp.MustCompile(`(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`)
	cleanMessage = re.ReplaceAllString(cleanMessage, "email")

	re = regexp.MustCompile(`\b\w{1,2}\b`)
	cleanMessage = re.ReplaceAllString(cleanMessage, " ")

	re = regexp.MustCompile(`[\W\d]`)
	cleanMessage = re.ReplaceAllString(cleanMessage, " ")

	re = regexp.MustCompile(`\s+`)
	cleanMessage = re.ReplaceAllString(cleanMessage, " ")

	cleanMessage = strings.TrimSpace(cleanMessage)
	parts := strings.Split(cleanMessage, " ")
	return parts
}


func copyChunk(chunk [][]string) [][]string {
	new := make([][] string , len(chunk))
	for index,data := range chunk {
		new[index] = make([]string, len(data))
		copy(new[index], data)
	}
	return new

}

func updateData(filename string, chunks int) (map[string]*Data, error) {
	records, err := csvReader(filename)
	if err != nil {
		return nil, err
	}

	var dataCount map[string]*Data
	if value, err := loadData("./count.csv"); err == nil {
		dataCount = value
	} else {
		dataCount = make(map[string]*Data)
	}


	offset := 0
	aux := len(records)/chunks

	var wg sync.WaitGroup
	var mutex = sync.Mutex{}
	var chunk [][]string
	
	for offset < len(records) {
	
		if offset+aux < len(records) {
			chunk = records[offset:offset+aux]
		} else {
			chunk = records[offset:]
		}

		wg.Add(1)
		go func(records [][]string, wg *sync.WaitGroup, mutex *sync.Mutex) {
			defer wg.Done()

			for i := 1; i < len(records); i++ {
				element := records[i]
				tokens := tokenize(element[1])

				for _, token := range tokens {		
					mutex.Lock()
					_, ok := dataCount[token]
					if !ok {
						data :=newData(token, 0, 0)
						dataCount[token] = &data
					}
					if element[0] == "spam" {
						dataCount[token].spamCount++
					} else {
						dataCount[token].hamCount++
					}
					mutex.Unlock()
				}
			}
		}(chunk, &wg, &mutex)
		offset += aux	
	}
	wg.Wait()
	return dataCount, err
}

func trainAux(data []*Data, class string, relativeLenght int) <-chan map[string]float64 {
	channel := make(chan map[string]float64)
	var mutex = sync.Mutex{}

	go func(mutex *sync.Mutex) {
		defer close(channel)
		weights := make(map[string]float64)

		mutex.Lock()
		lenght := len(data)
		mutex.Unlock()

		aux := float64(relativeLenght + len(data))
		for i := 0; i < lenght; i++ {
			mutex.Lock()
			value := data[i].copy()
			mutex.Unlock()

			count := 0.0
			if class == "spam" {
				count = value.spamCount
			} else {
				count = value.hamCount
			}
			weights[value.word] = (count + 1) / aux
		}
		channel <- weights
	}(&mutex)
	return channel
}

func train(data []*Data) {
	spamWords := 0.0
	hamWords := 0.0
	lenght := len(data)

	for _, value := range data {
		spamWords += value.spamCount
		hamWords += value.hamCount
	}

	spamChan, hamChan := trainAux(data, "spam", int(spamWords)), trainAux(data, "ham", int(hamWords))
	spam, ham := <-spamChan, <-hamChan

	result := make(map[string]*Data)
	for i := 0; i < lenght; i++ {
		word := data[i].word
		data := newData(word, spam[word], ham[word])
		result[word] = &data
	}

	resultList := mapToList(result, len(result))
	saveData(resultList, "./probabilities.csv")
}

func classify(data map[string]*Data, message string) int {
	tokens := tokenize(message)
	spamResult := 1.0
	hamResult := 1.0

	for _, word := range tokens {
		if val, ok := data[word]; ok {
			spamResult *= val.spamCount
			hamResult *= val.hamCount
		}
	}

	if spamResult > hamResult {
		return 1
	}
	return 0
}

func test(probabilities map[string]*Data, filename string) {
	confusionMatrix := [][]float64{{0, 0}, {0, 0}}
	records, _ := csvReader(filename)
	records = records[1:]

	for _, record := range records {
		result := classify(probabilities, record[1])
		if record[0] == "spam" {
			confusionMatrix[1][result]++
		} else {
			confusionMatrix[0][result]++
		}
	}

	accuaracy := 0.0

	for i := 0; i < len(confusionMatrix); i++ {
		for j := 0; j < len(confusionMatrix); j++ {
			accuaracy += confusionMatrix[i][j]
		}
	}

	accuaracy = (confusionMatrix[1][1] + confusionMatrix[0][0]) / accuaracy
	precision := confusionMatrix[1][1] / (confusionMatrix[1][1] + confusionMatrix[0][1])
	recall := confusionMatrix[1][1] / (confusionMatrix[1][1] + confusionMatrix[1][0])

	fmt.Printf("Accuaracy: %f\n", accuaracy)
	fmt.Printf("Precision: %f\n", precision)
	fmt.Printf("Recall: %f\n", recall)
	fmt.Println(confusionMatrix)
}
