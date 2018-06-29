// Code generated by counterfeiter. DO NOT EDIT.
package authfakes

import (
	"sync"

	"github.com/Peripli/service-manager-cli/internal/auth"
	"golang.org/x/oauth2"
)

type FakeAuthenticationStrategy struct {
	AuthenticateStub        func(issuerURL, user, password string) (*oauth2.Config, *oauth2.Token, error)
	authenticateMutex       sync.RWMutex
	authenticateArgsForCall []struct {
		issuerURL string
		user      string
		password  string
	}
	authenticateReturns struct {
		result1 *oauth2.Config
		result2 *oauth2.Token
		result3 error
	}
	authenticateReturnsOnCall map[int]struct {
		result1 *oauth2.Config
		result2 *oauth2.Token
		result3 error
	}
	RefreshTokenStub        func(oauth2.Config, oauth2.Token) (*oauth2.Token, error)
	refreshTokenMutex       sync.RWMutex
	refreshTokenArgsForCall []struct {
		arg1 oauth2.Config
		arg2 oauth2.Token
	}
	refreshTokenReturns struct {
		result1 *oauth2.Token
		result2 error
	}
	refreshTokenReturnsOnCall map[int]struct {
		result1 *oauth2.Token
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeAuthenticationStrategy) Authenticate(issuerURL string, user string, password string) (*oauth2.Config, *oauth2.Token, error) {
	fake.authenticateMutex.Lock()
	ret, specificReturn := fake.authenticateReturnsOnCall[len(fake.authenticateArgsForCall)]
	fake.authenticateArgsForCall = append(fake.authenticateArgsForCall, struct {
		issuerURL string
		user      string
		password  string
	}{issuerURL, user, password})
	fake.recordInvocation("Authenticate", []interface{}{issuerURL, user, password})
	fake.authenticateMutex.Unlock()
	if fake.AuthenticateStub != nil {
		return fake.AuthenticateStub(issuerURL, user, password)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.authenticateReturns.result1, fake.authenticateReturns.result2, fake.authenticateReturns.result3
}

func (fake *FakeAuthenticationStrategy) AuthenticateCallCount() int {
	fake.authenticateMutex.RLock()
	defer fake.authenticateMutex.RUnlock()
	return len(fake.authenticateArgsForCall)
}

func (fake *FakeAuthenticationStrategy) AuthenticateArgsForCall(i int) (string, string, string) {
	fake.authenticateMutex.RLock()
	defer fake.authenticateMutex.RUnlock()
	return fake.authenticateArgsForCall[i].issuerURL, fake.authenticateArgsForCall[i].user, fake.authenticateArgsForCall[i].password
}

func (fake *FakeAuthenticationStrategy) AuthenticateReturns(result1 *oauth2.Config, result2 *oauth2.Token, result3 error) {
	fake.AuthenticateStub = nil
	fake.authenticateReturns = struct {
		result1 *oauth2.Config
		result2 *oauth2.Token
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeAuthenticationStrategy) AuthenticateReturnsOnCall(i int, result1 *oauth2.Config, result2 *oauth2.Token, result3 error) {
	fake.AuthenticateStub = nil
	if fake.authenticateReturnsOnCall == nil {
		fake.authenticateReturnsOnCall = make(map[int]struct {
			result1 *oauth2.Config
			result2 *oauth2.Token
			result3 error
		})
	}
	fake.authenticateReturnsOnCall[i] = struct {
		result1 *oauth2.Config
		result2 *oauth2.Token
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeAuthenticationStrategy) RefreshToken(arg1 oauth2.Config, arg2 oauth2.Token) (*oauth2.Token, error) {
	fake.refreshTokenMutex.Lock()
	ret, specificReturn := fake.refreshTokenReturnsOnCall[len(fake.refreshTokenArgsForCall)]
	fake.refreshTokenArgsForCall = append(fake.refreshTokenArgsForCall, struct {
		arg1 oauth2.Config
		arg2 oauth2.Token
	}{arg1, arg2})
	fake.recordInvocation("RefreshToken", []interface{}{arg1, arg2})
	fake.refreshTokenMutex.Unlock()
	if fake.RefreshTokenStub != nil {
		return fake.RefreshTokenStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.refreshTokenReturns.result1, fake.refreshTokenReturns.result2
}

func (fake *FakeAuthenticationStrategy) RefreshTokenCallCount() int {
	fake.refreshTokenMutex.RLock()
	defer fake.refreshTokenMutex.RUnlock()
	return len(fake.refreshTokenArgsForCall)
}

func (fake *FakeAuthenticationStrategy) RefreshTokenArgsForCall(i int) (oauth2.Config, oauth2.Token) {
	fake.refreshTokenMutex.RLock()
	defer fake.refreshTokenMutex.RUnlock()
	return fake.refreshTokenArgsForCall[i].arg1, fake.refreshTokenArgsForCall[i].arg2
}

func (fake *FakeAuthenticationStrategy) RefreshTokenReturns(result1 *oauth2.Token, result2 error) {
	fake.RefreshTokenStub = nil
	fake.refreshTokenReturns = struct {
		result1 *oauth2.Token
		result2 error
	}{result1, result2}
}

func (fake *FakeAuthenticationStrategy) RefreshTokenReturnsOnCall(i int, result1 *oauth2.Token, result2 error) {
	fake.RefreshTokenStub = nil
	if fake.refreshTokenReturnsOnCall == nil {
		fake.refreshTokenReturnsOnCall = make(map[int]struct {
			result1 *oauth2.Token
			result2 error
		})
	}
	fake.refreshTokenReturnsOnCall[i] = struct {
		result1 *oauth2.Token
		result2 error
	}{result1, result2}
}

func (fake *FakeAuthenticationStrategy) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.authenticateMutex.RLock()
	defer fake.authenticateMutex.RUnlock()
	fake.refreshTokenMutex.RLock()
	defer fake.refreshTokenMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeAuthenticationStrategy) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ auth.AuthenticationStrategy = new(FakeAuthenticationStrategy)
