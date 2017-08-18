# go-stopandshop
> Golang library for interacting with the [Stop and Shop](https://stopandshop.com) API.

## Install

```
$ go get github.com/blaskovicz/go-stopandshop
```

## Use

```go
import (
  sns "github.com/blaskovicz/go-stopandshop"
)

// initialize a default client
client := sns.New()

// log in must be called to access the api.
// this fetches an access token for bearer auth.
if err := client.Login("youremail@somewhere.com", "yourpassword"); err != nil {
  panc(err)
}

// then fetch your profile info
profile, err := client.ReadProfile()
if err != nil {
  panic(err)
}

// with the profile info, we can then check for coupons!
coupons, err := cilent.ReadCoupons(profile.CardNumber)
if err != nil {
  panic(err)
}

fmt.Printf("Found %d coupons, maybe some free items?\n", len(coupons))

// more to come, including using refresh token.
```

## Test

```
$ go test ./...
```
