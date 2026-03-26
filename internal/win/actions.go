package win

import (
	"os/exec"
	"wike/internal/config"
	"wike/internal/logger"
)

func executeAction(action config.Action) {
	logger.Printf("Executing action: %+v\n", action)

	if len(action.Kb) > 0 {
		sendKeys(action.Kb, true, true)
	}

	if action.Cmd != nil {
		logger.Printf("Executing command: %s\n", *action.Cmd)
		exec.Command(*action.Cmd).Run()
	}

	if action.Launch != nil {
		logger.Printf("Launching application: %s\n", *action.Launch)
		openOrFocus(*action.Launch)
	}
}
