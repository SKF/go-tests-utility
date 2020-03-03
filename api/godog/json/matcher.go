package json

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/pkg/errors"
)

type term struct {
	jsonType int
	value    interface{}
}

type segment struct {
	key     string
	index   int64
	indexed bool
}

const (
	jsonNumber = iota
	jsonString = iota
	jsonBool   = iota
	jsonNull   = iota
	jsonArray  = iota
	jsonObject = iota
)

var null = term{jsonType: jsonNull, value: nil}

func MatchNull(json []byte, path string) error {
	t, err := resolve(json, path)
	if err != nil {
		return err
	}

	if t.value != nil {
		return errors.Errorf("Match error: Expected null got '%T' JSON: %s", t.value, string(json))
	}

	return nil
}

func Match(json []byte, path string, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Wrap(err, "Failed to compile regexp")
	}

	t, err := resolve(json, path)
	if err != nil {
		return err
	}

	value, ok := t.value.(string)
	if !ok {
		return errors.Errorf("Match error: Expected a scalar got '%T' JSON: %s", t.value, string(json))
	}

	if !re.MatchString(value) {
		return errors.Errorf("Match error: Values mismatch, pattern: '%s' value: '%s' JSON: %s", pattern, value, string(json))
	}

	return nil
}

func ArrayLen(json []byte, path string, length int) error {
	t, err := resolve(json, path)
	if err != nil {
		return err
	}

	value, ok := t.value.([]term)
	if !ok {
		return errors.Errorf("Match error: Expected an array got '%T' JSON: %s", t.value, string(json))
	}

	if len(value) != length {
		return errors.Errorf("Match error: Expected an array of length: %d got: %d JSON: %s", length, len(value), string(json))
	}
	return nil
}

func Read(json []byte, path string) (result string, err error) {
	t, err := resolve(json, path)
	if err != nil {
		return "", err
	}

	value, ok := t.value.(string)
	if !ok {
		return "", errors.Errorf("Match error: Expected a scalar got '%T' JSON: %s", t.value, string(json))
	}

	return value, nil
}

func resolve(json []byte, path string) (term, error) {
	t, err := parse(json)
	if err != nil {
		return t, err
	}

	segments, err := parseMatchSegments(path)
	if err != nil {
		return t, err
	}

	for _, seg := range segments {
		obj, ok := t.value.(map[string]term)
		if !ok {
			return t, errors.Errorf("Match error: Expected an object got '%T' JSON: %s", t.value, string(json))
		}

		t, ok = obj[seg.key]
		if !ok {
			keys := make([]string, 0)
			for key := range obj {
				keys = append(keys, key)
			}
			return t, errors.Errorf("Match error: Missing key '%s' in map, possible keys are: '%s' JSON: %s", seg.key, strings.Join(keys, ", "), string(json))
		}

		if seg.indexed {
			arr, ok := t.value.([]term)
			if !ok {
				return t, errors.Errorf("Match error: Expected an array, got: %T JSON: %s", t.value, string(json))
			}
			if len(arr) <= int(seg.index) {
				return t, errors.Errorf("Match error: Array out ouf bounds: %+v (%d) JSON: %s", arr, seg.index, string(json))
			}
			t = arr[seg.index]
		}
	}

	return t, nil
}

func parseMatchSegments(path string) ([]segment, error) {
	re := regexp.MustCompile(`([^["]+)(\[(\d+)\])?`)
	matchPath := strings.Split(strings.Trim(path, "."), ".")

	segments := make([]segment, 0)
	for _, p := range matchPath {
		ms := re.FindStringSubmatch(p)
		if ms[1] == "" {
			return segments, errors.Errorf("Parse error: Failed to parse match segment: %s", p)
		}
		if ms[3] != "" {
			idx, err := strconv.ParseInt(ms[3], 10, 64)
			if err != nil {
				return segments, errors.Wrap(err, "Parse error: Failed to parse array index")
			}
			segments = append(segments, segment{key: ms[1], index: idx, indexed: true})
		} else {
			segments = append(segments, segment{key: ms[1], index: -1, indexed: false})
		}
	}
	return segments, nil
}

func parse(json []byte) (term, error) {
	var lex scanner.Scanner
	lex.Init(bytes.NewReader(json))
	lex.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings

	t, err := parseTerm(&lex)
	if err != nil {
		return null, err
	}

	_, ok := t.value.(map[string]term)
	if !ok {
		return null, errors.Errorf("Match error: Expected a single toplevel object JSON: `%s`", string(json))
	}

	tok := lex.Scan()
	if tok != scanner.EOF {
		return null, errors.Errorf("Match error: Expected a single toplevel object JSON (%v): `%s`", tok, string(json))
	}

	return t, nil
}

func parseTerm(lex *scanner.Scanner) (term, error) {
	switch lex.Scan() {
	case '{':
		return parseObject(lex)
	case '[':
		return parseArray(lex)
	case '-':
		return parseNegativeNumber(lex)
	case scanner.Ident:
		return encodeIdent(lex)
	case scanner.String:
		return encodeString(lex)
	case scanner.Int, scanner.Float:
		return encodeNumber(lex)
	case scanner.EOF:
		return null, errors.Errorf("Parse error: Unexpected EOF")
	default:
		return null, errors.Errorf("Parse error: Unexpected token: '%c'", lex.Peek())
	}
}

func encodeIdent(lex *scanner.Scanner) (term, error) {
	text := lex.TokenText()

	switch text {
	case "null":
		return null, nil
	case "true", "false":
		return term{jsonType: jsonBool, value: text}, nil
	}

	return null, errors.Errorf("Parse error: Unexpected token: %s", text)
}

func encodeString(lex *scanner.Scanner) (term, error) {
	str, err := strconv.Unquote(lex.TokenText())

	if err != nil {
		return null, err
	}

	return term{jsonType: jsonString, value: str}, nil
}

func encodeNumber(lex *scanner.Scanner) (term, error) {
	text := lex.TokenText()
	return term{jsonType: jsonNumber, value: text}, nil
}

func parseObject(lex *scanner.Scanner) (term, error) {
	obj := make(map[string]term)

	for tok := lex.Scan(); tok != scanner.EOF && tok != '}'; tok = lex.Scan() {
		if tok != scanner.String {
			return null, errors.Errorf("Parse error: Expected a key got: %s", lex.TokenText())
		}

		key, err := strconv.Unquote(lex.TokenText())
		if err != nil {
			return null, err
		}

		tok = lex.Scan()
		if tok != ':' {
			return null, errors.Errorf("Parse error: Expected separator: %s", lex.TokenText())
		}

		val, err := parseTerm(lex)
		if err != nil {
			return null, err
		}

		if lex.Peek() == ',' {
			lex.Scan()
		}

		obj[key] = val
	}

	return term{jsonType: jsonObject, value: obj}, nil
}

func parseArray(lex *scanner.Scanner) (term, error) {
	arr := make([]term, 0)

	for lex.Peek() != ']' {
		val, err := parseTerm(lex)
		if err != nil {
			return null, err
		}

		arr = append(arr, val)

		if lex.Peek() == ',' {
			lex.Scan()
		}
	}
	lex.Scan() // Drop ']'

	return term{jsonType: jsonArray, value: arr}, nil
}

func parseNegativeNumber(lex *scanner.Scanner) (term, error) {
	tok := lex.Scan()
	if tok == scanner.Int || tok == scanner.Float {
		text := lex.TokenText()
		return term{jsonType: jsonNumber, value: "-" + text}, nil
	}
	return null, errors.Errorf("Parse error: Expected a number got: '%s'", lex.TokenText())
}
