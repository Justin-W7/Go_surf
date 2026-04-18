--
-- PostgreSQL database dump
--

\restrict LcjMXdp51hoC5tAr3xO21scEsK8teGhAWz2m64UCauQxR4BrOkUg6Yt0JBkavlo

-- Dumped from database version 16.13 (Ubuntu 16.13-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.13 (Ubuntu 16.13-0ubuntu0.24.04.1)

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

SET default_table_access_method = heap;

--
-- Name: buoys; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.buoys (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL
);


ALTER TABLE public.buoys OWNER TO postgres;

--
-- Name: cities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cities (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    country character varying(100) NOT NULL,
    state character varying(100) NOT NULL,
    county character varying(100) NOT NULL
);


ALTER TABLE public.cities OWNER TO postgres;

--
-- Name: current_surf_spot_conditions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.current_surf_spot_conditions (
    id integer NOT NULL,
    spot_id integer NOT NULL,
    recorded_at timestamp without time zone NOT NULL,
    dom_swell_height_m double precision,
    dom_swell_dir double precision,
    wind_speed_mph character varying(15),
    wind_direction character varying(5),
    air_temp_deg_c double precision,
    water_temp_deg_c double precision,
    precipitation double precision,
    cloud_coverage text,
    domwp_sec double precision
);


ALTER TABLE public.current_surf_spot_conditions OWNER TO postgres;

--
-- Name: current_surf_spot_conditions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.current_surf_spot_conditions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.current_surf_spot_conditions_id_seq OWNER TO postgres;

--
-- Name: current_surf_spot_conditions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.current_surf_spot_conditions_id_seq OWNED BY public.current_surf_spot_conditions.id;


--
-- Name: current_weather; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.current_weather (
    id integer NOT NULL,
    city_id integer NOT NULL,
    recorded_at timestamp with time zone NOT NULL,
    wind_speed text,
    wind_direction character varying(10),
    air_temp_c double precision,
    precipitation double precision,
    cloud_coverage text,
    observed_at timestamp with time zone NOT NULL
);


ALTER TABLE public.current_weather OWNER TO postgres;

--
-- Name: current_weather_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.current_weather_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.current_weather_id_seq OWNER TO postgres;

--
-- Name: current_weather_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.current_weather_id_seq OWNED BY public.current_weather.id;


--
-- Name: real_time_buoy_data_points; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.real_time_buoy_data_points (
    id integer NOT NULL,
    buoy_id integer NOT NULL,
    recorded_at timestamp without time zone DEFAULT now() NOT NULL,
    winddir_degt double precision,
    windspeed_m_pers double precision,
    windgust_m_pers double precision,
    waveh_m double precision,
    domwp_sec double precision,
    avgwavep_sec double precision,
    meanwavedir_degt double precision,
    airt_degc double precision,
    watert_degc double precision,
    inserted_at timestamp with time zone
);


ALTER TABLE public.real_time_buoy_data_points OWNER TO postgres;

--
-- Name: real_time_buoy_data_points_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.real_time_buoy_data_points_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.real_time_buoy_data_points_id_seq OWNER TO postgres;

--
-- Name: real_time_buoy_data_points_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.real_time_buoy_data_points_id_seq OWNED BY public.real_time_buoy_data_points.id;


--
-- Name: surfspot; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.surfspot (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    city_id integer NOT NULL,
    break_type character varying(100),
    orientation double precision,
    nearest_buoy integer
);


ALTER TABLE public.surfspot OWNER TO postgres;

--
-- Name: current_surf_spot_conditions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.current_surf_spot_conditions ALTER COLUMN id SET DEFAULT nextval('public.current_surf_spot_conditions_id_seq'::regclass);


--
-- Name: current_weather id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.current_weather ALTER COLUMN id SET DEFAULT nextval('public.current_weather_id_seq'::regclass);


--
-- Name: real_time_buoy_data_points id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.real_time_buoy_data_points ALTER COLUMN id SET DEFAULT nextval('public.real_time_buoy_data_points_id_seq'::regclass);


--
-- Name: buoys buoys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.buoys
    ADD CONSTRAINT buoys_pkey PRIMARY KEY (id);


--
-- Name: cities cities_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cities
    ADD CONSTRAINT cities_pkey PRIMARY KEY (id);


--
-- Name: current_surf_spot_conditions current_surf_spot_conditions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.current_surf_spot_conditions
    ADD CONSTRAINT current_surf_spot_conditions_pkey PRIMARY KEY (id);


--
-- Name: current_weather current_weather_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.current_weather
    ADD CONSTRAINT current_weather_pkey PRIMARY KEY (id);


--
-- Name: real_time_buoy_data_points real_time_buoy_data_points_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.real_time_buoy_data_points
    ADD CONSTRAINT real_time_buoy_data_points_pkey PRIMARY KEY (id);


--
-- Name: surfspot surfspot_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.surfspot
    ADD CONSTRAINT surfspot_pkey PRIMARY KEY (id);


--
-- Name: current_weather current_weather_city_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.current_weather
    ADD CONSTRAINT current_weather_city_id_fkey FOREIGN KEY (city_id) REFERENCES public.cities(id);


--
-- Name: surfspot fk_city; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.surfspot
    ADD CONSTRAINT fk_city FOREIGN KEY (city_id) REFERENCES public.cities(id);


--
-- Name: current_surf_spot_conditions fk_spot; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.current_surf_spot_conditions
    ADD CONSTRAINT fk_spot FOREIGN KEY (spot_id) REFERENCES public.surfspot(id);


--
-- Name: real_time_buoy_data_points real_time_buoy_data_points_buoy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.real_time_buoy_data_points
    ADD CONSTRAINT real_time_buoy_data_points_buoy_id_fkey FOREIGN KEY (buoy_id) REFERENCES public.buoys(id);


--
-- PostgreSQL database dump complete
--

\unrestrict LcjMXdp51hoC5tAr3xO21scEsK8teGhAWz2m64UCauQxR4BrOkUg6Yt0JBkavlo

