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
		"Come watch lietu.tv/live/ because I'm awesome",
		"https://www.lietu.tv/live/",
		"https://www.lietu.tv/live",
		"https://lietu.tv/live",
		"http://www.lietu.tv/live",
		"http://lietu.tv/live",
		"www.lietu.tv/live",
		"lietu.tv/live",
		"https://www.twitch.tv/lietu",
		"https://twitch.tv/lietu",
		"http://www.twitch.tv/lietu",
		"http://twitch.tv/lietu",
		"www.twitch.tv/lietu",
		"twitch.tv/lietu",
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

	strings = []string{
		"<Minin> blind playtrough of Fallout 2 https://www.twitch.tv/minin/",
	}

	for _, s := range strings {
		stream := getStreamerAdvertised(s)

		if stream != "minin" {
			t.Errorf("Did not match expected streamer: %s", s)
		} else {
			log.Printf("%s -> %s: OK!", s, stream)
		}
	}
}