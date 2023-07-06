--
-- PostgreSQL database cluster dump
-- This dump was exported using the following command:
-- pg_dumpall --no-role-password  -h SOME_HOST -U postgres --exclude-database=rdsadmin -f aws_aurora_pg14.sql

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
--
-- User Configurations
--

--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET "TimeZone" TO 'utc';
--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET log_statement TO 'all';
--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET log_min_error_statement TO 'debug5';
--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET log_min_messages TO 'panic';
--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET exit_on_error TO '0';
--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET statement_timeout TO '0';
--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET role TO 'rdsadmin';
--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET "auto_explain.log_min_duration" TO '-1';
--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET temp_file_limit TO '-1';
--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET search_path TO 'pg_catalog', 'public';
--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET "pg_hint_plan.enable_hint" TO 'off';
--
-- User Config "rdsadmin"
--

ALTER ROLE rdsadmin SET default_transaction_read_only TO 'off';


--
-- Role memberships
--

GRANT pg_monitor TO rds_superuser WITH ADMIN OPTION GRANTED BY rdsadmin;
GRANT pg_signal_backend TO rds_superuser WITH ADMIN OPTION GRANTED BY rdsadmin;
GRANT rds_password TO rds_superuser WITH ADMIN OPTION GRANTED BY rdsadmin;
GRANT rds_replication TO rds_superuser WITH ADMIN OPTION GRANTED BY rdsadmin;
GRANT rds_superuser TO postgres GRANTED BY rdsadmin;




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

-- Dumped from database version 14.6
-- Dumped by pg_dump version 14.8 (Ubuntu 14.8-0ubuntu0.22.04.1)

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
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM rdsadmin;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


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

-- Dumped from database version 14.6
-- Dumped by pg_dump version 14.8 (Ubuntu 14.8-0ubuntu0.22.04.1)

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
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM rdsadmin;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database cluster dump complete
--

