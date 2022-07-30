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

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: page_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.page_tokens (
    next_page_token character varying(200) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    is_used boolean DEFAULT false NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL
);


--
-- Name: videos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.videos (
    id integer NOT NULL,
    youtube_id character varying(20) NOT NULL,
    title character varying(200) NOT NULL,
    description character varying(5000),
    published_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    thumbnail_url character varying(500) NOT NULL
);


--
-- Name: videos_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.videos_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: videos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.videos_id_seq OWNED BY public.videos.id;


--
-- Name: videos id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.videos ALTER COLUMN id SET DEFAULT nextval('public.videos_id_seq'::regclass);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: videos videos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.videos
    ADD CONSTRAINT videos_pkey PRIMARY KEY (id);


--
-- Name: idx_videos_published_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_videos_published_at ON public.videos USING btree (published_at);


--
-- Name: idx_videos_title_description_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_videos_title_description_index ON public.videos USING btree (title, description);


--
-- Name: videos_youtube_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX videos_youtube_id_idx ON public.videos USING btree (youtube_id);


--
-- Name: videos_youtube_id_idx1; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX videos_youtube_id_idx1 ON public.videos USING btree (youtube_id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20220729152134');
