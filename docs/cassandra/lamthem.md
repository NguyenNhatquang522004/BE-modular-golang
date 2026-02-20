<!-- CREATE MATERIALIZED VIEW IF NOT EXISTS live_comments_by_user AS
    SELECT *
    FROM live_comments
    WHERE user_id IS NOT NULL AND created_at IS NOT NULL AND stream_id IS NOT NULL AND comment_id IS NOT NULL
    PRIMARY KEY (user_id, created_at, stream_id, comment_id)
    WITH CLUSTERING ORDER BY (created_at DESC); -->