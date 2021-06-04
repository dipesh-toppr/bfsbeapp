-- Table: slots

-- DROP TABLE slots;

CREATE TABLE slots
(
    id integer NOT NULL DEFAULT nextval('slots_id_seq'::regclass),
    teacher_id integer NOT NULL,
    available_slot integer NOT NULL,
    is_booked integer NOT NULL,
    CONSTRAINT slots_pkey PRIMARY KEY (id)
)
TABLESPACE pg_default;

ALTER TABLE slots
    OWNER to postgres;