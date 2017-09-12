package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/blaskovicz/go-stopandshop/models"
)

type StopAndShopAPI struct {
	server        *httptest.Server
	offersApplied map[string]interface{}
}

func stopAndShopMux(s *StopAndShopAPI) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/oauth/token", s.handleLogin)
	mux.HandleFunc("/auth/profile/SNS", s.handleProfile)
	mux.HandleFunc("/auth/api/private/synergy/coupons/offers/", s.handleOffers)
	return mux
}

func NewStopAndShopAPI() *StopAndShopAPI {
	s := &StopAndShopAPI{offersApplied: map[string]interface{}{}}
	s.server = httptest.NewServer(stopAndShopMux(s))
	return s
}

func (s *StopAndShopAPI) URL() string {
	return s.server.URL
}

func (s *StopAndShopAPI) Close() error {
	s.server.Close()
	return nil
}
func hasBearerToken(w http.ResponseWriter, req *http.Request) (ok bool) {
	auth := req.Header.Get("Authorization")
	if !strings.Contains(auth, "Bearer deadbeefeb2808001d182ebc24f31a44ac948ba1e20e3d7661104d8109ada6e3") {
		writeError(w, http.StatusUnauthorized, &models.ErrorResponse{ID: "unauthorized", Description: "Token required"})
		return false
	}
	return true
}
func writeError(w http.ResponseWriter, code int, e *models.ErrorResponse) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		fmt.Println(err)
	}
}
func (s *StopAndShopAPI) handleOffers(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else if !hasBearerToken(w, req) {
		return
	} else if !strings.HasSuffix(req.URL.Path, "/12345") {
		w.WriteHeader(http.StatusNotFound) // for credit card
	}
	if req.Method == "GET" {
		s.handleGetOffers(w, req)
	} else if req.Method == "PUT" {
		s.handlePutOffers(w, req)
	}
}
func (s *StopAndShopAPI) handlePutOffers(w http.ResponseWriter, req *http.Request) {
	var c models.CouponPayload
	if err := json.NewDecoder(req.Body).Decode(&c); err != nil || (err == nil && c.CouponID == "") {
		w.WriteHeader(400)
		w.Write([]byte(`{"exception":"org.springframework.web.client.HttpClientErrorException: 400 Bad Request","url":"https://stopandshop.com/auth/api/private/synergy/coupons/offers/12345","timestamp":"2017-09-11T21:15:39.967","status":400}`))
		return
	}
	if _, wasApplied := s.offersApplied[c.CouponID]; !wasApplied {
		s.offersApplied[c.CouponID] = struct{}{}
		w.Write([]byte(`{ "code" : "0", "description" : "Success: Customer 0660000000012 opted into customer group 44." }`))
	}
	// empty body otherwise
}
func (s *StopAndShopAPI) handleGetOffers(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`{
		"cardNumber": "12345",
		"offers": [
			{
					"id" : "0a06f213-298d-47cc-9260-99fc4450c0a4",
					"name" : "Fortify™",
					"description" : "On any Fortify™ 50 Billion Formula",
					"startDate" : "2017-07-10",
					"expirationDate" : "2017-08-31",
					"url" : "http://cdn.cpnscdn.com/static.coupons.com/ext/bussys/cpa/pod/94/267/494267_3638a9fc-c4b8-49b1-aa37-f2e1b0dab086.gif",
					"loaded" : false,
					"title" : "Save $4.00",
					"price" : 4.00,
					"couponSource" : "OMS",
					"couponCategory" : "Health & Wellness",
					"priceQualifier" : "0",
					"source" : "COUPONS"
			}
		]
	}`))
}
func (s *StopAndShopAPI) handleProfile(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else if !hasBearerToken(w, req) {
		return
	}
	w.Write([]byte(`{
		"card_number": "12345",
		"first_name": "test",
		"id": "999000",
		"preferred_store": "80239",
		"storeNumber": "00001",
		"login": "foo@somesite.com"
	}`))
}
func (s *StopAndShopAPI) handleLogin(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// TODO check username, password, and other fields if contract changes
	if _, _, ok := req.BasicAuth(); !ok {
		writeError(w, http.StatusUnauthorized, &models.ErrorResponse{ID: "unauthorized", Description: "Basic creds required"})
		return
	}
	w.Write([]byte(`
	{"access_token":"deadbeefeb2808001d182ebc24f31a44ac948ba1e20e3d7661104d8109ada6e3","token_type": "bearer","refresh_token":"deadbeef-cf80-4adb-8f3f-c51a18628cfd","expires_in":3599,"scope":"profile"}
	`))
}
