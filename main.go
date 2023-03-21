package main

import (
	"fmt"
	"math"

	"github.com/scosman/airplay-music-watcher/mdns"
)

const DeviceSupportsRelayBitmask = 0x800

func main() {
	// Hello world, the web server

	entriesCh := make(chan *mdns.AirplayFlagsEntry, 4)
	defer close(entriesCh)
	go func() {
		for entry := range entriesCh {
			//fmt.Printf("Airplay Event For: %v, flags raw: %v, flags parsed: %s\n", entry.Name, entry.RawFlags, fmt.Sprintf("%016x", entry.Flags))

			// We're interested if the device is playing or not. From experimentation the 11th bit
			// AKA DeviceSupportsRelay is the most relable way of checking this
			// It stays on after you pause for about 8 mins, but does flip off eventually.
			// If you manually disconnect airplay from a device, you get the off immediately
			// https://github.com/openairplay/airplay-spec/blob/master/src/status_flags.md
			isPlaying := (DeviceSupportsRelayBitmask & entry.Flags) > 0
			fmt.Printf("Device \"%s\" is playing: %t\n", entry.DeviceName, isPlaying)
		}
	}()

	// Start the lookup

	// Publishes a MDNS query for stereo.local.mdns. Not really needed, but not a problem either, and minimizing refactoring
	// of the mdns client.go so will just leave this in.
	params := mdns.DefaultParams("stereo")
	params.Entries = entriesCh
	// TODO can we remove this?
	params.DisableIPv6 = true
	// Timeout set to max time -- will never timeout
	params.Timeout = math.MaxInt64

	// Runs forever, does term on term signal so we'll call that a win
	mdns.Query(params)
}
