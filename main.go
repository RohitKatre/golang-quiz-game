package main

import (
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
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

func main() {
	filepath := "data/problems.csv"
	fileData, fileReadErr := read_csv(filepath)
	if fileReadErr != nil {
		log.Fatalf("failed to read the file data. %s", fileReadErr.Error())
	}

	questionCount := 0
	correctAnswerCount := 0

	println("Please provide the answers to the corrusponding questions...\n")

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
	fmt.Printf("Total Questions: %d", questionCount)
	fmt.Printf("Correct Answer: %d", correctAnswerCount)

}
