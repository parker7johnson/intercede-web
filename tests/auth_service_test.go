package tests

import (
	"errors"
	"reflect"
	"testing"

	"github.com/parkerjohnson/intercede/services"
	"github.com/parkerjohnson/intercede/tests/mocks"
	"github.com/supabase-community/gotrue-go/types"
)

func TestServiceLogin(t *testing.T) {

	cases := []struct {
		name    string
		authRes types.Session
		authErr error
	}{
		{"valid credentials", types.Session{}, nil},
		{"bad credentials", types.Session{}, errors.New("bad credentials provided")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			svc := services.NewAuthService(&mocks.MockSupabaseClient{
				SignInFn: func(s1, s2 string) (types.Session, error) {
					return tc.authRes, tc.authErr
				},
				AutoRefreshFn: func(s types.Session) {return},
			})

			session, err := svc.Login(testEmail, testPassword)

			if !reflect.DeepEqual(session, tc.authRes){
				t.Fatalf("expected session %v but got %v", tc.authRes, session)
			}

			if err != tc.authErr {
				t.Fatalf("expected error %v but got %v", tc.authRes, session)
			}

		})
	}
}
