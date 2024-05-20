CREATE TABLE IF NOT EXISTS "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "image_count" bigint NOT NULL DEFAULT 0,
  "subscribed" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("owner");


CREATE TABLE IF NOT EXISTS "images" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
  "url" varchar NOT NULL,
  "text" varchar NOT NULL DEFAULT '',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "images" ("user_id");

