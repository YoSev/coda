package coda

import (
	"encoding/json"
)

type OperationParameter struct {
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Mandatory   bool     `json:"mandatory" yaml:"mandatory"`
	Type        string   `json:"type" yaml:"type"`
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
}

// Get a list of all operations available in the Coda engine
func (c *Coda) GetOperations() map[string]*OperationHandler {
	return operations
}

var operations = map[string]*OperationHandler{
	"message.shoutrrr": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Shoutrrr(params)
		},
		Name:        "message.shoutrrr",
		Description: "Sends a message using the Shoutrrr notification system",
		Category:    OperationCategoryMessaging,
		Parameters: []OperationParameter{
			{Name: "urls", Description: "The shoutrrr targets", Type: "array", Mandatory: true},
			{Name: "message", Description: "The shoutrrr message to send", Mandatory: true},
			{Name: "parameters", Description: "Additional shoutrrr properties", Type: "object", Mandatory: false},
		},
	},
	"s3.upload": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.UploadToS3(params)
		},
		Name:        "s3.upload",
		Description: "Uploads a file to S3",
		Category:    OperationCategoryFile,
		Parameters: []OperationParameter{
			{Name: "endpoint", Description: "The S3 endpoint to use", Mandatory: true},
			{Name: "bucket", Description: "The S3 bucket to use", Mandatory: true},
			{Name: "region", Description: "The S3 region to use", Mandatory: true},
			{Name: "key_id", Description: "The S3 key ID to use", Mandatory: true},
			{Name: "key_secret", Description: "The S3 key secret to use", Mandatory: true},
			{Name: "local_path", Description: "The local path to upload", Mandatory: true},
			{Name: "remote_path", Description: "The remote path in the S3 bucket", Mandatory: false},
			{Name: "remote_prefix", Description: "The remote prefix in the S3 bucket (for recursive upload)", Mandatory: false},
			{Name: "invisible_files", Description: "If true, invisible files will be uploaded", Type: "boolean", Mandatory: false},
		},
	},
	"s3.download": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.DownloadFromS3(params)
		},
		Name:        "s3.download",
		Description: "Downloads a file from S3",
		Category:    OperationCategoryFile,
		Parameters: []OperationParameter{
			{Name: "endpoint", Description: "The S3 endpoint to use", Mandatory: true},
			{Name: "bucket", Description: "The S3 bucket to use", Mandatory: true},
			{Name: "region", Description: "The S3 region to use", Mandatory: true},
			{Name: "key_id", Description: "The S3 key ID to use", Mandatory: true},
			{Name: "key_secret", Description: "The S3 key secret to use", Mandatory: true},
			{Name: "local_path", Description: "The local path to download to", Mandatory: true},
			{Name: "remote_path", Description: "The remote path in the S3 bucket", Mandatory: false},
		},
	},
	"ai.openai": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.OpenAI(params)
		},
		Name:        "ai.openai",
		Description: "Performs an AI request",
		Category:    OperationCategoryAI,
		Parameters: []OperationParameter{
			{Name: "prompt", Description: "The actual prompt", Mandatory: true},
			{Name: "model", Description: "The modal to use", Mandatory: true},
			{Name: "api_key", Description: "The key to use", Mandatory: true},
			{Name: "system", Description: "The system query", Mandatory: false},
			{Name: "attachments", Description: "The attachments to include", Type: "array", Mandatory: false},
		},
	},
	"http.request": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.HttpReq(params)
		},
		Name:        "http.request",
		Description: "Performs an HTTP request",
		Category:    OperationCategoryHTTP,
		Parameters: []OperationParameter{
			{Name: "url", Description: "The url to query", Mandatory: true},
			{Name: "method", Description: "The HTTP method to use", Enum: []string{"GET", "POST", "PUT", "PATCH", "DELETE"}, Mandatory: true},
			{Name: "headers", Description: "The Headers to use", Type: "object", Mandatory: false},
			{Name: "body", Description: "The Body to use", Type: "any", Mandatory: false},
		},
	},
	"os.exec": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Exec(params)
		},
		Name:        "os.exec",
		Description: "Returns output of the command execution",
		Category:    OperationCategoryOS,
		Parameters: []OperationParameter{
			{Name: "command", Description: "The command to execute", Mandatory: true},
			{Name: "arguments", Description: "The arguments for the execution", Type: "array", Mandatory: true},
		},
	},
	"os.name": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.GetOS(params)
		},
		Name:        "os.name",
		Description: "Returns the name of the underlying operating system",
		Category:    OperationCategoryOS,
		Parameters:  []OperationParameter{},
	},
	"os.arch": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.GetArch(params)
		},
		Name:        "os.arch",
		Description: "Returns the name of the underlying operating architecture",
		Category:    OperationCategoryOS,
		Parameters:  []OperationParameter{},
	},
	"os.env.get": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.GetEnv(params)
		},
		Name:        "os.env.get",
		Description: "Returns the value of an environment variable",
		Category:    OperationCategoryOS,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The name of the environment variable", Mandatory: true},
		},
	},
	"io.stdout": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Stdout(params)
		},
		Name:        "io.stdout",
		Description: "Prints a value to stdout",
		Category:    OperationCategoryIO,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The value to print", Mandatory: true},
		},
	},
	"io.stderr": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Stderr(params)
		},
		Name:        "io.stderr",
		Description: "Prints a value to stderr",
		Category:    OperationCategoryIO,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The value to print", Mandatory: true},
		},
	},
	"string.upper": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.UpperCase(params)
		},
		Name:        "string.upper",
		Description: "Convert string to upper case",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to uppercase", Mandatory: true},
		},
	},
	"string.lower": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.LowerCase(params)
		},
		Name:        "string.lower",
		Description: "Convert string to lower case",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to lowercase", Mandatory: true},
		},
	},
	"string.camel": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.CamelCase(params)
		},
		Name:        "string.camel",
		Description: "Convert string to camel case",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to camelcase", Mandatory: true},
		},
	},
	"string.snake": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.SnakeCase(params)
		},
		Name:        "string.snake",
		Description: "Convert string to snake case",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to snakecase", Mandatory: true},
		},
	},
	"string.kebap": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.KebapCase(params)
		},
		Name:        "string.kebap",
		Description: "Convert string to kebap case",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to kebapcase", Mandatory: true},
		},
	},
	"string.trim": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.StringTrim(params)
		},
		Name:        "string.trim",
		Description: "Trim a string by delimiter (default: whitespace)",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to trim", Mandatory: true},
			{Name: "delimiter", Description: "The delimiter to trim", Mandatory: false},
		},
	},
	"string.reverse": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.StringReverse(params)
		},
		Name:        "string.reverse",
		Description: "Reverse a string",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to reverse", Mandatory: true},
		},
	},
	"string.split": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.StringSplit(params)
		},
		Name:        "string.split",
		Description: "Split a string by delimiter",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to split", Mandatory: true},
			{Name: "delimiter", Description: "The delimiter to use", Mandatory: false},
		},
	},
	"string.join": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.StringJoin(params)
		},
		Name:        "string.join",
		Description: "Join a string by delimiter",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to join", Mandatory: true},
			{Name: "delimiter", Description: "The delimiter to use", Mandatory: false},
		},
	},
	"string.resolve": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.StringResolve(params)
		},
		Name:        "string.resolve",
		Description: "Return a (resolved) string",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to resolve", Mandatory: true},
		},
	},
	"json.decode": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.JsonDecode(params)
		},
		Name:        "json.decode",
		Description: "Decode an object to a json string",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The object to decode", Mandatory: true},
		},
	},
	"json.encode": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.JsonEncode(params)
		},
		Name:        "json.encode",
		Description: "Encode a json string",
		Category:    OperationCategoryString,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to encode", Mandatory: true},
		},
	},
	"file.size": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Size(params)
		},
		Name:        "file.size",
		Description: "Get the size of a file",
		Category:    OperationCategoryFile,
		Parameters: []OperationParameter{
			{Name: "source", Description: "The path of the source", Mandatory: true},
		},
	},
	"file.modified": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Modified(params)
		},
		Name:        "file.modified",
		Description: "Get the modify date of a file as unix timestamp in milliseconds",
		Category:    OperationCategoryFile,
		Parameters: []OperationParameter{
			{Name: "source", Description: "The path of the source", Mandatory: true},
		},
	},
	"file.copy": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Copy(params)
		},
		Name:        "file.copy",
		Description: "Copy a file from source to destination",
		Category:    OperationCategoryFile,
		Parameters: []OperationParameter{
			{Name: "source", Description: "The path of the source", Mandatory: true},
			{Name: "destination", Description: "The path of the destination", Mandatory: true},
		},
	},
	"file.delete": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Delete(params)
		},
		Name:        "file.delete",
		Description: "Delete a file",
		Category:    OperationCategoryFile,
		Parameters: []OperationParameter{
			{Name: "source", Description: "The path of the file to delete", Mandatory: true},
		},
	},
	"file.move": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Move(params)
		},
		Name:        "file.move",
		Description: "Move a file",
		Parameters: []OperationParameter{
			{Name: "source", Description: "The path of the file to move", Mandatory: true},
			{Name: "destination", Description: "The path of the destination", Mandatory: true},
		},
	},
	"file.read": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Read(params)
		},
		Name:        "file.read",
		Description: "Read a files content",
		Category:    OperationCategoryFile,
		Parameters: []OperationParameter{
			{Name: "source", Description: "The path of the file to read", Mandatory: true},
		},
	},
	"file.write": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Write(params)
		},
		Name:        "file.write",
		Description: "Write content to a file",
		Category:    OperationCategoryFile,
		Parameters: []OperationParameter{
			{Name: "destination", Description: "The path of the file to write to", Mandatory: true},
			{Name: "value", Description: "The value to write", Mandatory: true},
		},
	},
	"time.datetime": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.GenerateDatetime(params)
		},
		Name:        "time.datetime",
		Description: "Generate a date time string",
		Category:    OperationCategoryTime,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The format to return", Mandatory: true},
		},
	},
	"time.timestamp.sec": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.GenerateTimestampSec(params)
		},
		Name:        "time.timestamp.sec",
		Description: "Generate a timestamp in seconds",
		Category:    OperationCategoryTime,
		Parameters:  []OperationParameter{},
	},
	"time.timestamp.milli": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.GenerateTimestampMilli(params)
		},
		Name:        "time.timestamp.milli",
		Description: "Generate a timestamp in milliseconds",
		Category:    OperationCategoryTime,
		Parameters:  []OperationParameter{},
	},
	"time.timestamp.micro": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.GenerateTimestampMicro(params)
		},
		Name:        "time.timestamp.micro",
		Description: "Generate a timestamp in microseconds",
		Category:    OperationCategoryTime,
		Parameters:  []OperationParameter{},
	},
	"time.timestamp.nano": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.GenerateTimestampNano(params)
		},
		Name:        "time.timestamp.nano",
		Description: "Generate a timestamp in nanoseconds",
		Category:    OperationCategoryTime,
		Parameters:  []OperationParameter{},
	},
	"time.sleep": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Sleep(params)
		},
		Name:        "time.sleep",
		Description: "Sleep for a given time in milliseconds",
		Category:    OperationCategoryTime,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The amount of milliseconds to sleep", Type: "integer", Mandatory: true},
		},
	},
	"hash.md5": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.HashMD5(params)
		},
		Name:        "hash.md5",
		Description: "Encode to MD5",
		Category:    OperationCategoryTime,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to hash", Mandatory: true},
		},
	},
	"hash.sha1": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.HashSha256(params)
		},
		Name:        "hash.sha1",
		Description: "Encode to Sha1",
		Category:    OperationCategoryTime,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to hash", Mandatory: true},
		},
	},
	"hash.sha256": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.HashSha256(params)
		},
		Name:        "hash.sha256",
		Description: "Encode to Sha256",
		Category:    OperationCategoryTime,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to hash", Mandatory: true},
		},
	},
	"hash.sha512": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.HashSha512(params)
		},
		Name:        "hash.sha512",
		Description: "Encode to Sha512",
		Category:    OperationCategoryTime,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to hash", Mandatory: true},
		},
	},
	"hash.base64.encode": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Base64Enc(params)
		},
		Name:        "hash.base64.encode",
		Description: "Encode to Base64",
		Category:    OperationCategoryTime,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to encode", Mandatory: true},
		},
	},
	"hash.base64.decode": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.Base64Dec(params)
		},
		Name:        "hash.base64.decode",
		Description: "Decode from Base64",
		Category:    OperationCategoryTime,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The string to decode", Mandatory: true},
		},
	},
	"math.inc": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.MathInc(params)
		},
		Name:        "math.inc",
		Description: "Increment a number",
		Category:    OperationCategoryMath,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The value to increment", Type: "number", Mandatory: true},
			{Name: "amount", Description: "The amount by which to increment (defaults to 1)", Type: "number", Mandatory: false},
		},
	},
	"math.dec": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.MathDec(params)
		},
		Name:        "math.dec",
		Description: "Decrement a number",
		Category:    OperationCategoryMath,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The value to decrement", Type: "number", Mandatory: true},
			{Name: "amount", Description: "The amount by which to decrement (defaults to 1)", Type: "number", Mandatory: false},
		},
	},
	"math.multiply": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.MathMultiply(params)
		},
		Name:        "math.multiply",
		Description: "Multiple a number",
		Category:    OperationCategoryMath,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The value to multiply", Type: "number", Mandatory: true},
			{Name: "amount", Description: "The amount by which to multiply (defaults to 1)", Type: "number", Mandatory: false},
		},
	},
	"math.divide": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.MathDivide(params)
		},
		Name:        "math.divide",
		Description: "Divide a number",
		Category:    OperationCategoryMath,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The value to divide", Type: "number", Mandatory: true},
			{Name: "amount", Description: "The amount by which to divide (defaults to 1)", Type: "number", Mandatory: false},
		},
	},
	"math.modulo": {
		Fn: func(c *Coda, params json.RawMessage) (json.RawMessage, error) {
			return c.fn.MathModulo(params)
		},
		Name:        "math.modulo",
		Description: "Get the floating-point remainder of value/amount",
		Category:    OperationCategoryMath,
		Parameters: []OperationParameter{
			{Name: "value", Description: "The source value", Type: "number", Mandatory: true},
			{Name: "amount", Description: "The amount by which to mod", Type: "number", Mandatory: false},
		},
	},
}

type OperationHandler struct {
	Fn          func(*Coda, json.RawMessage) (json.RawMessage, error)
	Name        string
	Description string
	Category    OperationCategory
	Parameters  []OperationParameter
}

type OperationCategory string

const (
	OperationCategoryFile      OperationCategory = "File"
	OperationCategoryString    OperationCategory = "String"
	OperationCategoryTime      OperationCategory = "Time"
	OperationCategoryIO        OperationCategory = "I/O"
	OperationCategoryMessaging OperationCategory = "Messaging"
	OperationCategoryOS        OperationCategory = "OS"
	OperationCategoryHTTP      OperationCategory = "HTTP"
	OperationCategoryHash      OperationCategory = "Hash"
	OperationCategoryMath      OperationCategory = "Math"
	OperationCategoryAI        OperationCategory = "AI"
)
