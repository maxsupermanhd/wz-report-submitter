package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func panicErr(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	logFilePath := flag.String("logfile", "init.0.log", "Path to log file")
	logGid := flag.Int("gid", -1, "Game id")
	submitURL := flag.String("url", "https://localhost:5555/b/end/", "Link to submit endpoint")
	flag.Parse()
	if logGid == nil || *logGid < 1 {
		log.Print("Invalid GID")
		return
	}
	log.Print("Log file: ", *logFilePath)
	log.Print("Game id: ", *logGid)
	log.Print("Submit url: ", *submitURL+strconv.Itoa(*logGid))
	logFileB, err := os.ReadFile(*logFilePath)
	panicErr(err)
	logFile := string(logFileB)
	re := regexp.MustCompile(`__REPORTextended__(?P<json>.*)__ENDREPORTextended__`)
	for i, match := range re.FindAllString(logFile, -1) {
		log.Print("Found extended report at ", i)
		toread := match[len("__REPORTextended__") : len(match)-len("__ENDREPORTextended__")]
		log.Print("Report len: ", len(toread))
		obj := map[string]interface{}{}
		err = json.Unmarshal([]byte(toread), &obj)
		panicErr(err)
		req, err := http.NewRequest("POST", *submitURL+strconv.Itoa(*logGid), bytes.NewBuffer([]byte(toread)))
		panicErr(err)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		panicErr(err)
		defer resp.Body.Close()
		fmt.Println("Response Status:", resp.Status)
		fmt.Println("Response Headers:", resp.Header)
		body, err := ioutil.ReadAll(resp.Body)
		panicErr(err)
		fmt.Println("Response Body:", string(body))
		break
	}
}
