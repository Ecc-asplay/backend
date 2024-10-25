CREATE TABLE "User" (
    "UserID" UUID UNIQUE PRIMARY KEY NOT NULL,
    "Name" VARCHAR NOT NULL,
    "Email" VARCHAR UNIQUE NOT NULL,
    "Age" INT NOT NULL,
    "Gender" VARCHAR NOT NULL,
    "Role" VARCHAR NOT NULL,
    "Disease" VARCHAR NOT NULL,
    "Condition" VARCHAR NOT NULL,
    "Certification" BOOLEAN DEFAULT FALSE,
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "HashPassword" VARCHAR NOT NULL,
    "ResetPasswordAt" TIMESTAMPTZ DEFAULT (NOW())
);

CREATE TABLE "Posts" (
    "UserID" UUID NOT NULL,
    "PostsID" UUID UNIQUE PRIMARY KEY NOT NULL,
    "ShowID" VARCHAR NOT NULL,
    "Title" VARCHAR NOT NULL,
    "Content" VARCHAR NOT NULL,
    "Media" BYTEA[],
    "IsSensitive" BOOLEAN DEFAULT FALSE,
    "Status" VARCHAR NOT NULL,
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "EditedAt" TIMESTAMPTZ DEFAULT (NOW())
);

CREATE TABLE "UserComment" (
    "UserID" UUID NOT NULL,
    "PostsID" UUID NOT NULL,
    "Comment" VARCHAR NOT NULL,
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

CREATE TABLE "Comment" (
    "PostsID" UUID NOT NULL,
    "Comment" VARCHAR NOT NULL,
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

CREATE TABLE "SearchRecord" (
    "SearchContent" VARCHAR NOT NULL,
    "isUser" BOOLEAN DEFAULT FALSE,
    "SearchedAt" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

CREATE TABLE "Block" (
    "UserID" UUID NOT NULL,
    "BlockUserID" UUID NOT NULL,
    "Reason" VARCHAR NOT NULL,
    "BlockAt" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "UnBlockAt" TIMESTAMPTZ DEFAULT (NOW())
);

CREATE TABLE "Tap" (
    "PostID" UUID PRIMARY KEY,
    "Tap" VARCHAR[] NOT NULL
);

CREATE TABLE "SavePost" (
    "UserID" UUID NOT NULL,
    "PostID" UUID NOT NULL,
    "SaveAt" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

CREATE TABLE "Token" (
    "TokenID" UUID UNIQUE NOT NULL,
    "Email" VARCHAR NOT NULL,
    "Status" VARCHAR NOT NULL,
    "Role" VARCHAR NOT NULL,
    "TakeAt" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "ExpiresAt" TIMESTAMPTZ NOT NULL
);

ALTER TABLE "Posts"
    ADD FOREIGN KEY (
        "UserID"
    )
        REFERENCES "User" (
            "UserID"
        );

ALTER TABLE "UserComment"
    ADD FOREIGN KEY (
        "UserID"
    )
        REFERENCES "User" (
            "UserID"
        );

ALTER TABLE "UserComment"
    ADD FOREIGN KEY (
        "PostsID"
    )
        REFERENCES "Posts" (
            "PostsID"
        );

ALTER TABLE "Comment"
    ADD FOREIGN KEY (
        "PostsID"
    )
        REFERENCES "Posts" (
            "PostsID"
        );

ALTER TABLE "Block"
    ADD FOREIGN KEY (
        "UserID"
    )
        REFERENCES "User" (
            "UserID"
        );

ALTER TABLE "Tap"
    ADD FOREIGN KEY (
        "PostID"
    )
        REFERENCES "Posts" (
            "PostsID"
        );

ALTER TABLE "SavePost"
    ADD FOREIGN KEY (
        "UserID"
    )
        REFERENCES "User" (
            "UserID"
        );

ALTER TABLE "SavePost"
    ADD FOREIGN KEY (
        "PostID"
    )
        REFERENCES "Posts" (
            "PostsID"
        );

ALTER TABLE "Token"
    ADD FOREIGN KEY (
        "Email"
    )
        REFERENCES "User" (
            "Email"
        );