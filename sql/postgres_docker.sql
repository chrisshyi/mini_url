CREATE TABLE mini_urls (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    visits BIGINT
);
CREATE INDEX ON mini_urls(url);