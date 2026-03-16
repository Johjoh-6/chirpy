package auth

import "testing"

func TestHashPassword(t *testing.T) {
	cases := []struct {
		password string
		error    bool
	}{
		{password: "password", error: false},
		{password: "", error: false},
	}

	for _, c := range cases {
		hash, err := HashPassword(c.password)
		if err != nil && !c.error {
			t.Errorf("HashPassword(%q) error: %v", c.password, err)
		}
		if hash != "" && c.error {
			t.Errorf("HashPassword(%q) = %q, want error", c.password, hash)
		}
	}
}

func TestCheckPasswordHash(t *testing.T) {
	cases := []struct {
		password    string
		checkedPass string
		expected    bool
	}{
		{password: "password", checkedPass: "password", expected: true},
		{password: "", checkedPass: "", expected: true},
		{password: "password", checkedPass: "", expected: false},
		{password: "mylovelypassword", checkedPass: "mylovelypassword", expected: true},
		{password: "password", checkedPass: "mylovelypassword", expected: false},
	}

	for _, c := range cases {
		// create the hash
		hash, err := HashPassword(c.password)
		if err != nil {
			t.Errorf("HashPassword(%q) error: %v", c.password, err)
		}
		ok, err := CheckPasswordHash(c.checkedPass, hash)
		if err != nil {
			t.Errorf("CheckPasswordHash(%q, %q) error: %v", c.password, hash, err)
		}
		if ok != c.expected {
			t.Errorf("CheckPasswordHash(%q, %q) = %v, want %v", c.password, hash, ok, c.expected)
		}
	}
}
