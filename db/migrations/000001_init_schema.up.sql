CREATE TABLE IF NOT EXISTS "accounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" varchar NOT NULL,
  "email" varchar NOT NULL UNIQUE,
  "google_id" varchar DEFAULT NULL,
  "facebook_id" varchar DEFAULT NULL,
  "image_count" bigint NOT NULL DEFAULT 0,
  "subscribed" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
  CHECK (
    (google_id IS NOT NULL AND facebook_id IS NULL) OR 
    (google_id IS NULL AND facebook_id IS NOT NULL)
  )
);

CREATE INDEX ON "accounts" ("owner");

COMMENT ON COLUMN "accounts"."google_id" IS 'google calls it sub';
COMMENT ON COLUMN "accounts"."facebook_id" IS 'facebook calls it ...';
