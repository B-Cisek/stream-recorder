package services

import "log"

type RecorderService interface {
	StartRecording(channel string) error
	StopRecording(channel string) error
}

type recorderService struct {
}

func NewRecorderService() RecorderService {
	return &recorderService{}
}

func (s *recorderService) StartRecording(channel string) error {
	log.Printf("Starting recording for channel: %s", channel)
	// Logic to start recording
	return nil
}

func (s *recorderService) StopRecording(channel string) error {
	log.Printf("Stopping recording for channel: %s", channel)
	// Logic to stop recording
	return nil
}
