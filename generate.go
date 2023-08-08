package main

import (
	"encoding/json"
	"fmt"
	"os"
)

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

func main() {
	var collectList []CollectList

	// Create 100 instances using a for loop
	for i := 0; i < 100; i++ {
		collectList = append(collectList, CollectList{
			FilePath:        fmt.Sprintf("/tmp/test0/test%d.log", i),
			LogGroupName:    "0",
			LogStreamName:   "{instance_id}",
			RetentionInDays: "1",
		})
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
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Write to file
	err = os.WriteFile("config.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("config.json file created successfully!")
}
