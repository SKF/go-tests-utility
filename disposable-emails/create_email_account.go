package disposableemail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const disposableEmailsBaseURL = "https://api.disposable-emails.enlight.skf.com"

func NewEmailAddress() (_ string, err error) {
	const url = disposableEmailsBaseURL + "/email-addresses/new"

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("Wrong status: %q", resp.Status)
		return
	}

	var respBody struct {
		Data struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return
	}

	return respBody.Data.EmailAddress, nil
}

func PollForMessageWithSubject(emailAddress, subject string, fromTimestamp time.Time) (string, error) {
	time.Sleep(2 * time.Second) //nolint: gomnd

	const timeOut = 5
	for n := 0; n < timeOut; n++ {
		time.Sleep(time.Second)

		messages, err := getAllMessages(emailAddress, fromTimestamp)
		if err != nil {
			continue
		}

		message, err := findMessageWithSubject(messages, subject)
		if err != nil {
			continue
		}

		return extractHTMLPart(message)
	}

	return "", errors.Errorf(
		"poll for message with subject [%s] to email [%s] timed out after %d seconds, search interval: [%s] - [%s]",
		subject, emailAddress, timeOut, fromTimestamp, time.Now().String(),
	)
}

func getAllMessages(emailAddress string, fromTimestamp time.Time) (messages []mail.Message, err error) {
	url := fmt.Sprintf(
		disposableEmailsBaseURL+"/email-addresses/%s/messages",
		emailAddress,
	)

	resp, err := http.Get(url) //nolint: gosec
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("Wrong status: %q", resp.Status)
		return
	}

	var respBody struct {
		Data [][]byte `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return
	}

	for _, rawMsg := range respBody.Data {
		r := bytes.NewReader(rawMsg)
		msg, innerErr := mail.ReadMessage(r)
		if innerErr != nil {
			continue
		}
		const dateFormat = "Mon, 2 Jan 2006 15:04:05 -0700"
		date, innerErr := time.Parse(dateFormat, msg.Header.Get("Date"))
		if innerErr != nil {
			continue
		}

		start := fromTimestamp.Truncate(time.Second)
		if msg.Header.Get("To") == emailAddress &&
			start.Before(date) || start.Equal(date) {
			messages = append(messages, *msg)
		}
	}

	return messages, err
}

func findMessageWithSubject(messages []mail.Message, subject string) (msg mail.Message, err error) {
	for _, msg = range messages {
		if msg.Header.Get("Subject") == subject {
			return
		}
	}

	err = errors.Errorf("Couldn't find message with subject: %q", subject)
	return
}

func extractHTMLPart(msg mail.Message) (string, error) {
	const htmlContentType = "text/html; charset=UTF-8"

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return "", errors.Wrap(err, "mime.ParseMediaType failed")
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				return "", nil
			}
			if err != nil {
				return "", errors.Wrap(err, "mr.NextPart failed")
			}
			if part.Header.Get("Content-Type") == htmlContentType {
				return parsePart(part)
			}
		}
	}

	return "", errors.Errorf("couldn't find part of type: %q", htmlContentType)
}

func parsePart(part *multipart.Part) (string, error) {
	defer part.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(part); err != nil {
		return "", errors.Wrap(err, "buf.ReadFrom failed")
	}

	re := regexp.MustCompile(`\r?\n`)
	bodyPart := re.ReplaceAllString(buf.String(), "\n")

	return strings.TrimSpace(bodyPart), nil
}
