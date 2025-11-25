package services

import (
	"fmt"
	"path/filepath"
	"time"
)

type PlatformStrategy interface {
	GetStreamlinkArgs(channel string, recordingDir string) []string
	GetStrategyName() string
}

type TwitchStrategy struct{}

func (t *TwitchStrategy) GetStreamlinkArgs(channel string, recordingDir string) []string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s-%s.mp4", channel, timestamp)
	filePath := filepath.Join(recordingDir, filename)
	return []string{"twitch.tv/" + channel, "best", "-o", filePath}
}

func (t *TwitchStrategy) GetStrategyName() string {
	return "twitch"
}

// GetPlatformStrategy returns the appropriate strategy based on the platform name
func GetPlatformStrategy(platform string) (PlatformStrategy, error) {
	switch platform {
	case "twitch":
		return &TwitchStrategy{}, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}
