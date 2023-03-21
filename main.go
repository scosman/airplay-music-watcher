package main

import (
	"fmt"
	"time"

	"github.com/scosman/airplay-music-watcher/mdns"
)

const DeviceSupportsRelayBitmask = 0x800

func main() {
	// Hello world, the web server

	entriesCh := make(chan *mdns.AirplayFlagsEntry, 4)
	go func() {
		for entry := range entriesCh {
			//fmt.Printf("Airplay Event For: %v, flags raw: %v, flags parsed: %s\n", entry.Name, entry.RawFlags, fmt.Sprintf("%016x", entry.Flags))

			// We're interested if the device is playing or not. From experimentation the 11th bit
			// AKA DeviceSupportsRelay is the most relable way of checking this
			// It stays on after you pause for about 8 mins, but does flip off
			// https://github.com/openairplay/airplay-spec/blob/master/src/status_flags.md
			isPlaying := (DeviceSupportsRelayBitmask & entry.Flags) > 0
			fmt.Printf("Device \"%s\" is playing: %t\n", entry.Name, isPlaying)
		}
	}()

	// Start the lookup
	//mdns.Lookup("_googlecast._tcp.local.", entriesCh)

	timeout := 300 * time.Second
	params := mdns.DefaultParams("stereo")
	params.Entries = entriesCh
	// TODO can we remove this
	params.DisableIPv6 = true
	params.Timeout = timeout

	mdns.Query(params)

	// TODO
	time.Sleep(timeout)

	close(entriesCh)
	time.Sleep(1 * time.Second)
}
