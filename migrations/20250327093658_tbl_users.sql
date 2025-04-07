-- +goose Up
CREATE TABLE tbl_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    login_session TEXT NULL,
    last_access TIMESTAMP WITHOUT TIME ZONE,
    profile_photo TEXT NULL,
    phone_number VARCHAR(20) NULL,
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
INSERT INTO tbl_users (username, password, email) VALUES
    ('ADMIN1', '123456', 'admin1@example.com'),
    ('ADMIN2', '123456', 'admin2@example.com'),
    ('ADMIN3', '123456', 'admin3@example.com'),
    ('ADMIN4', '123456', 'admin4@example.com'),
    ('ADMIN5', '123456', 'admin5@example.com'),
    ('ADMIN6', '123456', 'admin6@example.com'),
    ('ADMIN7', '123456', 'admin7@example.com'),
    ('ADMIN8', '123456', 'admin8@example.com'),
    ('ADMIN9', '123456', 'admin9@example.com');
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS tbl_users;
