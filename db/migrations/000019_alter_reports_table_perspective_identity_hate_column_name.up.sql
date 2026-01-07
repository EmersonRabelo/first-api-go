ALTER TABLE reports
RENAME COLUMN perspective_identity_hate TO perspective_identity_attack;

COMMENT ON COLUMN reports.perspective_identity_attack
IS 'Score de ataque de ódio por identidade retornado pela Perspective API';