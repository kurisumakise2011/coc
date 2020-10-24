package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

const fileDesc = "The filepath to input."
const wordDesc = "The word to count occurences."
const sensetiveDesc = "Set the flag for the case sensetive search. Default is false."
const green = "\u001B[32m"
const reset = "\u001B[0m"

var file = flag.String("file", "", fileDesc)
var word = flag.String("name", "", wordDesc)
var sensetive = flag.Bool("case", false, sensetiveDesc)

func main() {
	flag.StringVar(file, "f", "", "The short alias for -file")
	flag.StringVar(word, "i", "", "The short alias for -name")
	flag.BoolVar(sensetive, "c", false, "The short alias for -case")
	flag.Parse()

	if *word == "" {
		return
	}

	rd, err := os.OpenFile(*file, os.O_RDONLY, 0444)
	if err != nil {
		log.Fatal(err)
	}
	defer rd.Close()

	scan := bufio.NewScanner(rd)
	scan.Split(bufio.ScanLines)

	type strint struct {
		s string
		i []int
		n int
	}

	lines := make(chan strint)
	counter := make(chan strint)
	done := make(chan int)

	go func(rd *bufio.Scanner) {
		n := 0
		for scan.Scan() {
			line := scan.Text()
			n++
			if line == "" {
				continue
			}
			lines <- strint{s: line, n: n}
		}
		if err := scan.Err(); err != nil {
			log.Println(err)
			done <- 2
			return
		}
		close(lines)
	}(scan)

	go func() {
		for {
			line, ok := <-lines
			if !ok {
				close(counter)
				return
			}

			s := line.s
			strlen := len(*word)
			indexes := []int{}
			shift := 0
			for {
				var index int
				if *sensetive {
					index = strings.Index(s, *word)
				} else {
					index = strings.Index(strings.ToLower(s), strings.ToLower(*word))
				}
				if index == -1 {
					break
				}
				indexes = append(indexes, index+shift)
				shift += index + strlen
				s = s[index+strlen:]
			}
			if len(indexes) > 0 {
				counter <- strint{line.s, indexes, line.n}
			}
		}
	}()

	count := 0
	go func() {
		for {
			match, ok := <-counter

			if !ok {
				done <- 1
				return
			}

			indexes := match.i
			strlen := len(*word)
			line := match.s

			for j, i := range indexes {
				i += j * (len(green) + len(reset))
				count++
				line = line[0:i] + green + line[i:i+strlen] + reset + line[i+strlen:]
			}
			log.Println("row " + strconv.Itoa(match.n) + ", " + line)
		}
	}()

	<-done
	log.Printf("matches %d\n", count)
}
