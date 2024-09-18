CREATE SEQUENCE IF NOT EXISTS vault_id_seq;
CREATE TABLE IF NOT EXISTS vault
(
    id      BIGINT PRIMARY KEY DEFAULT nextval('vault_id_seq'),
    key     VARCHAR(256),
    value   oid,
    user_id BIGINT references users NOT NULL
);
