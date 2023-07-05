--
-- PostgreSQL database cluster dump
-- This dump was exported using the following command:
-- pg_dumpall --no-role-passwords -h SOME_HOST -U postgres --exclude-database=rdsadmin -f aws_pg15.sql
--

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Roles
--

CREATE ROLE postgres;
ALTER ROLE postgres WITH NOSUPERUSER INHERIT CREATEROLE CREATEDB LOGIN NOREPLICATION NOBYPASSRLS VALID UNTIL 'infinity';
CREATE ROLE rds_ad;
ALTER ROLE rds_ad WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB NOLOGIN NOREPLICATION NOBYPASSRLS;
CREATE ROLE rds_iam;
ALTER ROLE rds_iam WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB NOLOGIN NOREPLICATION NOBYPASSRLS;
CREATE ROLE rds_password;
ALTER ROLE rds_password WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB NOLOGIN NOREPLICATION NOBYPASSRLS;
CREATE ROLE rds_replication;
ALTER ROLE rds_replication WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB NOLOGIN NOREPLICATION NOBYPASSRLS;
CREATE ROLE rds_superuser;
ALTER ROLE rds_superuser WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB NOLOGIN NOREPLICATION NOBYPASSRLS;
CREATE ROLE rdsadmin;
ALTER ROLE rdsadmin WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS VALID UNTIL 'infinity';
CREATE ROLE rdsrepladmin;
ALTER ROLE rdsrepladmin WITH NOSUPERUSER NOINHERIT NOCREATEROLE NOCREATEDB NOLOGIN REPLICATION NOBYPASSRLS;
CREATE ROLE rdstopmgr;
ALTER ROLE rdstopmgr WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB LOGIN NOREPLICATION NOBYPASSRLS;

--
-- User Configurations
--

--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET "TimeZone" TO 'utc';
ALTER ROLE rdsadmin SET log_statement TO 'all';
ALTER ROLE rdsadmin SET log_min_error_statement TO 'debug5';
ALTER ROLE rdsadmin SET log_min_messages TO 'panic';
ALTER ROLE rdsadmin SET exit_on_error TO '0';
ALTER ROLE rdsadmin SET statement_timeout TO '0';
ALTER ROLE rdsadmin SET role TO 'rdsadmin';
ALTER ROLE rdsadmin SET "auto_explain.log_min_duration" TO '-1';
ALTER ROLE rdsadmin SET temp_file_limit TO '-1';
ALTER ROLE rdsadmin SET search_path TO 'pg_catalog', 'public';
ALTER ROLE rdsadmin SET synchronous_commit TO 'local';
ALTER ROLE rdsadmin SET default_tablespace TO '';
ALTER ROLE rdsadmin SET stats_fetch_consistency TO 'snapshot';
ALTER ROLE rdsadmin SET idle_session_timeout TO '0';
ALTER ROLE rdsadmin SET "pg_hint_plan.enable_hint" TO 'off';
ALTER ROLE rdsadmin SET default_transaction_read_only TO 'off';

--
-- User Config "rdstopmgr"
--

ALTER ROLE rdstopmgr SET log_statement TO 'all';
ALTER ROLE rdstopmgr SET log_min_error_statement TO 'debug5';
ALTER ROLE rdstopmgr SET log_min_messages TO 'panic';
ALTER ROLE rdstopmgr SET exit_on_error TO '0';
ALTER ROLE rdstopmgr SET statement_timeout TO '0';
ALTER ROLE rdstopmgr SET "TimeZone" TO 'utc';
ALTER ROLE rdstopmgr SET search_path TO 'pg_catalog', 'public';
ALTER ROLE rdstopmgr SET "auto_explain.log_min_duration" TO '-1';
ALTER ROLE rdstopmgr SET role TO 'rdstopmgr';
ALTER ROLE rdstopmgr SET temp_file_limit TO '-1';
ALTER ROLE rdstopmgr SET "pg_hint_plan.enable_hint" TO 'off';
ALTER ROLE rdstopmgr SET default_transaction_read_only TO 'off';
ALTER ROLE rdstopmgr SET idle_session_timeout TO '0';


--
-- Role memberships
--

GRANT pg_checkpoint TO rds_superuser WITH ADMIN OPTION GRANTED BY rdsadmin;
GRANT pg_checkpoint TO rdstopmgr GRANTED BY rdsadmin;
GRANT pg_monitor TO rds_superuser WITH ADMIN OPTION GRANTED BY rdsadmin;
GRANT pg_monitor TO rdstopmgr GRANTED BY rdsadmin;
GRANT pg_read_all_data TO rds_superuser WITH ADMIN OPTION GRANTED BY rdsadmin;
GRANT pg_signal_backend TO rds_superuser WITH ADMIN OPTION GRANTED BY rdsadmin;
GRANT pg_write_all_data TO rds_superuser WITH ADMIN OPTION GRANTED BY rdsadmin;
GRANT rds_password TO rds_superuser WITH ADMIN OPTION GRANTED BY rdsadmin;
GRANT rds_replication TO rds_superuser WITH ADMIN OPTION GRANTED BY rdsadmin;
GRANT rds_superuser TO postgres GRANTED BY rdsadmin;




--
-- Tablespaces
--

CREATE TABLESPACE rds_temp_tablespace OWNER rds_superuser LOCATION '/rdsdbdata/tmp/rds_temp_tablespace';
GRANT ALL ON TABLESPACE rds_temp_tablespace TO PUBLIC;


--
-- Databases
--

--
-- Database "template1" dump
--

\connect template1

--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3
-- Dumped by pg_dump version 15.3 (Debian 15.3-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- PostgreSQL database dump complete
--

--
-- Database "postgres" dump
--

\connect postgres

--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3
-- Dumped by pg_dump version 15.3 (Debian 15.3-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database cluster dump complete
--

