FROM caddy:builder AS builder
RUN xcaddy build --with github.com/honeybadger-io/caddy-token-checker

FROM caddy
COPY --from=builder /usr/bin/caddy /usr/bin/caddy
