-- +goose Up
CREATE TABLE tbl_users(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    user_name VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    login_session TEXT NULL,
    profile_photot TEXT NULL,
    user_alias VARCHAR(255) NULL,
    phone_number VARCHAR NULL,
    user_avatar_id INTEGER NULL,
    last_access TIMESTAMP WITHOUT TIME ZONE,
    status_id SMALLINT DEFAULT 1,
    "order" INTEGER NULL DEFAULT 1,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_by INTEGER,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- +goose Down
DROP TABLE IF EXISTS tbl_users;