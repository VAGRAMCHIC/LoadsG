
## Registration

curl -X POST http://localhost:8080/v1/users/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -d '{                
    "comment": "admin",
    "token": "12345678"
  }'

curl -X POST http://localhost:8080/v1/users/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -d '{                
    "token": "12345678"
  }'

## Login

curl -i -X POST http://localhost:8080/v1/auth/login   -H "Content-Type: application/json"   -d '{"uid":"UID","token":"12345678"}'

#### Refresh

curl -i -X POST http://localhost:8080/v1/auth/refresh \
  -H "Cookie: refresh_token=eyJhbGciOi..."

### Load Manager

#### Constant load

##### GET

curl -X POST http://localhost:8080/v1/manager/http/constant \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" -H "Content-Type: application/json" \
  -d '{
    "job_name": "constant-get-test",
    "type": "constant_http",
    "start_time": "2026-01-16T12:10:00Z",

    "count": 2000,
    "url": "http://target-service/api/orders?status=active",
    "method": "GET",
    "headers": {
      "Accept": "application/json"
    }
  }'

#### POST

curl -X POST http://localhost:8080/v1/manager/http/constant \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" -H "Content-Type: application/json" \
  -d '{
    "job_name": "constant-post-test-body",
    "type": "constant_http",
    "start_time": "2026-01-16T12:00:00Z",

    "count": 1000,
    "url": "http://target-service/api/orders",
    "method": "POST",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": {
      "user_id": 42,
      "amount": 100
    }
  }'

