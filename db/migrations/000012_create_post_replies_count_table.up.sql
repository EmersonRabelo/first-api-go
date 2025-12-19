CREATE TABLE IF NOT EXISTS post_replies_count (
    post_id uuid PRIMARY KEY NOT NULL,
    reply_count BIGINT NOT NULL DEFAULT 0
);

COMMENT ON TABLE post_replies_count IS 'Tabela de contagem de respostas por postagem';

COMMENT ON COLUMN post_replies_count.post_id IS 'Id único da postagem';
COMMENT ON COLUMN post_replies_count.reply_count IS 'Quantidade atual de respostas na postagem';
