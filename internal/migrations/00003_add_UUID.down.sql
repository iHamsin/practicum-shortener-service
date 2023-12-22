DROP INDEX "public"."uuid_idx";

ALTER TABLE "public"."links" 
  DROP COLUMN "uuid";