-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Create courses table
CREATE TABLE IF NOT EXISTS `courses` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL,
    `name` varchar(255) NOT NULL,
    `description` text,
    `created_date` datetime NOT NULL,
    `modified_date` datetime NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_courses_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create course_modules table
CREATE TABLE IF NOT EXISTS `course_modules` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `course_id` bigint(20) NOT NULL,
    `name` varchar(255) NOT NULL,
    `description` text,
    `content` longtext,
    `order_index` int(11) DEFAULT 0,
    `created_date` datetime NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_course_modules_course_id` (`course_id`),
    FOREIGN KEY (`course_id`) REFERENCES `courses`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create course_quizzes table
CREATE TABLE IF NOT EXISTS `course_quizzes` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `course_id` bigint(20) NOT NULL,
    `name` varchar(255) NOT NULL,
    `description` text,
    `order_index` int(11) DEFAULT 0,
    `passing_score` int(11) DEFAULT 70,
    `created_date` datetime NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_course_quizzes_course_id` (`course_id`),
    FOREIGN KEY (`course_id`) REFERENCES `courses`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create course_questions table
CREATE TABLE IF NOT EXISTS `course_questions` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `quiz_id` bigint(20) NOT NULL,
    `question` text NOT NULL,
    `order_index` int(11) DEFAULT 0,
    PRIMARY KEY (`id`),
    INDEX `idx_course_questions_quiz_id` (`quiz_id`),
    FOREIGN KEY (`quiz_id`) REFERENCES `course_quizzes`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create question_options table
CREATE TABLE IF NOT EXISTS `question_options` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `question_id` bigint(20) NOT NULL,
    `option` varchar(500) NOT NULL,
    `is_correct` tinyint(1) DEFAULT 0,
    `order_index` int(11) DEFAULT 0,
    PRIMARY KEY (`id`),
    INDEX `idx_question_options_question_id` (`question_id`),
    FOREIGN KEY (`question_id`) REFERENCES `course_questions`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create course_enrollments table
CREATE TABLE IF NOT EXISTS `course_enrollments` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL,
    `course_id` bigint(20) NOT NULL,
    `campaign_id` bigint(20) NULL,
    `enrolled_date` datetime NOT NULL,
    `started_date` datetime NULL,
    `completed_date` datetime NULL,
    `status` varchar(50) NOT NULL DEFAULT 'enrolled',
    `progress` int(11) DEFAULT 0,
    PRIMARY KEY (`id`),
    INDEX `idx_course_enrollments_user_id` (`user_id`),
    INDEX `idx_course_enrollments_course_id` (`course_id`),
    UNIQUE KEY `unique_user_course` (`user_id`, `course_id`),
    FOREIGN KEY (`course_id`) REFERENCES `courses`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create module_progress table
CREATE TABLE IF NOT EXISTS `module_progress` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `enrollment_id` bigint(20) NOT NULL,
    `module_id` bigint(20) NOT NULL,
    `completed_date` datetime NULL,
    `status` varchar(50) DEFAULT 'not_started',
    PRIMARY KEY (`id`),
    INDEX `idx_module_progress_enrollment_id` (`enrollment_id`),
    FOREIGN KEY (`enrollment_id`) REFERENCES `course_enrollments`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`module_id`) REFERENCES `course_modules`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create quiz_attempts table
CREATE TABLE IF NOT EXISTS `quiz_attempts` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `enrollment_id` bigint(20) NOT NULL,
    `quiz_id` bigint(20) NOT NULL,
    `attempt_date` datetime NOT NULL,
    `score` int(11) DEFAULT 0,
    `passed` tinyint(1) DEFAULT 0,
    PRIMARY KEY (`id`),
    INDEX `idx_quiz_attempts_enrollment_id` (`enrollment_id`),
    FOREIGN KEY (`enrollment_id`) REFERENCES `course_enrollments`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`quiz_id`) REFERENCES `course_quizzes`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create quiz_responses table
CREATE TABLE IF NOT EXISTS `quiz_responses` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `attempt_id` bigint(20) NOT NULL,
    `question_id` bigint(20) NOT NULL,
    `selected_option_id` bigint(20) NOT NULL,
    `is_correct` tinyint(1) DEFAULT 0,
    PRIMARY KEY (`id`),
    INDEX `idx_quiz_responses_attempt_id` (`attempt_id`),
    FOREIGN KEY (`attempt_id`) REFERENCES `quiz_attempts`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`question_id`) REFERENCES `course_questions`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`selected_option_id`) REFERENCES `question_options`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Add course_id to campaigns table for auto-enrollment
ALTER TABLE `campaigns` ADD COLUMN `course_id` bigint(20) NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- Remove course_id column from campaigns
ALTER TABLE `campaigns` DROP COLUMN `course_id`;

-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS `quiz_responses`;
DROP TABLE IF EXISTS `quiz_attempts`;
DROP TABLE IF EXISTS `module_progress`;
DROP TABLE IF EXISTS `course_enrollments`;
DROP TABLE IF EXISTS `question_options`;
DROP TABLE IF EXISTS `course_questions`;
DROP TABLE IF EXISTS `course_quizzes`;
DROP TABLE IF EXISTS `course_modules`;
DROP TABLE IF EXISTS `courses`;