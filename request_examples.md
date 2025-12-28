
## Registration

curl -X POST http://localhost:8080/v1/register   -H "Content-Type: application/json"   -d '{"uid":"admin","password":"12345678"}'


## Login

curl -X POST http://localhost:8080/v1/login   -H "Content-Type: application/json"   -d '{"uid":"admin","password":"12345678"}'