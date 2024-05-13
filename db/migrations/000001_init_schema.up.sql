-- Use to initialize empty DB and tables OR reset DB and tables

-- DO
-- $$
-- BEGIN
--     IF EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'userdb_user') THEN
--         DROP ROLE userdb_user;
--     END IF;
-- END
-- $$;
-- CREATE ROLE userdb_user WITH LOGIN PASSWORD 'test1234';
-- DROP DATABASE IF EXISTS userdb;
-- CREATE DATABASE userdb;
-- GRANT ALL PRIVILEGES ON DATABASE userdb TO userdb_user;

CREATE TABLE IF NOT EXISTS "accounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" varchar NOT NULL,
  "email" varchar NOT NULL,
  "google_id" varchar DEFAULT NULL,
  "facebook_id" varchar DEFAULT NULL,
  "image_count" bigint NOT NULL DEFAULT 0,
  "subscribed" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

COMMENT ON COLUMN "accounts"."google_id" IS 'google calls it sub';
COMMENT ON COLUMN "accounts"."facebook_id" IS 'facebook calls it ...';
