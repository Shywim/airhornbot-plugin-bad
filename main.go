package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/jonas747/dca"
)

const (
	soundPrefix = "@bad_"
	url         = "http://cdn.shywim.fr/antoine/snd/"
	urlSuffix   = ".ogg"
)

// Name of your plugin, must be unique from all other plugins
var Name = "com.matthieuharle.bad"

// Handle is called to check if your plugin can handle the sound
// if you return true here, GetSound is expected to return valid sound data
func Handle(name string) bool {
	if !strings.HasPrefix(name, soundPrefix) {
		return false
	}

	sound := strings.TrimLeft(name, soundPrefix)
	resp, err := http.Head(url + sound + urlSuffix)
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

// GetSound is called when one of your sound will be played
func GetSound(name string) (buffer [][]byte) {
	sound := strings.TrimLeft(name, soundPrefix)
	resp, err := http.Get(url + sound + urlSuffix)
	if err != nil {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	defer resp.Body.Close()

	dcaSessions, err := dca.EncodeMem(resp.Body, dca.StdEncodeOptions)
	defer dcaSessions.Cleanup()
	if err != nil {
		return nil
	}

	decoder := dca.NewDecoder(dcaSessions)
	for {
		frame, err := decoder.OpusFrame()
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return buffer
			}

			return nil
		}

		buffer = append(buffer, frame)
	}
}
