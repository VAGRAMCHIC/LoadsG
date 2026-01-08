
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
#### create fixed load
curl -i -X POST http://localhost:8080/v1/manager/http/fixed \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer zoElsCHl09-QBhbn00Z09pVncaCRCladl8bjPH-NPWc" \
  -d '{
    "job_name": "load_test", 
    "type": "fixed",
    "start_time": "2026-01-04T15:00:00Z",
    "rps": "1000", "duration":"10.0",
    "payload": {
      "method": "GET",
      "url": "https://example.com/api/test",
      "headers": "Authorization: Bearer test",
      "body": "{}"
    }
  }'
  
#### create ramp up load
curl -i -X POST http://localhost:8080/v1/manager/http/ramp-up \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer zoElsCHl09-QBhbn00Z09pVncaCRCladl8bjPH-NPW" \
        -d '{
      "job_name": "load_ramp", 
      "type": "ramp-up",
      "start_time": "2026-01-04T15:00:00Z",
      "rps_s": "10", "rps_f": "100", "duration":"20",
      "payload": {
        "method": "GET",
        "url": "https://example.com/api/test",
        "headers": "Authorization: Bearer test",
        "body": "{}"
      }
    }'

#### Delete fixed load
curl -i -X DELETE http://localhost:8080/v1/manager/http/fixed/cbdd89c0-52d9-4c7f-b5bb-2b0045a01261 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.XXX.YYY"
