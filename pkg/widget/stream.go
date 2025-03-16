package widget

type StreamType string

const (
	StreamTypeInput    StreamType = "input"
	StreamTypeOutput   StreamType = "output"
	StreamTypePlayback StreamType = "playback" // audio playback device
	StreamTypeCapture  StreamType = "capture"  // audio capture device
	StreamTypeMonitor  StreamType = "monitor"  // audio monitor
)

type StreamComponent struct {
	Type StreamType
}
