COMMENT ON COLUMN "accounts"."google_id" IS NULL;
COMMENT ON COLUMN "accounts"."facebook_id" IS NULL;
DROP INDEX IF EXISTS accounts_owner_idx;
DROP TABLE IF EXISTS "accounts";
REVOKE ALL PRIVILEGES ON DATABASE userdb FROM userdb_user;
DROP DATABASE IF EXISTS userdb;
DROP ROLE IF EXISTS userdb_user;
