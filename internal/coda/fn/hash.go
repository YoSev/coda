package fn

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type hashParams struct {
	Value string `json:"value" yaml:"value"`
}

func (f *Fn) HashMD5(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		hash := md5.Sum([]byte(params.Value))
		return returnAny(fmt.Sprintf("%x", hash)), nil
	})
}

func (f *Fn) HashSha1(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		hash := sha1.New()
		hash.Write([]byte(params.Value))
		hashBytes := hash.Sum(nil)
		return returnAny(fmt.Sprintf("%x", hashBytes)), nil
	})
}

func (f *Fn) HashSha256(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		hash := sha256.New()
		hash.Write([]byte(params.Value))
		hashBytes := hash.Sum(nil)
		return returnAny(fmt.Sprintf("%x", hashBytes)), nil
	})
}

func (f *Fn) HashSha512(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		hash := sha512.New()
		hash.Write([]byte(params.Value))
		hashBytes := hash.Sum(nil)
		return returnAny(fmt.Sprintf("%x", hashBytes)), nil
	})
}

func (f *Fn) Base64Enc(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		str := base64.StdEncoding.EncodeToString([]byte(params.Value))
		return returnAny(str), nil
	})
}

func (f *Fn) Base64Dec(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *hashParams) (json.RawMessage, error) {
		data, err := base64.StdEncoding.DecodeString(params.Value)
		if err != nil {
			return nil, err
		}
		return returnAny(data), nil
	})
}
