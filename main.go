package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Question struct {
	question, answer string
}

var (
	filePath  string
	timeLimit int
	randomize bool
)

func readCSV(filePath string) ([]Question, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	questionList, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	questions := make([]Question, 0)
	for _, question := range questionList {
		questions = append(questions, Question{strings.TrimSpace(question[0]), strings.TrimSpace(question[1])})
	}
	if randomize {
		for i := len(questions) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			questions[i], questions[j] = questions[j], questions[i]
		}
	}
	return questions, nil
}

func quiz(questions []Question, timeout int) (int, error) {
	reader := bufio.NewReader(os.Stdin)
	timer := time.NewTicker(time.Second * time.Duration(timeout))
	score := 0

	errs := make(chan error)
	completed := make(chan int)

	go func() {
		for _, question := range questions {
			fmt.Printf("What is %s?: ", question.question)
			guess, err := reader.ReadString('\n')
			if err != nil {
				errs <- err
			}
			if strings.TrimSpace(guess) == question.answer {
				score++
			}
		}
		completed <- score
	}()

	select {
	case <-errs:
		return score, <-errs
	case <-completed:
		return score, nil
	case <-timer.C:
		fmt.Println("\nTime is up!")
	}
	return score, nil
}

func init() {
	flag.StringVar(&filePath, "csv", "problems.csv", "File path for a .csv gile in the format of 'question, answer'")
	flag.IntVar(&timeLimit, "t", 30, "Time limit for answering questions in seconds, defaults to 30")
	flag.BoolVar(&randomize, "random", false, "Randomize order of questions, defaults to 'false'")
	flag.Parse()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press enter to start!")
	reader.ReadString('\n')
	questions, err := readCSV(filePath)
	if err != nil {
		log.Fatal(err)
	}
	score, err := quiz(questions, timeLimit)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Score: %d/%d\n", score, len(questions))
	reader.ReadString('\n')
}
