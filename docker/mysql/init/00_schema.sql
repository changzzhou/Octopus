-- Octopus Workflow Platform - Database Schema Stub
-- Full schema will be added in BE-1 (Workflow Definition CRUD)

-- Placeholder: workflows table will be generated via goctl model mysql ddl
-- This file serves as a template for database initialization

-- Example structure (to be replaced by goctl-generated migrations):
-- CREATE TABLE IF NOT EXISTS workflows (
--     id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
--     name VARCHAR(255) NOT NULL,
--     description TEXT,
--     definition JSON,
--     status TINYINT NOT NULL DEFAULT 0,
--     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
--     PRIMARY KEY (id),
--     INDEX idx_status (status),
--     INDEX idx_created_at (created_at)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SELECT 'Octopus database initialized' AS message;
