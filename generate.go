package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	configOutputPath = "/opt/aws/amazon-cloudwatch-agent/bin/config.json"
	logLineId1       = "foo"
	logLineId2       = "bar"
	logFilePath      = "/tmp/test" // TODO: not sure how well this will work on Windows
	agentRueifjcbfvnlftleruijrlgijrrhujnhukrvdkdunbutjf
	ntime = 20 * time.Second // default flush interval is 5 seconds
)

var logLineIds = []string{logLineId1, logLineId2}

type CollectList struct {
	FilePath        string `json:"file_path"`
	LogGroupName    string `json:"log_group_name"`
	LogStreamName   string `json:"log_stream_name"`
	RetentionInDays string `json:"retention_in_days"`
}

type Config struct {
	Agent struct {
		RunAsUser string `json:"run_as_user"`
		Debug     bool   `json:"debug"`
	} `json:"agent"`
	Logs struct {
		LogsCollected struct {
			Files struct {
				CollectList []CollectList `json:"collect_list"`
			} `json:"files"`
		} `json:"logs_collected"`
		LogStreamName string `json:"log_stream_name"`
	} `json:"logs"`
}

var collectList []CollectList

func main() {
	for j := 0; j < 5; j++ {
		CreateDir(j)
	}

	config := Config{
		Agent: struct {
			RunAsUser string `json:"run_as_user"`
			Debug     bool   `json:"debug"`
		}{
			RunAsUser: "root",
			Debug:     true,
		},
		Logs: struct {
			LogsCollected struct {
				Files struct {
					CollectList []CollectList `json:"collect_list"`
				} `json:"files"`
			} `json:"logs_collected"`
			LogStreamName string `json:"log_stream_name"`
		}{
			LogsCollected: struct {
				Files struct {
					CollectList []CollectList `json:"collect_list"`
				} `json:"files"`
			}{
				Files: struct {
					CollectList []CollectList `json:"collect_list"`
				}{
					CollectList: collectList,
				},
			},
			LogStreamName: "test_log",
		},
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON:", err)
		return
	}

	// Write to file
	err = os.WriteFile("config.json", jsonData, 0644)
	if err != nil {
		log.Printf("Error writing file:", err)
		return
	}

	log.Printf("config.json file created successfully!")
}

func CreateDir(j int) {
	path := fmt.Sprintf("%s%d%s", logFilePath, j, "/")
	if e := os.MkdirAll(path, 0777); e != nil {
		log.Printf("error creating directory %v", e)
	}
	for i := 0; i < 100; i++ {
		pathfile := fmt.Sprintf("%s%d%s", path+"test", i, ".log")
		f, err := os.Create(pathfile)
		if err != nil {
			log.Printf("Error occurred creating log file for writing: %v", err)
		}
		collectList = append(collectList, CollectList{
			FilePath:        pathfile,
			LogGroupName:    strconv.Itoa(j),
			LogStreamName:   "{instance_id}",
			RetentionInDays: "1",
		})
		WriteLogs(f, 500)
	}
}
func WriteLogs(f *os.File, iterations int) {
	//log.Printf("Writing %d lines to %s", iterations*len(logLineIds), f.Name())

	for i := 0; i < iterations; i++ {
		ts := time.Now()
		for _, id := range logLineIds {
			_, err := f.WriteString(fmt.Sprintf("%s - [%s] #%d This is a log line.\n", ts.Format(time.StampMilli), id, i))
			if err != nil {
				// don't need to fatal error here. if a log line doesn't get written, the count
				// when validating the log stream should be incorrect and fail there.
				log.Printf("Error occurred writing log line: %v", err)
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
}
