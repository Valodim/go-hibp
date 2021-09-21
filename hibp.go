package hibp

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Version represents the version of this package
const Version = "0.1.1"

// Client is the HIBP client object
type Client struct {
	hc *http.Client  // HTTP client to perform the API requests
	to time.Duration // HTTP client timeout
	ak string        // HIBP API key

	PwnedPassword *PwnedPassword // Reference to the PwnedPassword API
}

// Option is a function that is used for grouping of Client options.
type Option func(*Client)

// New creates and returns a new HIBP client object
func New(options ...Option) *Client {
	c := &Client{}

	// Set defaults
	c.to = time.Second * 5

	// Set additional options
	for _, opt := range options {
		opt(c)
	}

	// Add a http client to the Client object
	c.hc = httpClient(c.to)

	// Associate the different HIBP service APIs with the Client
	c.PwnedPassword = &PwnedPassword{hc: c}

	return c
}

// WithHttpTimeout overrides the default http client timeout
func WithHttpTimeout(t time.Duration) Option {
	return func(c *Client) {
		c.to = t
	}
}

// WithApiKey set the optional API key to the Client object
func WithApiKey(k string) Option {
	return func(c *Client) {
		c.ak = k
	}
}

// HttpReq performs an HTTP request to the corresponding API
func (c *Client) HttpReq(m, p string) (*http.Request, error) {
	u, err := url.Parse(p)
	if err != nil {
		return nil, err
	}

	hr, err := http.NewRequest(m, u.String(), nil)
	if err != nil {
		return nil, err
	}
	hr.Header.Set("Accept", "application/json")
	hr.Header.Set("User-Agent", fmt.Sprintf("go-hibp v%s - https://github.com/wneessen/go-hibp", Version))

	if c.ak != "" {
		hr.Header["hibp-api-key"] = []string{c.ak}
	}

	return hr, nil
}

// httpClient returns a custom http client for the HIBP Client object
func httpClient(to time.Duration) *http.Client {
	tlsConfig := &tls.Config{
		MaxVersion: tls.VersionTLS13,
		MinVersion: tls.VersionTLS12,
	}
	httpTransport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   5 * time.Second,
	}
	if to.Nanoseconds() > 0 {
		httpClient.Timeout = to
	}

	return httpClient
}
