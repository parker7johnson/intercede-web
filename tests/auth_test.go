package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/parkerjohnson/intercede/handlers/api"
	"github.com/parkerjohnson/intercede/services"
	"github.com/parkerjohnson/intercede/tests/mocks"
	"github.com/supabase-community/gotrue-go/types"
)



func newAuthHandler(authResponse types.Session, authErr error) *api.ApiHandlers {
	ah := services.NewAuthService(&mocks.MockSupabaseClient{
		SignInFn: func(s1, s2 string) (types.Session, error) {
			return authResponse, authErr
		},
		AutoRefreshFn: func(s types.Session) {
			return
		},
	})

	return api.NewApi(nil, ah)
}

func newAuthFormRequest(method string) *http.Request {
	form := url.Values{}

	form.Add("email", testEmail)
	form.Add("password", testPassword)

	request := httptest.NewRequest(method, adminAPILoginPath, strings.NewReader(form.Encode()))

	return request
}

func TestLogin(t *testing.T) {

	cases := []struct {
		name         string
		method       string
		wantStatus   int
		authResponse types.Session
		authErr      error
	}{
		{"valid login", http.MethodPost, http.StatusOK, types.Session{AccessToken: testToken}, nil},
		{"bad method", http.MethodGet, http.StatusMethodNotAllowed, types.Session{AccessToken: testToken}, nil},
		{"bad login", http.MethodPost, http.StatusUnauthorized, types.Session{AccessToken: testToken}, errors.New("bad login supplied")},
	}

  for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			newAuthHandler(tc.authResponse, tc.authErr).Login(w, newAuthFormRequest(tc.method))

			if w.Code != tc.wantStatus {
				t.Fatalf("Failed with status %v wanted %v", w.Code, tc.wantStatus)
			}

			if tc.authErr == nil {
				cookies := w.Result().Cookies()

				for _, cookie := range cookies {
					if cookie.Name == sessionCookieName && cookie.Value != testToken {
						t.Fatalf("Token not set correctly, wanted %v got %v", testToken, cookie.Value)
					}
				}
			}
		})
	}
}
