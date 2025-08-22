# Go-Service - Student Report Generation

Go service which will generate the Student Report in PDF format.
Not explicit auth added, internally utilizing the `/api/v1/auth/login` with demo user and password mentioned in the config by the `backend` service.

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

```

### API call using curl utility

- Login using the demo user mentioned in `backend` service and store the required cookies
```sh
curl -X POST http://localhost:5008/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin@school-admin.com","password":"3OU4zn3q6Zh9"}' \
  -c cookies.txt
```

- Use the cookie and get student report for a given ID(2)
```sh
curl -X GET http://localhost:5008/api/v1/students/2/report -b cookies.txt -o report.pdf
```

- Use the cookie and get student details for a given ID
```sh
curl -X GET http://localhost:5008/api/v1/students/2 -b cookies.txt
```
