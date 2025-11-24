package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type RecorderService interface {
	StartRecording(channel string) error
	StopRecording(channel string) error
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

func (s *recorderService) StartRecording(channel string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.activeRecordings[channel]; exists {
		return fmt.Errorf("recording already in progress for channel: %s", channel)
	}

	// Ensure recordings directory exists
	recordingsDir := "./recordings"
	if err := os.MkdirAll(recordingsDir, 0755); err != nil {
		return fmt.Errorf("failed to create recordings directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s-%s.mp4", channel, timestamp)
	filePath := filepath.Join(recordingsDir, filename)

	// Construct streamlink command
	// streamlink twitch.tv/<channel> best -o <filepath>
	cmd := exec.Command("streamlink", "twitch.tv/"+channel, "best", "-o", filePath)

	// Optional: set stdout/stderr to log or a file for debugging
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	log.Printf("Starting recording for channel: %s, output: %s", channel, filePath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start streamlink: %w", err)
	}

	s.activeRecordings[channel] = cmd

	// Monitor process in a goroutine
	go func() {
		err := cmd.Wait()
		s.mu.Lock()
		defer s.mu.Unlock()

		// Check if it was still in the map (might have been stopped manually)
		if _, exists := s.activeRecordings[channel]; exists {
			delete(s.activeRecordings, channel)
			if err != nil {
				log.Printf("Recording for channel %s finished with error: %v", channel, err)
			} else {
				log.Printf("Recording for channel %s finished successfully", channel)
			}
		}
	}()

	return nil
}

func (s *recorderService) StopRecording(channel string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cmd, exists := s.activeRecordings[channel]
	if !exists {
		return fmt.Errorf("no active recording found for channel: %s", channel)
	}

	log.Printf("Stopping recording for channel: %s", channel)

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
	delete(s.activeRecordings, channel)

	return nil
}
