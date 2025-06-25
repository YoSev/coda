package fn

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/yosev/coda/internal/utils"
)

type fnHash struct {
	category FnCategory
}

func (f *fnHash) init(fn *Fn) {
	fn.register("hash.md5", &FnEntry{
		Handler:     f.md5,
		Name:        "MD5 Hash",
		Description: "Calculate the MD5 hash of a string",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to hash", Mandatory: true},
		},
	})
	fn.register("hash.sha1", &FnEntry{
		Handler:     f.sha1,
		Name:        "SHA1 Hash",
		Description: "Calculate the SHA1 hash of a string",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to hash", Mandatory: true},
		},
	})
	fn.register("hash.sha256", &FnEntry{
		Handler:     f.sha256,
		Name:        "SHA256 Hash",
		Description: "Calculate the SHA256 hash of a string",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to hash", Mandatory: true},
		},
	})
	fn.register("hash.sha512", &FnEntry{
		Handler:     f.sha512,
		Name:        "SHA512 Hash",
		Description: "Calculate the SHA512 hash of a string",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to hash", Mandatory: true},
		},
	})
	fn.register("hash.base64.encode", &FnEntry{
		Handler:     f.b64enc,
		Name:        "Base64 Encode",
		Description: "Encode to Base64",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to encode", Mandatory: true},
		},
	})
	fn.register("hash.base64.decode", &FnEntry{
		Handler:     f.b64dec,
		Name:        "Base64 Decode",
		Description: "Decode from Base64",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to decode", Mandatory: true},
		},
	})
}

type hashParams struct {
	Value string `json:"value" yaml:"value"`
}

func (f *fnHash) md5(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		hash := md5.Sum([]byte(params.Value))
		return utils.ReturnRaw(fmt.Sprintf("%x", hash)), nil
	})
}

func (f *fnHash) sha1(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		hash := sha1.New()
		hash.Write([]byte(params.Value))
		hashBytes := hash.Sum(nil)
		return utils.ReturnRaw(fmt.Sprintf("%x", hashBytes)), nil
	})
}

func (f *fnHash) sha256(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		hash := sha256.New()
		hash.Write([]byte(params.Value))
		hashBytes := hash.Sum(nil)
		return utils.ReturnRaw(fmt.Sprintf("%x", hashBytes)), nil
	})
}

func (f *fnHash) sha512(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		hash := sha512.New()
		hash.Write([]byte(params.Value))
		hashBytes := hash.Sum(nil)
		return utils.ReturnRaw(fmt.Sprintf("%x", hashBytes)), nil
	})
}

func (f *fnHash) b64enc(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		str := base64.StdEncoding.EncodeToString([]byte(params.Value))
		return utils.ReturnRaw(str), nil
	})
}

func (f *fnHash) b64dec(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		data, err := base64.StdEncoding.DecodeString(params.Value)
		if err != nil {
			return nil, err
		}
		return utils.ReturnRaw(data), nil
	})
}
