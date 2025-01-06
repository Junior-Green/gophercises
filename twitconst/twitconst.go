package twitconst

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"time"

	"github.com/Junior-Green/gophercises/set"
	"github.com/Junior-Green/gophercises/twitter"
)

type Contest struct {
	Post        string
	BearerToken string
}

//1865710776717033901

// const apiKeySecret = "sB19UmYChlT5nuAj4dH64KIZ24jMkTTYUU9RPbL2WrEOl8qGgM"
// const apiKey = "3xGEcMOdejdr1PhYHlmrlnNjJ"
const bearerToken = "AAAAAAAAAAAAAAAAAAAAANGqxQEAAAAAwJCn%2FZMp38TTTxQiYtAD%2BmY8AmI%3D5bU02ipImuHvU1E5l3R4zGxCDy1QqeKcsF8pV7xSLSuKojwEgR"

func (c *Contest) Start() {
	var (
		pickAfter     time.Duration
		fetchInterval int
		numWinners    int
	)

	flag.DurationVar(&pickAfter, "pick", time.Minute*30, "sets when the winner will be pick from the time of execution")
	flag.IntVar(&fetchInterval, "i", 60, "sets interval of fetching new retweeters in seconds.")
	flag.IntVar(&numWinners, "w", 1, "sets amount of winners to choose")
	flag.Parse()

	if flag.Arg(0) == "" {
		fmt.Fprintln(os.Stderr, "missing arguments")
		os.Exit(1)
	}
	postId := flag.Arg(0)

	for {
		select {
		case <-time.Tick(time.Second * time.Duration(fetchInterval)):
			if err := updateEntries(postId); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		case <-time.After(pickAfter):
			winners, err := pickWinners(postId, numWinners)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			fmt.Println("Congrats to winners!")
			for _, winner := range winners {
				fmt.Printf("- %s\n", winner)
			}
			return
		}
	}
}

func pickWinners(postId string, numWinners int) ([]string, error) {
	ids, err := readCsv(postId)
	if err != nil {
		return nil, err
	}
	winners := make([]string, 0, numWinners)

	for ; numWinners > 0; numWinners-- {
		idx := rand.Intn(len(ids))
		winners = append(winners, ids[idx])
		ids[idx], ids[len(ids)-1] = ids[len(ids)-1], ids[idx]
		ids = ids[:len(ids)-1]
	}

	return winners, nil
}

func updateEntries(postId string) error {
	var err error

	client := twitter.TwitterClient{BearerToken: bearerToken}
	ids, err := client.GetRetweetedFromPostId(postId)
	if err != nil {
		return err
	}

	newIds, err := readCsv(postId)
	if err != nil {
		return err
	}

	uniqueIds := set.NewSet[string]()
	uniqueIds.AddAll(ids...)

	newIds = slices.DeleteFunc(newIds, func(id string) bool {
		return uniqueIds.Has(id)
	})

	return writeCsv(postId, newIds)
}

func writeCsv(postId string, ids []string) error {
	file, err := os.OpenFile(postId+".csv", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	newRecords := make([][]string, 0, len(ids))
	for _, id := range ids {
		newRecords = append(newRecords, []string{id})
	}

	return csv.NewWriter(file).WriteAll(newRecords)
}

func readCsv(postId string) ([]string, error) {
	file, err := os.OpenFile(postId+".csv", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(records))
	for _, record := range records {
		if len(record) != 1 {
			return nil, fmt.Errorf("corrupted csv file")
		}
		ids = append(ids, record[0])
	}

	return ids, nil
}
