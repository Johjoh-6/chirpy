package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSignAndVerify(t *testing.T) {
	cases := []struct {
		name           string
		secret         string
		validateSecret string
		expiresIn      time.Duration
		expectError    bool
	}{
		{
			name:           "valid input",
			secret:         "secret",
			validateSecret: "secret",
			expiresIn:      time.Hour,
			expectError:    false,
		},
		{
			name:           "invalid secret",
			secret:         "secret",
			validateSecret: "invalid",
			expiresIn:      time.Hour,
			expectError:    true,
		},
		{
			name:           "passed expired jwt",
			secret:         "passed",
			validateSecret: "passed",
			expiresIn:      -time.Hour,
			expectError:    true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id := uuid.New()
			jwtString, err := MakeJWT(id, c.secret, c.expiresIn)
			if err != nil {
				if !c.expectError {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}

			// check that the JWT is valid
			uuidJwt, err := ValidateJWT(jwtString, c.validateSecret)
			if err != nil {
				if !c.expectError {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}
			if c.expectError {
				t.Errorf("expected error, got nil")
			}
			if uuidJwt != id {
				t.Errorf("expected uuid %v, got %v", id, uuidJwt)
			}
		})
	}
}
