CREATE TABLE IF NOT EXISTS reports (
    report_id uuid PRIMARY KEY NOT NULL,
    post_id uuid NOT NULL,
    reporter_id uuid NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    perspective_toxicity DOUBLE PRECISION NULL,
    perspective_insult DOUBLE PRECISION NULL,
    perspective_profanity DOUBLE PRECISION NULL,
    perspective_threat DOUBLE PRECISION NULL,
    perspective_identity_hate DOUBLE PRECISION NULL,
    perspective_language VARCHAR(50) NULL,
    perspective_response_at TIMESTAMP NULL,

    CONSTRAINT fk_reports_post
        FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    CONSTRAINT fk_reports_reporter
        FOREIGN KEY (reporter_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_reports_id
ON reports(report_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_post_reports
ON reports(post_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_report
ON reports(reporter_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_reports_status
ON reports(status)
WHERE deleted_at IS NULL;

COMMENT ON TABLE reports IS 'Tabela de denúncias (reports) de posts e resultados de processamento (Perspective)';

COMMENT ON COLUMN reports.report_id IS 'Id único do report';
COMMENT ON COLUMN reports.post_id IS 'Id do post denunciado';
COMMENT ON COLUMN reports.reporter_id IS 'Usuário que realizou a denúncia';
COMMENT ON COLUMN reports.status IS 'Status do processamento do report (pending/processing/done/error)';
COMMENT ON COLUMN reports.created_at IS 'Data de criação do registro';
COMMENT ON COLUMN reports.updated_at IS 'Data da última atualização do registro';
COMMENT ON COLUMN reports.deleted_at IS 'Data da exclusão do registro (soft delete)';

COMMENT ON COLUMN reports.perspective_toxicity IS 'Score de toxicidade retornado pela Perspective API';
COMMENT ON COLUMN reports.perspective_insult IS 'Score de insulto retornado pela Perspective API';
COMMENT ON COLUMN reports.perspective_profanity IS 'Score de profanidade retornado pela Perspective API';
COMMENT ON COLUMN reports.perspective_threat IS 'Score de ameaça retornado pela Perspective API';
COMMENT ON COLUMN reports.perspective_identity_hate IS 'Score de ódio por identidade retornado pela Perspective API';
COMMENT ON COLUMN reports.perspective_language IS 'Idioma detectado/retornado pela Perspective API';
COMMENT ON COLUMN reports.perspective_response_at IS 'Data/hora do retorno do processamento pela Perspective API';