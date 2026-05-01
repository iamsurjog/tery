package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type StatusRule struct {
	State   string `json:"state"`
	Command string `json:"command"`
}

type RangeRule struct {
	Upper   int    `json:"upper"`
	Lower   int    `json:"lower"`
	Command string `json:"command"`
}

type Config struct {
	Status []StatusRule `json:"status"`
	Range  []RangeRule  `json:"range"`
}

func get_config_path() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "tery", "config.json")
}

func runCommand(cmdStr string) {
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Run()
}


func init() {
	exampleConfig := `{
    "status": [
        {
            "state": "Charging",
            "command": "notify-send -i ac-adapter-symbolic Charging"
        },
        {
            "state": "Discharging",
            "command": "notify-send -i ac-adapter-symbolic Discharging"
        },
        {
            "state": "Not charging",
            "command": "notify-send -i ac-adapter-symbolic \"Not Charging\""
        },
        {
            "state": "Full",
            "command": "notify-send -i ac-adapter-symbolic Full"
        }
    ],
    "range": [
        {
            "upper": 20,
            "lower": 10,
            "command": "notify-send -i dialog-error \"Battery Extremely Low\" \"Plug in Charger Immediately\""
        },
        {
            "upper": 10,
            "lower": 0,
            "command": "notify-send -i battery-low \"Battery Low\" \"Plug in Charger\""
        }
    ]
}`

	filePath := get_config_path()
	configDir := filepath.Dir(filePath)

	_ = os.MkdirAll(configDir, os.ModePerm)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_ = os.WriteFile(filePath, []byte(exampleConfig), 0644)
		fmt.Println("Created default config at:", filePath)
	}
}

func main() {
	// Load and parse the config
	configData, err := os.ReadFile(get_config_path())
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	var cfg Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	batteryPath := "/sys/class/power_supply/BAT0/"
	
	// Initial state setup
	statusBytes, _ := os.ReadFile(batteryPath + "status")
	lastStatus := strings.TrimSpace(string(statusBytes))
	
	// Track which range rules have been triggered to avoid spamming
	notifiedRanges := make(map[int]bool)
	for {
		// 1. Check Status Changes
		sBytes, _ := os.ReadFile(batteryPath + "status")
		currentStatus := strings.TrimSpace(string(sBytes))

		if currentStatus != lastStatus {
			for _, rule := range cfg.Status {
				if rule.State == currentStatus {
					runCommand(rule.Command)
					break
				}
			}
			lastStatus = currentStatus
		}

		// 2. Check Capacity Ranges
		cBytes, _ := os.ReadFile(batteryPath + "capacity")
		percent, _ := strconv.Atoi(strings.TrimSpace(string(cBytes)))

		rangeTriggered := false
		for i, rule := range cfg.Range {
			if percent <= rule.Upper && percent >= rule.Lower {
				if !notifiedRanges[i] {
					runCommand(rule.Command)
					notifiedRanges[i] = true
				}
				rangeTriggered = true
				break // Only trigger the first matching range
			}
		}

		// Reset notifications if we move out of all defined ranges
		if !rangeTriggered {
			for i := range notifiedRanges {
				notifiedRanges[i] = false
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
