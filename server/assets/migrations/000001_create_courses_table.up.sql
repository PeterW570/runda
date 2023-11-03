CREATE TABLE IF NOT EXISTS courses (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    last_updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1,
    archived_at timestamp(0) with time zone,
    name text NOT NULL,
    description text,
    location point NOT NULL,
    tags text [] NOT NULL DEFAULT '{}'::text [],
    website text
);
CREATE INDEX IF NOT EXISTS courses_name_idx ON courses USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS courses_tags_idx ON courses USING GIN (tags);