
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

curl -i -X POST http://localhost:8080/v1/login   -H "Content-Type: application/json"   -d '{"uid":"UID","token":"12345678"}'

#### Refresh

curl -i -X POST http://localhost:8080/v1/auth/refresh \
  -H "Cookie: refresh_token=eyJhbGciOi..."

### Load Manager
#### create fixed load
curl -i -X POST http://localhost:8080/v1/http-load/fixed \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.XXX.YYY" \
  -d '{
    "job_name": "fixed_load_test",
    "type": "fixed",
    "start_time": "2026-01-04T15:00:00Z",
    "request_count": "1000",
    "payload": {
      "method": "GET",
      "url": "https://example.com/api/test",
      "headers": "Authorization: Bearer test",
      "body": "{}"
    }
  }'
