package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

// printError displays an error message to stdout
func printError(e error) {
	fmt.Println("Error occurred:", e)
}

// printResults displays quiz results to stdout
func printResults(total, attempted, correct int) {

	var completedPercentage, attemptedScore, totalScore float64

	completedPercentage = float64(attempted) / float64(total) * 100
	attemptedScore = float64(correct) / float64(attempted) * 100
	totalScore = float64(correct) / float64(total) * 100

	fmt.Printf("You attempted %d questions out of %d total questions (%.2f%%). You answered %d correctly.\n", attempted, total, completedPercentage, correct)
	fmt.Printf("Your attempted score is %.2f%%. Your total score is %.2f%%.\n", attemptedScore, totalScore)
}

func main() {

	// command line flag definitions and parsing

	var questionFilename = flag.String("q", "problems.csv", "Filename with path for question list")
	var timeLimit = flag.Int("limit", 30, "Time limit for quiz in seconds")
	flag.Parse()

	fmt.Println("The question file provided is", *questionFilename)
	fmt.Printf("You will have %d seconds to complete the quiz.\n", *timeLimit)

	qFile, err := os.Open(*questionFilename)
	defer qFile.Close()

	if err != nil {
		printError(err)
	}

	questionList, err := csv.NewReader(qFile).ReadAll()
	if err != nil {
		printError(err)
	}

	totalQuestions := len(questionList)
	attemptedQuestions, correctQuestions := 0, 0
	var answer string

	userInput := bufio.NewScanner(os.Stdin)
	inputCh := make(chan string)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	// Using an anonymous function to break out of the quiz loop when the select statement
	// is broken by the timer. Gophercises solution also shows use of a label and break <label>
	// Maybe refactoring this as a separate function would be cleaner?
	func() {
		for _, problem := range questionList {
			fmt.Printf("%s=?\tYour answer: ", problem[0])
			attemptedQuestions += 1
			go func() {
				userInput.Scan()
				inputCh <- userInput.Text()
			}()
			select {
			case answer = <-inputCh:
				if answer == problem[1] {
					correctQuestions += 1
				}
			case <-timer.C:
				return
			}
		}
	}()
	printResults(totalQuestions, attemptedQuestions, correctQuestions)
}
