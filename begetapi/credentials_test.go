package begetapi

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type credentialTestTransport func(*http.Request) (*http.Response, error)

func (f credentialTestTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestErrorsDoNotExposeCredentials(t *testing.T) {
	credentials := Credentials{Login: "private-login", Passwd: "private+password&"}
	transportErr := errors.New("connection failed")
	for _, method := range []string{"getData", "changeRecords"} {
		for _, scenario := range []string{"transport", "http", "api"} {
			if method == "getData" && scenario == "api" {
				continue
			}
			t.Run(method+"/"+scenario, func(t *testing.T) {
				endpoint, _ := url.Parse("https://example.invalid")
				client := NewApiClient(endpoint)
				client.client.Transport = credentialTestTransport(func(r *http.Request) (*http.Response, error) {
					if r.URL.Query().Get("passwd") != credentials.Passwd {
						t.Fatal("missing request credentials")
					}
					if scenario == "transport" {
						return nil, transportErr
					}
					status := http.StatusOK
					body := `{"answer":{"result":false,"status":"` + credentials.Passwd + `"}}`
					if scenario == "http" {
						status = http.StatusForbidden
						body = r.URL.String()
					}
					return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
				})
				var err error
				if method == "getData" {
					_, err = client.GetData("example.com", credentials)
				} else {
					err = client.ChangeRecords("example.com", Records{}, credentials)
				}
				if err == nil {
					t.Fatal("expected failure")
				}
				for _, secret := range []string{credentials.Login, credentials.Passwd, url.QueryEscape(credentials.Passwd)} {
					if strings.Contains(err.Error(), secret) {
						t.Fatal("error exposes credentials")
					}
				}
				if scenario == "transport" && !errors.Is(err, transportErr) {
					t.Fatal("transport cause lost")
				}
			})
		}
	}
}
