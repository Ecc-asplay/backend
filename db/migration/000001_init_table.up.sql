CREATE TABLE "users" (
    "user_id" UUID UNIQUE PRIMARY KEY NOT NULL,
    "username" VARCHAR NOT NULL,
    "email" VARCHAR UNIQUE NOT NULL,
    "birth" DATE NOT NULL,
    "gender" VARCHAR NOT NULL,
    "is_privacy" BOOLEAN NOT NULL DEFAULT FALSE,
    "disease" VARCHAR NOT NULL,
    "condition" VARCHAR NOT NULL,
    "hashpassword" VARCHAR NOT NULL,
    "certification" BOOLEAN DEFAULT FALSE,
    "reset_password_at" TIMESTAMP DEFAULT (NOW()),
    "created_at" TIMESTAMP NOT NULL DEFAULT (NOW()),
    "updated_at" TIMESTAMP NOT NULL DEFAULT (NOW())
);

CREATE TABLE "posts" (
    "user_id" UUID NOT NULL,
    "post_id" UUID UNIQUE PRIMARY KEY NOT NULL,
    "show_id" VARCHAR NOT NULL,
    "title" VARCHAR NOT NULL,
    "feel" VARCHAR NOT NULL,
    "content" VARCHAR NOT NULL,
    "reaction" INT NOT NULL,
    "image" BYTEA[],
    "is_sensitive" BOOLEAN DEFAULT FALSE,
    "status" VARCHAR NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT (NOW()),
    "updated_at" TIMESTAMP DEFAULT (NOW())
);

CREATE TABLE "comments" (
    "comment_id" UUID UNIQUE PRIMARY KEY NOT NULL,
    "user_id" UUID NOT NULL,
    "post_id" UUID NOT NULL,
    "status" VARCHAR NOT NULL,
    "is_public" BOOLEAN NOT NULL,
    "comment" VARCHAR NOT NULL,
    "reaction" INT NOT NULL,
    "is_censored" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP NOT NULL DEFAULT (NOW())
);

CREATE TABLE "bookmarks" (
    "post_id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT (NOW())
);

CREATE TABLE "searchrecord" (
    "search_content" VARCHAR NOT NULL,
    "is_user" BOOLEAN NOT NULL DEFAULT FALSE,
    "searched_at" TIMESTAMP NOT NULL DEFAULT (NOW())
);

CREATE TABLE "block" (
    "user_id" UUID NOT NULL,
    "blockuser_id" UUID NOT NULL,
    "reason" VARCHAR NOT NULL,
    "block_at" TIMESTAMP NOT NULL DEFAULT (NOW()),
    "unblock_at" TIMESTAMP DEFAULT (NOW())
);

CREATE TABLE "adminuser" (
    "email" VARCHAR UNIQUE NOT NULL,
    "hash_password" VARCHAR NOT NULL,
    "staff_name" VARCHAR NOT NULL,
    "department" VARCHAR NOT NULL
);

CREATE TABLE "tap" (
    "post_id" UUID PRIMARY KEY,
    "tap" VARCHAR[] NOT NULL
);

CREATE TABLE "token" (
    "token" VARCHAR UNIQUE NOT NULL,
    "email" VARCHAR NOT NULL,
    "role" VARCHAR NOT NULL,
    "status" VARCHAR NOT NULL,
    "take_at" TIMESTAMP NOT NULL DEFAULT (NOW()),
    "expires_at" TIMESTAMP NOT NULL
);

ALTER TABLE "posts"
    ADD FOREIGN KEY (
        "user_id"
    )
        REFERENCES "users" (
            "user_id"
        );

ALTER TABLE "comments"
    ADD FOREIGN KEY (
        "user_id"
    )
        REFERENCES "users" (
            "user_id"
        );

ALTER TABLE "comments"
    ADD FOREIGN KEY (
        "post_id"
    )
        REFERENCES "posts" (
            "post_id"
        );

ALTER TABLE "bookmarks"
    ADD FOREIGN KEY (
        "post_id"
    )
        REFERENCES "posts" (
            "post_id"
        );

ALTER TABLE "bookmarks"
    ADD FOREIGN KEY (
        "user_id"
    )
        REFERENCES "users" (
            "user_id"
        );

ALTER TABLE "block"
    ADD FOREIGN KEY (
        "user_id"
    )
        REFERENCES "users" (
            "user_id"
        );

ALTER TABLE "tap"
    ADD FOREIGN KEY (
        "post_id"
    )
        REFERENCES "posts" (
            "post_id"
        );

ALTER TABLE "token"
    ADD FOREIGN KEY (
        "email"
    )
        REFERENCES "users" (
            "email"
        );