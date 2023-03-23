[![Build and Test](https://github.com/scosman/airplay-music-watcher/actions/workflows/test.yml/badge.svg)](https://github.com/scosman/airplay-music-watcher/actions/workflows/test.yml)
[![Format and Vet](https://github.com/scosman/airplay-music-watcher/actions/workflows/format_check.yml/badge.svg)](https://github.com/scosman/airplay-music-watcher/actions/workflows/format_check.yml)
[![Release Binary Build](https://github.com/scosman/airplay-music-watcher/actions/workflows/release.yml/badge.svg)](https://github.com/scosman/airplay-music-watcher/actions/workflows/release.yml)

# Airplay Music Watcher

A golang service which watches for devices starting or stopping playing airplay audio. When a device starts or stops, you can issue system command line tasks.

This is primarily useful for home automations. For example: detecting music playing on an airplay dongle, and turning on your receiver (via a smart switch, or IR command). I use it with homebridge and an IR blaster to turn a nice vintage receiver on an off for airplay.

## Note on event timing

If you start or stop playing it sends the commands immediately. If you simply pause, the stop command comes eventually but not instantly. From experimentation, pausing triggers a stop event after approximately 8 minutes.

## Setup

1) Download the [appriopiate binary](https://github.com/scosman/airplay-music-watcher/blob/main/README.md#which-binary-should-i-download) for your system from the latest release.
2) Create a config file using JSON format below. You can add as many entries as you want. The device_name and action are used to identify the trigger, and the command is run when the event is triggered
3) Run the server via command line and ensure it works as you want it to: `./airplay-music-watcher ./your_json_config.json`
4) Setup the command to run on boot using systemd, launchd, or your daemon manager of choice!

## Which Binary Should I Download?

Releases can be found here: https://github.com/scosman/airplay-music-watcher/releases 

Releases have several builds, and you'll need the right one for your machine:

 - Raspberry Pi: linux-arm
 - Mac with M1/M2/M* processor: darwin-arm64
 - Older Mac with non M processor: darwin-amd64
 - Windows: windows-amd64, and if it's old maybe windows-386
 - Windows with arm processor: I don't believe you
 - Linux: you got this

## JSON Config File Example

This example json file would say "stereo starting" when the Airplay device named "Stereo" starts playing, and say "stereo stopping" when it stops.

Fields:

 - Device name: the name of the Airplay device, as you have setup in the Apple Home app. Case sensitive.
 - Action: either "start_playing" or "end_playing"
 - Command: a command line that will be run when the event triggets (examples: "echo hello", "curl http://...")

```
{
    "actions": [
        {
            "device_name": "Stereo",
            "action": "start_playing",
            "command": "say 'stereo starting'"
        },
        {
            "device_name": "Stereo",
            "action": "end_playing",
            "command": "say 'stereo stopping'"
        }
    ]
}
```

For homebridge users, this works well with the [homebridge-config-ui-x API](https://github.com/oznu/homebridge-config-ui-x/wiki/API-Reference). Use curl to call the API and enable/disable smart home devices (like a smart plug controlling your vintage audio amp).

## How this works

This process monitors UDP traffic on your network for mdns records from Airplay devices. Certain airplay mdns TXT records include a [bitmask](https://github.com/openairplay/airplay-spec/blob/master/src/status_flags.md) which let's us infer the device's state (playing or not).

It doesn't require any special permissions, just network access. It's simply reading un-encrypted multicast UDP (mdns) traffic. It needs to be run on the same network as the device in question.

The mdns monitoring and parsing code was extracted from: https://github.com/hashicorp/mdns

## Licence

MIT Licence
