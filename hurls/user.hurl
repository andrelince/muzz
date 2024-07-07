# create user
POST http://localhost:3000/user/create
{
 "email": "a@a.com",
 "password": "pword",
 "name": "a",
 "gender": "M",
 "dob": "2000-01-01" 
}
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.email" == "a@a.com"
jsonpath "$.password" == "pword"
jsonpath "$.name" == "a"
jsonpath "$.age" == 24