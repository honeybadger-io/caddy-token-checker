# Caddy Token Check

A Caddy module for restricting access to resources using Redis for API key storage.

## Usage

Load this module with your Caddy build, point it at your Redis server, and protect your endpoints with an API key check for every request.

### Build with Caddy

```shell
xcaddy build --with github.com/honeybadger-io/caddy-token-checker
```

### Docker

```shell
docker run --net host -v ./Caddyfile:/etc/caddy/Caddyfile -p 2020:2020 honeybadger-io/caddy-token-checker
```

## Credits

Thanks to [caddy-geofence](https://github.com/circa10a/caddy-geofence) for a clear example of how to build middleware for Caddy.
