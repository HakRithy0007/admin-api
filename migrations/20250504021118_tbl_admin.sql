-- +goose Up
CREATE TABLE tbl_admin (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(20) DEFAULT NULL,
    last_name VARCHAR(20) DEFAULT NULL,
    admin_name VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(15) DEFAULT NULL,
    is_admin SMALLINT DEFAULT 0,
    login_session TEXT NULL,
    last_access TIMESTAMP WITHOUT TIME ZONE,
    role_id INTEGER DEFAULT 0,
    profile_photo TEXT NULL,
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
INSERT INTO tbl_admin (first_name, last_name, admin_name, password, email) VALUES
    ('Admin', 'One', 'ADMIN1', '123456', 'admin1@example.com'),
    ('Admin', 'Two', 'ADMIN2', '123456', 'admin2@example.com'),
    ('Admin', 'Three', 'ADMIN3', '123456', 'admin3@example.com'),
    ('Admin', 'Four', 'ADMIN4', '123456', 'admin4@example.com'),
    ('Admin', 'Five', 'ADMIN5', '123456', 'admin5@example.com'),
    ('Admin', 'Six', 'ADMIN6', '123456', 'admin6@example.com'),
    ('Admin', 'Seven', 'ADMIN7', '123456', 'admin7@example.com'),
    ('Admin', 'Eight', 'ADMIN8', '123456', 'admin8@example.com'),
    ('Admin', 'Nine', 'ADMIN9', '123456', 'admin9@example.com');
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS tbl_admin;
