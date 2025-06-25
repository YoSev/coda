package fn

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/yosev/coda/internal/utils"
)

type fnFile struct {
	category FnCategory
}

func (f *fnFile) init(fn *Fn) {
	fn.register("file.size", &FnEntry{
		Handler:     f.size,
		Name:        "File size",
		Description: "Get the size of a file",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "source", Description: "The path of the source", Mandatory: true},
		},
	})
	fn.register("file.modified", &FnEntry{
		Handler:     f.modified,
		Name:        "File modified",
		Description: "Get the modify date of a file as unix timestamp in milliseconds",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "source", Description: "The path of the source", Mandatory: true},
		},
	})
	fn.register("file.copy", &FnEntry{
		Handler:     f.copy,
		Name:        "Copy file",
		Description: "Copy a file",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "source", Description: "The path of the source", Mandatory: true},
			{Name: "destination", Description: "The path of the destination", Mandatory: true},
		},
	})
	fn.register("file.delete", &FnEntry{
		Handler:     f.delete,
		Name:        "Delete file",
		Description: "Delete a file",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "source", Description: "The path of the file to delete", Mandatory: true},
		},
	})
	fn.register("file.move", &FnEntry{
		Handler:     f.move,
		Name:        "Move file",
		Description: "Move a file",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "source", Description: "The path of the file to move", Mandatory: true},
			{Name: "destination", Description: "The path of the destination", Mandatory: true},
		},
	})
	fn.register("file.read", &FnEntry{
		Handler:     f.read,
		Name:        "Read file",
		Description: "Read a file",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "source", Description: "The path of the file to read", Mandatory: true},
		},
	})
	fn.register("file.write", &FnEntry{
		Handler:     f.write,
		Name:        "Write file",
		Description: "Write content to a file",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "destination", Description: "The path of the file to write to", Mandatory: true},
			{Name: "value", Description: "The value to write", Mandatory: true},
		},
	})
}

type sourceFileParams struct {
	Source string `json:"source" yaml:"source"`
}

func (f *fnFile) size(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *sourceFileParams) (json.RawMessage, error) {
		fileInfo, err := os.Stat(params.Source)
		if err != nil {
			return nil, err
		}

		return utils.ReturnRaw(fileInfo.Size()), nil
	})
}

func (f *fnFile) modified(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *sourceFileParams) (json.RawMessage, error) {
		fileInfo, err := os.Stat(params.Source)
		if err != nil {
			return nil, err
		}

		return utils.ReturnRaw(fileInfo.ModTime().UnixMilli()), nil
	})
}

func (f *fnFile) delete(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *sourceFileParams) (json.RawMessage, error) {
		os.Remove(params.Source)

		return nil, nil
	})
}

type copyMoveParams struct {
	Source      string `json:"source" yaml:"source"`
	Destination string `json:"destination" yaml:"destination"`
}

func (f *fnFile) copy(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *copyMoveParams) (json.RawMessage, error) {
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

		return utils.ReturnRaw(params.Destination), err
	})
}

func (f *fnFile) move(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *copyMoveParams) (json.RawMessage, error) {
		err := os.Rename(params.Source, params.Destination)
		if err != nil {
			return nil, err
		}
		return utils.ReturnRaw(params.Destination), nil
	})
}

func (f *fnFile) read(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *sourceFileParams) (json.RawMessage, error) {
		file, err := os.ReadFile(params.Source)
		if err != nil {
			return nil, err
		}

		return utils.ReturnRaw(file), nil
	})
}

type writeFileParams struct {
	Destination string `json:"destination" yaml:"destination"`
	Value       string `json:"value" yaml:"value"`
}

func (f *fnFile) write(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *writeFileParams) (json.RawMessage, error) {
		err := os.WriteFile(params.Destination, []byte(params.Value), 0644)
		if err != nil {
			return nil, err
		}
		return utils.ReturnRaw(params.Destination), nil
	})
}
