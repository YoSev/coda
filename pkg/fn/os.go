package fn

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/yosev/coda/internal/utils"
)

type fnOs struct {
	category FnCategory
}

func (f *fnOs) init(fn *Fn) {
	fn.register("os.name", &FnEntry{
		Handler:     f.os,
		Name:        "Get OS",
		Description: "Get the current operating system",
		Category:    f.category,
		Parameters:  []FnParameter{},
	})
	fn.register("os.arch", &FnEntry{
		Handler:     f.arch,
		Name:        "Get Architecture",
		Description: "Get the current architecture",
		Category:    f.category,
		Parameters:  []FnParameter{},
	})
	fn.register("os.exec", &FnEntry{
		Handler:     f.exec,
		Name:        "Execute Command",
		Description: "Returns output of the command execution",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "command", Description: "The command to execute", Mandatory: true},
			{Name: "arguments", Description: "The arguments for the execution", Type: "array", Mandatory: true},
		},
	})
	fn.register("os.env.get", &FnEntry{
		Handler:     f.env,
		Name:        "Get Environment Variable",
		Description: "Get the value of an environment variable",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The name of the environment variable", Mandatory: true},
		},
	})
}

func (f *fnOs) os(j json.RawMessage) (json.RawMessage, error) {
	return utils.ReturnRaw(runtime.GOOS), nil
}

func (f *fnOs) arch(j json.RawMessage) (json.RawMessage, error) {
	return utils.ReturnRaw(runtime.GOARCH), nil
}

type cmdParams struct {
	Command   string   `json:"command" yaml:"command"`
	Arguments []string `json:"arguments" yaml:"arguments"`
}

func (f *fnOs) exec(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *cmdParams) (json.RawMessage, error) {
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

		return utils.ReturnRaw(res), nil
	})
}

func (f *fnOs) env(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stdParams) (json.RawMessage, error) {
		return utils.ReturnRaw(os.Getenv(params.Value)), nil
	})
}
