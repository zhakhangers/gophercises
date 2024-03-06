package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var filename string
	var durationInSeconds int
	var shuffle bool
	flag.StringVar(&filename, "f", "problems.csv", "file containing list of problems")
	flag.IntVar(&durationInSeconds, "t", 2, "Amount of time in seconds for answering all the questions")
	flag.BoolVar(&shuffle, "sh", false, "Shuffle the card")
	flag.Parse()
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("error openning the file:", err)
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("error reading from file: ", err)
		return
	}

	if shuffle {
		rand.Shuffle(len(records), func(i int, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}

	counter := 0
	questionCounter := 0

	timer := time.NewTimer(time.Duration(durationInSeconds) * time.Second)
	inputChan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputChan <- scanner.Text()
		}

	}()

	for _, record := range records {
		fmt.Print(record[0], "=")
		select {
		case <-timer.C:
			fmt.Printf("\nNumber of correct answers: %d \nNumber of total questions: %d ", counter, questionCounter)
			return
		case providedResult := <-inputChan:
			userAnswer, _ := strconv.Atoi(providedResult)
			correctAnswer, _ := strconv.Atoi(strings.TrimLeft(record[1], " "))
			if userAnswer == correctAnswer {
				counter++
			}
			questionCounter++
		}
	}

	fmt.Printf("\nNumber of correct answers: %d \nNumber of total questions: %d ", counter, questionCounter)

}
