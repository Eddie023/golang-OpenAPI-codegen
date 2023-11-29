-- create "transactions" table
CREATE TABLE "transactions" ("id" uuid NOT NULL, "date" timestamptz NOT NULL, "amount_in_usd" double precision NOT NULL, "description" character varying NOT NULL, PRIMARY KEY ("id"));
