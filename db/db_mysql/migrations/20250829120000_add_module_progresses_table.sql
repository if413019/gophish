-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE module_progresses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    enrollment_id BIGINT NOT NULL,
    module_id BIGINT NOT NULL,
    started_date DATETIME NULL,
    completed_date DATETIME NULL,
    time_spent INT DEFAULT 0,
    progress_percentage INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'not_started',
    completion_data TEXT,
    created_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_date DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Add indexes for better performance
CREATE INDEX idx_module_progresses_enrollment ON module_progresses(enrollment_id);
CREATE INDEX idx_module_progresses_module ON module_progresses(module_id);
CREATE UNIQUE INDEX idx_module_progresses_enrollment_module ON module_progresses(enrollment_id, module_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX idx_module_progresses_enrollment_module ON module_progresses;
DROP INDEX idx_module_progresses_module ON module_progresses;
DROP INDEX idx_module_progresses_enrollment ON module_progresses;
DROP TABLE module_progresses;