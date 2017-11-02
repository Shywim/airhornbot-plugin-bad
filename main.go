package main

import (
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jonas747/dca"
)

const (
	soundPrefix = "@bad_"
	url         = "http://cdn.shywim.fr/antoine/snd/"
	urlSuffix   = ".ogg"
)

// Name of your plugin, must be unique from all other plugins
var Name = "com.matthieuharle.bad"

func randomRange(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

// Handle is called to check if your plugin can handle the sound
// if you return true here, GetSound is expected to return valid sound data
func Handle(name string) bool {
	if !strings.HasPrefix(name, soundPrefix) {
		return false
	}

	sound := strings.TrimPrefix(name, soundPrefix)
	if sound == "salut" {
		return true
	}
	if sound == "_random" {
		return true
	}

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
	sound := strings.TrimPrefix(name, soundPrefix)
	if sound == "salut" {
		sound = sound + strconv.Itoa(randomRange(1, 37))
	} else if sound == "_random" {
		resp, err := http.Get(url)
		if err != nil {
			return nil
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
		index := string(body)
		resp.Body.Close()

		re := regexp.MustCompile(">(.*?).ogg")
		allSounds := re.FindAllStringSubmatch(index, -1)
		if allSounds == nil {
			return nil
		}
		sound = allSounds[randomRange(0, len(allSounds))][1]
	}

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
