# create user1
POST http://localhost:3000/user/create
{
 "email": "a@a.com",
 "password": "pword",
 "name": "a",
 "gender": "M",
 "dob": "2000-01-01",
 "locationLat": 38.72583954689884,
 "locationLong": -9.245569404061271
}
HTTP 200

# create user2
POST http://localhost:3000/user/create
{
 "email": "b@b.com",
 "password": "pword",
 "name": "b",
 "gender": "F",
 "dob": "2000-01-01",
 "locationLat": 38.827564251494366,
 "locationLong": -9.456503302312559
}
HTTP 200

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
GET http://localhost:3000/discover?gender=F&min_age=15&max_age=25
Authorization: Bearer {{user1token}}
HTTP 200