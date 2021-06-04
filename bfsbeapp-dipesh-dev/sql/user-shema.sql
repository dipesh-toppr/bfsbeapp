-- postgres database credentials:

-- # psql -U postgres
-- # postgres
-- # CREATE DATABASE users;
-- \c users; 

-- Table: users

drop table if exists users cascade;

create table users
(
  id serial,
  firstname character varying(255) not null,
  lastname character varying(255) not null,
  email character varying(255) not null,
  identity character varying(255) not null,
  password character varying(255) not null,
  constraint users_pkey primary key (id)
)
with (
  oids=false
);

alter table users owner to postgres;
drop index if exists email_index;
create unique index email_index on users using btree (email);
