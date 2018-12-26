/*
 * @Author: Allen Flickinger (allen.flickinger@gmail.com)
 * @Date: 2018-01-07 12:40:31
 * @Last Modified by: FuzzyStatic
 * @Last Modified time: 2018-01-16 19:32:55
 */

// Package blizzard is the top level library needed to use the
// API calls for Config games within the subpackages.
package blizzard

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// For testing
var c *Config

// Config regional API URLs, locale, access token, api key
type Config struct {
	client   *http.Client
	oauth    OAuth
	region   Region
	oauthURL string
	apiURL   string
	locale   string
}

// Region type
type Region int

// Region constants
const (
	CN Region = iota
	EU
	KR
	SEA
	TW
	US
)

// Path constants
const (
	localeQuery = "locale="
)

// New create new WorldOfWarcraft structure
func New(clientID, clientSecret string, region Region) *Config {
	var c = Config{
		client: &http.Client{
			Timeout: time.Second * time.Duration(60),
		},
		oauth: OAuth{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			ExpiresAt:    time.Now(),
		},
		region: region,
	}

	switch c.region {
	case EU:
		c.oauthURL = "https://eu.battle.net"
		c.apiURL = "https://eu.api.blizzard.com"
		c.locale = "en_GB"
	case KR:
		c.oauthURL = "https://kr.battle.net"
		c.apiURL = "https://kr.api.blizzard.com"
		c.locale = "ko_KR"
	case SEA:
		c.oauthURL = "https://sea.battle.net"
		c.apiURL = "https://sea.api.blizzard.com"
		//c.locale = "zh_CN"
	case TW:
		c.oauthURL = "https://tb.battle.net"
		c.apiURL = "https://tb.api.blizzard.com"
		c.locale = "zh_TW"
	case US:
		c.oauthURL = "https://us.battle.net"
		c.apiURL = "https://us.api.blizzard.com"
		c.locale = "en_US"
	default: // USA! USA!
		c.oauthURL = "https://us.battle.net"
		c.apiURL = "https://us.api.blizzard.com"
		c.locale = "en_US"
	}

	return &c
}

// GetURLBody processes simple GET request based on URL
func (c *Config) GetURLBody(url string) ([]byte, error) {
	var (
		req  *http.Request
		res  *http.Response
		body []byte
		err  error
	)

	err = c.UpdateAccessTokenIfExp()
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.oauth.AccessTokenRequest.AccessToken)
	req.Header.Set("Accept", "application/json")

	res, err = c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			return
		}
	}()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		return nil, errors.New(res.Status)
	}

	return body, nil
}
