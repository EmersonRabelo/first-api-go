CREATE TABLE IF NOT EXISTS likes (
    like_id uuid PRIMARY KEY NOT NULL,
    user_id uuid NOT NULL,
    post_id uuid NOT NULL,
    quantity BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL, -- Se o delete_at for <> de null, o usuário removeu o like
    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_posts FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_likes
ON likes(user_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_post_likes
ON likes(post_id)
WHERE deleted_at IS NULL;

COMMENT ON TABLE likes IS 'Tabela de likes, que também contabiliza o total de likes por post';

COMMENT ON COLUMN likes.like_id IS 'Id único do like';
COMMENT ON COLUMN likes.user_id IS 'Usuário que curtiu o post';
COMMENT ON COLUMN likes.post_id IS 'Post que foi curtido';
COMMENT ON COLUMN likes.quantity IS 'Total de likes no post (incremento e decremento)';
COMMENT ON COLUMN likes.created_at IS 'Data de criação do registro';
COMMENT ON COLUMN likes.updated_at IS 'Data da última atualização do registro';
COMMENT ON COLUMN likes.deleted_at IS 'Data da exclusão do registro (soft delete). Usuário removeu o like';