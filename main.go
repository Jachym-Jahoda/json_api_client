package main

import (
	"encoding/json"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/kardianos/service"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const version = "2023.2.5.16"
const serviceName = "json api client"
const serviceDescription = "json api client"

var (
	serviceIsRunning = false
	serviceSync      sync.Mutex
)

type SalaryPerUser struct {
	FirstName     string
	LastName      string
	Email         string
	Age           int
	MonthlySalary []MonthlySalary
}

type MonthlySalary struct {
	Basic int
	HRA   int
	TA    int
}

type program struct{}

func main() {
	fmt.Println(color.Ize(color.Green, "INF [SYSTEM] "+serviceName+" ["+version+"] starting..."))
	fmt.Println(color.Ize(color.Green, "INF [SYSTEM] © "+strconv.Itoa(time.Now().Year())+" Jachym Jahoda"))
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		fmt.Println(color.Ize(color.Red, "ERR [SYSTEM] Cannot start: "+err.Error()))
	}

	err = s.Run()
	if err != nil {
		fmt.Println(color.Ize(color.Red, "ERR [SYSTEM] Cannot start: "+err.Error()))
	}
}

func (p *program) Start(service.Service) error {
	fmt.Println(color.Ize(color.Green, "INF [SYSTEM] "+serviceName+" ["+version+"] started"))
	go p.run()
	serviceSync.Lock()
	serviceIsRunning = true
	serviceSync.Unlock()
	return nil
}

func (p *program) Stop(service.Service) error {
	serviceSync.Lock()
	serviceIsRunning = false
	serviceSync.Unlock()
	fmt.Println(color.Ize(color.Green, "INF [SYSTEM] "+serviceName+" ["+version+"] stopped"))
	return nil
}

func (p *program) run() {
	for {
		serviceSync.Lock()
		serviceRunning := serviceIsRunning
		serviceSync.Unlock()
		if !serviceRunning {
			break
		}
		start := time.Now()

		jsonfile := readJson()
		dataForPage, err := json.Marshal(jsonfile)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(dataForPage))
		newFile, err := os.Create("after.json")
		if err != nil {
			log.Println(err.Error())
		}
		defer newFile.Close()
		_, err2 := newFile.WriteString(string(dataForPage))
		if err2 != nil {
			log.Println(err2.Error())
		}
		fmt.Println(color.Ize(color.Green, "INF File written successfully"))
		if time.Since(start) < (60 * time.Second) {
			sleepTime := 60 * time.Second
			fmt.Println(color.Ize(color.Green, "INF [MAIN] Sleeping for "+sleepTime.String()))
			time.Sleep(sleepTime)
		}
	}
}

func readJson() []SalaryPerUser {
	file, err := os.Open("test.json")
	if err != nil {
		log.Println("Error opening json file:", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading json data:", err)
	}
	var jsonfile []SalaryPerUser
	err = json.Unmarshal(data, &jsonfile)
	if err != nil {
		log.Println("Error unmarshalling json data:", err)
	}
	return jsonfile
}
