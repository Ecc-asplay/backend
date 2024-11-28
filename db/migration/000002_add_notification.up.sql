CREATE TABLE "notification" (
    "user_id" UUID NOT NULL,
    "content" VARCHAR NOT NULL,
    "icon" BYTEA NOT NULL,
    "is_read" BOOL NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP(0) NOT NULL DEFAULT NOW()
);