package stopandshop

import "strings"

// Reports true if the another access_token needs to be
// generated from the refresh_token
func IsAccessTokenExpired(err error) bool {
	return strings.Contains(err.Error(), "Invalid access token")
}

// Returns true if the user must log in again.
func IsRefreshTokenExpired(err error) bool {
	return strings.Contains(err.Error(), "Invalid refresh token")
}
