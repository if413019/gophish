-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Add columns for local file storage
ALTER TABLE course_modules 
ADD COLUMN video_file_path TEXT,
ADD COLUMN presentation_file_path TEXT,
ADD COLUMN original_filename VARCHAR(255),
ADD COLUMN file_size BIGINT,
ADD COLUMN mime_type VARCHAR(100),
ADD COLUMN uploaded_date DATETIME;

-- Add index for file lookups
CREATE INDEX idx_course_modules_video_file ON course_modules(video_file_path(255));
CREATE INDEX idx_course_modules_presentation_file ON course_modules(presentation_file_path(255));

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX idx_course_modules_presentation_file ON course_modules;
DROP INDEX idx_course_modules_video_file ON course_modules;

-- Remove added columns
ALTER TABLE course_modules 
DROP COLUMN uploaded_date,
DROP COLUMN mime_type,
DROP COLUMN file_size,
DROP COLUMN original_filename,
DROP COLUMN presentation_file_path,
DROP COLUMN video_file_path;