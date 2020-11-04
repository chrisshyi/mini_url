CREATE TABLE mini_urls (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    visits BIGINT
);
CREATE INDEX ON mini_urls(url);

GRANT ALL PRIVILEGES ON TABLE mini_urls TO mini_url;
GRANT SELECT, USAGE ON SEQUENCE mini_urls_id_seq TO mini_url;