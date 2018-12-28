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
	fileName := flag.String("csv", "problems.csv",
		"a csv file in the format of question,answer")
	timeLimit := flag.Int("limit", 30, "time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Printf("filed to open the csv file: %s\n", *fileName)
		os.Exit(1)
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		fmt.Printf("filed to parse the provided csv file")
		os.Exit(1)
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemLoop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	out := make([]problem, len(lines))
	for i, line := range lines {
		out[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return out
}

type problem struct {
	q string
	a string
}
