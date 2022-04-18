package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const version = 0.1
const monitoringTimes = 4
const delaySeconds = 5
const webSitesFileName = "sites.txt"
const logFileName = "log.txt"

func intro(version float32) {
	fmt.Println("1- Start monitoring")
	fmt.Println("2- Show logs")
	fmt.Println("3- Exit")
	fmt.Print("Command: ")
}

func readCommand() int {
	var command int
	fmt.Scan(&command)
	fmt.Println("You chose command", command)
	fmt.Println("")
	return command
}

func logRegister(webSite string, status bool, statusCode int) {
	file, error := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if error != nil {
		fmt.Println("ERROR when registering logs")
		fmt.Println("Error", error)
	}

	var statusString string

	if status {
		statusString = "- ONLINE"
	} else {
		statusString = "- OFFLINE"
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " --- " + webSite + "- status code:" + strconv.Itoa(statusCode) + statusString + "\n")
	file.Close()
}

func startMonitoring(webSites []string) {
	fmt.Println("Starting monitors...")
	fmt.Println("")
	for i, webSite := range webSites {
		fmt.Println("Strating monitor on web site number", i, ":", webSite)
		response, error := http.Get(webSite)

		if error != nil {
			fmt.Println("Requisition ERROR")
			fmt.Println(error)
		}

		if response.StatusCode == 200 {
			fmt.Println("-> SUCCESS when monioring:", webSite, ". Status code:", response.StatusCode)
			fmt.Println("")
			logRegister(webSite, true, response.StatusCode)
		} else {
			fmt.Println("-> ERROR when monitoring:", webSite, ". Status code:", response.StatusCode)
			fmt.Println("")
			logRegister(webSite, false, response.StatusCode)

		}
	}

}

func showLogs() {
	fmt.Println("Showing logs...")
	fmt.Println("")
	//read file and returns a bytes array
	// ioutil closes the file auto.
	file, error := ioutil.ReadFile(logFileName)

	if error != nil {
		fmt.Println("Error reading", logFileName)
		fmt.Println("Error:", error)
	}
	//string(file) converts bytes array to string

	fmt.Println(string(file))
}

func readSitesFromTxt(fileName string) []string {
	sites := []string{}

	file, errorOpen := os.Open(fileName)

	if errorOpen != nil {
		fmt.Println("Error when opening", fileName, "file")
		fmt.Println("Error:", errorOpen)
		os.Exit(-1)
	}

	reader := bufio.NewReader(file) // start file reader

	for {
		line, errorRead := reader.ReadString('\n') // reads until gets to \n
		line = strings.TrimSpace(line)             //cuts all spaces and \n from the readed line

		fmt.Println(line)
		sites = append(sites, line)

		if errorRead == io.EOF {
			break
		}

	}

	file.Close()
	return sites
}
func printLine(length int) {
	for l := 0; l < length; l++ {
		fmt.Print("-")
	}
}

func main() {
	fmt.Println("App version:", version)
	fmt.Println("")
	for {
		intro(version)
		command := readCommand()

		switch command {
		case 1:
			webSites := readSitesFromTxt(webSitesFileName)
			for i := 0; i < monitoringTimes; i++ {
				startMonitoring(webSites)
				printLine(100)
				fmt.Println("")
				time.Sleep(delaySeconds * time.Second)
			}
		case 2:
			showLogs()
		case 3:
			fmt.Println("Exit")
			os.Exit(0)
		default:
			fmt.Println("Invalid command:", command)
		}
	}

}
