package utility

import (
	"fmt"
	"time"
)

func GetTimestampedName(experimentName string) string {
	// Get the location for JST
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Printf("failed to load location. Abort using timestamp. %s", err)
		return fmt.Sprintf("%s", experimentName)
	}

	// Get the current time in JST
	now := time.Now().In(loc)

	// Format the time in the format "mm-dd-hh-mm-ss"
	timeString := now.Format("01-02-15-04-05")
	return fmt.Sprintf("%s-%s", experimentName, timeString)
}

func GetS3Key(experimentName, deploymentName string) string {
	return fmt.Sprintf("%s/%s", experimentName, deploymentName)
}

func ParseDurationWithDefault(durationStr string, defaultDuration time.Duration) time.Duration {
	if duration, err := time.ParseDuration(durationStr); err == nil {
		return duration
	}
	return defaultDuration
}
