-- +goose Up
CREATE TABLE tbl_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    login_session TEXT NULL,
    last_access TIMESTAMP WITHOUT TIME ZONE,
    profile_photo TEXT NULL,
    phone_number VARCHAR NULL,
    status_id SMALLINT DEFAULT 1,
    "order" INTEGER NULL DEFAULT 1,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_by INTEGER,
    deleted_at TIMESTAMP WITHOUT TIME ZONE

);

-- +goose StatementBegin
INSERT INTO tbl_users (
    username, password, email
) VALUES ('ADMIN', '123456', 'admin@example.com');
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS tbl_users;
