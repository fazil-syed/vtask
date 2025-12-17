package handlers

// implement a interface with a stub struct and a stub Transcribe method that returns a fixed string or the the transcription user provided
type STTService interface {
	Transcribe(filePath string) (string, error)
}
type StubSTTService struct {
	Transcription string
}
func (s *StubSTTService) Transcribe(filePath string) (string, error) {
	return s.Transcription, nil
}