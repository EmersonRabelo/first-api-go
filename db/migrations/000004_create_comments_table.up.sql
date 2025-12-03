CREATE TABLE IF NOT EXISTS replies (
    reply_id uuid PRIMARY KEY,
    reply_body VARCHAR(280),
    user_id uuid NOT NULL,
    post_id uuid NOT NULL,
    quantity BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL, -- Se o delete_at for <> de null, o usuário removeu a resposta
    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_posts FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_replies_id
ON replies(reply_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_replies
ON replies(user_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_post_replies 
ON replies(post_id)
WHERE deleted_at IS NULL;

COMMENT ON TABLE replies IS 'Tabela de respostas, que também contabiliza o total de respostas por post';

COMMENT ON COLUMN replies.reply_id IS 'Id único da resposta';
COMMENT ON COLUMN replies.reply_body IS 'Contéudo da resposta';
COMMENT ON COLUMN replies.user_id IS 'Usuário que respondeu';
COMMENT ON COLUMN replies.post_id IS 'Id do post que foi respondido';
COMMENT ON COLUMN replies.quantity IS 'Total de respostas no post (incremento e decremento)';
COMMENT ON COLUMN replies.created_at IS 'Data de criação do registro';
COMMENT ON COLUMN replies.updated_at IS 'Data da última atualização do registro';
COMMENT ON COLUMN replies.deleted_at IS 'Data da exclusão do registro (soft delete). Usuário removeu o like';