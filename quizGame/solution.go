package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "CSV file that reads the data")
	timeLimit := flag.Int("timeLimit", 30, "Timi limit in seconds for the quiz")

	flag.Parse()

	file, error := os.Open(*csvFileName)

	if error != nil {

		exit(fmt.Sprintf("Failed to open CSV File Name: %s\n", *csvFileName))
	}

	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the CSV")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	score := 0
problemLoop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d, %s = ", i+1, problem.q)
		answerChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer

		}()
		select {
		case <-timer.C:
			fmt.Println("Timeout !!!")
			break problemLoop
		case answer := <-answerChannel:
			if answer == problem.a {
				score++
			}
		}
	}

	fmt.Printf("\nYou scored %d points out of %d.\n", score, len(problems))
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))

	for i, line := range lines {
		problems[i] = problem{q: line[0], a: strings.TrimSpace(line[1])}
	}

	return problems
}

type problem struct {
	q string
	a string
}

func exit(message string) {
	fmt.Printf(message)
	os.Exit(1)
}
