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
INSERT INTO tbl_admin (first_name, last_name, admin_name, password, email, phone, is_admin, role_id, profile_photo, status_id, created_by) VALUES
('ADMIN', 'SUPER1', 'ADMIN1', '123456', 'ADMIN1@example.com', '1234567890', 1, 1, 'admin1.jpg', 1, 1),
('ADMIN', 'SUPER2', 'ADMIN2', '123456', 'ADMIN2@example.com', '1234567890', 1, 1, 'admin2.jpg', 1, 1),
('ADMIN', 'SUPER14', 'ADMIN3', '123456', 'ADMIN3@example.com', '0987654321', 1, 1, 'admin3.jpg', 1, 1);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS tbl_admin;
