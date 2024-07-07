# login
POST http://localhost:3000/login
{
 "email": "a@a.com",
 "password": "pword"
}
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.expires" > 0
jsonpath "$.token" matches /^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$/
