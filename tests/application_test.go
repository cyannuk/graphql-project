package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"graphql-project/core"
	"graphql-project/domain/model"
)

func TestJwtAuthentication(t *testing.T) {
	tokens, err := core.NewJwt(&model.User{ID: 1, Email: "borer-hudson@yahoo.com", Name: "Hudson Borer", Role: model.RoleUser},
		0, 0, Cfg.JwtSecret())
	if err != nil {
		t.Errorf("create token %v", err)
		return
	}
	_, err = doTestRequest(tokens.AccessToken, t.Name(), nil)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestLogin(t *testing.T) {
	token, err := core.JwtAnon(time.Hour, Cfg.JwtSecret())
	if err != nil {
		t.Errorf("create token %v", err)
		return
	}
	data, err := doTestRequest(token, t.Name(), func(m map[string]any) any { return getTokens(m, "login") })
	if err != nil {
		t.Error(err)
		return
	}
	newTokens := getTokens(data, "login")
	assert.NotNil(t, newTokens, "invalid jwt")
	claims, ok := core.JwtVerify(newTokens.AccessToken, Cfg.JwtSecret())
	assert.True(t, ok, "jwt not verified")
	assert.Equal(t, "Hudson Borer", claims.Name)
	assert.Equal(t, "borer-hudson@yahoo.com", claims.Email)
	assert.Equal(t, int64(1), claims.Uid)
	assert.Equal(t, model.RoleUser, claims.Role)

	claims, ok = core.JwtVerify(newTokens.RefreshToken, Cfg.JwtSecret())
	assert.True(t, ok, "JWT not verified")
	assert.Equal(t, int64(1), claims.Uid)
	assert.Equal(t, model.RoleRefresh, claims.Role)
}

func TestLoginFail(t *testing.T) {
	token, err := core.JwtAnon(time.Hour, Cfg.JwtSecret())
	if err != nil {
		t.Errorf("create token %v", err)
		return
	}
	data, err := doTestRequest(token, t.Name(), nil)
	if err != nil {
		t.Error(err)
		return
	}
	newTokens := getTokens(data, "login")
	assert.Nil(t, newTokens, "returned JWT")
}

func TestRefreshToken(t *testing.T) {
	tokens, err := core.NewJwt(&model.User{ID: 1, Email: "borer-hudson@yahoo.com", Name: "Hudson Borer", Role: model.RoleUser},
		Cfg.JwtExpiration(), Cfg.JwtRefreshExpiration(), Cfg.JwtSecret())
	if err != nil {
		t.Errorf("create token %v", err)
		return
	}
	data, err := doTestRequest(tokens.RefreshToken, t.Name(), func(m map[string]any) any {
		p := map[string]any{"OldRefreshToken": tokens.RefreshToken}
		if t := getTokens(m, "refreshToken"); t != nil {
			p["AccessToken"] = t.AccessToken
			p["RefreshToken"] = t.RefreshToken
		}
		return p
	})
	if err != nil {
		t.Error(err)
		return
	}
	newTokens := getTokens(data, "refreshToken")
	assert.NotNil(t, newTokens, "invalid JWT")
	claims, ok := core.JwtVerify(newTokens.AccessToken, Cfg.JwtSecret())
	assert.True(t, ok, "jwt not verified")
	assert.Equal(t, "Hudson Borer", claims.Name)
	assert.Equal(t, "borer-hudson@yahoo.com", claims.Email)
	assert.Equal(t, int64(1), claims.Uid)
	assert.Equal(t, model.RoleUser, claims.Role)

	claims, ok = core.JwtVerify(newTokens.RefreshToken, Cfg.JwtSecret())
	assert.True(t, ok, "JWT not verified")
	assert.Equal(t, int64(1), claims.Uid)
	assert.Equal(t, model.RoleRefresh, claims.Role)
}

func TestRefreshTokenFailed(t *testing.T) {
	tokens, err := core.NewJwt(&model.User{ID: 1, Email: "borer-hudson@yahoo.com", Name: "Hudson Borer", Role: model.RoleUser},
		Cfg.JwtExpiration(), Cfg.JwtRefreshExpiration(), Cfg.JwtSecret())
	if err != nil {
		t.Errorf("create token %v", err)
		return
	}
	data, err := doTestRequest(tokens.RefreshToken, t.Name(), nil)
	if err != nil {
		t.Error(err)
		return
	}
	newTokens := getTokens(data, "refreshToken")
	assert.Nil(t, newTokens, "returned JWT")
}

func TestGetUser(t *testing.T) {
	tokens, err := core.NewJwt(&model.User{ID: 1, Email: "borer-hudson@yahoo.com", Name: "Hudson Borer", Role: model.RoleUser},
		Cfg.JwtExpiration(), Cfg.JwtRefreshExpiration(), Cfg.JwtSecret())
	if err != nil {
		t.Errorf("create token %v", err)
		return
	}
	_, err = doTestRequest(tokens.AccessToken, t.Name(), nil)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCreateUser(t *testing.T) {
	token, err := core.JwtAnon(time.Hour, Cfg.JwtSecret())
	if err != nil {
		t.Errorf("create token %v", err)
		return
	}
	_, err = doTestRequest(token, t.Name(), nil)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestUpdateUser(t *testing.T) {
	tokens, err := core.NewJwt(&model.User{ID: 1, Email: "borer-hudson@yahoo.com", Name: "Hudson Borer", Role: model.RoleUser},
		Cfg.JwtExpiration(), Cfg.JwtRefreshExpiration(), Cfg.JwtSecret())
	if err != nil {
		t.Errorf("create token %v", err)
		return
	}
	_, err = doTestRequest(tokens.AccessToken, t.Name(), nil)
	if err != nil {
		t.Error(err)
		return
	}
}
