CREATE USER muzz WITH PASSWORD 'api';
CREATE DATABASE api OWNER muzz;
GRANT ALL PRIVILEGES ON DATABASE api TO muzz;
REVOKE ALL ON DATABASE api FROM PUBLIC;