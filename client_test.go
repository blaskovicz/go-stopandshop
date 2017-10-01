package stopandshop

import (
	"testing"

	"github.com/blaskovicz/go-stopandshop/mocks"
	"github.com/blaskovicz/go-stopandshop/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	s := mocks.NewStopAndShopAPI()
	defer s.Close()
	username := "testuser"
	password := "testpass"
	var err error
	var client *Client
	var profile *models.Profile
	t.Run("Init", func(t *testing.T) {
		client = New()
		client.SetRootURI(s.URL())
		require.NotNil(t, client, "client was nil")
	})
	t.Run("Login", func(t *testing.T) {
		err = client.Login(username, password)
		require.NoError(t, err, "login failed")
	})
	t.Run("ReadProfile", func(t *testing.T) {
		require.NotNil(t, client)
		profile, err = client.ReadProfile()
		require.NoError(t, err, "read profile failed")
		require.NotNil(t, profile, "profile was nil")
		require.Equal(t, "test", profile.FirstName)
		require.Equal(t, "12345", profile.CardNumber)
		require.Equal(t, "foo@somesite.com", profile.Login)
	})
	t.Run("ReadCoupons", func(t *testing.T) {
		require.NotNil(t, client)
		require.NotNil(t, profile)
		var cs []models.Coupon
		cs, err = client.ReadCoupons(profile.CardNumber)
		require.NoError(t, err, "read coupons failed")
		require.NotNil(t, cs, "coupons were nil")
		require.Equal(t, 2, len(cs))
		c := cs[0]
		assert.Equal(t, "0a06f213-298d-47cc-9260-99fc4450c0a4", c.ID, "id mismatch")
		assert.Equal(t, "Fortify™", c.Name, "name mistmatch")
		assert.Equal(t, "On any Fortify™ 50 Billion Formula", c.Description, "description mismatch")
		assert.False(t, c.Loaded, "coupon loaded")

		c = cs[1]
		assert.Equal(t, "2b9f19ca-7655-499a-8284-491c039068b4", c.ID, "id mismatch")
		assert.Equal(t, "Duncan Hines®", c.Name, "name mistmatch")
		assert.Equal(t, "Duncan Hines® Perfect Size for 1® Cake Mix (Max. Value $3.49)", c.Description, "description mismatch")
		assert.True(t, c.Loaded, "coupon not loaded")
	})
	t.Run("LoadCouponToProfile", func(t *testing.T) {
		require.NotNil(t, client)
		require.NotNil(t, profile)
		err := client.LoadCoupon(profile.CardNumber, "0a06f213-298d-47cc-9260-99fc4450c0a4")
		require.NoError(t, err, "load coupon failed")

		err = client.LoadCoupon(profile.CardNumber, "0a06f213-298d-47cc-9260-99fc4450c0a4")
		require.NoError(t, err, "load coupon failed")

		err = client.LoadCoupon(profile.CardNumber, "")
		require.Error(t, err, "load coupon failed")
	})
}
