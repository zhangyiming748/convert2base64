package main

import (
	"os/exec"
	"testing"
)

func TestGetBase64(t *testing.T) {
	cmd := exec.Command("ffmpeg")
	err := cmd.Run()
	if err != nil {
		return
	}
}
