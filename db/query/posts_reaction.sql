-- name: CreatePostsReaction :one
INSERT INTO POSTS_REACTION (
    USER_ID,
    POST_ID,
    P_REACTION_THANKS,
    P_REACTION_HELPFUL,
    P_REACTION_USEFUL,
    P_REACTION_HEART
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetPostsReaction :many
SELECT * FROM POSTS_REACTION 
WHERE POST_ID = $1;


-- name: UpdatePostsReactionThanksPlusOne :one
UPDATE POSTS_REACTION
SET 
    P_REACTION_THANKS = TRUE
WHERE 
    USER_ID = $1 AND POST_ID = $2
RETURNING *;

-- name: UpdatePostsReactionHelpfulPlusOne :one
UPDATE POSTS_REACTION
SET 
    P_REACTION_HELPFUL = TRUE
WHERE 
    USER_ID = $1 AND POST_ID = $2
RETURNING *;

-- name: UpdatePostsReactionUsefulPlusOne :one
UPDATE POSTS_REACTION
SET 
    P_REACTION_USEFUL = TRUE
WHERE 
    USER_ID = $1 AND POST_ID = $2
RETURNING *;

-- name: UpdatePostsReactionHeartPlusOne :one
UPDATE POSTS_REACTION
SET 
    P_REACTION_HEART = TRUE
WHERE 
    USER_ID = $1 AND POST_ID = $2
RETURNING *;

-- name: DeletePostsReaction :exec
DELETE FROM POSTS_REACTION
WHERE 
    USER_ID = $1 AND POST_ID = $2
    AND NOT (P_REACTION_THANKS OR P_REACTION_HEART OR P_REACTION_HELPFUL OR P_REACTION_USEFUL);