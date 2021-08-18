package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Number of times to run the loop to monitor
const nTimesToMonitor = 2

// The delay between the monitor, in seconds
const delay = 5

func main() {
	showIntro()
	for {
		showMenu()
		input := readInput()

		switch input {
		case 1:
			fmt.Println("Start monitoring...")
			fmt.Println("")
			startMonitoring()
		case 2:
			fmt.Println("Showing logs...")
			fmt.Println("")
			showLogs()
		case 0:
			fmt.Println("Exiting the program...")
			os.Exit(0)
		default:
			fmt.Println("Unknown option, exiting the program!")
			os.Exit(-1)
		}
	}
}

func showIntro() {
	version := 1.0
	fmt.Println("-------------------------------------")
	fmt.Println("- Simple Website Heal Check Monitor -")
	fmt.Println("------------ Version", version, " -------------")
	fmt.Println("-------------------------------------")
	fmt.Println("")
}

func showMenu() {
	fmt.Println("--------------- Menu ----------------")
	fmt.Println("1 - Start to monitor")
	fmt.Println("2 - Show logs")
	fmt.Println("0 - Exit the program")
}

func readInput() int {
	var inputRead int
	_, err := fmt.Scan(&inputRead)
	if err != nil {
		fmt.Println("There was the following problem on reading the input:", err)
		os.Exit(-1)
	}
	return inputRead
}

func startMonitoring() {

	sites := readFileWithSitsList()

	for i := 0; i < nTimesToMonitor; i++ {
		for i, site := range sites {
			fmt.Println("Checking website n", i, ":", site)
			checkSite(site)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}

	fmt.Println("")
}

func checkSite(site string) {
	res, err := http.Get(site)

	if err != nil {
		fmt.Println("There was the following problem on checking the website health:", err)
		os.Exit(-1)
	}

	if res.StatusCode == 200 {
		fmt.Println("The website", site, "is healthy!")
		saveLogToFile(site, true)
	} else {
		fmt.Println("The website", site, "is unhealthy, with the response status code: ", res.StatusCode, "!")
		saveLogToFile(site, false)
	}
}

func readFileWithSitsList() []string {
	file, err := os.Open("sites.txt")
	var sites []string

	if err != nil {
		fmt.Println("There was the following error while trying to open the websites' file:", err)
		os.Exit(-1)
	}

	fileReader := bufio.NewReader(file)

	for {
		line, err := fileReader.ReadString('\n')

		line = strings.TrimSpace(line)

		sites = append(sites, line)

		if err == io.EOF {
			break
		} else if err != nil && err != io.EOF {
			fmt.Println("There was the following error while trying to read the websites' file:", err)
			os.Exit(-1)
		}
	}

	err = file.Close()
	if err != nil {
		fmt.Println("There was the following error while trying to open the websites' file:", err)
		os.Exit(-1)
	}

	return sites
}

func saveLogToFile(site string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("There was the following error while trying to open the log's file:", err)
		os.Exit(-1)
	}

	timestamp := time.Now().Format("02/01/2006 15:04:05")

	if status == true {
		_, err := file.WriteString("[INFO][ONLINE][" + timestamp + "] " + site + "\n")
		if err != nil {
			fmt.Println("There was the following error while trying to write to the log's file:", err)
			os.Exit(-1)
		}
	} else {
		_, err := file.WriteString("[ERR][OFFLINE][" + timestamp + "] " + site + "\n")
		if err != nil {
			fmt.Println("There was the following error while trying to write to the log's file:", err)
			os.Exit(-1)
		}
	}

	err = file.Close()
	if err != nil {
		fmt.Println("There was the following error while trying to open the websites' file:", err)
		os.Exit(-1)
	}
}

func showLogs() {
	file, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println("There was the following error while trying to open the log's file:", err)
		os.Exit(-1)
	}

	fmt.Println(string(file))
}
