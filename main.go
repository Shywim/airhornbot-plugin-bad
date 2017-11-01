package main

import (
	"io"
	"net/http"
	"strings"
)

const (
	soundPrefix = "@bad_"
	url         = "http//cdn.shywim.fr/antoine/snd/"
	urlSuffix   = ".ogg"
)

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

	var buf []byte
	_, err = io.ReadFull(resp.Body, buf)
	if err != nil {
		return nil
	}

	buffer = append(buffer, buf)
	return
}
