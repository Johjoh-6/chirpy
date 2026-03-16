package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	cases := []struct {
		name    string
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name:    "no headers",
			headers: http.Header{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "bearer token",
			headers: http.Header{"Authorization": []string{"Bearer token123"}},
			want:    "token123",
			wantErr: false,
		},
		{
			name:    "bearer token with space",
			headers: http.Header{"Authorization": []string{"Bearer "}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "basic auth",
			headers: http.Header{"Authorization": []string{"Basic token123"}},
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetBearerToken(tc.headers)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("GetBearerToken() got = %v, want %v", got, tc.want)
			}
		})
	}
}
