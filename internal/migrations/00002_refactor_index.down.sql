DROP INDEX "public"."short_link";

CREATE INDEX "original_link" ON "public"."links" USING btree (
  "original_link" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "whort_link" ON "public"."links" USING btree (
  "short_link" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);