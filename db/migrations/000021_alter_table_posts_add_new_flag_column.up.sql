ALTER TABLE posts
ADD COLUMN flag VARCHAR(48) NOT NULL DEFAULT 'visible';

COMMENT ON COLUMN posts.flag
IS 'São flags do estado de visualização e moderação de cada post';