# Insights Setup

Everything you need to get started and running! :)

## Prerequisites

- [Go](https://golang.org/dl/)
- [Git](https://git-scm.com/)

## Clone the Repository

```sh
git clone https://github.com/winjw7/Insights.git
```

## Install Dependencies

```sh
go mod tidy
```

## Build and Run

```sh
go build -o insights
./insights
```

The server will start on port **3000**.

## API Endpoints

- `POST /api/login/new` — Store a login event
- `GET /api/login/suspicious` — Retrieve suspicious login origins (requires API key)

## Testing

Here's some example `curl` requests:

```sh
curl -X POST http://localhost:3000/api/login/new \
  -H "Content-Type: application/json" \
  -d '{"tenant":"amazon","user":"bobby","origin":"127.0.0.1","status":"failure"}'
```

```sh
curl -X GET "http://localhost:3000/api/login/suspicious?threshold=5&minutes=3" \
  -H "X-API-Key: amazon"
```