DROP INDEX "public"."original_link";

DROP INDEX "public"."whort_link";

CREATE UNIQUE INDEX "short_link_idx" ON "public"."links" (
  "short_link" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);