--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.11
-- Dumped by pg_dump version 10.5

-- Started on 2019-03-13 14:19:34 UTC

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 1 (class 3079 OID 12393)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;

SET default_with_oids = false;

--
-- TOC entry 185 (class 1259 OID 18344)
-- Name: picsure_query; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.picsure_query (
    uuid uuid NOT NULL,
    metadata bytea,
    query text,
    readytime date,
    resourceresultid character varying(255),
    starttime date,
    status integer,
    resourceid uuid
);


--
-- TOC entry 186 (class 1259 OID 18352)
-- Name: picsure_resource; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.picsure_resource (
    uuid uuid NOT NULL,
    description character varying(8192),
    name character varying(255),
    resourcerspath character varying(255),
    targeturl character varying(255),
    token character varying(8192)
);


--
-- TOC entry 187 (class 1259 OID 18360)
-- Name: picsure_user; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.picsure_user (
    uuid uuid NOT NULL,
    roles character varying(255),
    subject character varying(255),
    userid character varying(255)
);


--
-- TOC entry 2015 (class 2606 OID 18351)
-- Name: picsure_query picsure_query_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.picsure_query
    ADD CONSTRAINT picsure_query_pkey PRIMARY KEY (uuid);


--
-- TOC entry 2017 (class 2606 OID 18359)
-- Name: picsure_resource picsure_resource_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.picsure_resource
    ADD CONSTRAINT picsure_resource_pkey PRIMARY KEY (uuid);


--
-- TOC entry 2019 (class 2606 OID 18367)
-- Name: picsure_user picsure_user_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.picsure_user
    ADD CONSTRAINT picsure_user_pkey PRIMARY KEY (uuid);


--
-- TOC entry 2021 (class 2606 OID 18371)
-- Name: picsure_user uk_il7dyn2pdljx2rcedtav6goth; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.picsure_user
    ADD CONSTRAINT uk_il7dyn2pdljx2rcedtav6goth UNIQUE (userid);


--
-- TOC entry 2023 (class 2606 OID 18369)
-- Name: picsure_user uk_s0ineuvin1jbw7fecsn3wokme; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.picsure_user
    ADD CONSTRAINT uk_s0ineuvin1jbw7fecsn3wokme UNIQUE (subject);


--
-- TOC entry 2024 (class 2606 OID 18372)
-- Name: picsure_query fkcfdw400825df2fjoxoj1erdc1; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.picsure_query
    ADD CONSTRAINT fkcfdw400825df2fjoxoj1erdc1 FOREIGN KEY (resourceid) REFERENCES public.picsure_resource(uuid);


-- Completed on 2019-03-13 14:19:34 UTC

--
-- PostgreSQL database dump complete
--

