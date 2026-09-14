-- Octopus Workflow Platform - Database Schema
-- BE-1: Workflow Definition CRUD

CREATE TABLE IF NOT EXISTS `workflows` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `status` VARCHAR(20) NOT NULL DEFAULT 'draft' COMMENT 'draft|enabled|disabled',
    `version` INT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Incremented on each save for optimistic concurrency',
    `nodes` JSON COMMENT 'Array of workflow nodes',
    `edges` JSON COMMENT 'Array of workflow edges',
    `entry_node_id` VARCHAR(64) COMMENT 'ID of the entry node',
    `variables_schema` JSON COMMENT 'JSON schema for workflow variables',
    `canvas_meta` JSON COMMENT 'Canvas metadata (viewport, zoom, etc.)',
    `created_by` VARCHAR(64),
    `updated_by` VARCHAR(64),
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Node structure (stored in nodes JSON column):
-- {
--   "id": "node_1",
--   "type": "script|http|human",
--   "name": "Node Name",
--   "position": {"x": 100, "y": 200},
--   "config": {},
--   "credential_ref": "optional_credential_id"
-- }

-- Edge structure (stored in edges JSON column):
-- {
--   "id": "edge_1",
--   "source": "node_1",
--   "target": "node_2",
--   "outlet": "success|failure",
--   "label": "optional label"
-- }

-- Workflow runs table (BE-2)
CREATE TABLE IF NOT EXISTS `workflow_runs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `workflow_id` BIGINT UNSIGNED NOT NULL,
    `workflow_version` INT UNSIGNED NOT NULL COMMENT 'Version of workflow at run time',
    `definition_snapshot` JSON NOT NULL COMMENT 'Frozen workflow definition',
    `status` VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT 'pending|running|succeeded|failed|cancelled',
    `trigger_type` VARCHAR(32) NOT NULL DEFAULT 'manual' COMMENT 'manual|schedule|webhook',
    `started_at` TIMESTAMP NULL,
    `finished_at` TIMESTAMP NULL,
    `error_message` TEXT,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_workflow_id` (`workflow_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Workflow run steps table (BE-2)
CREATE TABLE IF NOT EXISTS `workflow_run_steps` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `run_id` BIGINT UNSIGNED NOT NULL,
    `node_id` VARCHAR(64) NOT NULL COMMENT 'Node ID from definition',
    `node_type` VARCHAR(32) NOT NULL COMMENT 'script|http|human|etc',
    `node_config` JSON COMMENT 'Node configuration snapshot',
    `status` VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT 'pending|running|succeeded|failed|waiting_human|skipped',
    `input_data` JSON COMMENT 'Input data for this step',
    `output_data` JSON COMMENT 'Output data from this step',
    `error_message` TEXT,
    `started_at` TIMESTAMP NULL,
    `finished_at` TIMESTAMP NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_run_id` (`run_id`),
    INDEX `idx_node_id` (`node_id`),
    INDEX `idx_status` (`status`),
    UNIQUE KEY `uk_run_node` (`run_id`, `node_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SELECT 'Octopus database initialized with workflows, workflow_runs, workflow_run_steps tables' AS message;
