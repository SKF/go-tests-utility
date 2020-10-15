package json

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
)

func MatchNull(json []byte, path string) error {
	if bytes.Equal(bytes.TrimSpace(json), []byte("null")) && (path == "" || path == ".") {
		return nil
	}

	path = revertLegacySyntax(path)

	result := gjson.Get(string(json), path)

	if result.Exists() && !result.IsArray() && !result.IsObject() && result.Value() == nil {
		return nil
	}

	return errors.Errorf("Match error: Expected null got '%T' JSON: %s", result.String(), string(json))
}

func Match(json []byte, path string, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Wrap(err, "Failed to compile regexp")
	}

	path = revertLegacySyntax(path)

	result := gjson.Get(string(json), path)

	if !re.MatchString(result.String()) {
		return errors.Errorf("Match error: Values mismatch, pattern: '%s' value: '%s' JSON: %s", pattern, result.String(), string(json))
	}

	return nil
}

func revertLegacySyntax(path string) string {
	// Remove brackets earlier they were used for array indexing
	bracketRegexp := regexp.MustCompile(`\.?\[(\d+)]`)
	path = bracketRegexp.ReplaceAllString(path, ".$1")

	if strings.HasPrefix(path, ".") {
		path = string([]rune(path)[1:])
	}

	return path
}

func ArrayLen(json []byte, path string, length int) error {
	path = revertLegacySyntax(path)

	res := gjson.Get(string(json), path)


	if !res.IsArray() {
		return errors.Errorf("Match error: Expected an array got: %s, JSON: %s", res.String(), string(json))
	}

	if len(res.Array()) != length {
		return errors.Errorf("Match error: Expected an array of length: %d got: %d JSON: %s", length, len(res.Array()), string(json))
	}

	return nil
}

func Read(json []byte, path string) (result string, err error) {
	path = revertLegacySyntax(path)

	res := gjson.Get(string(json), path)

	if !res.Exists() || res.IsObject() || res.IsArray() {
		return "", errors.Errorf("Match error: Expected a scalar got '%T' JSON: %s", res.String(), string(json))
	}

	return res.String(), nil
}

func ReadStringArr(json []byte, path string) (result []string, err error) {
	path = revertLegacySyntax(path)

	res := gjson.Get(string(json), path)



	if !res.IsArray() {
		return []string{}, errors.Errorf("Match error: Expected an array got: %s", res.String())
	}

	for _, v := range res.Array() {
		stringValue := v.String()
		if !(v.IsObject() || !v.IsArray()) {
			return []string{}, errors.Errorf("Match error: Expected a scalar got: %s", v.String())
		}
		result = append(result, stringValue)
	}

	return result, nil
}
