CREATE TABLE "images" (
  "post_id" uuid NOT NULL,
  "page" int NOT NULL,
  "img1" bytea NOT NULL,
  "img2" bytea,
  "img3" bytea,
  "img4" bytea,
  "img5" bytea,
  "created_at" TIMESTAMP(0) NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP(0)
);

ALTER TABLE "images" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("post_id");
