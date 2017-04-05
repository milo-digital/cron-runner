package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"time"
	"bytes"
	"os/exec"
	"net/url"
	"strings"
	"strconv"
)

type jobDefinition struct {
	Schedule string `json:"schedule"`
	Url string `json:"url"`
}


func main() {
	exec.Command("crond")
	process()
	c := time.Tick(60 * time.Second)
	for range c {
		process()
	}
}

func process() {
	os.Remove("/tmp/cron")
	f, _ := os.OpenFile("/tmp/cron", os.O_WRONLY|os.O_CREATE, 0666)

	jobsMap := readFiles()
	for siteName,jobs := range jobsMap {
		f.WriteString(fmt.Sprintf("### %s ###\n", siteName))
		for _, job := range jobs {
			f.WriteString(fmt.Sprintf("%s curl %s\n", job.Schedule, job.Url))
		}
		f.WriteString("\n\n")

	}
	f.Close()

	newFile, _ := ioutil.ReadFile("/tmp/cron")
	existing, err := ioutil.ReadFile("/etc/crontabs/root")
	if err != nil || !bytes.Equal(newFile, existing) {
		ret := exec.Command("crontab", "/tmp/cron")
		err = ret.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
	os.Remove("/tmp/cron")
}

func readFiles() (jobs map[string][]jobDefinition) {
	fmt.Printf("Polling at %s\n", time.Now())
	jobs = make(map[string][]jobDefinition)
	deployEnv := os.Getenv("DEPLOY_ENV")
	files, _ := ioutil.ReadDir("/mnt")
	OUTER:
	for _, f := range files {
		if f.IsDir() {
			cronFileName := fmt.Sprintf("/mnt/%v/code/live/cron/%v.cron", f.Name(), deployEnv)
			if _, err := os.Stat(cronFileName); err == nil {
				fmt.Printf("Processing file: %v\n", cronFileName)
				file, e := ioutil.ReadFile(cronFileName)
				if e != nil {
					fmt.Printf("File read error: %v\n", e)
					continue
				}
				var thisSiteJobs []jobDefinition
				err := json.Unmarshal(file, &thisSiteJobs)
				if err != nil {
					fmt.Printf("File format error: %v %v\n", cronFileName, err)
					continue
				}
				for i, line := range thisSiteJobs{
					_, err := url.Parse(line.Url)
					if err != nil {
						fmt.Printf("Line %v of %v has invalid URL\n", i, cronFileName)
						continue OUTER
					}
					if valid,err := scheduleValid(line.Schedule); !valid{
						fmt.Printf("Line %v of %v: %v\n", i, cronFileName, err)
						continue OUTER
					}
				}
				jobs[f.Name()] = thisSiteJobs
			}
		}
	}
	return
}

func scheduleValid (schedule string) (valid bool, errMsg string){
	scheduleArr := strings.Split(schedule, " ")
	if len(scheduleArr) != 5 {
		valid = false
		fmt.Println(scheduleArr)
		errMsg = "Schedule does not contain 5 parts"
		return
	}
	for i, value := range scheduleArr {
		if value == "*"{
			continue
		}
		ival, err := strconv.Atoi(value)
		if err != nil{
			valid = false
			errMsg = fmt.Sprintf("Element #%v of schedule is invalid", i)
			return
		}
		var thisValid bool
		switch (i){
		case 0:
			thisValid = ival >=0 && ival <=59
		case 1:
			thisValid = ival >=0 && ival <=23
		case 2:
			thisValid = ival >=1 && ival <=31
		case 3:
			thisValid = ival >=1 && ival <=12
		case 4:
			thisValid = ival >=0 && ival <=6
		}

		if !thisValid {
			valid = false
			errMsg = fmt.Sprintf("Element #%v of schedule is invalid", i)
			return
		}
	}
	valid = true
	return
}