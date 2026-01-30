-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Create elearning_settings table (singleton - only one row with id=1)
CREATE TABLE IF NOT EXISTS `elearning_settings` (
    `id` int(11) NOT NULL,
    `smtp_id` bigint(20) DEFAULT NULL,
    `email_subject` varchar(255) NOT NULL DEFAULT 'Security Awareness Training Required',
    `email_html` text,
    `base_url` varchar(500) NOT NULL DEFAULT 'https://localhost:3333',
    `company_name` varchar(255) NOT NULL DEFAULT 'Your Organization',
    `modified_date` datetime NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_elearning_settings_smtp_id` (`smtp_id`),
    CONSTRAINT `fk_elearning_settings_smtp` FOREIGN KEY (`smtp_id`) REFERENCES `smtp` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS `elearning_settings`;
