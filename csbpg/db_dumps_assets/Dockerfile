ARG PG_VERSION
FROM postgres:${PG_VERSION}
HEALTHCHECK --interval=10s --timeout=1s CMD pg_isready || exit 1

# Copy dump script to the right folder so it'll get executed automatically
ARG DUMP_FILE
COPY "${DUMP_FILE}" /docker-entrypoint-initdb.d/1_dump.sql

# AWS dump tries to write to the following path and fails if it doesn't exist
RUN mkdir -p /rdsdbdata/tmp/rds_temp_tablespace && chown -R postgres:postgres /rdsdbdata

# Dumps don't include any passwords. Let's assign a test one to the admin user
RUN echo "ALTER USER postgres WITH PASSWORD 'test-password';"   > /docker-entrypoint-initdb.d/2_password.sql

ENV POSTGRES_USER="restoredump" \
    POSTGRES_PASSWORD="restoredump" \
    POSTGRES_DB="restoredump"

# =================================================================================
# We should find better ways to restore dumps and get rid of the instructions below
# =================================================================================

# Required to be able to login with restored users
ENV POSTGRES_HOST_AUTH_METHOD=trust

# Dumps don't seem to be replicating original permissions for some objects.
# This is something PostgresDB should be able to do this automatically.
# Until we find how to do this properly, lets hardcode some manual fixes.
RUN echo '\
ALTER DATABASE postgres OWNER TO postgres; \n\
\\connect postgres;                        \n\
ALTER SCHEMA public OWNER TO postgres;     \n' > /docker-entrypoint-initdb.d/3_fix_ownership.sql

