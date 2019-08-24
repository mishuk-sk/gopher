package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type problem struct {
	Q string
	A string
}

func main() {
	//#############################
	//Initializing all io thing
	l, err := os.Create("logs.out")
	if err != nil {
		fmt.Println("Something wrong with logs. Contact me. Err: ", err)
	}
	defer l.Close()
	log.SetOutput(l)
	fn := flag.String("f", "quiz.csv", "Quiz questions csv file in format 'question,answer'")
	timeLimit := flag.Int("t", 30, "Time limit for the quiz")
	flag.Parse()
	f, err := os.Open(*fn)
	defer f.Close()
	if err != nil {
		fmt.Printf("Can't open file with file name: %s\n", *fn)
		log.Fatalf("Cant open file with name %s, with err: %s\n", *fn, err)
	}
	//############################
	quiz := readQuiz(f)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	answers := handleQuiz(quiz, timer)
	fmt.Printf("\nYou answered %d of %d questions\n", answers, len(quiz))
}

func handleQuiz(quiz []problem, timer *time.Timer) int {
	var answers int
	r := bufio.NewReader(os.Stdin)
	ansChan := make(chan string, 1)
	for _, qst := range quiz {
		fmt.Print(qst.Q, "=")
		go readInput(r, ansChan)
		select {
		case <-timer.C:
			return answers
		case ans := <-ansChan:
			if ans == qst.A {
				answers++
			}
		}
	}
	return answers
}

func readInput(r *bufio.Reader, answers chan<- string) {
	ans, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("Program ran into an internal err.")
		log.Fatalf("Can't read from stdin. Err: %s", err)
	}
	ans = strings.TrimSpace(ans[:len(ans)-1])
	ans = strings.ToLower(ans)
	answers <- ans
}

func readQuiz(f io.Reader) []problem {
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error reading records. Err - %s\n", err)
	}
	problems := make([]problem, 0, len(records))
	for _, record := range records {
		if len(record) != 2 {
			log.Printf("Record (%v) is not valid length\n", record)
			continue
		}
		answer := strings.TrimSpace(record[1])
		answer = strings.ToLower(answer)
		problems = append(problems, problem{
			Q: record[0],
			A: answer,
		})
	}
	return problems
}
