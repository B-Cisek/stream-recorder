package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

type RecorderService interface {
	StartRecording(platform, channel string) error
	StopRecording(platform, channel string) error
}

type recorderService struct {
	activeRecordings map[string]*exec.Cmd
	mu               sync.Mutex
}

func NewRecorderService() RecorderService {
	return &recorderService{
		activeRecordings: make(map[string]*exec.Cmd),
	}
}

func (s *recorderService) StartRecording(platform, channel string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s-%s", platform, channel)

	if _, exists := s.activeRecordings[key]; exists {
		return fmt.Errorf("recording already in progress for %s channel: %s", platform, channel)
	}

	strategy, err := GetPlatformStrategy(platform)
	if err != nil {
		return err
	}

	// Ensure recordings directory exists
	recordingDir := "./recordings"
	if err := os.MkdirAll(recordingDir, 0755); err != nil {
		return fmt.Errorf("failed to create recordings directory: %w", err)
	}

	args := strategy.GetStreamlinkArgs(channel, recordingDir)
	cmd := exec.Command("streamlink", args...)

	// Optional: set stdout/stderr to log or a file for debugging
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	log.Printf("Starting recording for %s channel: %s, output: %s", strategy.GetStrategyName(), channel, args)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start streamlink: %w", err)
	}

	s.activeRecordings[key] = cmd

	// Monitor process in a goroutine
	go func() {
		err := cmd.Wait()
		s.mu.Lock()
		defer s.mu.Unlock()

		if _, exists := s.activeRecordings[key]; exists {
			delete(s.activeRecordings, key)
			if err != nil {
				log.Printf("Recording for %s channel %s finished with error: %v", strategy.GetStrategyName(), channel, err)
			} else {
				log.Printf("Recording for %s channel %s finished successfully", strategy.GetStrategyName(), channel)
			}
		}
	}()

	return nil
}

func (s *recorderService) StopRecording(platform, channel string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s-%s", platform, channel)

	cmd, exists := s.activeRecordings[key]
	if !exists {
		return fmt.Errorf("no active recording found for %s channel: %s", platform, channel)
	}

	log.Printf("Stopping recording for %s channel: %s", platform, channel)

	// Send interrupt signal to gracefully stop streamlink
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		// Fallback to Kill if Interrupt fails
		log.Printf("Failed to send interrupt to process, killing: %v", err)
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
	}

	// We remove it from the map here to prevent further operations,
	// but the goroutine will also try to remove it.
	// The goroutine checks existence, so it's safe.
	delete(s.activeRecordings, key)

	return nil
}
