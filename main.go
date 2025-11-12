package main

import (
	"embed"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

//go:embed data
var staticFiles embed.FS // embed data.csv into an embed.FS

type Quiz struct {
	Question   string
	Answer     string
	UserAnswer string
}

// read_csv - read a csv file and return the data
func read_csv(file_path string) ([]Quiz, error) {
	var file_data []Quiz

	f, fErr := staticFiles.Open(file_path)
	if fErr != nil {
		log.Fatalf("failed to read the provided file. %s", fErr.Error())
	}

	defer f.Close()
	csvr := csv.NewReader(f)
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return file_data, err
		}

		quiz := Quiz{
			Question: strings.TrimSpace(row[0]),
			Answer:   strings.TrimSpace(row[1]),
		}
		file_data = append(file_data, quiz)
	}
}

func readAllAndExecuetLineByLine(quizTime *int, shuffel *bool) {

	filepath := "data/problems.csv"
	fileData, fileReadErr := read_csv(filepath)
	if fileReadErr != nil {
		log.Fatalf("failed to read the file data. %s", fileReadErr.Error())
	}

	questionCount := 0
	correctAnswerCount := 0

	println("Please provide the answers to the corrusponding questions...\n")

	quizTimer := time.NewTimer(time.Duration(*quizTime) * time.Second)

	if *shuffel {
		// Shuffle the slice
		rand.Shuffle(len(fileData), func(i, j int) {
			fileData[i], fileData[j] = fileData[j], fileData[i]
		})
	}

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

	for _, quiz := range fileData {
		var userAnswer string
		fmt.Println(fmt.Sprintf("Question: %s", quiz.Question))
		fmt.Scan(&userAnswer)
		questionCount = questionCount + 1
		quiz.UserAnswer = userAnswer
		if strings.TrimSpace(userAnswer) == quiz.Answer {
			correctAnswerCount = correctAnswerCount + 1
		}
	}

	done <- true // stopping the timer
	fmt.Printf("Total Questions: %d\n", questionCount)
	fmt.Printf("Correct Answers: %d\n", correctAnswerCount)
}

func main() {
	quizTime := flag.Int("timer", 30, "duration of the quiz, default is 30s")
	shuffel := flag.Bool("shuffel", false, "Shuffel the quiz question, defalult is false")
	flag.Parse()
	fmt.Println("provided timer is", *quizTime)

	readAll := 0
	if readAll == 1 {
		readAllAndExecuetLineByLine(quizTime, shuffel)
	} else {
		ReadAndExecuteLineByline(quizTime)
	}
}
