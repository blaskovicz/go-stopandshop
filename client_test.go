package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	var username string
	var password string
	var firstName string
	var err error
	var client *Client
	var profile *Profile
	t.Run("Init", func(t *testing.T) {
		client = New()
		require.NotNil(t, client, "client was nil")
	})
	// TODO mock out the backend for testing once i'm satisfied with real data
	t.Run("Login", func(t *testing.T) {
		err = client.Login(username, password)
		require.NoError(t, err, "login failed")
	})
	t.Run("ReadProfile", func(t *testing.T) {
		profile, err = client.ReadProfile()
		require.NoError(t, err, "read profile failed")
		require.NotNil(t, profile, "profile was nil")
		require.Equal(t, firstName, profile.FirstName)
	})
	t.Run("ReadCoupons", func(t *testing.T) {
		var cs []Coupon
		cs, err = client.ReadCoupons(profile.CardNumber)
		require.NoError(t, err, "read coupons failed")
		require.NotNil(t, cs, "coupons were nil")
		require.NotEqual(t, 0, len(cs))

		/* now you kno
		for _, c := range cs {
			// nasty for now
			maybeFree := strings.Contains(strings.ToLower(fmt.Sprintf("%#v", c)), "free")
			if maybeFree {
				fmt.Printf("MAYBE %#v?\n", c)
				if !c.Loaded {
					// TODO -> do the load
					fmt.Printf("\tbut first load it!\n")
				}
			}
		}
		*/
	})
}
