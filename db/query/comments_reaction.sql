-- name: CreateCommentsReaction :one
INSERT INTO COMMENTS_REACTION (
    USER_ID,
    COMMENT_ID,
    C_REACTION_THANKS,
    C_REACTION_HELPFUL,
    C_REACTION_USEFUL,
    C_REACTION_HEART
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetCommentsReaction :many
SELECT * FROM COMMENTS_REACTION 
WHERE COMMENT_ID = $1;


-- name: UpdateCommentsReactionThanks :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_THANKS = NOT C_REACTION_THANKS
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING *;

-- name: UpdateCommentsReactionHelpful :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_HELPFUL = NOT C_REACTION_HELPFUL
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING *;

-- name: UpdateCommentsReactionUseful :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_USEFUL = NOT C_REACTION_USEFUL
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING *;

-- name: UpdateCommentsReactionHeart :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_HEART = NOT C_REACTION_HEART
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING *;

-- name: DeleteCommentsReaction :exec
DELETE FROM COMMENTS_REACTION
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
    AND NOT (C_REACTION_THANKS OR C_REACTION_HEART OR C_REACTION_HELPFUL OR C_REACTION_USEFUL);


-- name: GetCommentsHeartOfTrue :one
SELECT COUNT(*)
FROM COMMENTS_REACTION
WHERE C_REACTION_HEART = TRUE AND COMMENT_ID = $1;

-- name: GetCommentsThanksOfTrue :one
SELECT COUNT(*)
FROM COMMENTS_REACTION
WHERE C_REACTION_THANKS = TRUE AND COMMENT_ID = $1;

-- name: GetCommentsHelpfulOfTrue :one
SELECT COUNT(*)
FROM COMMENTS_REACTION
WHERE C_REACTION_HELPFUL = TRUE AND COMMENT_ID = $1;

-- name: GetCommentsUsefulOfTrue :one
SELECT COUNT(*)
FROM COMMENTS_REACTION
WHERE C_REACTION_USEFUL = TRUE AND COMMENT_ID = $1;
