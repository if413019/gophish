-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Drop the old course_modules table and recreate with enhanced structure
DROP TABLE IF EXISTS `course_modules`;

-- Create enhanced course_modules table with module types
CREATE TABLE IF NOT EXISTS `course_modules` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `course_id` bigint(20) NOT NULL,
    `name` varchar(255) NOT NULL,
    `description` text,
    `module_type` varchar(50) NOT NULL DEFAULT 'html',
    `content` longtext,
    `video_url` varchar(500),
    `video_duration` int(11),
    `must_complete` tinyint(1) DEFAULT 1,
    `min_time_spent` int(11) DEFAULT 0,
    `order_index` int(11) DEFAULT 0,
    `created_date` datetime NOT NULL,
    `modified_date` datetime,
    PRIMARY KEY (`id`),
    INDEX `idx_course_modules_course_id_type` (`course_id`, `module_type`),
    INDEX `idx_course_modules_order` (`course_id`, `order_index`),
    FOREIGN KEY (`course_id`) REFERENCES `courses`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Update course_quizzes to be standalone with enhanced options
ALTER TABLE `course_quizzes` ADD COLUMN `module_id` bigint(20) NULL;
ALTER TABLE `course_quizzes` ADD COLUMN `time_limit` int(11) NULL;
ALTER TABLE `course_quizzes` ADD COLUMN `max_attempts` int(11) DEFAULT 3;
ALTER TABLE `course_quizzes` ADD COLUMN `shuffle_questions` tinyint(1) DEFAULT 0;
ALTER TABLE `course_quizzes` ADD COLUMN `show_results` tinyint(1) DEFAULT 1;

-- Create enhanced module_progress table
DROP TABLE IF EXISTS `module_progress`;
CREATE TABLE IF NOT EXISTS `module_progress` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `enrollment_id` bigint(20) NOT NULL,
    `module_id` bigint(20) NOT NULL,
    `started_date` datetime NULL,
    `completed_date` datetime NULL,
    `time_spent` int(11) DEFAULT 0,
    `progress_percentage` int(11) DEFAULT 0,
    `status` varchar(50) DEFAULT 'not_started',
    `completion_data` text,
    PRIMARY KEY (`id`),
    INDEX `idx_module_progress_enrollment_module` (`enrollment_id`, `module_id`),
    INDEX `idx_module_progress_status` (`status`),
    UNIQUE KEY `unique_enrollment_module` (`enrollment_id`, `module_id`),
    FOREIGN KEY (`enrollment_id`) REFERENCES `course_enrollments`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`module_id`) REFERENCES `course_modules`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Enhanced quiz_attempts with more tracking
ALTER TABLE `quiz_attempts` ADD COLUMN `time_taken` int(11) NULL;
ALTER TABLE `quiz_attempts` ADD COLUMN `attempt_number` int(11) DEFAULT 1;
ALTER TABLE `quiz_attempts` ADD COLUMN `started_date` datetime NULL;
ALTER TABLE `quiz_attempts` ADD COLUMN `completed_date` datetime NULL;

-- Add question types and enhanced options
ALTER TABLE `course_questions` ADD COLUMN `question_type` varchar(50) DEFAULT 'multiple_choice';
ALTER TABLE `course_questions` ADD COLUMN `points` int(11) DEFAULT 1;
ALTER TABLE `course_questions` ADD COLUMN `explanation` text NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- Remove the added columns
ALTER TABLE `course_quizzes` DROP COLUMN IF EXISTS `module_id`;
ALTER TABLE `course_quizzes` DROP COLUMN IF EXISTS `time_limit`;
ALTER TABLE `course_quizzes` DROP COLUMN IF EXISTS `max_attempts`;
ALTER TABLE `course_quizzes` DROP COLUMN IF EXISTS `shuffle_questions`;
ALTER TABLE `course_quizzes` DROP COLUMN IF EXISTS `show_results`;

ALTER TABLE `quiz_attempts` DROP COLUMN IF EXISTS `time_taken`;
ALTER TABLE `quiz_attempts` DROP COLUMN IF EXISTS `attempt_number`;
ALTER TABLE `quiz_attempts` DROP COLUMN IF EXISTS `started_date`;
ALTER TABLE `quiz_attempts` DROP COLUMN IF EXISTS `completed_date`;

ALTER TABLE `course_questions` DROP COLUMN IF EXISTS `question_type`;
ALTER TABLE `course_questions` DROP COLUMN IF EXISTS `points`;
ALTER TABLE `course_questions` DROP COLUMN IF EXISTS `explanation`;

-- Recreate original tables
DROP TABLE IF EXISTS `module_progress`;
CREATE TABLE IF NOT EXISTS `module_progress` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `enrollment_id` bigint(20) NOT NULL,
    `module_id` bigint(20) NOT NULL,
    `completed_date` datetime NULL,
    `status` varchar(50) DEFAULT 'not_started',
    PRIMARY KEY (`id`),
    FOREIGN KEY (`enrollment_id`) REFERENCES `course_enrollments`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`module_id`) REFERENCES `course_modules`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `course_modules`;
CREATE TABLE IF NOT EXISTS `course_modules` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `course_id` bigint(20) NOT NULL,
    `name` varchar(255) NOT NULL,
    `description` text,
    `content` text,
    `order_index` int(11) DEFAULT 0,
    `created_date` datetime NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`course_id`) REFERENCES `courses`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;