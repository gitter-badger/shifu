package main


import (
	"github.com/docopt/docopt-go"
	"fmt"
	"net/http"
	//"log"
	//"io"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

var version = "0.0.1"

const usage = `
  Usage:
    shifu run [--testSeriesfile=<file>]
    shifu -h | --help
    shifu --version
  Options:
    -t, --testSeriesfile file   test file to execute
    -h, --help          		output help information
    -v, --version       		output version
  Examples:
    output tasks
    $ shifu
    run a test series
    $ shifu run -t test1.yaml
`


type TestDescriber struct {
    Test_name string
    Comment string
    Command_sequence []Command
}


type Command struct {
    Order_id int
    Comment string
    Type string
    Method string
    Url string
    Data string
    Expect Expect
    Headers map[string]string
    Repeat_times int
    Waiting_time int
}

type Expect struct {
	Value string
	Respones_code int
	Type string
}

func assertEquealString(actual string, expected string) {
	if actual != expected {
		fmt.Printf("❌ didnt match expected result, expected: %v ,found: %v\n" , expected, actual)
		os.Exit(1)
	} else {
		fmt.Printf("✅ matched result: %v\n", expected)
	}
}


func assertEquealInt(actual int, expected int) {
	if actual != expected {
		fmt.Printf("❌ didnt match expected result, expected: %v ,found: %v\n" , expected, actual)
		os.Exit(1)
	} else {
		fmt.Printf("✅ matched result: %v\n", expected)
	}
}


func processTestSeriesFile(file string) {
	fmt.Println("the file with instruction: " , file)

	yamlFile, _ := ioutil.ReadFile(file)
		
	testDescriber := TestDescriber{}

    yaml.Unmarshal([]byte(yamlFile), &testDescriber)
    
    fmt.Println("🐼 Running Test Describer " + testDescriber.Test_name)
    fmt.Println(testDescriber.Comment)


    tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	for _, command := range testDescriber.Command_sequence {
		fmt.Println("🐯 " + command.Comment)
		req, _ := http.NewRequest(command.Method, command.Url, nil)
		for key, val := range command.Headers {
			req.Header.Add(key,val)
		}

		resp, _ := client.Do(req)
	    defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)
		assertEquealInt(resp.StatusCode, command.Expect.Respones_code)
		assertEquealString(string(byteArray[:]), command.Expect.Value)
	}
}

func main() {

	args, _ := docopt.Parse(usage, nil, false, version, false)
	files, _ := args["--testSeriesfile"].([]string)
	
	for _,file := range files {
		processTestSeriesFile(file)
	}
}

