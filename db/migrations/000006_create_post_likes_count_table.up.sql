CREATE TABLE IF NOT EXISTS post_likes_count (
    post_id uuid PRIMARY KEY NOT NULL,
    like_count BIGINT NOT NULL DEFAULT 0
);

COMMENT ON TABLE post_likes_count IS 'Tabela de contagem de likes por postagem';

COMMENT ON COLUMN post_likes_count.post_id IS 'Id único da postagem';
COMMENT ON COLUMN post_likes_count.like_count IS 'Quantidade atual de curtidas na postagem';
