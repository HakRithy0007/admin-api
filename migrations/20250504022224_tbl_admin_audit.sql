-- +goose Up
CREATE TABLE tbl_admin_audit (
    id SERIAL PRIMARY KEY,
    admin_id INTEGER NOT NULL,
    admin_audit_context VARCHAR(255) NOT NULL,
    admin_audit_desc VARCHAR(255) NOT NULL,
    audit_type_id INTEGER NOT NULL,
    admin_agent VARCHAR(255) NOT NULL,
    operator VARCHAR(100) NOT NULL,
    ip VARCHAR(50) NOT NULL,
    status_id SMALLINT DEFAULT 1,
    "order" INTEGER DEFAULT 1,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NULL,
    deleted_by INTEGER NULL,
    deleted_at TIMESTAMP WITHOUT TIME ZONE NULL
);

-- +goose Down
DROP TABLE IF EXISTS tbl_admin_audit;