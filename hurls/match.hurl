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

# login user 2
POST http://localhost:3000/login
{
 "email": "b@b.com",
 "password": "pword"
}
HTTP 200
[Captures]
user2token: jsonpath "$['token']"

# swipe user1
POST http://localhost:3000/swipe
Authorization: Bearer {{user1token}}
Content-Type: application/json
{
 "user_id": {{user2id}},
 "preference": "yes"
}
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.matched" == false

# swipe user2
POST http://localhost:3000/swipe
Authorization: Bearer {{user2token}}
Content-Type: application/json
{
 "user_id": {{user1id}},
 "preference": "yes"
}
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.matched" == true