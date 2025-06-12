package coda

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

func (c *Coda) resolveVariables(in json.RawMessage) (json.RawMessage, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(in) == 0 {
		return in, nil // No input to resolve
	}
	defer func() {
		c.Stats.VariablesTotal++
	}()

	// Marshal `c` so we can use gjson to query it
	codaJSON, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Coda struct: %w", err)
	}

	// Unmarshal the input into an interface{}
	var input any
	if err := json.Unmarshal(in, &input); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input: %w", err)
	}

	// Recursively resolve variables in the data
	resolved := resolveValue(input, codaJSON)

	// Re-marshal to JSON
	out, err := json.Marshal(resolved)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resolved result: %w", err)
	}

	c.Stats.VariablesSuccessfulTotal++
	return json.RawMessage(out), nil
}

func resolveValue(val any, codaJSON []byte) any {
	switch v := val.(type) {
	case map[string]any:
		for key, value := range v {
			v[key] = resolveValue(value, codaJSON)
		}
		return v
	case []any:
		for i, item := range v {
			v[i] = resolveValue(item, codaJSON)
		}
		return v
	case string:
		return resolveString(v, codaJSON)
	default:
		return val
	}
}

func resolveString(input string, codaJSON []byte) any {
	// Updated regex: capture variable path and any filter string (starting with a pipe) until "}"
	re := regexp.MustCompile(`\${\s*([^}\|]+?)\s*((?:\|[^}]+)+)?\s*}`)

	matches := re.FindAllStringSubmatch(input, -1)

	if len(matches) == 1 && strings.TrimSpace(input) == matches[0][0] {
		path := matches[0][1] // variable path trimmed
		filters := parseFilters(matches[0][2])
		val := gjson.GetBytes(codaJSON, path)
		return applyFilterChain(val, filters)
	}

	result := input
	for _, match := range matches {
		full := match[0]
		path := match[1]
		filters := parseFilters(match[2])

		val := gjson.GetBytes(codaJSON, path)
		replacement := fmt.Sprintf("%v", applyFilterChain(val, filters))
		result = strings.ReplaceAll(result, full, replacement)
	}

	return result
}

func parseFilters(raw string) []Filter {
	var filters []Filter
	if raw == "" {
		return filters
	}
	// Remove the leading pipe if present.
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "|") {
		raw = raw[1:]
	}
	// Split on pipes and trim each part.
	parts := strings.Split(raw, "|")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// Split filter name and optional argument (e.g., "join" or "join:/")
		sub := strings.SplitN(part, ":", 2)
		f := Filter{Name: strings.TrimSpace(sub[0])}
		if len(sub) == 2 {
			f.Arg = strings.TrimSpace(sub[1])
		}
		filters = append(filters, f)
	}
	return filters
}

type Filter struct {
	Name string
	Arg  string
}

func applyFilterChain(val gjson.Result, filters []Filter) any {
	current := parseRaw(val)

	for _, filter := range filters {
		current = applySingleFilter(current, filter)
	}
	return current
}

func applySingleFilter(val any, filter Filter) any {
	switch filter.Name {
	case "string":
		return fmt.Sprintf("%v", val)
	case "join":
		del := ","
		if filter.Arg != "" {
			del = filter.Arg
		}
		// Use reflection to support slices of any type
		rVal := reflect.ValueOf(val)
		if rVal.Kind() != reflect.Slice {
			fmt.Println("join: value is not a slice:", val)
			return val
		}
		parts := make([]string, rVal.Len())
		for i := 0; i < rVal.Len(); i++ {
			parts[i] = fmt.Sprintf("%v", rVal.Index(i).Interface())
		}
		return strings.Join(parts, del)
	case "upper":
		if s, ok := val.(string); ok {
			return strings.ToUpper(s)
		}
	case "lower":
		if s, ok := val.(string); ok {
			return strings.ToLower(s)
		}
	case "trim":
		if s, ok := val.(string); ok {
			return strings.TrimSpace(s)
		}
	case "split":
		if s, ok := val.(string); ok {
			del := "."
			if filter.Arg != "" {
				del = filter.Arg
			}
			return strings.Split(s, del)
		}
	case "md5":
		if s, ok := val.(string); ok {
			hash := md5.Sum([]byte(s))
			return fmt.Sprintf("%x", hash)
		}
	case "sha1":
		if s, ok := val.(string); ok {
			hash := sha1.New()
			hash.Write([]byte(s))
			hashBytes := hash.Sum(nil)
			return fmt.Sprintf("%x", hashBytes)
		}
	case "sha256":
		if s, ok := val.(string); ok {
			hash := sha256.New()
			hash.Write([]byte(s))
			hashBytes := hash.Sum(nil)
			return fmt.Sprintf("%x", hashBytes)
		}
	case "sha512":
		if s, ok := val.(string); ok {
			hash := sha512.New()
			hash.Write([]byte(s))
			hashBytes := hash.Sum(nil)
			return fmt.Sprintf("%x", hashBytes)
		}
	case "jsonDecode":
		if s, ok := val.(string); ok {
			var b = make(map[string]interface{})
			err := json.Unmarshal([]byte(s), &b)
			if err != nil {
				return s
			}
			return b
		}
	case "jsonEncode":
		if s, ok := val.(any); ok {
			b, err := json.Marshal(s)
			if err != nil {
				return s
			}
			return string(b)
		}
	case "base64Decode":
		if s, ok := val.(string); ok {
			b, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				return s
			}
			return string(b)
		}
	case "base64Encode":
		if s, ok := val.(string); ok {
			return base64.StdEncoding.EncodeToString([]byte(s))
		}
	}
	return val
}

func parseRaw(v gjson.Result) any {
	if v.Type == gjson.String {
		return v.String()
	}
	var out any
	err := json.Unmarshal([]byte(v.Raw), &out)
	if err != nil {
		return v.String()
	}
	return out
}
