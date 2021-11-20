package main

import (
	"encoding/json"
	"fmt"
	"github.com/roessland/csgostate/pkg/csgostate"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("specify input file as first argument")
	}

	fileName := os.Args[1]
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	msgsStr := string(buf)
	msgs := strings.Split(msgsStr, "\n\n")

	firstTimestamp := 0
	prevTimestamp := 0
	offset := 1612210217 - 1612207556

	for _, msg := range msgs {
		//
		if len(msg) == 0 {
			continue
		}
		change := csgostate.State{}
		err := json.Unmarshal([]byte(msg), &change)
		if err != nil {
			fmt.Println("Json is", string(msg))
			panic(err)
		}
		if firstTimestamp == 0 {
			firstTimestamp = change.Provider.Timestamp
			prevTimestamp = change.Provider.Timestamp
			fmt.Println("setting first timetsmp to", firstTimestamp)
		}


		if change.Provider.Timestamp < firstTimestamp + offset  {
			prevTimestamp = change.Provider.Timestamp
			continue
		}

		fmt.Println("sleeping for ", change.Provider.Timestamp - prevTimestamp, "seconds)")
		// time.Sleep(time.Duration(change.Provider.Timestamp - prevTimestamp)*time.Second/3)

		resp, err := http.Post("http://127.0.0.1:3528/api/push", "text/html", strings.NewReader(msg))
		if err != nil {
			log.Print("fuck", resp, err)
		}

		prevTimestamp = change.Provider.Timestamp
	}
}