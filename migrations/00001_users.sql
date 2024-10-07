CREATE SEQUENCE IF NOT EXISTS users_id_seq;
CREATE TABLE IF NOT EXISTS users
(
    id       BIGINT PRIMARY KEY DEFAULT nextval('users_id_seq'),
    login    VARCHAR(200),
    password VARCHAR(200)
);
