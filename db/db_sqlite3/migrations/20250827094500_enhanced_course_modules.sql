-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Drop the old course_modules table and recreate with enhanced structure
DROP TABLE IF EXISTS course_modules;

-- Create enhanced course_modules table with module types
CREATE TABLE IF NOT EXISTS "course_modules" (
    "id" integer primary key autoincrement,
    "course_id" bigint NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" text,
    "module_type" varchar(50) NOT NULL DEFAULT 'html', -- 'html', 'video', 'quiz'
    "content" text, -- HTML content for html modules, video URL/embed for video modules
    "video_url" varchar(500), -- Direct video URL
    "video_duration" integer, -- Duration in seconds
    "must_complete" boolean DEFAULT 1, -- Whether module must be completed to progress
    "min_time_spent" integer DEFAULT 0, -- Minimum time in seconds user must spend
    "order_index" integer DEFAULT 0,
    "created_date" datetime NOT NULL,
    "modified_date" datetime,
    FOREIGN KEY ("course_id") REFERENCES "courses"("id")
);

-- Update course_quizzes to be standalone (not embedded in modules)
-- Add more quiz configuration options
ALTER TABLE course_quizzes ADD COLUMN "module_id" bigint;
ALTER TABLE course_quizzes ADD COLUMN "time_limit" integer; -- Time limit in minutes, 0 = no limit
ALTER TABLE course_quizzes ADD COLUMN "max_attempts" integer DEFAULT 3; -- Maximum attempts allowed
ALTER TABLE course_quizzes ADD COLUMN "shuffle_questions" boolean DEFAULT 0; -- Randomize question order
ALTER TABLE course_quizzes ADD COLUMN "show_results" boolean DEFAULT 1; -- Show results after completion

-- Create module_progress table with enhanced tracking
DROP TABLE IF EXISTS module_progress;
CREATE TABLE IF NOT EXISTS "module_progress" (
    "id" integer primary key autoincrement,
    "enrollment_id" bigint NOT NULL,
    "module_id" bigint NOT NULL,
    "started_date" datetime,
    "completed_date" datetime,
    "time_spent" integer DEFAULT 0, -- Time spent in seconds
    "progress_percentage" integer DEFAULT 0, -- 0-100
    "status" varchar(50) DEFAULT 'not_started', -- 'not_started', 'in_progress', 'completed'
    "completion_data" text, -- JSON data for module-specific completion info
    FOREIGN KEY ("enrollment_id") REFERENCES "course_enrollments"("id"),
    FOREIGN KEY ("module_id") REFERENCES "course_modules"("id"),
    UNIQUE("enrollment_id", "module_id")
);

-- Enhanced quiz_attempts with more tracking
ALTER TABLE quiz_attempts ADD COLUMN "time_taken" integer; -- Time taken in seconds
ALTER TABLE quiz_attempts ADD COLUMN "attempt_number" integer DEFAULT 1;
ALTER TABLE quiz_attempts ADD COLUMN "started_date" datetime;
ALTER TABLE quiz_attempts ADD COLUMN "completed_date" datetime;

-- Add question types and enhanced options
ALTER TABLE course_questions ADD COLUMN "question_type" varchar(50) DEFAULT 'multiple_choice'; -- 'multiple_choice', 'true_false', 'text'
ALTER TABLE course_questions ADD COLUMN "points" integer DEFAULT 1; -- Points for this question
ALTER TABLE course_questions ADD COLUMN "explanation" text; -- Explanation shown after answering

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS "idx_course_modules_course_id_type" ON "course_modules"("course_id", "module_type");
CREATE INDEX IF NOT EXISTS "idx_course_modules_order" ON "course_modules"("course_id", "order_index");
CREATE INDEX IF NOT EXISTS "idx_module_progress_enrollment_module" ON "module_progress"("enrollment_id", "module_id");
CREATE INDEX IF NOT EXISTS "idx_module_progress_status" ON "module_progress"("status");

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- Remove the added columns
ALTER TABLE course_quizzes DROP COLUMN IF EXISTS "module_id";
ALTER TABLE course_quizzes DROP COLUMN IF EXISTS "time_limit";
ALTER TABLE course_quizzes DROP COLUMN IF EXISTS "max_attempts";
ALTER TABLE course_quizzes DROP COLUMN IF EXISTS "shuffle_questions";
ALTER TABLE course_quizzes DROP COLUMN IF EXISTS "show_results";

ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS "time_taken";
ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS "attempt_number";
ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS "started_date";
ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS "completed_date";

ALTER TABLE course_questions DROP COLUMN IF EXISTS "question_type";
ALTER TABLE course_questions DROP COLUMN IF EXISTS "points";
ALTER TABLE course_questions DROP COLUMN IF EXISTS "explanation";

-- Recreate original tables
DROP TABLE IF EXISTS "module_progress";
CREATE TABLE IF NOT EXISTS "module_progress" (
    "id" integer primary key autoincrement,
    "enrollment_id" bigint NOT NULL,
    "module_id" bigint NOT NULL,
    "completed_date" datetime,
    "status" varchar(50) DEFAULT 'not_started',
    FOREIGN KEY ("enrollment_id") REFERENCES "course_enrollments"("id"),
    FOREIGN KEY ("module_id") REFERENCES "course_modules"("id")
);

DROP TABLE IF EXISTS "course_modules";
CREATE TABLE IF NOT EXISTS "course_modules" (
    "id" integer primary key autoincrement,
    "course_id" bigint NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" text,
    "content" text,
    "order_index" integer DEFAULT 0,
    "created_date" datetime NOT NULL,
    FOREIGN KEY ("course_id") REFERENCES "courses"("id")
);