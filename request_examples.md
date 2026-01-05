
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
