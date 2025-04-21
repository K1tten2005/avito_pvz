CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TYPE user_role AS ENUM ('employee', 'moderator');
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,                     
    role user_role NOT NULL DEFAULT 'employee',                              
    password_hash BYTEA NOT NULL                   
);

CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY,
    registration_date DATE NOT NULL DEFAULT now(),
    city TEXT NOT NULL
);

CREATE TYPE reception_status AS ENUM ('in_progress', 'close');
CREATE TABLE IF NOT EXISTS reception (
    id UUID PRIMARY KEY,
    reception_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    status reception_status NOT NULL
);

CREATE TYPE product_category AS ENUM ('электроника', 'одежда', 'обувь');
CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY,
    reception_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    reception_id UUID NOT NULL REFERENCES reception(id) ON DELETE CASCADE,
    category product_category NOT NULL
);

INSERT INTO users (id, email, role, password_hash)
VALUES
(uuid_generate_v4(), 'nick@mail.ru', 'employee', decode('ff936a28b3fd98ea01207aa8b6b0662c3c6c3fd68a49241960bdf4d89e91003748d497862c7bbe48', 'hex'));

