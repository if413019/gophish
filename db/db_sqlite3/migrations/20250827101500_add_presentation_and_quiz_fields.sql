-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Add presentation_url and quiz_id fields to course_modules table
ALTER TABLE course_modules ADD COLUMN "presentation_url" varchar(500);
ALTER TABLE course_modules ADD COLUMN "quiz_id" bigint;

-- Create index for quiz_id lookups
CREATE INDEX IF NOT EXISTS "idx_course_modules_quiz_id" ON "course_modules"("quiz_id");

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- Remove the added columns
ALTER TABLE course_modules DROP COLUMN IF EXISTS "presentation_url";
ALTER TABLE course_modules DROP COLUMN IF EXISTS "quiz_id";

-- Drop the index
DROP INDEX IF EXISTS "idx_course_modules_quiz_id";