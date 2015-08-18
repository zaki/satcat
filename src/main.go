package main

import (
	"io/ioutil"
	"strings"
	"github.com/zaki/satcat"
)

func main() {
	for satellite := range processSatellite(parseSatellites()) {
		if !satellite.IsDebris && satellite.PayloadFlag && !satellite.HasDecayed {
			satellite.Print()
		}
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func parseSatellites() <-chan string {
	satellites := make(chan string)

	dat, err := ioutil.ReadFile("dat/satcat.txt")
	checkError(err)

	lines := strings.Split(string(dat), "\n")

	go func() {
		for _, line := range lines {
			satellites <- line
		}
		close(satellites)
	}()

	return satellites
}

func processSatellite(satellites <-chan string) <-chan satcat.SatelliteEntry {
	entries := make(chan satcat.SatelliteEntry)

	go func() {
		for satellite := range satellites {
			if len(satellite) > 0 {
				entries <- parseSatelliteEntry(satellite)
			}
		}
		close(entries)
	}()

	return entries
}

func parseSatelliteEntry(line string) satcat.SatelliteEntry {
	entry := satcat.SatelliteEntry{}
	entry.FromString(line)

	return entry
}
