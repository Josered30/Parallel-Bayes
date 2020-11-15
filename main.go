package main

import (
	"os"
	"fmt"
)

func main() {
	os.Remove("./count.csv")
	os.Remove("./probabilities.csv")
	
	dataMap, _ := updateData("./spamham.csv",4)
	dataList := mapToList(dataMap, 2000)
	saveData(dataList, "./count.csv")
	
	//dataList, _ := loadDataList("./count.csv")
	train(dataList)

	probabilities, _ := loadData("./probabilities.csv")
	test(probabilities,"./spamham.csv")
	fmt.Println("Done")
}