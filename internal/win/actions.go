package win

import (
	"fmt"
	"os/exec"
	"wike/internal/settings"
)

func executeAction(action settings.Action) {
	fmt.Printf("Executing action: %+v\n", action)

	if len(action.Kb) > 0 {
		sendKeys(action.Kb, true, true)
	}

	if action.Cmd != nil {
		fmt.Printf("Executing command: %s\n", *action.Cmd)
		exec.Command(*action.Cmd).Run()
	}

	if action.Launch != nil {
		fmt.Printf("Launching application: %s\n", *action.Launch)
		openOrFocus(*action.Launch)
	}
}
