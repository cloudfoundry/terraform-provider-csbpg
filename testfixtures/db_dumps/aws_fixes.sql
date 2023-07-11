-- Dumps don't include any passwords. Let's assign a test one to the admin user
ALTER USER testuser WITH PASSWORD 'test-password';

-- Dumps don't seem to be replicating original permissions for some objects.
-- This is something PostgresDB should be able to do automatically.
-- Until we find how to do this properly, lets hardcode some manual fixes.
\connect postgres
ALTER SCHEMA PUBLIC OWNER TO testuser;

\connect template1
ALTER SCHEMA PUBLIC OWNER TO testuser;

\connect testdb
ALTER SCHEMA PUBLIC OWNER TO testuser;

