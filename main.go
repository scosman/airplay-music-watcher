package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/scosman/airplay-music-watcher/actions"
	"github.com/scosman/airplay-music-watcher/mdns"
)

const DeviceSupportsRelayBitmask = 0x800

func main() {
	args := os.Args
	if len(args) != 2 {
		// first arg is program path, so 2 == 1...
		log.Fatal("command requires exactly 1 arg -- the path to the json config file")
	}
	jsonFilePath := args[1]
	actionRunner, err := actions.NewAirplayMusicActionRunner(jsonFilePath)
	if err != nil {
		log.Fatalf("Error parsing json config file: %v", err)
	}

	entriesCh := make(chan *mdns.AirplayFlagsEntry, 4)
	defer close(entriesCh)
	go func() {
		for entry := range entriesCh {
			// We're interested if the device is playing or not. From experimentation the 11th bit
			// AKA DeviceSupportsRelay is the most relable way of checking this
			// It stays on after you pause for about 8 mins, but does flip off eventually.
			// If you manually disconnect airplay from a device, you get the off immediately
			// https://github.com/openairplay/airplay-spec/blob/master/src/status_flags.md
			isPlaying := (DeviceSupportsRelayBitmask & entry.Flags) > 0
			fmt.Printf("Airplay Device \"%s\" event, is playing: %t\n", entry.DeviceName, isPlaying)
			actionRunner.RunActionForDeviceState(entry.DeviceName, isPlaying)
		}
	}()

	// Publishes a MDNS query for stereo.local.mdns. Not really needed, but not a problem either, and minimizing refactoring
	// of the mdns client.go so will just leave this in.
	params := mdns.DefaultParams("stereo")
	params.Entries = entriesCh
	// Timeout set to max time -- will never timeout
	params.Timeout = math.MaxInt64

	// Runs forever, does term on term signal so we'll call that a win
	mdns.Query(params)
}
