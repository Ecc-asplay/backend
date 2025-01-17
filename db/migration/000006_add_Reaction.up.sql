CREATE TABLE "posts_reaction" (
    "user_id" UUID NOT NULL,
    "post_id" UUID NOT NULL,
    "p_reaction_thanks" BOOLEAN NOT NULL DEFAULT FALSE,
    "p_reaction_heart" BOOLEAN NOT NULL DEFAULT FALSE,
    "p_reaction_helpful" BOOLEAN NOT NULL DEFAULT FALSE,
    "p_reaction_useful" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP(0) NOT NULL DEFAULT (now()),
    PRIMARY KEY ("user_id", "post_id")
);

CREATE TABLE "comments_reaction" (
    "user_id" UUID NOT NULL,
    "comment_id" UUID NOT NULL,
    "c_reaction_thanks" BOOLEAN NOT NULL DEFAULT FALSE,
    "c_reaction_heart" BOOLEAN NOT NULL DEFAULT FALSE,
    "c_reaction_helpful" BOOLEAN NOT NULL DEFAULT FALSE,
    "c_reaction_useful" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP(0) NOT NULL DEFAULT (now()),
    PRIMARY KEY ("user_id", "comment_id")
);

ALTER TABLE "posts_reaction" 
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id"),
ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("post_id");

ALTER TABLE "comments_reaction" 
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id"),
ADD FOREIGN KEY ("comment_id") REFERENCES "comments" ("comment_id");