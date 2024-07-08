# create user1
POST http://localhost:3000/user/create
{
 "email": "a@a.com",
 "password": "pword",
 "name": "a",
 "gender": "M",
 "dob": "2000-01-01" 
}
HTTP 200
[Captures]
user1id: jsonpath "$['id']"

# create user2
POST http://localhost:3000/user/create
{
 "email": "b@b.com",
 "password": "pword",
 "name": "b",
 "gender": "F",
 "dob": "2000-01-01" 
}
HTTP 200
[Captures]
user2id: jsonpath "$['id']"

# login user 1
POST http://localhost:3000/login
{
 "email": "a@a.com",
 "password": "pword"
}
HTTP 200
[Captures]
user1token: jsonpath "$['token']"

# discover for user 1
GET http://localhost:3000/discover
Authorization: Bearer {{user1token}}
HTTP 200