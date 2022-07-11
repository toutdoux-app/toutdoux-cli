package v1

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

const (
	APIVersion = "v1"

	LoginEndpoint  = "/auth/new"
	LogoutEndpoint = "/auth/logout"

	csrfTokenHeaderName = "X-Csrf-Token"
)

type Client interface {
	Initialize() error
	GetTodoLists(opts GetTodoListsOptions) (TodoLists, error)
	GetTodoListByID(todoListID string, opts GetTodoListByIDOptions) (TodoList, error)
	GetTodoListEntriesByID(todoListID string, opts GetTodoListEntriesByIDOptions) (TodoEntries, error)
}

type ClientOptions struct {
	Endpoint string
	Username string
	Password string
}

type httpDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type client struct {
	opts ClientOptions

	httpClient httpDoer
}

func NewClientWithOptions(opts ClientOptions) (*client, error) {
	c := &client{
		opts: opts,
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, errors.Wrap(err, "initializing cookie jar")
	}

	c.httpClient = &http.Client{
		Timeout: 60 * time.Second,
		Jar:     jar,
	}

	return c, nil
}

func (c *client) Initialize() error {
	if err := c.loginWithPassword(); err != nil {
		return errors.Wrap(err, "logging in with password")
	}

	return nil
}

func (c *client) apiURLPathJoin(path string) string {
	return c.opts.Endpoint + "/api/" + APIVersion + path
}

func (c *client) loginWithPassword() error {
	csrfToken, err := c.getCSRFToken()
	if err != nil {
		return errors.Wrap(err, "getting CSRF token")
	}

	formData := url.Values{}
	formData.Set("Email", c.opts.Username)
	formData.Set("Password", c.opts.Password)

	loginReq, err := http.NewRequest(http.MethodPost, c.opts.Endpoint+LoginEndpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return errors.Wrap(err, "creating login request")
	}

	loginReq.Header.Add(csrfTokenHeaderName, csrfToken)
	loginReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(loginReq)
	if err != nil {
		return errors.Wrap(err, "performing login request")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status code %d", resp.StatusCode)
	}

	return nil
}

func (c *client) getCSRFToken() (string, error) {
	csrfRequest, err := http.NewRequest(http.MethodGet, c.opts.Endpoint+LoginEndpoint, nil)
	if err != nil {
		return "", errors.Wrap(err, "creating HTTP request")
	}

	resp, err := c.httpClient.Do(csrfRequest)
	if err != nil {
		return "", errors.Wrap(err, "performing HTTP request")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad HTTP status code %d returned by API", resp.StatusCode)
	}

	csrfHeaderValue := resp.Header.Get(csrfTokenHeaderName)
	if csrfHeaderValue == "" {
		return "", fmt.Errorf("empty csrf token header value")
	}

	return csrfHeaderValue, nil
}
