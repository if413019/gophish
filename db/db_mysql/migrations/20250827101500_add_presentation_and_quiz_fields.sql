-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Add presentation_url and quiz_id fields to course_modules table
ALTER TABLE `course_modules` ADD COLUMN `presentation_url` varchar(500);
ALTER TABLE `course_modules` ADD COLUMN `quiz_id` bigint(20);

-- Create index for quiz_id lookups
CREATE INDEX `idx_course_modules_quiz_id` ON `course_modules`(`quiz_id`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- Remove the added columns
ALTER TABLE `course_modules` DROP COLUMN `presentation_url`;
ALTER TABLE `course_modules` DROP COLUMN `quiz_id`;

-- Drop the index
DROP INDEX `idx_course_modules_quiz_id` ON `course_modules`;