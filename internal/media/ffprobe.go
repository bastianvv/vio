package media

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

// FFProbeOutput models the subset of ffprobe JSON we care about.
type FFProbeOutput struct {
	Format struct {
		Filename string `json:"filename"`
		Duration string `json:"duration"` // seconds as string
		Size     string `json:"size"`     // bytes as string
	} `json:"format"`
	Streams []struct {
		Index     int    `json:"index"`
		CodecType string `json:"codec_type"` // "video", "audio", "subtitle"
		CodecName string `json:"codec_name"`

		// Video
		Width  int `json:"width"`
		Height int `json:"height"`

		// Audio
		Channels int `json:"channels"`

		// Subtitles
		Tags struct {
			Language string `json:"language"`
			Title    string `json:"title"`
		} `json:"tags"`

		Disposition struct {
			Default int `json:"default"`
			Forced  int `json:"forced"`
		} `json:"disposition"`
	} `json:"streams"`
}

// RunFFProbe executes ffprobe and returns parsed JSON.
func RunFFProbe(path string) (*FFProbeOutput, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		path,
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var result FFProbeOutput
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
