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
    "certification" BOOLEAN NOT NULL DEFAULT FALSE,
    "reset_password_at" TIMESTAMP(0),
    "created_at" TIMESTAMP(0) NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP(0)
);

CREATE TABLE "posts" (
    "user_id" UUID NOT NULL,
    "post_id" UUID UNIQUE PRIMARY KEY NOT NULL,
    "show_id" VARCHAR NOT NULL,
    "title" VARCHAR NOT NULL,
    "feel" VARCHAR NOT NULL,
    "content" JSONB NOT NULL,
    "images" JSONB,
    "reaction" INT NOT NULL,
    "is_sensitive" BOOLEAN NOT NULL DEFAULT FALSE,
    "status" VARCHAR NOT NULL,
    "created_at" TIMESTAMP(0) NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP(0)
);

CREATE TABLE "comments" (
    "comment_id" UUID UNIQUE PRIMARY KEY NOT NULL,
    "user_id" UUID NOT NULL,
    "post_id" UUID NOT NULL,
    "status" VARCHAR NOT NULL,
    "is_public" BOOLEAN NOT NULL,
    "comments" VARCHAR NOT NULL,
    "reaction" INT NOT NULL,
    "is_censored" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP(0) NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP(0)
);

CREATE TABLE "bookmarks" (
    "user_id" UUID NOT NULL,
    "post_id" UUID NOT NULL,
    "created_at" TIMESTAMP(0) NOT NULL DEFAULT NOW()
);

CREATE TABLE "searchrecord" (
    "search_content" VARCHAR NOT NULL,
    "is_user" BOOLEAN NOT NULL DEFAULT FALSE,
    "searched_at" TIMESTAMP(0) NOT NULL DEFAULT NOW()
);

CREATE TABLE "blockuser" (
    "user_id" UUID NOT NULL,
    "blockuser_id" UUID NOT NULL,
    "reason" VARCHAR NOT NULL,
    "status" VARCHAR NOT NULL,
    "block_at" TIMESTAMP(0) NOT NULL DEFAULT NOW(),
    "unblock_at" TIMESTAMP(0)
);

CREATE TABLE "adminuser" (
    "email" VARCHAR UNIQUE NOT NULL,
    "hashpassword" VARCHAR NOT NULL,
    "staff_name" VARCHAR NOT NULL,
    "department" VARCHAR NOT NULL,
    "joined_at" TIMESTAMP(0) NOT NULL DEFAULT NOW()
);

CREATE TABLE "tag" (
    "post_id" UUID PRIMARY KEY,
    "tag_comments" VARCHAR NOT NULL
);

CREATE TABLE "token" (
    "id" UUID UNIQUE PRIMARY KEY NOT NULL,
    "user_id" UUID NOT NULL,
    "access_token" VARCHAR UNIQUE NOT NULL,
    "roles" VARCHAR NOT NULL,
    "status" VARCHAR NOT NULL,
    "take_at" TIMESTAMP(0) NOT NULL DEFAULT NOW(),
    "expires_at" TIMESTAMP(0) NOT NULL
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
        "user_id"
    )
        REFERENCES "users" (
            "user_id"
        );

ALTER TABLE "bookmarks"
    ADD FOREIGN KEY (
        "post_id"
    )
        REFERENCES "posts" (
            "post_id"
        );

ALTER TABLE "blockuser"
    ADD FOREIGN KEY (
        "user_id"
    )
        REFERENCES "users" (
            "user_id"
        );

ALTER TABLE "tag"
    ADD FOREIGN KEY (
        "post_id"
    )
        REFERENCES "posts" (
            "post_id"
        );