package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

type question struct {
	Q string
	A int
}

func main() {
	l, err := os.Create("logs.out")
	if err != nil {
		fmt.Println("Something wrong with logs. Contact me. Err: ", err)
	}
	defer l.Close()
	log.SetOutput(l)
	fn := flag.String("f", "filename", "Quiz questions file")
	timeLimit := flag.Int("t", 30, "Time limit to quiz")
	flag.Parse()
	ctx, close := context.WithTimeout(context.Background(), time.Duration(*timeLimit)*time.Second)
	defer close()
	quiz := readQuiz(*fn)
	var ta int64
	go func() {
		r := bufio.NewReader(os.Stdin)
		for _, qst := range quiz {
			fmt.Print(qst.Q, "=")

			// TODO Qst: How to cancel reading from io to terminate goroutine correctly
			text, err := r.ReadString('\n')
			if err != nil {
				fmt.Println("Program ran on internal err.")
				log.Fatalf("Can't read from stdin. Err: %s", err)
			}
			text = text[:len(text)-1]
			ans, err := strconv.Atoi(text)
			if err != nil {
				log.Printf("Can't parse user input to record (%s=%d). Answer is: %s. Err: %s\n", qst.Q, qst.A, text, err)
				continue
			}
			if qst.A == ans {
				atomic.AddInt64(&ta, 1)
			}
		}
		close()
	}()
	<-ctx.Done()
	fmt.Printf("You answered %d of %d questions\n", ta, len(quiz))
}

func readQuiz(filename string) []question {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		fmt.Printf("Can't open file with filename: %s\n", filename)
		log.Fatalf("Cant open file with name %s, with err: %s\n", filename, err)
	}
	questions := make([]question, 0)
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error reading records: %s\n", err)
		}
		if len(record) != 2 {
			log.Printf("Record (%v) is not valid length\n", record)
			continue
		}
		qst := record[0]
		ans, err := strconv.Atoi(record[1])
		if err != nil {
			log.Printf("Can't convert answer to integer. Err: %s\n", err)
			continue
		}
		questions = append(questions, question{qst, ans})
	}
	return questions
}
