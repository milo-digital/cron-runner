package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"time"
	"bytes"
	"os/exec"
)

type jobDefinition struct {
	Schedule string `json:"schedule"`
	Url string `json:"url"`
}


func main() {
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
	fmt.Println(jobsMap)
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
	if err != nil || ! bytes.Equal(newFile, existing) {
		exec.Command("crontab", "/tmp/cron")
		os.Remove("/tmp/cron")
	}
}

func readFiles() (jobs map[string][]jobDefinition) {
	jobs = make(map[string][]jobDefinition)
	deployEnv := os.Getenv("DEPLOY_ENV")
	files, _ := ioutil.ReadDir("/mnt")
	for _, f := range files {
		fmt.Println(f.Name())

		if f.IsDir() {
			cronFileName := fmt.Sprintf("/mnt/%v/code/live/cron/%v.cron", f.Name(), deployEnv)
			fmt.Println(cronFileName)
			if _, err := os.Stat(cronFileName); err == nil {

				file, e := ioutil.ReadFile(cronFileName)
				if e != nil {
					fmt.Printf("File error: %v\n", e)
					os.Exit(1)
				}
				var thisSiteJobs []jobDefinition
				err := json.Unmarshal(file, &thisSiteJobs)
				if err != nil {
					continue
				}
				jobs[f.Name()] = thisSiteJobs
			}
		}
	}
	return
}
