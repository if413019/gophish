-- Dummy E-Learning Courses Data for Testing
-- This script creates sample security awareness courses

-- Course 1: Email Security Basics
INSERT INTO courses (id, user_id, name, description, created_date, modified_date) VALUES 
(1, 1, 'Email Security Fundamentals', 'Learn the basics of email security and how to identify phishing attempts', datetime('now'), datetime('now'));

-- Modules for Email Security Course
INSERT INTO course_modules (id, course_id, name, description, content, order_index, created_date) VALUES 
(1, 1, 'Introduction to Email Threats', 'Understanding different types of email-based attacks', '<h3>Email Threats Overview</h3><p>Email is one of the most common attack vectors used by cybercriminals. In this module, you will learn about:</p><ul><li>Phishing attacks</li><li>Malware distribution</li><li>Social engineering tactics</li><li>Business Email Compromise (BEC)</li></ul><p>Understanding these threats is the first step in protecting yourself and your organization.</p>', 0, datetime('now')),
(2, 1, 'Identifying Suspicious Emails', 'How to spot red flags in emails', '<h3>Red Flags to Watch For</h3><p>Learn to identify suspicious emails by looking for these warning signs:</p><ul><li>Urgent or threatening language</li><li>Requests for personal information</li><li>Suspicious sender addresses</li><li>Poor grammar and spelling</li><li>Unexpected attachments</li><li>Generic greetings</li></ul><p>Always verify the sender through alternative means when in doubt.</p>', 1, datetime('now')),
(3, 1, 'Best Practices for Email Security', 'Implement security measures for safe email usage', '<h3>Email Security Best Practices</h3><p>Follow these guidelines to stay secure:</p><ul><li>Never click suspicious links</li><li>Verify requests for sensitive information</li><li>Use strong, unique passwords</li><li>Enable two-factor authentication</li><li>Keep software updated</li><li>Report suspicious emails</li></ul><p>Remember: When in doubt, dont click!</p>', 2, datetime('now'));

-- Quiz for Email Security Course
INSERT INTO course_quizzes (id, course_id, name, description, order_index, passing_score, created_date) VALUES 
(1, 1, 'Email Security Knowledge Check', 'Test your understanding of email security concepts', 0, 75, datetime('now'));

-- Questions for Email Security Quiz
INSERT INTO course_questions (id, quiz_id, question, order_index) VALUES 
(1, 1, 'Which of the following is a common sign of a phishing email?', 0),
(2, 1, 'What should you do if you receive an unexpected email asking for your password?', 1),
(3, 1, 'What is the best way to verify a suspicious email from your bank?', 2);

-- Options for Question 1
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(1, 1, 'Urgent language demanding immediate action', 1, 0),
(2, 1, 'Professional formatting and company logo', 0, 1),
(3, 1, 'Detailed contact information in signature', 0, 2),
(4, 1, 'Personalized greeting with your full name', 0, 3);

-- Options for Question 2
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(5, 2, 'Reply immediately with your password', 0, 0),
(6, 2, 'Delete the email and report it as suspicious', 1, 1),
(7, 2, 'Forward it to your colleagues for verification', 0, 2),
(8, 2, 'Click any links to verify its legitimate', 0, 3);

-- Options for Question 3
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(9, 3, 'Reply to the email asking for confirmation', 0, 0),
(10, 3, 'Call the bank using the number on your card or statement', 1, 1),
(11, 3, 'Click the link in the email to verify', 0, 2),
(12, 3, 'Forward the email to friends for opinions', 0, 3);

-- Course 2: Password Security
INSERT INTO courses (id, user_id, name, description, created_date, modified_date) VALUES 
(2, 1, 'Password Security Mastery', 'Master the art of creating and managing secure passwords', datetime('now'), datetime('now'));

-- Modules for Password Security Course
INSERT INTO course_modules (id, course_id, name, description, content, order_index, created_date) VALUES 
(4, 2, 'Why Passwords Matter', 'Understanding the importance of strong passwords', '<h3>The Password Problem</h3><p>Passwords are your first line of defense against cyber attacks. Weak passwords can lead to:</p><ul><li>Account takeovers</li><li>Identity theft</li><li>Data breaches</li><li>Financial losses</li></ul><p>In this course, youll learn how to create and manage strong, unique passwords.</p>', 0, datetime('now')),
(5, 2, 'Creating Strong Passwords', 'Learn techniques for creating uncrackable passwords', '<h3>Strong Password Guidelines</h3><p>A strong password should be:</p><ul><li>At least 12 characters long</li><li>Include uppercase and lowercase letters</li><li>Contain numbers and special characters</li><li>Avoid dictionary words</li><li>Be unique for each account</li></ul><p>Consider using passphrases like "Coffee$Sunrise#2024!" instead of complex passwords.</p>', 1, datetime('now')),
(6, 2, 'Password Managers', 'Simplify password security with password managers', '<h3>Password Manager Benefits</h3><p>Password managers help you:</p><ul><li>Generate strong, unique passwords</li><li>Store passwords securely</li><li>Auto-fill login forms</li><li>Sync across devices</li><li>Monitor for breached passwords</li></ul><p>Popular options include: Bitwarden, 1Password, LastPass, and Dashlane.</p>', 2, datetime('now'));

-- Quiz for Password Security Course
INSERT INTO course_quizzes (id, course_id, name, description, order_index, passing_score, created_date) VALUES 
(2, 2, 'Password Security Assessment', 'Evaluate your password security knowledge', 0, 80, datetime('now'));

-- Questions for Password Security Quiz
INSERT INTO course_questions (id, quiz_id, question, order_index) VALUES 
(4, 2, 'What makes a password strong?', 0),
(5, 2, 'How often should you change your passwords?', 1);

-- Options for Password Question 1
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(13, 4, 'Length, complexity, and uniqueness', 1, 0),
(14, 4, 'Using your pets name with numbers', 0, 1),
(15, 4, 'Including your birthday and address', 0, 2),
(16, 4, 'Using the same secure password everywhere', 0, 3);

-- Options for Password Question 2
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(17, 5, 'Every 30 days regardless of circumstances', 0, 0),
(18, 5, 'Only when there is evidence of compromise', 1, 1),
(19, 5, 'Never, once set they should remain the same', 0, 2),
(20, 5, 'Every time you log in', 0, 3);

-- Course 3: Social Engineering Awareness
INSERT INTO courses (id, user_id, name, description, created_date, modified_date) VALUES 
(3, 1, 'Social Engineering Defense', 'Recognize and defend against social engineering attacks', datetime('now'), datetime('now'));

-- Module for Social Engineering Course
INSERT INTO course_modules (id, course_id, name, description, content, order_index, created_date) VALUES 
(7, 3, 'Understanding Social Engineering', 'What is social engineering and how it works', '<h3>Social Engineering Tactics</h3><p>Social engineering exploits human psychology rather than technical vulnerabilities. Common tactics include:</p><ul><li>Pretexting (creating fake scenarios)</li><li>Baiting (offering something enticing)</li><li>Tailgating (following someone into secure areas)</li><li>Quid pro quo (offering a service for information)</li></ul><p>The best defense is awareness and verification.</p>', 0, datetime('now'));

-- Quiz for Social Engineering Course
INSERT INTO course_quizzes (id, course_id, name, description, order_index, passing_score, created_date) VALUES 
(3, 3, 'Social Engineering Recognition Test', 'Can you spot social engineering attempts?', 0, 70, datetime('now'));

-- Question for Social Engineering Quiz
INSERT INTO course_questions (id, quiz_id, question, order_index) VALUES 
(6, 3, 'Someone calls claiming to be from IT and asks for your password to "update the system". What should you do?', 0);

-- Options for Social Engineering Question
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(21, 6, 'Provide the password since they said theyre from IT', 0, 0),
(22, 6, 'Hang up and contact IT through official channels to verify', 1, 1),
(23, 6, 'Ask them to prove their identity by telling you your password first', 0, 2),
(24, 6, 'Give them a fake password to test if theyre legitimate', 0, 3);