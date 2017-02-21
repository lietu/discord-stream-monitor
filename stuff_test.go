package main

import (
	"testing"
	"log"
)

func TestStuff(t *testing.T) {
	strings := []string{
		"Come watch https://www.twitch.tv/lietu because I'm awesome",
		"Come watch https://twitch.tv/lietu because I'm awesome",
		"Come watch http://www.twitch.tv/lietu because I'm awesome",
		"Come watch http://twitch.tv/lietu because I'm awesome",
		"Come watch www.twitch.tv/lietu because I'm awesome",
		"Come watch twitch.tv/lietu because I'm awesome",
		"Come watch https://www.lietu.tv/live because I'm awesome",
		"Come watch https://lietu.tv/live because I'm awesome",
		"Come watch http://www.lietu.tv/live because I'm awesome",
		"Come watch http://lietu.tv/live because I'm awesome",
		"Come watch www.lietu.tv/live because I'm awesome",
		"Come watch lietu.tv/live because I'm awesome",
	}

	for _, s := range strings {
		stream := getStreamerAdvertised(s)

		if stream != "lietu" {
			t.Errorf("Did not match expected streamer: %s", s)
		} else {
			log.Printf("%s -> %s: OK!", s, stream)
		}
	}

	strings = []string{
		"Come watch me because I'm awesome",
		"Test",
		"Test,",
		"lulz",
		"http://youtube.com/aasdasd",
		"https://stackoverflow.com/questions/24613271/golang-is-conversion-between-different-struct-types-possible",
		"twitch.tv/directory/following/live",
	}

	for _, s := range strings {
		stream := getStreamerAdvertised(s)

		if stream != "" {
			t.Errorf("False positive match %s: %s", stream, s)
		} else {
			log.Printf("%s -> %s: OK!", s, stream)
		}
	}
}