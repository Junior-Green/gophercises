package quiz

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func shuffleQuestions(slice [][]string) {
	rand.Shuffle(len(slice), func(i int, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func printResults(correct, total int) {
	fmt.Printf("You got %d correct and %d incorrect\n", correct, total-correct)
}

func startQuiz(done chan<- bool, questions [][]string, correct *int) {

	for i := range questions {
		var answer string

		fmt.Printf("%d. %s: ", i+1, questions[i][0])
		fmt.Scanln(&answer)

		if answer == questions[i][1] {
			*correct++
		}
	}

	done <- true
	close(done)
}

func Init() {
	var (
		filename  string
		shuffle   bool
		duration  int
		questions [][]string
		correct   int
		err       error
	)

	flag.StringVar(&filename, "f", "problems.csv", "Denote which .csv to open")
	flag.IntVar(&duration, "t", 300, "Denotes total test duration")
	flag.BoolVar(&shuffle, "shuffle", false, "Choose whether to shuffle questions using pseudo random algorithm")
	flag.Parse()

	file, err := os.Open(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured opening csv file: %s\n", filename)
		os.Exit(-1)
	}

	reader := csv.NewReader(file)
	questions, err = reader.ReadAll()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured parsing csv file: %s\n", filename)
		os.Exit(-1)
	}

	if shuffle {
		shuffleQuestions(questions)
	}

	fmt.Print("Press ENTER key when you are ready to start the quiz")
	fmt.Scanln()

	done := make(chan bool)
	timer := time.After(time.Duration(duration) * time.Second)
	go startQuiz(done, questions, &correct)

	for {
		select {
		case <-timer:
			fmt.Println("\nTime limit exceeded.")
			printResults(correct, len(questions))
			return
		case <-done:
			printResults(correct, len(questions))
			return
		}
	}
}
