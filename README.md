### Search for, listen to and download sermons from yt

git clone git@github.com:miceremwirigi/journey-family-sermons.git
go mod init
go mod tidy

# Configure database
Install postgresql, then if on linux, do -> sudo -u postgres psql 
CREATE DATABASE jfk_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE USER jfk_user WITH SUPERUSER CREATEDB CREATEROLE LOGIN PASSWORD 'jfk_pass';
ALTER ROLE jfk_user SET client_encoding TO 'utf8';
ALTER ROLE jfk_user SET default_transaction_isolation TO 'read committed';
ALTER ROLE jfk_user SET timezone TO 'UTC';
GRANT ALL PRIVILEGES ON DATABASE jfk_db TO jfk_user;
CREATE DATABASE jfk_test;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
GRANT ALL PRIVILEGES ON DATABASE jfk_test TO jfk_user;
\q

# Load env
source env.sh