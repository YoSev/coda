package fn

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"runtime"
	"strings"
)

func (f *Fn) GetOS(j json.RawMessage) (json.RawMessage, error) {
	return returnString(runtime.GOOS), nil
}

func (f *Fn) GetArch(j json.RawMessage) (json.RawMessage, error) {
	return returnString(runtime.GOARCH), nil
}

type cmdParams struct {
	Command   string   `json:"command" yaml:"command"`
	Arguments []string `json:"arguments" yaml:"arguments"`
}

func (f *Fn) Exec(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *cmdParams) (json.RawMessage, error) {
		cmd := exec.Command(params.Command, params.Arguments...)
		var stdOutBuffer bytes.Buffer
		var stdErrBuffer bytes.Buffer
		cmd.Stdout = &stdOutBuffer
		cmd.Stderr = &stdErrBuffer

		err := cmd.Run()
		if err != nil {
			return nil, err
		}

		res := &struct {
			StdOut string `json:"stdout" yaml:"stdout"`
			StdErr string `json:"stderr" yaml:"stderr"`
			Code   int    `json:"code" yaml:"code"`
		}{StdOut: strings.Trim(stdOutBuffer.String(), "\n"), StdErr: strings.Trim(stdErrBuffer.String(), "\n"), Code: cmd.ProcessState.ExitCode()}

		return returnAny(res), nil
	})
}
