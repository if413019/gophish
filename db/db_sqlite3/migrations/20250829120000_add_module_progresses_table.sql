-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE module_progresses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    enrollment_id INTEGER NOT NULL,
    module_id INTEGER NOT NULL,
    started_date DATETIME,
    completed_date DATETIME,
    time_spent INTEGER DEFAULT 0,
    progress_percentage INTEGER DEFAULT 0,
    status TEXT DEFAULT 'not_started',
    completion_data TEXT,
    created_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_date DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for better performance
CREATE INDEX idx_module_progresses_enrollment ON module_progresses(enrollment_id);
CREATE INDEX idx_module_progresses_module ON module_progresses(module_id);
CREATE UNIQUE INDEX idx_module_progresses_enrollment_module ON module_progresses(enrollment_id, module_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX idx_module_progresses_enrollment_module;
DROP INDEX idx_module_progresses_module;
DROP INDEX idx_module_progresses_enrollment;
DROP TABLE module_progresses;