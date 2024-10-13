package ffmpeg_go

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Probe Run ffprobe on the specified file and return a JSON representation of the output.
func Probe(probePath, fileName string, kwargs ...KwArgs) (string, error) {
	return ProbeWithTimeout(probePath, fileName, 0, MergeKwArgs(kwargs))
}

func ProbeWithTimeout(probePath, fileName string, timeOut time.Duration, kwargs KwArgs) (string, error) {
	args := KwArgs{
		"show_format":  "",
		"show_streams": "",
		"of":           "json",
	}

	return ProbeWithTimeoutExec(probePath, fileName, timeOut, MergeKwArgs([]KwArgs{args, kwargs}))
}

func ProbeWithTimeoutExec(probePath, fileName string, timeOut time.Duration, kwargs KwArgs) (string, error) {
	args := ConvertKwargsToCmdLineArgs(kwargs)
	args = append(args, fileName)
	ctx := context.Background()
	if timeOut > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(context.Background(), timeOut)
		defer cancel()
	}
	if probePath == "" {
		probePath = "ffprobe"
	}
	cmd := exec.CommandContext(ctx, probePath, args...)
	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf
	fmt.Println(cmd.String())
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}
