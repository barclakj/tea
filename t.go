// t -a <action> -d <dte> -t <topic> -p 1
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	Due       string   `json:"due"`
	CreatedTs int      `json:"createdTs"`
	DueTs     int      `json:"dueTs"`
	Priority  int      `json:"priority"`
	Topics    []string `json:"topics"`
}

type TaskQueryResponse []Task

const T_URL = "http://localhost:1643/t"

func main() {
	httpClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	if len(os.Args) < 3 {
		printInvalidSyntaxResponse()
		log.Fatal("Insuffient arguments")
	}

	switch action := os.Args[1]; action {
	case "-a":
		postTask(httpClient, argsToTask(os.Args))
	case "-d":
		item, _ := strconv.Atoi(os.Args[2])
		deleteTask(httpClient, item)
	case "-r":
		item, _ := strconv.Atoi(os.Args[2])
		task := getTask(httpClient, item)
		printTask(&task)
	default:
		printInvalidSyntaxResponse()
	}
}

func argsToTask(args []string) Task {
	var t Task

	t.Name = args[2]
	if len(args) > 3 {
		for i := 3; i < len(args); i += 2 {
			// fmt.Printf("%d = %s\n", i, args[i])
			if len(args) > i {
				switch args[i] {
				case "-p":
					t.Priority, _ = strconv.Atoi(args[i+1])
				case "-d":
					t.Due = args[i+1]
				case "-t":
					t.Topics = strings.Split(args[i+1], ",")
				}
			}
		}
	}

	return t
}

func printInvalidSyntaxResponse() {
	fmt.Printf("Invalid command syntax\n")
}

func printTask(task *Task) {
	fmt.Printf("P%d (%d) %s [%d] %s\n", task.Priority, task.Id, task.Name, task.DueTs, task.Topics)
}

func printResponse(queryResponse TaskQueryResponse) {
	fmt.Printf("Results,%d\n", len(queryResponse))
	for _, result := range queryResponse {
		fmt.Printf("%s,%d,%d,%d, %s\n", result.Id, result.CreatedTs, result.DueTs, result.Priority, result.Due)
	}
}

func postTask(httpClient http.Client, task Task) Task {
	path := T_URL + "/"

	data, err := json.Marshal(task)
	fmt.Printf("%s\n", data)

	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	body := execAndReturnBody(httpClient, req)

	postResult := Task{}
	jsonErr := json.Unmarshal(body, &postResult)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return postResult
}

func getTask(httpClient http.Client, id int) Task {
	path := T_URL + "/" + strconv.Itoa(id)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		log.Fatal(err)
	}
	body := execAndReturnBody(httpClient, req)

	getResult := Task{}
	jsonErr := json.Unmarshal(body, &getResult)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return getResult
}

func deleteTask(httpClient http.Client, id int) {
	path := T_URL + "/" + strconv.Itoa(id)

	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		log.Fatal(err)
	}
	execAndReturnBody(httpClient, req)
}

func execAndReturnBody(httpClient http.Client, req *http.Request) []byte {
	populateHeaders(req)

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	// fmt.Printf("%d", res.StatusCode)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body
}

func populateHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "topcat")
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
}
