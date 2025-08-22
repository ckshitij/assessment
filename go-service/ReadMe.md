# Go-Service - Student Report Generation

Go service which will generate the Student Report in PDF format.

## 🚀 Quick Start

### Prerequisites
- go 1.24.2
- set up the `backend` service to get the user information for report generation.
- use same auth mechanism or cookie `backend` code.
- `make` for easy access

### Installation & Setup

- Use Make

```bash
make start
```
- Without Make
```bash
# Install dependencies
go mod tidy

# Start server
go run cmd/server/main.go

# Build binary
go build -o report_srv cmd/server/main.go
```

### Config

- For Demo using the similar config used in the `backend` service.
```yaml
# current surver port
server:
  port: 5008

# backend service config
# using the same demo account using the login cred to use login flow.
backend:
  baseURL: "http://localhost:5007"
  username: "admin@school-admin.com"
  password: "3OU4zn3q6Zh9"

```


