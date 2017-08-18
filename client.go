package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/blaskovicz/go-stopandshop/models"
)

const (
	DefaultRootURI = "https://stopandshop.com"

	// To get these:
	// 1) go here -> https://stopandshop.com/
	// 2) find base.min.<sha>.js script (eg: https://images.stopandshop.com/static/common/js/bundle/base.min.c95120e11fc413ef.js)
	// 3) grab any cid: and tokenAuth: pair (SNS is below)
	DefaultClientID  = "54a8c012-9956-44a5-961b-20f93521c15e"
	DefaultTokenAuth = "NTRhOGMwMTItOTk1Ni00NGE1LTk2MWItMjBmOTM1MjFjMTVlOjczNzk3ZTQ3YmRhMGQ5NzQ3ZjlmOGI2NTg5YTc2YzgzNjhjODI3ZmRkNjU0YmFhNjQ5MjYwZWY0OGI0MmM1YTY="
)

type Client struct {
	rootURI   string
	tokenAuth string // for client_credentials grant; base64-encoded oauth client_id:<somepass>
	clientID  string // for oauth
	token     *models.Token
}

func New() *Client {
	tokenAuth := os.Getenv("STOP_AND_SHOP_TOKEN_AUTH")
	if tokenAuth == "" {
		tokenAuth = DefaultTokenAuth
	}
	clientID := os.Getenv("STOP_AND_SHOP_CLIENT_ID")
	if clientID == "" {
		clientID = DefaultClientID
	}
	rootURI := os.Getenv("STOP_AND_SHOP_ROOT_URI")
	if rootURI == "" {
		rootURI = DefaultRootURI
	}
	return &Client{rootURI: rootURI, tokenAuth: tokenAuth, clientID: clientID}
}

func (c *Client) uri(path string) string {
	return fmt.Sprintf("%s%s", c.rootURI, path)
}

func (c *Client) assertToken() error {
	if c.token == nil {
		return fmt.Errorf("client not logged in")
	}
	// TODO, keep track of expires/time and refresh the token if needed
	return nil
}

// TODO common c.call function
func (c *Client) ReadCoupons(cardNumber string) ([]models.Coupon, error) {
	if err := c.assertToken(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", c.uri(fmt.Sprintf("/auth/api/private/synergy/coupons/offers/%s?pageIndex=0&numRecords=2000&categories&brands", cardNumber)), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var o models.OfferResponse
		err = json.NewDecoder(resp.Body).Decode(&o)
		if err != nil {
			return nil, fmt.Errorf("failed to decode coupons: %s", err)
		}
		return o.Offers, nil
	} else {
		var e models.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			// TODO print payload substring here
			return nil, fmt.Errorf("failed to decode error payload: %s", err)
		}
		return nil, fmt.Errorf("coupon fetch failed: %s", e.Description)
	}
}
func (c *Client) ReadProfile() (*models.Profile, error) {
	if err := c.assertToken(); err != nil {
		return nil, err
	}
	// TODO figure out what /auth/profile (since it doesnt have card num)
	req, err := http.NewRequest("GET", c.uri("/auth/profile/SNS"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var p models.Profile
		err = json.NewDecoder(resp.Body).Decode(&p)
		if err != nil {
			return nil, fmt.Errorf("failed to decode profile: %s", err)
		}
		return &p, nil
	} else {
		var e models.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return nil, fmt.Errorf("failed to decode error payload: %s", err)
		}
		return nil, fmt.Errorf("profile fetch failed: %s", e.Description)
	}
}

// on stop-and-shop web, the flow is:
// 1) POST form to /auth/oauth/token (grant_type=client_credentials, scope=profile) with Authorization: Basic <SNS.tokenAuth> ->
//		{"access_token":"deadbeef2e06e1473542e3a58987c34140494592a52228844eed8508ec76a133","token_type":"Bearer","expires_in":2514,"scope":"profile"}
// 2) POST from to /auth/oauth/token (grant_type=password, username=<email>, password=<password>, client_id=<SNS.cid> with Authorization: Basic <SNS.tokenAuth> ->
//    {"access_token":"deadbeefeb2808001d182ebc24f31a44ac948ba1e20e3d7661104d8109ada6e3","token_type":"bearer","refresh_token":"deadbeef-cf80-4adb-8f3f-c51a18628cfd","expires_in":3599,"scope":"profile"}
// I think phase 1 is an interm token for other bearer requests, so skip it and just do number 2. Note we must spoof the client id (until I can figure out where to register new clients).
func (c *Client) Login(username, password string) error {
	payload := url.Values{
		"username":   []string{username},
		"password":   []string{password},
		"grant_type": []string{"password"},
		"client_id":  []string{c.clientID},
	}
	req, err := http.NewRequest("POST", c.uri("/auth/oauth/token"), strings.NewReader(payload.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// this is needed, unfortunately
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", c.tokenAuth))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var t models.Token
		err = json.NewDecoder(resp.Body).Decode(&t)
		if err != nil {
			return fmt.Errorf("failed to decode token: %s", err)
		}
		c.token = &t
	} else {
		var e models.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return fmt.Errorf("failed to decode error payload: %s", err)
		}
		return fmt.Errorf("login failed: %s", e.Description)
	}
	return nil
}
