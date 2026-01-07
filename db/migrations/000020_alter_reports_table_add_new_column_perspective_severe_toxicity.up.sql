ALTER TABLE reports
ADD COLUMN perspective_severe_toxicity DOUBLE PRECISION;

COMMENT ON COLUMN reports.perspective_severe_toxicity
IS 'Score de toxicidade severa retornado pela Perspective API';