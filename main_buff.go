package main

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"
)

// read_csv - read a csv file and return the data
func createFileReder(file_path string) (fs.File, error) {

	f, fErr := staticFiles.Open(file_path)
	if fErr != nil {
		log.Fatalf("failed to read the provided file. %s", fErr.Error())
		return nil, fErr
	}

	return f, nil
}

func ReadAndExecuteLineByline(quizTime *int) {

	filepath := "data/problems.csv"
	f, readerErrr := createFileReder(filepath)
	if readerErrr != nil {
		fmt.Print("failed to read file")
		os.Exit(1)
	}

	reader := csv.NewReader(f)

	questionCount := 0
	correctAnswerCount := 0

	println("Please provide the answers to the corrusponding questions...\n")

	quizTimer := time.NewTimer(time.Duration(*quizTime) * time.Second)

	var done = make(chan bool)

	fmt.Println("Starting a timer for the quiz")

	go func() {

		select {
		case <-quizTimer.C:
			fmt.Printf("Quiz time %ds is over, clossing the quiz\n", *quizTime)
			fmt.Printf("Total Questions: %d\n", questionCount)
			fmt.Printf("Correct Answer: %d\n", correctAnswerCount)
			quizTimer.Stop()
			os.Exit(0)
		case <-done:
			quizTimer.Stop() // stop the timer safely
		}

		// <-time.After(time.Duration(*quizTime) * time.Second)
		// fmt.Printf("Quiz time %ds is over, clossing the quiz\n", *quizTime)
		// fmt.Printf("Total Questions: %d\n", questionCount)
		// fmt.Printf("Correct Answer: %d\n", correctAnswerCount)
		// os.Exit(0)
	}()

	for {
		var userAnswer string
		question, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
		}

		fmt.Println(fmt.Sprintf("Question: %s", question[0]))
		fmt.Scan(&userAnswer)
		questionCount = questionCount + 1
		if strings.TrimSpace(userAnswer) == question[1] {
			correctAnswerCount = correctAnswerCount + 1
		}
	}

	done <- true // stopping the timer
	fmt.Printf("Total Questions: %d\n", questionCount)
	fmt.Printf("Correct Answers: %d\n", correctAnswerCount)

}
