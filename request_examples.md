
## Registration

curl -X POST http://localhost:8080/v1/register   -H "Content-Type: application/json"   -d '{"comment":"admin","token":"12345678"}'

curl -X POST http://localhost:8080/v1/register   -H "Content-Type: application/json"   -d '{"token":"12345678"}'


## Login

curl -X POST http://localhost:8080/v1/login   -H "Content-Type: application/json"   -d '{"uid":"UID","token":"12345678"}'