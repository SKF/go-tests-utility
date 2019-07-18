# HTTP utilities for godog test suits
A simple abstraction on top of the go http client.

## Usage Example:
Feature:
```
Feature: It should be posible to create and retreive users from the service
    
    Scenario: If should be possible to create a user
        When given the following json dump:
        """
            {"key-0" : "value", "key-1" : {"key-key-1": ["value-0", "value-1"]}}
        """
        Then the path ".key-1.key-key-1[1]" should match "value-1"
```

Backing code:
```
import (
	"github.com/DATA-DOG/godog/gherkin"
    "github.com/SKF/go-tests-utility/api/godog/json"
)

import (
	"fmt"
	"strings"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

type state struct {
    dump []byte
}

func (st *state) setJsonDump(input *gherkin.DocString) error {
	fmt.Printf("DocString %+v", *input)
    st.dump = []byte(strings.TrimSpace(input.Content))
	return nil
}

func (st *state) matchJson(key, pattern string) error {
	return Match(st.dump, key, pattern)
}

func FeatureContext(s *godog.Suite) {
    st := state{}
	s.Step(`^given the following json dump:$`, st.setJsonDump)
	s.Step(`^the path "([^"]*)" should match "([^"]*)"$`, st.matchJson)
}
    
```
