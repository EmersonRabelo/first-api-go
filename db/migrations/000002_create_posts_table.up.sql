CREATE TABLE IF NOT EXISTS posts (
    post_id uuid PRIMARY KEY NOT NULL,
    user_id uuid NOT NULL,
    post_body VARCHAR(280),
    is_active boolean NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_posts_id
ON posts(post_id)
WHERE deleted_at IS NULL AND is_active IS TRUE;

CREATE INDEX IF NOT EXISTS idx_posts_created_at
ON posts(created_at)
WHERE deleted_at IS NULL AND is_active IS TRUE;

CREATE INDEX IF NOT EXISTS idx_posts_user
ON posts(user_id)
WHERE deleted_at IS NULL AND is_active IS TRUE;

COMMENT ON TABLE posts IS 'Tabela de posts do sistema';

COMMENT ON COLUMN posts.post_id IS 'Id único do post';
COMMENT ON COLUMN posts.post_body IS 'Corpo do post';
COMMENT ON COLUMN posts.is_active IS 'Informa se o post está ativo ou inativo';
COMMENT ON COLUMN posts.created_at IS 'Data de criação do registro';
COMMENT ON COLUMN posts.updated_at IS 'Data da última atualização do registro';
COMMENT ON COLUMN posts.deleted_at IS 'Data da exclusão do registro (soft delete)';
COMMENT ON COLUMN posts.user_id IS 'Usuário que fez a postagem';