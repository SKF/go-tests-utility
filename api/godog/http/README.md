# HTTP utilities for godog test suits
A simple abstraction on top of the go http client.

## Usage Example:
Feature:
```
Feature: It should be posible to create and retreive users from the service
    
    Scenario: If should be possible to create a user
        When a create user request with email "test-user@example.com", firstname "test-user" and lastname "example" is POSTED
        Then the returned status should be "200 OK"
        And the response header "Content-Type" should match "application/json"

    Scenario: If should be possible to get a user
        When a get user request with email "aladdin@example.com" is sent
        Then the returned status should be "200 OK"
        And the response header "Content-Type" should match "application/json"
```

Backing code:
```
import (
	"github.com/DATA-DOG/godog"
	"github.com/pkg/errors"
    "github.com/SKF/go-tests-utility/api/godog/json"
)

const (
    loginEmail = "aladdin@example.com"
    loginPassword = "simsalabim"
)

type state struct {
    client       *http.HttpClient
	requestURL    string
    requestBody   interface{}
	httpResponse *http.HttpResponse
	result        interface{}
}

type user struct {
    Email     string
    Firstname string
    Lastname  string
}

func (st *state) createUser(email, firstname, lastname string) error {
    st.requestURL = "https://localhost/users"
    st.requestBody = user{firstname, lastname, email} 
	resp, err := st.client.Post(st.requestURL, st.requestBody, nil)
	if err != nil {
		return err
	}
	st.httpResponse = resp
    return nil
}

func (st *state) getUser(email string) error {
    responseBody := struct {
        user user
    }{}

	resp, err := st.client.Get(st.requestURL, &responseBody)
	if err != nil {
		return err
	}

	st.httpResponse = resp
    st.result = responseBody.user
	return nil
}


func (st *state) httpStatusMatch(status string) error {
    if st.httpResponse == nil {
        return errors.Errorf("Expected a HTTP response")
    }

    if st.httpResponse.Status != status {
        return errors.Errorf("HTTP Status mismatch: expected %s got %s", status, st.httpResponse.Status)
    }

    return nil
}

func (st *state) httpHeaderMatch(header, pattern string) error {
    if st.httpResponse == nil {
        return errors.Errorf("Expected a HTTP response")
    }

    h, ok := st.httpResponse.Headers[header]
    if !ok {
        return errors.Errorf("Missing a %s HTTP header", header)
    }

    for _, v := range h {
        if v == pattern {
            return nil
        }
    }

    return errors.Errorf("No HTTP header value matched for %s (%+v)", header, h)
}

func FeatureContext(s *godog.Suite) {
    st := state{}
    st.client = http.New()

    if err := st.client.FetchToken("sandbox", loginEmail, loginPassword); err != nil {
        panic(err)
    }

    s.Step(`^a get user request with email "([^"]+)" is sent$`, st.getUser)
    s.Step(`^a create user request with email "([^"]+)", firstname "([^"]+)" and lastname "([^"]+)" is sent$`, st.createUser)
    s.Step(`^the returned status should be "([^"]+)"$`, st.httpStatusMatch)
    s.Step(`^the response header "([^"]+)" should match "([^"]+)"$`, st.httpHeaderMatch)
}

```
