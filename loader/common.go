package loader

import (
	"bytes"
	"go.dedis.ch/onet/v3/log"
	"io"
	"os"
	"os/exec"
)

// DBSettings stores the database settings
type DBSettings struct {
	DBhost     string
	DBport     int
	DBuser     string
	DBpassword string
	DBname     string
}

// ExecuteScript executes a script with a specific path
func ExecuteScript(command string, args ...string) error {
	// Display just the stderr if an error occurs
	cmd := exec.Command(command, args...)
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	// Execute the command
	if err := cmd.Run(); err != nil {
		log.Lvl1("Error when running command.  Error log:", cmd.Stderr)
		log.Lvl1("Got command status:", err.Error())
		return err
	}

	return nil
}
