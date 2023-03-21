package actions

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestActionsParseValid(t *testing.T) {
	workingDirectory, _ := os.Getwd()
	testFile := filepath.Join(workingDirectory, "../test_files/valid_actions_1.json")
	log.Printf("wd: %v", testFile)

	actions, err := NewAirplayMusicActionRunner(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions.Actions) != 2 {
		t.Fatal("parse error -- wrong count")
	}
	first := actions.Actions[0]
	if first.DeviceName != "device1" || first.Command != "echo" || first.ActionName != "start_playing" {
		t.Fatal("didn't parse first action")
	}
	second := actions.Actions[1]
	if second.DeviceName != "Stereo" || second.Command != "echo" || second.ActionName != "end_playing" {
		t.Fatal("didn't parse second action")
	}
}
