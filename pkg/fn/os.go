package fn

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type FnOs struct {
	FnHandler
	*Fn

	Category FnCategory
}

func (f *FnOs) Register() {
	f.fns["os.name"] = &FnEntry{
		Handler:     f.os,
		Name:        "Get OS",
		Description: "Get the current operating system",
		Category:    f.Category,
		Parameters:  []FnParameter{},
	}
	f.fns["os.arch"] = &FnEntry{
		Handler:     f.arch,
		Name:        "Get Architecture",
		Description: "Get the current architecture",
		Category:    f.Category,
		Parameters:  []FnParameter{},
	}
	f.fns["os.exec"] = &FnEntry{
		Handler:     f.exec,
		Name:        "Execute Command",
		Description: "Returns output of the command execution",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "command", Description: "The command to execute", Mandatory: true},
			{Name: "arguments", Description: "The arguments for the execution", Type: "array", Mandatory: true},
		},
	}
	f.fns["os.env.get"] = &FnEntry{
		Handler:     f.env,
		Name:        "Get Environment Variable",
		Description: "Get the value of an environment variable",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The name of the environment variable", Mandatory: true},
		},
	}
}

func (f *FnOs) os(j json.RawMessage) (json.RawMessage, error) {
	return returnRaw(runtime.GOOS), nil
}

func (f *FnOs) arch(j json.RawMessage) (json.RawMessage, error) {
	return returnRaw(runtime.GOARCH), nil
}

type cmdParams struct {
	Command   string   `json:"command" yaml:"command"`
	Arguments []string `json:"arguments" yaml:"arguments"`
}

func (f *FnOs) exec(j json.RawMessage) (json.RawMessage, error) {
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

		return returnRaw(res), nil
	})
}

func (f *FnOs) env(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stdParams) (json.RawMessage, error) {
		return returnRaw(os.Getenv(params.Value)), nil
	})
}
