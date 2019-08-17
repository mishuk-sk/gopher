package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	fn := flag.String("f", "filename", "Quiz questions file")
	flag.Parse()
	f, err := os.Open(*fn)
	defer f.Close()
	if err != nil {
		fmt.Printf("Can't open file with filename: %s\n", *fn)
		log.Fatalf("Cant open file with name %s, with err: %s\n", *fn, err)
	}
	l, err := os.Create("logs.out")
	if err != nil {
		fmt.Println("Something wrong with logs. Contact me. Err: ", err)
	}
	defer l.Close()
	log.SetOutput(l)
	r := csv.NewReader(f)
	tq := 0
	ta := 0
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
		fmt.Print(qst, "=")
		tq++
		rStdin := bufio.NewReader(os.Stdin)
		text, err := rStdin.ReadString('\n')
		if err != nil {
			fmt.Println("Program ran on internal err.")
			log.Fatalf("Can't read from stdin. Err: %s", err)
		}
		usrAns, err := strconv.Atoi(text[:len(text)-1])
		if err != nil {
			log.Printf("Can't parse user input to record (%v). Answer is: %s. Err: %s\n", record, text, err)
			continue
		}
		if usrAns == ans {
			ta++
		}
	}
	fmt.Printf("You answered %d of %d questions\n", ta, tq)
}
