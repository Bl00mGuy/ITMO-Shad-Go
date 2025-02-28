//go:build !solution

package main

import (
	"bytes"
	"flag"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
)

var (
	serviceAddress = flag.String("service-addr", "localhost:8080", "service address")
	address        = flag.String("addr", "localhost:8081", "firewall listen address")
	configuration  = flag.String("conf", "./firewall/configs/example.yaml", "configuration file path")
)

type Rule struct {
	Endpoint               string   `yaml:"endpoint"`
	ForbiddenUserAgents    []string `yaml:"forbidden_user_agents"`
	ForbiddenHeaders       []string `yaml:"forbidden_headers"`
	RequiredHeaders        []string `yaml:"required_headers"`
	MaxRequestLengthBytes  int64    `yaml:"max_request_length_bytes"`
	MaxResponseLengthBytes int64    `yaml:"max_response_length_bytes"`
	ForbiddenResponseCodes []int    `yaml:"forbidden_response_codes"`
	ForbiddenRequestRe     []string `yaml:"forbidden_request_re"`
	ForbiddenResponseRe    []string `yaml:"forbidden_response_re"`
}

type RulesConfig struct {
	Rules []Rule `yaml:"rules"`
}

type Firewall struct {
	Tripper http.RoundTripper
	Config  RulesConfig
}

func main() {
	flag.Parse()

	config, err := loadConfig(*configuration)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	targetHost := parseHost(*serviceAddress)

	firewall := &Firewall{
		Tripper: &http.Transport{},
		Config:  config,
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = targetHost
		},
		Transport: firewall,
	}

	log.Printf("Starting firewall on %s, forwarding to %s", *address, *serviceAddress)
	log.Fatal(http.ListenAndServe(*address, proxy))
}

func loadConfig(path string) (RulesConfig, error) {
	var config RulesConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	return config, err
}

func parseHost(serviceAddress string) string {
	serviceAddress = strings.TrimPrefix(serviceAddress, "http://")
	serviceAddress = strings.TrimPrefix(serviceAddress, "https://")
	if idx := strings.Index(serviceAddress, "/"); idx != -1 {
		serviceAddress = serviceAddress[:idx]
	}
	return serviceAddress
}

func (firewall Firewall) RoundTrip(request *http.Request) (*http.Response, error) {
	if firewall.isRequestBlocked(request) {
		return forbiddenResponse(), nil
	}

	response, err := firewall.Tripper.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if firewall.isResponseBlocked(request, response) {
		return forbiddenResponse(), nil
	}
	return response, nil
}

func (firewall Firewall) isRequestBlocked(request *http.Request) bool {
	for _, rule := range firewall.Config.Rules {
		if !matchEndpoint(request.URL.String(), rule.Endpoint) {
			continue
		}

		if exceedsMaxRequestLength(request, rule) ||
			matchesForbiddenUserAgent(request, rule) ||
			matchesForbiddenHeaders(request, rule) ||
			missingRequiredHeaders(request, rule) ||
			matchesForbiddenRequestRegex(request, rule) {
			return true
		}
	}
	return false
}

func (firewall Firewall) isResponseBlocked(request *http.Request, response *http.Response) bool {
	for _, rule := range firewall.Config.Rules {
		if !matchEndpoint(request.URL.String(), rule.Endpoint) {
			continue
		}

		if exceedsMaxResponseLength(response, rule) ||
			matchesForbiddenResponseCode(response, rule) ||
			matchesForbiddenResponseRegex(response, rule) {
			return true
		}
	}
	return false
}

func matchEndpoint(url, endpoint string) bool {
	matched, err := regexp.MatchString(endpoint, url)
	if err != nil {
		log.Printf("error matching endpoint %s: %v", endpoint, err)
		return false
	}
	return matched
}

func exceedsMaxRequestLength(request *http.Request, rule Rule) bool {
	return rule.MaxRequestLengthBytes > 0 && request.ContentLength > 0 && request.ContentLength > rule.MaxRequestLengthBytes
}

func matchesForbiddenUserAgent(request *http.Request, rule Rule) bool {
	userAgent := request.Header.Get("User-Agent")
	for _, pattern := range rule.ForbiddenUserAgents {
		if matched, _ := regexp.MatchString(pattern, userAgent); matched {
			return true
		}
	}
	return false
}

func matchesForbiddenHeaders(request *http.Request, rule Rule) bool {
	for _, pattern := range rule.ForbiddenHeaders {
		for name, values := range request.Header {
			for _, value := range values {
				header := name + ": " + value
				if matched, _ := regexp.MatchString(pattern, header); matched {
					return true
				}
			}
		}
	}
	return false
}

func missingRequiredHeaders(request *http.Request, rule Rule) bool {
	for _, header := range rule.RequiredHeaders {
		if request.Header.Get(header) == "" {
			return true
		}
	}
	return false
}

func matchesForbiddenRequestRegex(request *http.Request, rule Rule) bool {
	body, _ := io.ReadAll(request.Body)
	request.Body = io.NopCloser(bytes.NewReader(body))
	for _, pattern := range rule.ForbiddenRequestRe {
		if matched, _ := regexp.Match(pattern, body); matched {
			return true
		}
	}
	return false
}

func exceedsMaxResponseLength(response *http.Response, rule Rule) bool {
	return rule.MaxResponseLengthBytes > 0 && response.ContentLength > 0 && response.ContentLength > rule.MaxResponseLengthBytes
}

func matchesForbiddenResponseCode(response *http.Response, rule Rule) bool {
	for _, code := range rule.ForbiddenResponseCodes {
		if response.StatusCode == code {
			return true
		}
	}
	return false
}

func matchesForbiddenResponseRegex(response *http.Response, rule Rule) bool {
	body, _ := io.ReadAll(response.Body)
	response.Body = io.NopCloser(bytes.NewReader(body))
	for _, pattern := range rule.ForbiddenResponseRe {
		if matched, _ := regexp.Match(pattern, body); matched {
			return true
		}
	}
	return false
}

func forbiddenResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusForbidden,
		Body:       io.NopCloser(strings.NewReader("Forbidden")),
	}
}
