-- CREATE TABLE IF NOT EXISTS  "public"."links" (
--     "original_link" text COLLATE "pg_catalog"."default" NOT NULL,
--     "short_link" text COLLATE "pg_catalog"."default" NOT NULL,
--     UNIQUE(original_link)
-- );

CREATE TABLE "public"."links" (
  "original_link" text COLLATE "pg_catalog"."default" NOT NULL,
  "short_link" text COLLATE "pg_catalog"."default" NOT NULL
)
;
ALTER TABLE "public"."links" OWNER TO "yp";

CREATE INDEX "original_link" ON "public"."links" USING btree (
  "original_link" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "whort_link" ON "public"."links" USING btree (
  "short_link" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);