CREATE ROLE mini_url LOGIN ENCRYPTED PASSWORD 'mini_pass';

CREATE TABLE mini_urls (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    visits BIGINT
);

GRANT ALL PRIVILEGES ON TABLE mini_urls TO mini_url;
GRANT SELECT, USAGE ON SEQUENCE mini_urls_id_seq TO mini_url;
