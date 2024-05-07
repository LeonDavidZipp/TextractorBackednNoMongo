CREATE TABLE "accounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" varchar NOT NULL,
  "image_count" bigint NOT NULL,
  "subscribed" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "images" (
  "id" BIGSERIAL PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "data_base64" text NOT NULL,
  "text" text,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "images" ("account_id");

COMMENT ON COLUMN "images"."data_base64" IS 'can be negative or positive';

COMMENT ON COLUMN "images"."text" IS 'transcribed text of image';

ALTER TABLE "images" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
