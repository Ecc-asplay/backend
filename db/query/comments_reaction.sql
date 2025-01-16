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


-- name: UpdateCommentsReactionThanksPlusOne :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_THANKS = TRUE
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING *;

-- name: UpdateCommentsReactionHelpfulPlusOne :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_HELPFUL = TRUE
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING *;

-- name: UpdateCommentsReactionUsefulPlusOne :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_USEFUL = TRUE
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING *;

-- name: UpdateCommentsReactionHeartPlusOne :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_HEART = TRUE
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING *;

-- name: DeleteCommentsReaction :exec
DELETE FROM COMMENTS_REACTION
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
    AND NOT (C_REACTION_THANKS OR C_REACTION_HEART OR C_REACTION_HELPFUL OR C_REACTION_USEFUL);