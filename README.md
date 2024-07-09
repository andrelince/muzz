# muzz

### Setup and running

`make up` -  Run the api available on `0.0.0.0:3000`

`make down` to shut down the api

The whole project is composed of an api service, a postgres database and a redis cache (currently unused).
The api container is built using `air` which is a hot-reload go docker image used solely for development

### Design choices

The project is split between 3 layers (packages) each with its own purpose. Each layer is independent of the layers below and its dependencies could be abstracted with mocks:

- rest/ 
    - Here lies all logic regarding request/response processing, payload validation and middlewares such as authentication

- service/
    - Here lies all the business logic of the application where for instance, upon performing a login  (`/login`) the user email is validated, the password is compared against the hashed password which is stored and the token is generated

- respository/
    - Here lies all data retrieval functionality where the data access is abstracted using interfaces to obscure the type of datasource. For instance `UserRepository` currenlty uses a postgres database but it could be changed to MySQL / MariaDB / MongoDB without compromising the layers above it with changes


For the `/discover` endpoint logic the following assumptions were made:
    
`attractiveness_score` is calculated based on positive swipes from a user representing the likelyhood of a future match. For ex. if a user has the tendency to positively swipe across other profiles then it is more likely that it has partial match to the current user

`distance` is calculated using the postgres `earthdistance` extension (https://www.postgresql.org/docs/current/earthdistance.html)

The returing profiles are sorted by attractiveness_score and closer distance to the current user.

### Available routes

- `/swagger`: auto generated api docs 

- `/healthz`: for checking service is healthy

- `/user/create`: for creating a profile

- `/login`: for authenticating a user

- `/swipe`: for simulating a user swipe over a profile

- `/discover`: for returing interesting profiles for a user with the following optional parameters:
    - `min_age`: number detailing minimum age for a prospective profile
    - `max_age`: number detailing minimum age for a prospective profile
    - `gender`: the profile gender (M | F)

All requests go through a layer of validation using the `https://github.com/go-playground/validator` package


### Points of improvement

- Add unit tests: due to lack of time I mostly focused on developing the features and setting only partial e2e tests using `hurl` (https://hurl.dev) available on `/hurl` folder of the repo. Unit test would provide an additional layer of safety to the source code.

- Add `/user/delete` endpoint to make each e2e test self sufficient. Currently we need to clear the db after each hurl test run as the user would fail the email validation upon creation

- Optimize discover query: The discover query would perform badly in large datasets, specially the portion where the `attractiveness_score` is calculated. This calculation could be abstracted in a db view.

- Add ci/cd pipeline to run unit tests and hurl e2e tests

- Improve `/healthz` by pinging db and cache for ensuring connections/repositories are up