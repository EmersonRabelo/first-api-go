CREATE TABLE IF NOT EXISTS users (
    user_id uuid PRIMARY KEY NOT NULL,
    user_name VARCHAR(50) UNIQUE NOT NULL,
    user_email VARCHAR(255) UNIQUE NOT NULL,
    is_active boolean NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_id
ON users(user_id)
WHERE deleted_at IS NULL AND is_active IS TRUE;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name
ON users(user_name)
WHERE deleted_at IS NULL AND is_active IS TRUE;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email
ON users(user_email)
WHERE deleted_at IS NULL AND is_active IS TRUE;

COMMENT ON TABLE users IS 'Tabela de usuários do sistema';

COMMENT ON COLUMN users.user_id IS 'Id único do usuário';
COMMENT ON COLUMN users.user_name IS 'Nickname/Nome do usuário (único)';
COMMENT ON COLUMN users.user_email IS 'E-mail do usuário (único)';
COMMENT ON COLUMN users.is_active IS 'Informa se o usuário está ativo ou inativo';
COMMENT ON COLUMN users.created_at IS 'Data de criação do registro';
COMMENT ON COLUMN users.updated_at IS 'Data da última atualização do registro';
COMMENT ON COLUMN users.deleted_at IS 'Data da exclusão do registro (soft delete)';