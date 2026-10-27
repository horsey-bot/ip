# ip

A minimal self-hosted client-IP service.

## Endpoints

- `/` — small HTML page
- `/text` — client IP followed by a newline
- `/json` — `{ "ip": "..." }`

The service has no JavaScript, database, GeoIP dependency, cookies, analytics, or external requests.

## Run

```sh
go run .
```

The server listens on `:8080`.

## Reverse proxy

Put it behind Caddy and keep port 8080 private:

```caddy
ip.example.com {
    reverse_proxy 127.0.0.1:8080
}
```

The application uses `X-Forwarded-For`, which is appropriate when the service is reachable only through the trusted reverse proxy.
