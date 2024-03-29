package tests

import (
	"testing"
	"time"

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
	_, err = doTestRequest(token, t.Name(), map[string]any{
		"JwtExpiration":        Cfg.JwtExpiration(),
		"JwtRefreshExpiration": Cfg.JwtRefreshExpiration(),
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestLoginFail(t *testing.T) {
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

func TestRefreshToken(t *testing.T) {
	tokens, err := core.NewJwt(&model.User{ID: 1, Email: "borer-hudson@yahoo.com", Name: "Hudson Borer", Role: model.RoleUser},
		Cfg.JwtExpiration(), Cfg.JwtRefreshExpiration(), Cfg.JwtSecret())
	if err != nil {
		t.Errorf("create token %v", err)
		return
	}
	_, err = doTestRequest(tokens.RefreshToken, t.Name(), map[string]any{
		"OldRefreshToken":      tokens.RefreshToken,
		"JwtExpiration":        Cfg.JwtExpiration(),
		"JwtRefreshExpiration": Cfg.JwtRefreshExpiration(),
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRefreshTokenFailed(t *testing.T) {
	tokens, err := core.NewJwt(&model.User{ID: 1, Email: "borer-hudson@yahoo.com", Name: "Hudson Borer", Role: model.RoleUser},
		Cfg.JwtExpiration(), Cfg.JwtRefreshExpiration(), Cfg.JwtSecret())
	if err != nil {
		t.Errorf("create token %v", err)
		return
	}
	_, err = doTestRequest(tokens.RefreshToken, t.Name(), nil)
	if err != nil {
		t.Error(err)
		return
	}
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
