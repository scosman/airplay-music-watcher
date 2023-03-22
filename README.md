[![Build and Test](https://github.com/scosman/airplay-music-watcher/actions/workflows/test.yml/badge.svg)](https://github.com/scosman/airplay-music-watcher/actions/workflows/test.yml)
[![Format and Vet](https://github.com/scosman/airplay-music-watcher/actions/workflows/format_check.yml/badge.svg)](https://github.com/scosman/airplay-music-watcher/actions/workflows/format_check.yml)

# Airplay Music Watcher

A golang service which watches for devices starting or stopping playing airplay audio. When a device starts or stops, you can issue system command line tasks.

This is primarily useful for home automations. For example: detecting music playing on an airplay dongle, and turning on your receiver (via a smart switch, or IR command). I use it with homebridge and an IR blaster to turn a nice vintage receiver on an off for airplay.

## Note on event timing

If you start or stop playing it sends the commands immediately. If you simply pause, the stop command comes eventually but not instantly. From experimentation, pausing triggers a stop event after approximately 8 minutes.

## Setup

1) Download the appriopiate binary for your system from the latest release: https://github.com/scosman/airplay-music-watcher/releases 
2) Create a config file using JSON format below. You can add as many entries as you want. The device_name and action are used to identify the trigger, and the command and command_args are what is run when the event is triggered
3) Run the server via command line and ensure it works as you want it to: `./airplay-music-watcher ./your_json_config.json`
4) Setup the command to run on boot using systemd, launchd, or your daemon manager of choice!

## Which Binary Should I Download?

Releases have several builds, and you'll need the right one for your machine.

 - Raspberry Pi: linux-arm
 - Mac with M1/M2/M* processor: darwin-arm64
 - Older Mac with non M processor: darwin-amd64
 - Windows: windows-amd64, and if it's old maybe windows-386
 - Windows with arm processor: I don't believe you
 - Linux: you got this
 
 Releases can be found here: https://github.com/scosman/airplay-music-watcher/releases 

## JSON Config File Example

This json file would say "stereo starting" when the Airplay device named "Stereo" starts playing, and say "stereo stopping" when it stops.

```
{
    "actions": [
        {
            "device_name": "Stereo",
            "action": "start_playing",
            "command": "say",
            "command_args": "'stereo starting'"
        },
        {
            "device_name": "Stereo",
            "action": "end_playing",
            "command": "say",
            "command_args": "'stereo stopping'"
        }
    ]
}
```
