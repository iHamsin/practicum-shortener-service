ALTER TABLE "public"."links" 
  ADD COLUMN "uuid" varchar(255);

CREATE INDEX "uuid_idx" ON "public"."links" (
  "uuid"
);