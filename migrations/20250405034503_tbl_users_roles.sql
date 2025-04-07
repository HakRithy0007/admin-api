-- +goose Up
CREATE TABLE tbl_users_roles (
    id SERIAL PRIMARY KEY,
    user_role_name VARCHAR(255) NOT NULL,
    user_role_desc VARCHAR(255) NOT NULL,
    status_id SMALLINT DEFAULT 1,
    "order" INTEGER DEFAULT 1,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_by INTEGER,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- +goose StatementBegin
INSERT INTO tbl_users_roles (
    user_role_name, 
    user_role_desc, 
    status_id, 
    "order",
    created_by,
    created_at
) VALUES 
    ('Admin', 'Administrator role', 1, 1, 1, CURRENT_TIMESTAMP),
    ('Moderator', 'Standard moderator role',1, 1, 1, CURRENT_TIMESTAMP),
    ('Operator', 'Standard operator role', 1, 1, 1,  CURRENT_TIMESTAMP);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS tbl_users_roles;
