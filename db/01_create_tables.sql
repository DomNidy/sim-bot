--
-- PostgreSQL database cluster dump
--

-- Started on 2024-09-30 15:07:50 UTC

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Roles
--

CREATE ROLE saint;
ALTER ROLE saint WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS;

--
-- User Configurations
--






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

-- Dumped from database version 16.4 (Debian 16.4-1.pgdg120+1)
-- Dumped by pg_dump version 16.4

-- Started on 2024-09-30 15:07:50 UTC

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

-- Completed on 2024-09-30 15:07:50 UTC

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

-- Dumped from database version 16.4 (Debian 16.4-1.pgdg120+1)
-- Dumped by pg_dump version 16.4

-- Started on 2024-09-30 15:07:50 UTC

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

-- Completed on 2024-09-30 15:07:50 UTC

--
-- PostgreSQL database dump complete
--

--
-- Database "saint" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.4 (Debian 16.4-1.pgdg120+1)
-- Dumped by pg_dump version 16.4

-- Started on 2024-09-30 15:07:50 UTC

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
-- TOC entry 3368 (class 1262 OID 16384)
-- Name: saint; Type: DATABASE; Schema: -; Owner: -
--

CREATE DATABASE saint WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


\connect saint

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
-- TOC entry 219 (class 1255 OID 17605)
-- Name: notify_sim_result(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.notify_sim_result() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
	PERFORM pg_notify('new_sim_result', NEW.id::text);
	RETURN NEW;
END;
$$;


--
-- TOC entry 220 (class 1255 OID 17616)
-- Name: notify_simulation_data(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.notify_simulation_data() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
	PERFORM pg_notify('new_simulation_data', NEW.from_request::text);
	RETURN NEW;
END;
$$;


--
-- TOC entry 216 (class 1259 OID 16441)
-- Name: simulation_data_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.simulation_data_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


SET default_table_access_method = heap;

--
-- TOC entry 218 (class 1259 OID 17570)
-- Name: simulation_data; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.simulation_data (
    id integer DEFAULT nextval('public.simulation_data_id_seq'::regclass) NOT NULL,
    from_request uuid NOT NULL,
    sim_result text NOT NULL
);


--
-- TOC entry 217 (class 1259 OID 17563)
-- Name: simulation_request; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.simulation_request (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    received_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    options jsonb NOT NULL
);


--
-- TOC entry 215 (class 1259 OID 16385)
-- Name: simulaton_data_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.simulaton_data_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3216 (class 2606 OID 17577)
-- Name: simulation_data simulation_data_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.simulation_data
    ADD CONSTRAINT simulation_data_pkey PRIMARY KEY (id);


--
-- TOC entry 3214 (class 2606 OID 17569)
-- Name: simulation_request simulation_request_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.simulation_request
    ADD CONSTRAINT simulation_request_pkey PRIMARY KEY (id);


--
-- TOC entry 3218 (class 2620 OID 17606)
-- Name: simulation_data new_sim_result; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER new_sim_result AFTER INSERT ON public.simulation_data FOR EACH ROW EXECUTE FUNCTION public.notify_sim_result();


--
-- TOC entry 3219 (class 2620 OID 17617)
-- Name: simulation_data new_simulation_data; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER new_simulation_data AFTER INSERT ON public.simulation_data FOR EACH ROW EXECUTE FUNCTION public.notify_simulation_data();


--
-- TOC entry 3217 (class 2606 OID 17578)
-- Name: simulation_data simulation_data_from_request_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.simulation_data
    ADD CONSTRAINT simulation_data_from_request_fkey FOREIGN KEY (from_request) REFERENCES public.simulation_request(id) ON DELETE CASCADE;


-- Completed on 2024-09-30 15:07:50 UTC

--
-- PostgreSQL database dump complete
--

-- Completed on 2024-09-30 15:07:50 UTC

--
-- PostgreSQL database cluster dump complete
--

