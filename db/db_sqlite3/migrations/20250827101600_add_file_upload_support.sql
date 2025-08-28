-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Add columns for local file storage
ALTER TABLE course_modules ADD COLUMN video_file_path TEXT;
ALTER TABLE course_modules ADD COLUMN presentation_file_path TEXT;
ALTER TABLE course_modules ADD COLUMN original_filename TEXT;
ALTER TABLE course_modules ADD COLUMN file_size INTEGER;
ALTER TABLE course_modules ADD COLUMN mime_type TEXT;
ALTER TABLE course_modules ADD COLUMN uploaded_date DATETIME;

-- Add index for file lookups
CREATE INDEX idx_course_modules_video_file ON course_modules(video_file_path);
CREATE INDEX idx_course_modules_presentation_file ON course_modules(presentation_file_path);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX idx_course_modules_presentation_file;
DROP INDEX idx_course_modules_video_file;

-- Remove added columns
ALTER TABLE course_modules DROP COLUMN uploaded_date;
ALTER TABLE course_modules DROP COLUMN mime_type;
ALTER TABLE course_modules DROP COLUMN file_size;
ALTER TABLE course_modules DROP COLUMN original_filename;
ALTER TABLE course_modules DROP COLUMN presentation_file_path;
ALTER TABLE course_modules DROP COLUMN video_file_path;