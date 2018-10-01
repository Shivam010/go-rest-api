SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

CREATE SCHEMA todolist_management;


ALTER SCHEMA todolist_management OWNER TO postgres;

CREATE SCHEMA user_management;


ALTER SCHEMA user_management OWNER TO postgres;

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

CREATE TABLE public.users (
    fname text NOT NULL,
    lname text,
    dob text,
    email text NOT NULL,
    phone_no bigint NOT NULL,
    id integer NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


CREATE TABLE todolist_management.todo_items (
    id integer NOT NULL,
    value text NOT NULL,
    list_id integer NOT NULL,
    completed boolean
);


ALTER TABLE todolist_management.todo_items OWNER TO postgres;

CREATE SEQUENCE todolist_management.todo_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE todolist_management.todo_items_id_seq OWNER TO postgres;

ALTER SEQUENCE todolist_management.todo_items_id_seq OWNED BY todolist_management.todo_items.id;


CREATE TABLE todolist_management.todo_lists (
    id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE todolist_management.todo_lists OWNER TO postgres;

CREATE SEQUENCE todolist_management.todolist_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE todolist_management.todolist_id_seq OWNER TO postgres;

ALTER SEQUENCE todolist_management.todolist_id_seq OWNED BY todolist_management.todo_lists.id;


CREATE TABLE user_management.users (
    fname text NOT NULL,
    lname text,
    dob text,
    email text NOT NULL,
    phone_no bigint NOT NULL,
    id integer DEFAULT nextval('public.users_id_seq'::regclass) NOT NULL
);


ALTER TABLE user_management.users OWNER TO postgres;

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


ALTER TABLE ONLY todolist_management.todo_items ALTER COLUMN id SET DEFAULT nextval('todolist_management.todo_items_id_seq'::regclass);


ALTER TABLE ONLY todolist_management.todo_lists ALTER COLUMN id SET DEFAULT nextval('todolist_management.todolist_id_seq'::regclass);


SELECT pg_catalog.setval('public.users_id_seq', 36, true);


SELECT pg_catalog.setval('todolist_management.todo_items_id_seq', 16, true);


SELECT pg_catalog.setval('todolist_management.todolist_id_seq', 8, true);


ALTER TABLE ONLY todolist_management.todo_items
    ADD CONSTRAINT todo_items_pkey PRIMARY KEY (id);


ALTER TABLE ONLY todolist_management.todo_lists
    ADD CONSTRAINT todolist_pkey PRIMARY KEY (id);


