package logics

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

const (
	// Rate to limit to, e.g., "4000kbit" for 4Mbit/s
	limitRate = "4000kbit"
	// Time period when the speed limit is REMOVED (peak hours)
	peakStartTime = "17:15"
	peakEndTime   = "23:30"
)

// Logger can be replaced with your project's logger, e.g., zap
var shaperLogger = log.New(log.Writer(), "[tcShaper] ", log.LstdFlags)

// GetDefaultInterface retrieves the default network interface name.
func GetDefaultInterface() (string, error) {
	cmd := exec.Command("ip", "route")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'ip route': %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "default via") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				return fields[4], nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning 'ip route' output: %w", err)
	}

	return "", fmt.Errorf("default interface not found")
}

// runTCCommand executes a tc command.
func runTCCommand(args ...string) error {
	cmd := exec.Command("sudo", append([]string{"tc"}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	shaperLogger.Printf("Executing TC command: sudo tc %s", strings.Join(args, " "))
	err := cmd.Run()
	if err != nil {
		// tc qdisc del often returns an error if the qdisc doesn't exist, which is fine.
		if !(strings.Contains(strings.Join(args, " "), "qdisc del") && strings.Contains(stderr.String(), "No such file or directory")) {
			return fmt.Errorf("tc command 'sudo tc %s' failed: %w, stderr: %s", strings.Join(args, " "), err, stderr.String())
		}
		shaperLogger.Printf("TC warning (ignored): %s", stderr.String())
	}
	return nil
}

// ClearTCRules removes all tc qdisc rules from the specified interface.
func ClearTCRules(ifaceName string) error {
	shaperLogger.Printf("Clearing TC rules on interface %s", ifaceName)
	return runTCCommand("qdisc", "del", "dev", ifaceName, "root")
}

// ApplyTCRules applies HTB-based rate limiting to the specified interface.
func ApplyTCRules(ifaceName string, rate string) error {
	shaperLogger.Printf("Applying TC rate limit %s on interface %s", rate, ifaceName)
	// First, try to delete any existing qdisc to avoid errors if it already exists with a different handle.
	_ = ClearTCRules(ifaceName) // Ignore error, as it might not exist

	err := runTCCommand("qdisc", "add", "dev", ifaceName, "root", "handle", "1:", "htb", "default", "10")
	if err != nil {
		return fmt.Errorf("failed to add root htb qdisc: %w", err)
	}

	err = runTCCommand("class", "add", "dev", ifaceName, "parent", "1:", "classid", "1:10", "htb", "rate", rate)
	if err != nil {
		// If class add fails, try to clean up the qdisc we just added.
		_ = ClearTCRules(ifaceName)
		return fmt.Errorf("failed to add htb class: %w", err)
	}
	shaperLogger.Printf("Successfully applied TC rate limit %s on interface %s", rate, ifaceName)
	return nil
}

// ManageBandwidth applies or clears bandwidth limits based on the current time.
func ManageBandwidth() {
	iface, err := GetDefaultInterface()
	if err != nil {
		shaperLogger.Printf("Error getting default interface: %v", err)
		return
	}
	shaperLogger.Printf("Default interface: %s", iface)

	now := time.Now()
	currentTimeStr := now.Format("15:04")

	parsedStartTime, err := time.Parse("15:04", peakStartTime)
	if err != nil {
		shaperLogger.Printf("Error parsing peak start time '%s': %v", peakStartTime, err)
		return
	}
	parsedEndTime, err := time.Parse("15:04", peakEndTime)
	if err != nil {
		shaperLogger.Printf("Error parsing peak end time '%s': %v", peakEndTime, err)
		return
	}

	// Normalize times to the same date (today) for comparison
	// This handles time ranges that do not cross midnight.
	// For ranges crossing midnight, the logic needs to be adjusted.
	currentDayTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	peakStartDateTime := time.Date(now.Year(), now.Month(), now.Day(), parsedStartTime.Hour(), parsedStartTime.Minute(), 0, 0, now.Location())
	peakEndDateTime := time.Date(now.Year(), now.Month(), now.Day(), parsedEndTime.Hour(), parsedEndTime.Minute(), 0, 0, now.Location())

	// Handle cases where peakEndTime is on the next day (e.g., 22:00 - 02:00)
	if peakEndDateTime.Before(peakStartDateTime) {
		// If current time is after peak start (e.g. 23:00 for 22:00-02:00) OR
		// if current time is before peak end on the next day (e.g. 01:00 for 22:00-02:00)
		if currentDayTime.After(peakStartDateTime) || currentDayTime.Before(peakEndDateTime.Add(24*time.Hour)) {
			// If current time is after peak start on current day (e.g. 23:00 for 22:00-02:00)
			// OR if current time is before peak end (which is on the next day, so add 24h to peakEndDateTime for comparison if current time is also on next day relative to start)
			// More robust: check if current time is between start and end, considering wrap-around
			var inPeakTime bool
			if currentDayTime.After(peakStartDateTime) { // e.g. current 23:00, start 22:00
				inPeakTime = true
			} else { // current time is on the next day relative to start time, e.g. current 01:00, start 22:00 (previous day)
				// peakEndDateTime is already set to today, so if current is before it, it's in range
				if currentDayTime.Before(peakEndDateTime) {
					inPeakTime = true
				}
			}

			if inPeakTime {
				shaperLogger.Printf("Current time %s is within peak hours (cross-day %s - %s). Removing speed limit.", currentTimeStr, peakStartTime, peakEndTime)
				if err := ClearTCRules(iface); err != nil {
					shaperLogger.Printf("Error clearing TC rules: %v", err)
				}
			} else {
				shaperLogger.Printf("Current time %s is outside peak hours (cross-day %s - %s). Applying speed limit %s.", currentTimeStr, peakStartTime, peakEndTime, limitRate)
				if err := ApplyTCRules(iface, limitRate); err != nil {
					shaperLogger.Printf("Error applying TC rules: %v", err)
				}
			}
		} else { // Normal non-crossing day range
			if currentDayTime.After(peakStartDateTime) && currentDayTime.Before(peakEndDateTime) {
				shaperLogger.Printf("Current time %s is within peak hours (%s - %s). Removing speed limit.", currentTimeStr, peakStartTime, peakEndTime)
				if err := ClearTCRules(iface); err != nil {
					shaperLogger.Printf("Error clearing TC rules: %v", err)
				}
			} else {
				shaperLogger.Printf("Current time %s is outside peak hours (%s - %s). Applying speed limit %s.", currentTimeStr, peakStartTime, peakEndTime, limitRate)
				if err := ApplyTCRules(iface, limitRate); err != nil {
					shaperLogger.Printf("Error applying TC rules: %v", err)
				}
			}
		}
	} else { // Peak time does not cross midnight
		if currentDayTime.After(peakStartDateTime) && currentDayTime.Before(peakEndDateTime) {
			shaperLogger.Printf("Current time %s is within peak hours (%s - %s). Removing speed limit.", currentTimeStr, peakStartTime, peakEndTime)
			if err := ClearTCRules(iface); err != nil {
				shaperLogger.Printf("Error clearing TC rules: %v", err)
			}
		} else {
			shaperLogger.Printf("Current time %s is outside peak hours (%s - %s). Applying speed limit %s.", currentTimeStr, peakStartTime, peakEndTime, limitRate)
			if err := ApplyTCRules(iface, limitRate); err != nil {
				shaperLogger.Printf("Error applying TC rules: %v", err)
			}
		}
	}
}

/*
// Example of how to use it in main or a service:
func main() {
    // For a one-time run:
    // ManageBandwidth()

    // For continuous monitoring (e.g., run every 5 minutes):
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            ManageBandwidth()
        }
    }
}
*/
