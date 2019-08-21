package oidc

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/Peripli/service-manager-cli/pkg/auth"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OIDC Client Suite")
}

var _ = Describe("OIDC Client", func() {
	options := &auth.Options{
		ClientID:      "client-id",
		ClientSecret:  "client-secret",
		User:          "user",
		Password:      "password",
		TokenEndpoint: "http://token-endpoint",
	}
	token := newToken(1 * time.Hour)
	tokenNoRefreshToken := &auth.Token{
		AccessToken:  "access-token",
		TokenType:    "token-type",
		RefreshToken: "",
		ExpiresIn:    time.Now().Add(-1 * time.Hour),
	}

	DescribeTable("NewClient",
		func(options *auth.Options, token *auth.Token, expectedErrMsg string, expetedToken *auth.Token) {
			client := NewClient(options, token)
			t, err := client.Token()
			if expectedErrMsg == "" {
				Expect(err).To(BeNil())
				Expect(*t).To(Equal(*token))
			} else {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring(expectedErrMsg))
			}
		},
		Entry("Valid token - reuses the token", options, token, "", token),
		Entry("No client credentials and valid token - reuses the token",
			&auth.Options{}, token, "", token),
		Entry("No client credentials and expired token - returns error to login",
			&auth.Options{},
			newToken(-1*time.Hour),
			"smctl login",
			nil),
		Entry("With client credentials and refresh token - refreshes the token",
			options,
			newToken(-1*time.Hour),
			options.TokenEndpoint,
			nil),
		Entry("With client credentials and no refresh token - fetches a new token using client credentials flow",
			&auth.Options{
				ClientID:      "client-id",
				ClientSecret:  "client-secret",
				TokenEndpoint: "http://token-endpoint",
			},
			tokenNoRefreshToken,
			"http://token-endpoint",
			nil),
		Entry("With client and user credentials and no refresh token - returns error to login",
			options,
			tokenNoRefreshToken,
			"smctl login",
			nil),
	)
})

func newToken(validity time.Duration) *auth.Token {
	return &auth.Token{
		AccessToken:  "access-token",
		TokenType:    "token-type",
		RefreshToken: "refresh-token",
		ExpiresIn:    time.Now().Add(validity),
	}
}
