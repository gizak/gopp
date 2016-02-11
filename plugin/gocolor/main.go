package goplugin

import (
	"io"
	"os"
	"os/exec"
)

var PLUGIN = aplugin{}

type aplugin struct {
	*exec.Cmd
}

func (ap aplugin) RewriteOutput(in io.Reader, out io.Writer) error {
	ap.Cmd = exec.Command("colorgo", os.Args[1:]...)
	ap.Stdin = in
	ap.Stdout = out
	ap.Stderr = out

	if len(os.Args) > 1 && (os.Args[1] == "build" || os.Args[1] == "test") {
		return ap.Run()
	}

	_, err := io.Copy(out, in)
	return err
}
