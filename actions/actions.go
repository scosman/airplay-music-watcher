package actions

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os/exec"
)

type ActionName string

type AirplayCommandLineAction struct {
	DeviceName  string     `json:"device_name"`
	Command     string     `json:"command"`
	CommandArgs string     `json:"command_args"`
	ActionName  ActionName `json:"action"`
}

type AirplayMusicActionRunner struct {
	Actions []*AirplayCommandLineAction `json:"actions"`
}

const (
	ACTION_NAME_START_PLAYING ActionName = "start_playing"
	ACTION_NAME_END_PLAYING   ActionName = "end_playing"
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

	return &parsedRunner, nil
}

func (r *AirplayMusicActionRunner) RunActionForDeviceState(deviceName string, isPlaying bool) {
	for _, action := range r.Actions {
		if action.DeviceName == deviceName {
			if (isPlaying && action.ActionName == ACTION_NAME_START_PLAYING) || (!isPlaying && action.ActionName == ACTION_NAME_END_PLAYING) {
				log.Printf("Running command: %v %v", action.Command, action.CommandArgs)
				cmd := exec.Command(action.Command, action.CommandArgs)
				if err := cmd.Run(); err != nil {
					log.Printf("Error running command: %v %v", action.Command, action.CommandArgs)
				}
			}
		}
	}
}
