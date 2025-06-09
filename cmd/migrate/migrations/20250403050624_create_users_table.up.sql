CREATE TYPE "role" AS ENUM (
    'PET_OWNER',
    'VET',
    'PET_SHOP'
    );


CREATE TABLE IF NOT EXISTS "Users" (
    "user_id"       BIGSERIAL PRIMARY KEY,
    "username"      varchar,
    "email"         varchar,
    "address"       varchar,
    "role"          role,
    "avatar"        varchar,
    "password_hash" varchar,
    "created_at"    timestamp DEFAULT (now()),
    "updated_at"    timestamp,
    "deleted_at"    timestamp
);
