package actions

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
)

type ActionName string

type AirplayCommandLineAction struct {
	DeviceName string     `json:"device_name"`
	Command    string     `json:"command"`
	ActionName ActionName `json:"action"`
}

type AirplayMusicActionRunner struct {
	Actions                []*AirplayCommandLineAction `json:"actions"`
	lastKnownStateOfDevice map[string]LastKnownState
}

const (
	ACTION_NAME_START_PLAYING ActionName = "start_playing"
	ACTION_NAME_END_PLAYING   ActionName = "end_playing"
)

type LastKnownState string

const (
	LAST_KNOWN_STATE_UNKNOWN LastKnownState = "unknown"
	LAST_KNOWN_STATE_PLAYING LastKnownState = "playing"
	LAST_KNOWN_STATE_STOPPED LastKnownState = "stopped"
)

// func DefaultParams(service string) *QueryParam {
func NewAirplayMusicActionRunner(configFilePath string) (*AirplayMusicActionRunner, error) {
	configBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	var parsedRunner AirplayMusicActionRunner
	err = json.Unmarshal(configBytes, &parsedRunner)
	if err != nil {
		return nil, err
	}

	for _, action := range parsedRunner.Actions {
		if action.ActionName != ACTION_NAME_START_PLAYING && action.ActionName != ACTION_NAME_END_PLAYING {
			return nil, errors.New("Invalid action name")
		}
	}
	parsedRunner.lastKnownStateOfDevice = make(map[string]LastKnownState)

	return &parsedRunner, nil
}

func (r *AirplayMusicActionRunner) RunActionForDeviceState(deviceName string, isPlaying bool) {
	priorState := r.lastKnownStateOfDevice[deviceName]
	if (priorState == LAST_KNOWN_STATE_PLAYING && isPlaying) || (priorState == LAST_KNOWN_STATE_STOPPED && !isPlaying) {
		// we already sent this, can skip
		return
	}
	if isPlaying {
		r.lastKnownStateOfDevice[deviceName] = LAST_KNOWN_STATE_PLAYING
	} else {
		r.lastKnownStateOfDevice[deviceName] = LAST_KNOWN_STATE_STOPPED
	}

	for _, action := range r.Actions {
		if action.DeviceName == deviceName {
			if (isPlaying && action.ActionName == ACTION_NAME_START_PLAYING) || (!isPlaying && action.ActionName == ACTION_NAME_END_PLAYING) {
				r.runActionForDevice(deviceName, isPlaying, *action)
			}
		}
	}
}

func (r *AirplayMusicActionRunner) runActionForDevice(deviceName string, isPlaying bool, action AirplayCommandLineAction) {
	log.Printf("Running command: %s\n", action.Command)
	cmd := exec.Command("sh", "-c", action.Command)
	if runtime.GOOS == "windows" {
		// Need someome to test this. From stack overflow, but no windows box...
		cmd = exec.Command("C:\\Windows\\System32\\cmd.exe", "/c", action.Command)
	}

	if err := cmd.Run(); err != nil {
		log.Printf("Error running command: %s\n", action.Command)
	}
}
