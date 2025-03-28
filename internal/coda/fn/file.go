package fn

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type sourceFileParams struct {
	Source string `json:"source" yaml:"source"`
}

func (f *Fn) Size(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *sourceFileParams) (json.RawMessage, error) {
		fileInfo, err := os.Stat(params.Source)
		if err != nil {
			return nil, err
		}

		return returnAny(fileInfo.Size()), nil
	})
}

func (f *Fn) Modified(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *sourceFileParams) (json.RawMessage, error) {
		fileInfo, err := os.Stat(params.Source)
		if err != nil {
			return nil, err
		}

		return returnInt64(fileInfo.ModTime().UnixMilli()), nil
	})
}

func (f *Fn) Delete(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *sourceFileParams) (json.RawMessage, error) {
		os.Remove(params.Source)

		return nil, nil
	})
}

type copyMoveParams struct {
	Source      string `json:"source" yaml:"source"`
	Destination string `json:"destination" yaml:"destination"`
}

func (f *Fn) Copy(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *copyMoveParams) (json.RawMessage, error) {
		sourceFileStat, err := os.Stat(params.Source)
		if err != nil {
			return nil, err
		}

		if !sourceFileStat.Mode().IsRegular() {
			return nil, fmt.Errorf("%s is not a regular file", params.Source)
		}

		source, err := os.Open(params.Source)
		if err != nil {
			return nil, err
		}
		defer source.Close()

		destination, err := os.Create(params.Destination)
		if err != nil {
			return nil, err
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)
		if err != nil {
			return nil, err
		}

		return returnAny(params.Destination), err
	})
}

func (f *Fn) Move(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *copyMoveParams) (json.RawMessage, error) {
		err := os.Rename(params.Source, params.Destination)
		if err != nil {
			return nil, err
		}
		return returnAny(params.Destination), nil
	})
}

func (f *Fn) Read(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *sourceFileParams) (json.RawMessage, error) {
		file, err := os.ReadFile(params.Source)
		if err != nil {
			return nil, err
		}

		return returnAny(file), nil
	})
}

type writeFileParams struct {
	Destination string `json:"destination" yaml:"destination"`
	Value       string `json:"value" yaml:"value"`
}

func (f *Fn) Write(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *writeFileParams) (json.RawMessage, error) {
		err := os.WriteFile(params.Destination, []byte(params.Value), 0644)
		if err != nil {
			return nil, err
		}
		return returnAny(params.Destination), nil
	})
}
