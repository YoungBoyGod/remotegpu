--
-- PostgreSQL database dump
--

\restrict IA5SsACGgKtsCB0nheG64BVP6kmXmVaRvKOad8JNCpaEvyHyIyfNyCuwOJ1kOab

-- Dumped from database version 17.7 (Debian 17.7-3.pgdg13+1)
-- Dumped by pg_dump version 17.7 (Debian 17.7-3.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: remotegpu_user
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO remotegpu_user;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: remotegpu_user
--

COMMENT ON SCHEMA public IS '';


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: remotegpu_user
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO remotegpu_user;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: active_alerts; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.active_alerts (
    id bigint NOT NULL,
    rule_id bigint NOT NULL,
    host_id character varying(64) NOT NULL,
    value numeric NOT NULL,
    message text,
    triggered_at timestamp with time zone,
    acknowledged boolean DEFAULT false
);


ALTER TABLE public.active_alerts OWNER TO remotegpu_user;

--
-- Name: active_alerts_id_seq; Type: SEQUENCE; Schema: public; Owner: remotegpu_user
--

CREATE SEQUENCE public.active_alerts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.active_alerts_id_seq OWNER TO remotegpu_user;

--
-- Name: active_alerts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: remotegpu_user
--

ALTER SEQUENCE public.active_alerts_id_seq OWNED BY public.active_alerts.id;


--
-- Name: alert_rules; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.alert_rules (
    id bigint NOT NULL,
    name character varying(128) NOT NULL,
    metric_type character varying(64) NOT NULL,
    threshold numeric NOT NULL,
    condition character varying(10) NOT NULL,
    severity character varying(20) DEFAULT 'warning'::character varying,
    enabled boolean DEFAULT true,
    created_at timestamp with time zone
);


ALTER TABLE public.alert_rules OWNER TO remotegpu_user;

--
-- Name: alert_rules_id_seq; Type: SEQUENCE; Schema: public; Owner: remotegpu_user
--

CREATE SEQUENCE public.alert_rules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.alert_rules_id_seq OWNER TO remotegpu_user;

--
-- Name: alert_rules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: remotegpu_user
--

ALTER SEQUENCE public.alert_rules_id_seq OWNED BY public.alert_rules.id;


--
-- Name: allocations; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.allocations (
    id character varying(64) NOT NULL,
    customer_id bigint NOT NULL,
    host_id character varying(64) NOT NULL,
    workspace_id bigint,
    start_time timestamp with time zone NOT NULL,
    end_time timestamp with time zone NOT NULL,
    actual_end_time timestamp with time zone,
    status character varying(32) DEFAULT 'active'::character varying,
    remark text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.allocations OWNER TO remotegpu_user;

--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.audit_logs (
    id bigint NOT NULL,
    customer_id bigint,
    username character varying(128),
    ip_address character varying(64),
    method character varying(10),
    path character varying(512),
    action character varying(128) NOT NULL,
    resource_type character varying(64),
    resource_id character varying(128),
    detail jsonb,
    status_code bigint,
    created_at timestamp with time zone
);


ALTER TABLE public.audit_logs OWNER TO remotegpu_user;

--
-- Name: audit_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: remotegpu_user
--

CREATE SEQUENCE public.audit_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.audit_logs_id_seq OWNER TO remotegpu_user;

--
-- Name: audit_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: remotegpu_user
--

ALTER SEQUENCE public.audit_logs_id_seq OWNED BY public.audit_logs.id;


--
-- Name: customers; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.customers (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    username character varying(64) NOT NULL,
    email character varying(128) NOT NULL,
    password_hash character varying(256) NOT NULL,
    display_name character varying(128),
    full_name character varying(256),
    company character varying(256),
    phone character varying(32),
    avatar_url character varying(512),
    role character varying(32) DEFAULT 'customer_owner'::character varying,
    user_type character varying(32) DEFAULT 'external'::character varying,
    account_type character varying(32) DEFAULT 'individual'::character varying,
    status character varying(32) DEFAULT 'active'::character varying,
    email_verified boolean DEFAULT false,
    phone_verified boolean DEFAULT false,
    balance numeric(10,4) DEFAULT 0,
    currency character varying(10) DEFAULT 'CNY'::character varying,
    last_login_at timestamp with time zone,
    must_change_password boolean DEFAULT false,
    company_code character varying(64)
);


ALTER TABLE public.customers OWNER TO remotegpu_user;

--
-- Name: customers_id_seq; Type: SEQUENCE; Schema: public; Owner: remotegpu_user
--

CREATE SEQUENCE public.customers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.customers_id_seq OWNER TO remotegpu_user;

--
-- Name: customers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: remotegpu_user
--

ALTER SEQUENCE public.customers_id_seq OWNED BY public.customers.id;


--
-- Name: dataset_mounts; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.dataset_mounts (
    id bigint NOT NULL,
    dataset_id bigint NOT NULL,
    host_id character varying(64) NOT NULL,
    mount_path character varying(256) NOT NULL,
    read_only boolean DEFAULT true,
    status character varying(20) DEFAULT 'mounting'::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.dataset_mounts OWNER TO remotegpu_user;

--
-- Name: dataset_mounts_id_seq; Type: SEQUENCE; Schema: public; Owner: remotegpu_user
--

CREATE SEQUENCE public.dataset_mounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.dataset_mounts_id_seq OWNER TO remotegpu_user;

--
-- Name: dataset_mounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: remotegpu_user
--

ALTER SEQUENCE public.dataset_mounts_id_seq OWNED BY public.dataset_mounts.id;


--
-- Name: datasets; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.datasets (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    customer_id bigint NOT NULL,
    workspace_id bigint,
    name character varying(256) NOT NULL,
    description text,
    storage_path character varying(512) NOT NULL,
    storage_type character varying(32) DEFAULT 'minio'::character varying,
    total_size bigint DEFAULT 0,
    file_count bigint DEFAULT 0,
    status character varying(20) DEFAULT 'uploading'::character varying,
    visibility character varying(20) DEFAULT 'private'::character varying
);


ALTER TABLE public.datasets OWNER TO remotegpu_user;

--
-- Name: datasets_id_seq; Type: SEQUENCE; Schema: public; Owner: remotegpu_user
--

CREATE SEQUENCE public.datasets_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.datasets_id_seq OWNER TO remotegpu_user;

--
-- Name: datasets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: remotegpu_user
--

ALTER SEQUENCE public.datasets_id_seq OWNED BY public.datasets.id;


--
-- Name: gpus; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.gpus (
    id bigint NOT NULL,
    host_id character varying(64) NOT NULL,
    index bigint NOT NULL,
    uuid character varying(128),
    name character varying(128) NOT NULL,
    memory_total_mb bigint NOT NULL,
    brand character varying(64),
    status character varying(20) DEFAULT 'available'::character varying,
    health_status character varying(20) DEFAULT 'healthy'::character varying,
    allocated_to character varying(64),
    updated_at timestamp with time zone
);


ALTER TABLE public.gpus OWNER TO remotegpu_user;

--
-- Name: gpus_id_seq; Type: SEQUENCE; Schema: public; Owner: remotegpu_user
--

CREATE SEQUENCE public.gpus_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.gpus_id_seq OWNER TO remotegpu_user;

--
-- Name: gpus_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: remotegpu_user
--

ALTER SEQUENCE public.gpus_id_seq OWNED BY public.gpus.id;


--
-- Name: hosts; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.hosts (
    id character varying(64) NOT NULL,
    name character varying(128) NOT NULL,
    hostname character varying(256),
    region character varying(64) DEFAULT 'default'::character varying,
    ip_address character varying(64) NOT NULL,
    public_ip character varying(64),
    ssh_port bigint DEFAULT 22,
    agent_port bigint DEFAULT 8080,
    os_type character varying(20) DEFAULT 'linux'::character varying,
    os_version character varying(64),
    cpu_info character varying(256),
    total_cpu bigint NOT NULL,
    total_memory_gb bigint NOT NULL,
    total_disk_gb bigint,
    status character varying(20) DEFAULT 'offline'::character varying,
    health_status character varying(20) DEFAULT 'unknown'::character varying,
    deployment_mode character varying(20) DEFAULT 'traditional'::character varying,
    last_heartbeat timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    ssh_username character varying(128),
    ssh_password text,
    ssh_key text,
    needs_collect boolean DEFAULT false
);


ALTER TABLE public.hosts OWNER TO remotegpu_user;

--
-- Name: images; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.images (
    id bigint NOT NULL,
    name character varying(256) NOT NULL,
    display_name character varying(256),
    description text,
    category character varying(64),
    framework character varying(64),
    cuda_version character varying(32),
    registry_url character varying(512),
    is_official boolean DEFAULT false,
    customer_id bigint,
    status character varying(20) DEFAULT 'active'::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.images OWNER TO remotegpu_user;

--
-- Name: images_id_seq; Type: SEQUENCE; Schema: public; Owner: remotegpu_user
--

CREATE SEQUENCE public.images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.images_id_seq OWNER TO remotegpu_user;

--
-- Name: images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: remotegpu_user
--

ALTER SEQUENCE public.images_id_seq OWNED BY public.images.id;


--
-- Name: machine_enrollments; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.machine_enrollments (
    id bigint NOT NULL,
    customer_id bigint,
    name character varying(128),
    hostname character varying(256),
    region character varying(64),
    address character varying(256) NOT NULL,
    ssh_port bigint DEFAULT 22,
    ssh_username character varying(128),
    ssh_password text,
    ssh_key text,
    status character varying(20) DEFAULT 'pending'::character varying,
    error_message text,
    host_id character varying(64),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.machine_enrollments OWNER TO remotegpu_user;

--
-- Name: TABLE machine_enrollments; Type: COMMENT; Schema: public; Owner: remotegpu_user
--

COMMENT ON TABLE public.machine_enrollments IS '用户机器添加任务表';


--
-- Name: COLUMN machine_enrollments.status; Type: COMMENT; Schema: public; Owner: remotegpu_user
--

COMMENT ON COLUMN public.machine_enrollments.status IS '状态: pending, success, failed';


--
-- Name: machine_enrollments_id_seq; Type: SEQUENCE; Schema: public; Owner: remotegpu_user
--

CREATE SEQUENCE public.machine_enrollments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.machine_enrollments_id_seq OWNER TO remotegpu_user;

--
-- Name: machine_enrollments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: remotegpu_user
--

ALTER SEQUENCE public.machine_enrollments_id_seq OWNED BY public.machine_enrollments.id;


--
-- Name: ssh_keys; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.ssh_keys (
    id bigint NOT NULL,
    customer_id bigint NOT NULL,
    name character varying(64) NOT NULL,
    fingerprint character varying(128),
    public_key text NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.ssh_keys OWNER TO remotegpu_user;

--
-- Name: ssh_keys_id_seq; Type: SEQUENCE; Schema: public; Owner: remotegpu_user
--

CREATE SEQUENCE public.ssh_keys_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ssh_keys_id_seq OWNER TO remotegpu_user;

--
-- Name: ssh_keys_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: remotegpu_user
--

ALTER SEQUENCE public.ssh_keys_id_seq OWNED BY public.ssh_keys.id;


--
-- Name: tasks; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.tasks (
    id character varying(64) NOT NULL,
    customer_id bigint NOT NULL,
    host_id character varying(64),
    name character varying(256) NOT NULL,
    type character varying(32) DEFAULT 'shell'::character varying NOT NULL,
    image_id bigint,
    command text NOT NULL,
    env_vars jsonb,
    status character varying(20) DEFAULT 'pending'::character varying,
    exit_code bigint,
    error_msg text,
    started_at timestamp with time zone,
    finished_at timestamp with time zone,
    created_at timestamp with time zone,
    process_id bigint DEFAULT 0,
    args jsonb,
    work_dir character varying(500),
    timeout bigint DEFAULT 3600,
    priority bigint DEFAULT 5,
    retry_count bigint DEFAULT 0,
    retry_delay bigint DEFAULT 60,
    max_retries bigint DEFAULT 3,
    machine_id character varying(64),
    group_id character varying(64),
    parent_id character varying(64),
    assigned_agent_id character varying(64),
    lease_expires_at timestamp with time zone,
    attempt_id character varying(64),
    assigned_at timestamp with time zone,
    ended_at timestamp with time zone
);


ALTER TABLE public.tasks OWNER TO remotegpu_user;

--
-- Name: workspaces; Type: TABLE; Schema: public; Owner: remotegpu_user
--

CREATE TABLE public.workspaces (
    id bigint NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    owner_id bigint NOT NULL,
    name character varying(128) NOT NULL,
    description text,
    type character varying(32) DEFAULT 'personal'::character varying,
    status character varying(32) DEFAULT 'active'::character varying,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.workspaces OWNER TO remotegpu_user;

--
-- Name: workspaces_id_seq; Type: SEQUENCE; Schema: public; Owner: remotegpu_user
--

CREATE SEQUENCE public.workspaces_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.workspaces_id_seq OWNER TO remotegpu_user;

--
-- Name: workspaces_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: remotegpu_user
--

ALTER SEQUENCE public.workspaces_id_seq OWNED BY public.workspaces.id;


--
-- Name: active_alerts id; Type: DEFAULT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.active_alerts ALTER COLUMN id SET DEFAULT nextval('public.active_alerts_id_seq'::regclass);


--
-- Name: alert_rules id; Type: DEFAULT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.alert_rules ALTER COLUMN id SET DEFAULT nextval('public.alert_rules_id_seq'::regclass);


--
-- Name: audit_logs id; Type: DEFAULT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.audit_logs ALTER COLUMN id SET DEFAULT nextval('public.audit_logs_id_seq'::regclass);


--
-- Name: customers id; Type: DEFAULT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.customers ALTER COLUMN id SET DEFAULT nextval('public.customers_id_seq'::regclass);


--
-- Name: dataset_mounts id; Type: DEFAULT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.dataset_mounts ALTER COLUMN id SET DEFAULT nextval('public.dataset_mounts_id_seq'::regclass);


--
-- Name: datasets id; Type: DEFAULT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.datasets ALTER COLUMN id SET DEFAULT nextval('public.datasets_id_seq'::regclass);


--
-- Name: gpus id; Type: DEFAULT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.gpus ALTER COLUMN id SET DEFAULT nextval('public.gpus_id_seq'::regclass);


--
-- Name: images id; Type: DEFAULT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.images ALTER COLUMN id SET DEFAULT nextval('public.images_id_seq'::regclass);


--
-- Name: machine_enrollments id; Type: DEFAULT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.machine_enrollments ALTER COLUMN id SET DEFAULT nextval('public.machine_enrollments_id_seq'::regclass);


--
-- Name: ssh_keys id; Type: DEFAULT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.ssh_keys ALTER COLUMN id SET DEFAULT nextval('public.ssh_keys_id_seq'::regclass);


--
-- Name: workspaces id; Type: DEFAULT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.workspaces ALTER COLUMN id SET DEFAULT nextval('public.workspaces_id_seq'::regclass);


--
-- Data for Name: active_alerts; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.active_alerts (id, rule_id, host_id, value, message, triggered_at, acknowledged) FROM stdin;
\.


--
-- Data for Name: alert_rules; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.alert_rules (id, name, metric_type, threshold, condition, severity, enabled, created_at) FROM stdin;
\.


--
-- Data for Name: allocations; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.allocations (id, customer_id, host_id, workspace_id, start_time, end_time, actual_end_time, status, remark, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: audit_logs; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.audit_logs (id, customer_id, username, ip_address, method, path, action, resource_type, resource_id, detail, status_code, created_at) FROM stdin;
1	1	admin	127.0.0.1	POST	/api/v1/admin/machines	create	machine		\N	200	2026-02-05 09:11:41.009506+00
2	1	admin	127.0.0.1	POST	/api/v1/admin/machines	create	machine		\N	200	2026-02-05 09:11:42.801934+00
3	1	admin	127.0.0.1	POST	/api/v1/admin/machines	create	machine		\N	200	2026-02-05 09:12:08.154646+00
4	1	admin	127.0.0.1	POST	/api/v1/admin/machines	create	machine		\N	200	2026-02-05 10:29:16.014897+00
5	1	admin	::1	POST	/api/v1/admin/machines	create	machine		\N	200	2026-02-05 10:54:23.604991+00
6	1	admin	::1	POST	/api/v1/admin/machines	create	machine		\N	200	2026-02-05 11:41:43.16782+00
7	1	admin	::1	POST	/api/v1/admin/machines	create	machine		\N	200	2026-02-05 11:42:08.201646+00
8	1	admin	127.0.0.1	POST	/api/v1/admin/machines//allocate	create	machine		\N	200	2026-02-05 14:37:17.809485+00
9	1	admin	127.0.0.1	POST	/api/v1/admin/customers	create	customer		\N	200	2026-02-06 06:57:49.551476+00
10	1	admin	127.0.0.1	POST	/api/v1/admin/customers	create	customer		\N	200	2026-02-06 06:59:58.630285+00
11	1	admin	127.0.0.1	POST	/api/v1/admin/machines/test-gpu-01/allocate	create	machine	test-gpu-01	\N	200	2026-02-06 09:27:06.716071+00
12	1	admin	127.0.0.1	POST	/api/v1/admin/machines/192.168.10.210/allocate	create	machine	192.168.10.210	\N	200	2026-02-06 09:27:06.980648+00
13	1	admin	127.0.0.1	POST	/api/v1/admin/machines/test-gpu-01/allocate	create	machine	test-gpu-01	\N	200	2026-02-06 09:35:56.586936+00
14	1	admin	127.0.0.1	POST	/api/v1/admin/machines/192.168.10.210/allocate	create	machine	192.168.10.210	\N	200	2026-02-06 09:35:56.586936+00
15	1	admin	127.0.0.1	POST	/api/v1/admin/machines/test-gpu-01/allocate	create	machine	test-gpu-01	\N	200	2026-02-06 09:36:02.710487+00
\.


--
-- Data for Name: customers; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.customers (id, created_at, updated_at, deleted_at, uuid, username, email, password_hash, display_name, full_name, company, phone, avatar_url, role, user_type, account_type, status, email_verified, phone_verified, balance, currency, last_login_at, must_change_password, company_code) FROM stdin;
2	2026-02-06 06:57:49.515095+00	2026-02-06 06:57:49.515095+00	\N	72db317f-39cb-4764-86f1-efa22bf174f1	test	luoyangtest@163.com	$2a$10$2NzHQ9IObUddCfkw5lUke.6jFKWly4pfmQN2VroFI6zflJ8.MQEYe			测试	18362493779		customer_owner	external	individual	active	f	f	0.0000	CNY	\N	f	test
3	2026-02-06 06:59:58.621604+00	2026-02-06 06:59:58.621604+00	\N	b91f9e5d-2e87-4bfe-9134-4f5038a705d1	dev	dev@dev.com	$2a$10$c9YFGz8BLPZi9..Oq9bUvuq9pWig/bS7hUiSchbYXrFqrLy.n5AN.			dev			customer_owner	external	individual	active	f	f	0.0000	CNY	\N	f	dev
1	2026-02-03 11:22:44.886294+00	2026-02-06 09:26:37.034546+00	\N	2709d78b-cb2b-4190-a23a-82afa00f6e8c	admin	admin@remotegpu.com	$2a$10$VU.r.gcsoKCXX2sW5Us0DuB0l3eCdakh3S4uP0KFm0HcZifKxpL/K	Administrator					admin	admin	individual	active	f	f	0.0000	CNY	2026-02-06 09:26:37.034205+00	f	\N
\.


--
-- Data for Name: dataset_mounts; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.dataset_mounts (id, dataset_id, host_id, mount_path, read_only, status, created_at) FROM stdin;
\.


--
-- Data for Name: datasets; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.datasets (id, created_at, updated_at, deleted_at, uuid, customer_id, workspace_id, name, description, storage_path, storage_type, total_size, file_count, status, visibility) FROM stdin;
\.


--
-- Data for Name: gpus; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.gpus (id, host_id, index, uuid, name, memory_total_mb, brand, status, health_status, allocated_to, updated_at) FROM stdin;
\.


--
-- Data for Name: hosts; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.hosts (id, name, hostname, region, ip_address, public_ip, ssh_port, agent_port, os_type, os_version, cpu_info, total_cpu, total_memory_gb, total_disk_gb, status, health_status, deployment_mode, last_heartbeat, created_at, updated_at, ssh_username, ssh_password, ssh_key, needs_collect) FROM stdin;
	192.168.10.210		南京			22	8080	linux			0	0	0	online	unknown	traditional	\N	2026-02-05 07:47:03.722456+00	2026-02-05 07:47:03.722456+00	\N	\N	\N	f
192.168.10.210	192.168.10.210		南京	192.168.10.210		22	8080	linux			0	0	0	offline	unknown	traditional	\N	2026-02-05 10:29:16.002239+00	2026-02-05 10:29:16.002239+00	luo	luoyang@767831939		f
test-gpu-01	test-gpu-01	test-gpu-01	default	192.168.1.100		22	8080	linux			0	0	0	offline	unknown	traditional	\N	2026-02-05 10:54:23.597144+00	2026-02-05 10:54:23.597144+00	root			f
gpu-201	gpu-201	gpu-201	default	192.168.10.201		22	8080	linux	6.8.0-90-generic	8 cores	8	3	292	idle	healthy	traditional	\N	2026-02-05 11:41:43.16347+00	2026-02-05 11:41:43.16347+00	luo	luo		f
\.


--
-- Data for Name: images; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.images (id, name, display_name, description, category, framework, cuda_version, registry_url, is_official, customer_id, status, created_at) FROM stdin;
\.


--
-- Data for Name: machine_enrollments; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.machine_enrollments (id, customer_id, name, hostname, region, address, ssh_port, ssh_username, ssh_password, ssh_key, status, error_message, host_id, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: ssh_keys; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.ssh_keys (id, customer_id, name, fingerprint, public_key, created_at) FROM stdin;
\.


--
-- Data for Name: tasks; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.tasks (id, customer_id, host_id, name, type, image_id, command, env_vars, status, exit_code, error_msg, started_at, finished_at, created_at, process_id, args, work_dir, timeout, priority, retry_count, retry_delay, max_retries, machine_id, group_id, parent_id, assigned_agent_id, lease_expires_at, attempt_id, assigned_at, ended_at) FROM stdin;
\.


--
-- Data for Name: workspaces; Type: TABLE DATA; Schema: public; Owner: remotegpu_user
--

COPY public.workspaces (id, uuid, owner_id, name, description, type, status, created_at, updated_at) FROM stdin;
\.


--
-- Name: active_alerts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: remotegpu_user
--

SELECT pg_catalog.setval('public.active_alerts_id_seq', 1, false);


--
-- Name: alert_rules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: remotegpu_user
--

SELECT pg_catalog.setval('public.alert_rules_id_seq', 1, false);


--
-- Name: audit_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: remotegpu_user
--

SELECT pg_catalog.setval('public.audit_logs_id_seq', 15, true);


--
-- Name: customers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: remotegpu_user
--

SELECT pg_catalog.setval('public.customers_id_seq', 3, true);


--
-- Name: dataset_mounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: remotegpu_user
--

SELECT pg_catalog.setval('public.dataset_mounts_id_seq', 1, false);


--
-- Name: datasets_id_seq; Type: SEQUENCE SET; Schema: public; Owner: remotegpu_user
--

SELECT pg_catalog.setval('public.datasets_id_seq', 1, false);


--
-- Name: gpus_id_seq; Type: SEQUENCE SET; Schema: public; Owner: remotegpu_user
--

SELECT pg_catalog.setval('public.gpus_id_seq', 1, false);


--
-- Name: images_id_seq; Type: SEQUENCE SET; Schema: public; Owner: remotegpu_user
--

SELECT pg_catalog.setval('public.images_id_seq', 1, false);


--
-- Name: machine_enrollments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: remotegpu_user
--

SELECT pg_catalog.setval('public.machine_enrollments_id_seq', 1, false);


--
-- Name: ssh_keys_id_seq; Type: SEQUENCE SET; Schema: public; Owner: remotegpu_user
--

SELECT pg_catalog.setval('public.ssh_keys_id_seq', 1, false);


--
-- Name: workspaces_id_seq; Type: SEQUENCE SET; Schema: public; Owner: remotegpu_user
--

SELECT pg_catalog.setval('public.workspaces_id_seq', 1, false);


--
-- Name: active_alerts active_alerts_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.active_alerts
    ADD CONSTRAINT active_alerts_pkey PRIMARY KEY (id);


--
-- Name: alert_rules alert_rules_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.alert_rules
    ADD CONSTRAINT alert_rules_pkey PRIMARY KEY (id);


--
-- Name: allocations allocations_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.allocations
    ADD CONSTRAINT allocations_pkey PRIMARY KEY (id);


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- Name: dataset_mounts dataset_mounts_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.dataset_mounts
    ADD CONSTRAINT dataset_mounts_pkey PRIMARY KEY (id);


--
-- Name: datasets datasets_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.datasets
    ADD CONSTRAINT datasets_pkey PRIMARY KEY (id);


--
-- Name: gpus gpus_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.gpus
    ADD CONSTRAINT gpus_pkey PRIMARY KEY (id);


--
-- Name: hosts hosts_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.hosts
    ADD CONSTRAINT hosts_pkey PRIMARY KEY (id);


--
-- Name: images images_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.images
    ADD CONSTRAINT images_pkey PRIMARY KEY (id);


--
-- Name: machine_enrollments machine_enrollments_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.machine_enrollments
    ADD CONSTRAINT machine_enrollments_pkey PRIMARY KEY (id);


--
-- Name: ssh_keys ssh_keys_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.ssh_keys
    ADD CONSTRAINT ssh_keys_pkey PRIMARY KEY (id);


--
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (id);


--
-- Name: gpus uni_gpus_uuid; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.gpus
    ADD CONSTRAINT uni_gpus_uuid UNIQUE (uuid);


--
-- Name: images uni_images_name; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.images
    ADD CONSTRAINT uni_images_name UNIQUE (name);


--
-- Name: workspaces workspaces_pkey; Type: CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.workspaces
    ADD CONSTRAINT workspaces_pkey PRIMARY KEY (id);


--
-- Name: idx_allocations_customer_id; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE INDEX idx_allocations_customer_id ON public.allocations USING btree (customer_id);


--
-- Name: idx_allocations_host_id; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE INDEX idx_allocations_host_id ON public.allocations USING btree (host_id);


--
-- Name: idx_allocations_status; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE INDEX idx_allocations_status ON public.allocations USING btree (status);


--
-- Name: idx_customers_deleted_at; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE INDEX idx_customers_deleted_at ON public.customers USING btree (deleted_at);


--
-- Name: idx_customers_email; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE UNIQUE INDEX idx_customers_email ON public.customers USING btree (email);


--
-- Name: idx_customers_role; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE INDEX idx_customers_role ON public.customers USING btree (role);


--
-- Name: idx_customers_username; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE UNIQUE INDEX idx_customers_username ON public.customers USING btree (username);


--
-- Name: idx_customers_uuid; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE UNIQUE INDEX idx_customers_uuid ON public.customers USING btree (uuid);


--
-- Name: idx_datasets_deleted_at; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE INDEX idx_datasets_deleted_at ON public.datasets USING btree (deleted_at);


--
-- Name: idx_datasets_uuid; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE UNIQUE INDEX idx_datasets_uuid ON public.datasets USING btree (uuid);


--
-- Name: idx_gpus_host_id; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE INDEX idx_gpus_host_id ON public.gpus USING btree (host_id);


--
-- Name: idx_host_gpu_index; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE UNIQUE INDEX idx_host_gpu_index ON public.gpus USING btree (host_id, index);


--
-- Name: idx_machine_enrollments_customer; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE INDEX idx_machine_enrollments_customer ON public.machine_enrollments USING btree (customer_id);


--
-- Name: idx_machine_enrollments_customer_id; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE INDEX idx_machine_enrollments_customer_id ON public.machine_enrollments USING btree (customer_id);


--
-- Name: idx_machine_enrollments_status; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE INDEX idx_machine_enrollments_status ON public.machine_enrollments USING btree (status);


--
-- Name: idx_ssh_keys_customer_id; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE INDEX idx_ssh_keys_customer_id ON public.ssh_keys USING btree (customer_id);


--
-- Name: idx_workspaces_uuid; Type: INDEX; Schema: public; Owner: remotegpu_user
--

CREATE UNIQUE INDEX idx_workspaces_uuid ON public.workspaces USING btree (uuid);


--
-- Name: machine_enrollments update_machine_enrollments_updated_at; Type: TRIGGER; Schema: public; Owner: remotegpu_user
--

CREATE TRIGGER update_machine_enrollments_updated_at BEFORE UPDATE ON public.machine_enrollments FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: active_alerts fk_active_alerts_rule; Type: FK CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.active_alerts
    ADD CONSTRAINT fk_active_alerts_rule FOREIGN KEY (rule_id) REFERENCES public.alert_rules(id);


--
-- Name: allocations fk_allocations_workspace; Type: FK CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.allocations
    ADD CONSTRAINT fk_allocations_workspace FOREIGN KEY (workspace_id) REFERENCES public.workspaces(id);


--
-- Name: allocations fk_customers_allocations; Type: FK CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.allocations
    ADD CONSTRAINT fk_customers_allocations FOREIGN KEY (customer_id) REFERENCES public.customers(id);


--
-- Name: ssh_keys fk_customers_ssh_keys; Type: FK CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.ssh_keys
    ADD CONSTRAINT fk_customers_ssh_keys FOREIGN KEY (customer_id) REFERENCES public.customers(id);


--
-- Name: workspaces fk_customers_workspaces; Type: FK CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.workspaces
    ADD CONSTRAINT fk_customers_workspaces FOREIGN KEY (owner_id) REFERENCES public.customers(id);


--
-- Name: dataset_mounts fk_datasets_dataset_mounts; Type: FK CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.dataset_mounts
    ADD CONSTRAINT fk_datasets_dataset_mounts FOREIGN KEY (dataset_id) REFERENCES public.datasets(id);


--
-- Name: allocations fk_hosts_allocations; Type: FK CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.allocations
    ADD CONSTRAINT fk_hosts_allocations FOREIGN KEY (host_id) REFERENCES public.hosts(id);


--
-- Name: gpus fk_hosts_gp_us; Type: FK CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.gpus
    ADD CONSTRAINT fk_hosts_gp_us FOREIGN KEY (host_id) REFERENCES public.hosts(id);


--
-- Name: tasks fk_tasks_image; Type: FK CONSTRAINT; Schema: public; Owner: remotegpu_user
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT fk_tasks_image FOREIGN KEY (image_id) REFERENCES public.images(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: remotegpu_user
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;


--
-- PostgreSQL database dump complete
--

\unrestrict IA5SsACGgKtsCB0nheG64BVP6kmXmVaRvKOad8JNCpaEvyHyIyfNyCuwOJ1kOab

