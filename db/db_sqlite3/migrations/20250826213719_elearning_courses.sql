-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Create courses table
CREATE TABLE IF NOT EXISTS "courses" (
    "id" integer primary key autoincrement,
    "user_id" bigint NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" text,
    "created_date" datetime NOT NULL,
    "modified_date" datetime NOT NULL
);

-- Create course_modules table
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

-- Create course_quizzes table
CREATE TABLE IF NOT EXISTS "course_quizzes" (
    "id" integer primary key autoincrement,
    "course_id" bigint NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" text,
    "order_index" integer DEFAULT 0,
    "passing_score" integer DEFAULT 70,
    "created_date" datetime NOT NULL,
    FOREIGN KEY ("course_id") REFERENCES "courses"("id")
);

-- Create course_questions table
CREATE TABLE IF NOT EXISTS "course_questions" (
    "id" integer primary key autoincrement,
    "quiz_id" bigint NOT NULL,
    "question" text NOT NULL,
    "order_index" integer DEFAULT 0,
    FOREIGN KEY ("quiz_id") REFERENCES "course_quizzes"("id")
);

-- Create question_options table
CREATE TABLE IF NOT EXISTS "question_options" (
    "id" integer primary key autoincrement,
    "question_id" bigint NOT NULL,
    "option" varchar(500) NOT NULL,
    "is_correct" boolean DEFAULT 0,
    "order_index" integer DEFAULT 0,
    FOREIGN KEY ("question_id") REFERENCES "course_questions"("id")
);

-- Create course_enrollments table
CREATE TABLE IF NOT EXISTS "course_enrollments" (
    "id" integer primary key autoincrement,
    "user_id" bigint NOT NULL,
    "course_id" bigint NOT NULL,
    "campaign_id" bigint,
    "enrolled_date" datetime NOT NULL,
    "started_date" datetime,
    "completed_date" datetime,
    "status" varchar(50) NOT NULL DEFAULT 'enrolled',
    "progress" integer DEFAULT 0,
    FOREIGN KEY ("course_id") REFERENCES "courses"("id"),
    UNIQUE("user_id", "course_id")
);

-- Create module_progress table
CREATE TABLE IF NOT EXISTS "module_progress" (
    "id" integer primary key autoincrement,
    "enrollment_id" bigint NOT NULL,
    "module_id" bigint NOT NULL,
    "completed_date" datetime,
    "status" varchar(50) DEFAULT 'not_started',
    FOREIGN KEY ("enrollment_id") REFERENCES "course_enrollments"("id"),
    FOREIGN KEY ("module_id") REFERENCES "course_modules"("id")
);

-- Create quiz_attempts table
CREATE TABLE IF NOT EXISTS "quiz_attempts" (
    "id" integer primary key autoincrement,
    "enrollment_id" bigint NOT NULL,
    "quiz_id" bigint NOT NULL,
    "attempt_date" datetime NOT NULL,
    "score" integer DEFAULT 0,
    "passed" boolean DEFAULT 0,
    FOREIGN KEY ("enrollment_id") REFERENCES "course_enrollments"("id"),
    FOREIGN KEY ("quiz_id") REFERENCES "course_quizzes"("id")
);

-- Create quiz_responses table
CREATE TABLE IF NOT EXISTS "quiz_responses" (
    "id" integer primary key autoincrement,
    "attempt_id" bigint NOT NULL,
    "question_id" bigint NOT NULL,
    "selected_option_id" bigint NOT NULL,
    "is_correct" boolean DEFAULT 0,
    FOREIGN KEY ("attempt_id") REFERENCES "quiz_attempts"("id"),
    FOREIGN KEY ("question_id") REFERENCES "course_questions"("id"),
    FOREIGN KEY ("selected_option_id") REFERENCES "question_options"("id")
);

-- Add course_id to campaigns table for auto-enrollment
ALTER TABLE "campaigns" ADD COLUMN "course_id" bigint;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS "idx_courses_user_id" ON "courses"("user_id");
CREATE INDEX IF NOT EXISTS "idx_course_modules_course_id" ON "course_modules"("course_id");
CREATE INDEX IF NOT EXISTS "idx_course_quizzes_course_id" ON "course_quizzes"("course_id");
CREATE INDEX IF NOT EXISTS "idx_course_questions_quiz_id" ON "course_questions"("quiz_id");
CREATE INDEX IF NOT EXISTS "idx_question_options_question_id" ON "question_options"("question_id");
CREATE INDEX IF NOT EXISTS "idx_course_enrollments_user_id" ON "course_enrollments"("user_id");
CREATE INDEX IF NOT EXISTS "idx_course_enrollments_course_id" ON "course_enrollments"("course_id");
CREATE INDEX IF NOT EXISTS "idx_module_progress_enrollment_id" ON "module_progress"("enrollment_id");
CREATE INDEX IF NOT EXISTS "idx_quiz_attempts_enrollment_id" ON "quiz_attempts"("enrollment_id");
CREATE INDEX IF NOT EXISTS "idx_quiz_responses_attempt_id" ON "quiz_responses"("attempt_id");

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- Remove course_id column from campaigns
ALTER TABLE "campaigns" DROP COLUMN "course_id";

-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS "quiz_responses";
DROP TABLE IF EXISTS "quiz_attempts";
DROP TABLE IF EXISTS "module_progress";
DROP TABLE IF EXISTS "course_enrollments";
DROP TABLE IF EXISTS "question_options";
DROP TABLE IF EXISTS "course_questions";
DROP TABLE IF EXISTS "course_quizzes";
DROP TABLE IF EXISTS "course_modules";
DROP TABLE IF EXISTS "courses";