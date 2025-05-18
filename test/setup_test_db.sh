#!/bin/bash

set -e

DB_NAME="sika_test"
DB_USER="sika_test_user"
DB_PASS="sika_test_pass"

echo "Setting up test database..."

psql -v ON_ERROR_STOP=1 -c "DO
\$do\$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = '$DB_USER') THEN
      CREATE USER $DB_USER WITH PASSWORD '$DB_PASS';
   END IF;
END
\$do\$;"

psql -v ON_ERROR_STOP=1 -c "SELECT 'CREATE DATABASE $DB_NAME'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME')\gexec"

psql -v ON_ERROR_STOP=1 -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"

echo "Test database setup completed successfully!" 