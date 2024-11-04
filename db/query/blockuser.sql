--name: CreateBlock: one
INSERT INTO BLOCKUSER (
    USER_ID,
    BLOCK_USER_ID,
    REASON,
    STATUS
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING *;

--name: GetAllBlockUsersList:many
SELECT
    *
FROM
    BLOCKUSER
ORDER BY
    BLOCK_AT DESC RETURNING *;

--name: GetBlockUserlist: many
SELECT
    *
FROM
    BLOCKUSER
WHERE
    USER_ID = $1
ORDER BY
    BLOCK_AT DESC RETURNING *;

--name: UnBlockUser: one
UPDATE TABLE BLOCKUSER
SET
    STATUS,
    UNBLOCK_AT
WHERE
    USER_ID = $1
    AND BLOCK_USER_ID = $2 RETURNING *;