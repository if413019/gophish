-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Create elearning_settings table (singleton - only one row with id=1)
CREATE TABLE IF NOT EXISTS "elearning_settings" (
    "id" integer primary key,
    "smtp_id" bigint,
    "email_subject" varchar(255) NOT NULL DEFAULT 'Security Awareness Training Required',
    "email_html" text,
    "base_url" varchar(500) NOT NULL DEFAULT 'https://localhost:3333',
    "company_name" varchar(255) NOT NULL DEFAULT 'Your Organization',
    "modified_date" datetime NOT NULL,
    FOREIGN KEY ("smtp_id") REFERENCES "smtp"("id")
);

-- Create index for smtp_id lookup
CREATE INDEX IF NOT EXISTS "idx_elearning_settings_smtp_id" ON "elearning_settings"("smtp_id");

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS "elearning_settings";
