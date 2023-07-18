ARG PG_VERSION
FROM postgres:${PG_VERSION}
HEALTHCHECK --interval=10s --timeout=1s CMD pg_isready || exit 1

# Copy dump script to the right folder so it'll get executed automatically
ARG DUMP_FILE
COPY "db_dumps/${DUMP_FILE}" /docker-entrypoint-initdb.d/1_dump.sql


ENV POSTGRES_USER="restoredump" \
    POSTGRES_PASSWORD="restoredump" \
    POSTGRES_DB="restoredump"

COPY ssl_postgres/ /mnt/
RUN chown -R postgres:postgres /mnt/ \
 && chmod 0600 /mnt/keys/server.key

# =================================================================================
# We should find better ways to restore dumps and get rid of the instructions below
# =================================================================================

# Required to be able to login with restored users
ENV POSTGRES_HOST_AUTH_METHOD=trust

COPY "db_dumps/aws_fixes.sql" /docker-entrypoint-initdb.d/2_aws_fixes.sql

# AWS dump tries to write to the following path and fails if it doesn't exist
RUN mkdir -p /rdsdbdata/tmp/rds_temp_tablespace
RUN mkdir -p /rdsdbdata/tmp/aurora_temp_tablespace
RUN chown -R postgres:postgres /rdsdbdata

