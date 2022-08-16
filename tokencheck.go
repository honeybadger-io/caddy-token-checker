package tokencheck

import (
	"net/http"
	"regexp"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/go-redis/redis"
)

func init() {
	caddy.RegisterModule(TokenCheck{})
	httpcaddyfile.RegisterHandlerDirective("token_check", parseCaddyfile)
}

type TokenCheck struct {
	RedisUrl       string
	RedisClient    *redis.Client
	KeyPrefix      string
	QueryParameter string
}

func (TokenCheck) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.token_check",
		New: func() caddy.Module { return new(TokenCheck) },
	}
}

func (tc *TokenCheck) Provision(ctx caddy.Context) error {
	if tc.RedisUrl == "" {
		tc.RedisUrl = "redis://localhost:6379/0"
	}

	opt, err := redis.ParseURL(tc.RedisUrl)
	if err != nil {
		return err
	}

	tc.RedisClient = redis.NewClient(opt)

	if tc.QueryParameter == "" {
		tc.QueryParameter = "token"
	}

	return nil
}

func (tc TokenCheck) Validate() error {
	return nil
}

func (tc TokenCheck) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	token := r.FormValue(tc.QueryParameter)

	// If it's not a valid token, don't even bother with redis
	matched, err := regexp.MatchString(`^\w{6,40}$`, token)
	if err != nil {
		return err
	}
	if !matched {
		return caddyhttp.Error(403, nil)
	}

	exists, err := tc.RedisClient.Exists(tc.KeyPrefix + token).Result()
	if err != nil { // Something bad happened with redis, so fail open
		return next.ServeHTTP(w, r)
	}

	if exists == 0 { // No key found
		return caddyhttp.Error(403, nil)
	} else {
		return next.ServeHTTP(w, r)
	}
}

func (tc *TokenCheck) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "key_prefix":
				if !d.NextArg() {
					return d.ArgErr()
				}

				tc.KeyPrefix = d.Val()

			case "query_parameter":
				if !d.NextArg() {
					return d.ArgErr()
				}

				if d.Val() == "" {
					return d.ArgErr()
				}

				tc.QueryParameter = d.Val()

			case "redis_url":
				if !d.NextArg() {
					return d.ArgErr()
				}

				tc.RedisUrl = d.Val()
			}
		}
	}
	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var tc TokenCheck
	err := tc.UnmarshalCaddyfile(h.Dispenser)
	return tc, err
}

var (
	_ caddy.Provisioner           = (*TokenCheck)(nil)
	_ caddy.Validator             = (*TokenCheck)(nil)
	_ caddyhttp.MiddlewareHandler = (*TokenCheck)(nil)
	_ caddyfile.Unmarshaler       = (*TokenCheck)(nil)
)
