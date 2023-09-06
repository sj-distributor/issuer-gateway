package middleware

import (
	"cert-gateway/cert/internal/config"
	"cert-gateway/cert/internal/errs"
	"cert-gateway/utils"
	xhttp "github.com/zeromicro/x/http"
	"net/http"
	"strings"
)

type AuthorizationMiddleware struct {
	Config *config.Config
}

func NewAuthorizationMiddleware(c *config.Config) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		Config: c,
	}
}

func (m *AuthorizationMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := m.Config
		token := r.Header.Get("Authorization")
		_, after, _ := strings.Cut(token, "Bearer ")
		if after != "" {
			token = after
		}
		if token != m.Config.Secret {
			err := utils.ParseJwt(token, c.JWTSecret, c.Secret, c.User.Pass, c.User.Name)
			if err != nil {
				xhttp.JsonBaseResponseCtx(r.Context(), w, errs.UnAuthorizationException)
				return
			}
		}
		next(w, r)
	}
}
