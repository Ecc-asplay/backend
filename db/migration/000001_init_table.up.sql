CREATE TABLE "User" (
    "UserID" UUID UNIQUE PRIMARY KEY NOT NULL,
    "Name" VARCHAR NOT NULL,
    "Email" VARCHAR UNIQUE NOT NULL,
    "Birth" DATE NOT NULL,
    "Gender" ENUM NOT NULL,
    "IsPrivacy" BOOL NOT NULL DEFAULT FALSE,
    "Disease" VARCHAR NOT NULL,
    "Condition" VARCHAR NOT NULL,
    "Password" VARCHAR NOT NULL,
    "Certification" BOOLEAN DEFAULT FALSE,
    "ResetPasswordAt" TIMESTAMP DEFAULT (NOW()),
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT (NOW()),
    "UpdatedAt" TIMESTAMP NOT NULL DEFAULT (NOW())
);

CREATE TABLE "Posts" (
    "UserID" UUID NOT NULL,
    "PostsID" UUID UNIQUE PRIMARY KEY NOT NULL,
    "ShowID" VARCHAR NOT NULL,
    "Title" VARCHAR NOT NULL,
    "Feel" ENUM NOT NULL,
    "Content" VARCHAR NOT NULL,
    "Reaction" INT NOT NULL,
    "Image" VARBINARY,
    "IsSensitive" BOOLEAN DEFAULT FALSE,
    "Status" VARCHAR NOT NULL,
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT (NOW()),
    "UpdatedAt" TIMESTAMP DEFAULT (NOW())
);

CREATE TABLE "Comments" (
    "CommentID" UUID UNIQUE PRIMARY KEY NOT NULL,
    "UserID" UUID NOT NULL,
    "PostsID" UUID NOT NULL,
    "Status" ENUM NOT NULL,
    "IsPublic" BOOL NOT NULL,
    "Comment" VARCHAR NOT NULL,
    "Reaction" INT NOT NULL,
    "IsCensored" BOOL NOT NULL DEFAULT FALSE,
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT (NOW())
);

CREATE TABLE "Bookmarks" (
    "PostID" UUID NOT NULL,
    "UserID" UUID NOT NULL,
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT (NOW())
);

CREATE TABLE "SearchRecord" (
    "SearchContent" VARCHAR NOT NULL,
    "isUser" BOOLEAN NOT NULL DEFAULT FALSE,
    "SearchedAt" TIMESTAMP NOT NULL DEFAULT (NOW())
);

CREATE TABLE "Block" (
    "UserID" UUID NOT NULL,
    "BlockUserID" UUID NOT NULL,
    "Reason" VARCHAR NOT NULL,
    "BlockAt" TIMESTAMP NOT NULL DEFAULT (NOW()),
    "UnBlockAt" TIMESTAMP DEFAULT (NOW())
);

CREATE TABLE "AdminUser" (
    "Email" VARCHAR UNIQUE NOT NULL,
    "Password" VARCHAR NOT NULL,
    "Name" VARCHAR NOT NULL,
    "Department" VARCHAR NOT NULL
);

CREATE TABLE "Tap" (
    "PostID" UUID PRIMARY KEY,
    "Tap" VARCHAR[] NOT NULL
);

CREATE TABLE "Token" (
    "Token" VARCHAR UNIQUE NOT NULL,
    "Email" VARCHAR NOT NULL,
    "Role" ENUM NOT NULL,
    "Status" VARCHAR NOT NULL,
    "TakeAt" TIMESTAMP NOT NULL DEFAULT (NOW()),
    "ExpiresAt" TIMESTAMP NOT NULL
);

ALTER TABLE "Posts"
    ADD FOREIGN KEY (
        "UserID"
    )
        REFERENCES "User" (
            "UserID"
        );

ALTER TABLE "Comments"
    ADD FOREIGN KEY (
        "UserID"
    )
        REFERENCES "User" (
            "UserID"
        );

ALTER TABLE "Comments"
    ADD FOREIGN KEY (
        "PostsID"
    )
        REFERENCES "Posts" (
            "PostsID"
        );

ALTER TABLE "Bookmarks"
    ADD FOREIGN KEY (
        "PostID"
    )
        REFERENCES "Posts" (
            "PostsID"
        );

ALTER TABLE "Bookmarks"
    ADD FOREIGN KEY (
        "UserID"
    )
        REFERENCES "User" (
            "UserID"
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

ALTER TABLE "Token"
    ADD FOREIGN KEY (
        "Email"
    )
        REFERENCES "User" (
            "Email"
        );