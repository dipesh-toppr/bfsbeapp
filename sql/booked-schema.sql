-- Table: bookeds

-- DROP TABLE bookeds;

CREATE TABLE bookeds
(
    id integer NOT NULL DEFAULT nextval('bookeds_id_seq'::regclass),
    student_id integer NOT NULL,
    slot_id integer NOT NULL,
    CONSTRAINT bookeds_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE bookeds
    OWNER to postgres;