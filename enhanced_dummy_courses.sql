-- Enhanced E-Learning Courses Data with Different Module Types
-- This script creates sample courses with HTML, Video, and Quiz modules

-- Clear existing data first
DELETE FROM quiz_responses;
DELETE FROM quiz_attempts;
DELETE FROM question_options;
DELETE FROM course_questions;
DELETE FROM course_quizzes;
DELETE FROM module_progress;
DELETE FROM course_enrollments;
DELETE FROM course_modules;
DELETE FROM courses;

-- Course 1: Comprehensive Email Security Training
INSERT INTO courses (id, user_id, name, description, created_date, modified_date) VALUES 
(1, 1, 'Advanced Email Security Training', 'Complete email security course with interactive content, videos, and assessments', datetime('now'), datetime('now'));

-- Module 1: HTML Content Module
INSERT INTO course_modules (id, course_id, name, description, module_type, content, must_complete, min_time_spent, order_index, created_date, modified_date) VALUES 
(1, 1, 'Email Threats Introduction', 'Learn about various email-based security threats', 'html', 
'<div class="module-content">
<h2>Welcome to Email Security Training</h2>
<p>Email remains one of the most common attack vectors for cybercriminals. In this comprehensive training, you will learn:</p>
<ul>
<li><strong>Phishing Attacks:</strong> How attackers trick users into revealing sensitive information</li>
<li><strong>Malware Distribution:</strong> How malicious software spreads through email attachments</li>
<li><strong>Business Email Compromise (BEC):</strong> Sophisticated attacks targeting organizations</li>
<li><strong>Social Engineering:</strong> Psychological manipulation techniques used in emails</li>
</ul>
<div class="alert alert-info">
<i class="fa fa-info-circle"></i> <strong>Did you know?</strong> 91% of cyber attacks begin with a phishing email!
</div>
<h3>Types of Email Threats</h3>
<ol>
<li><strong>Spam:</strong> Unwanted bulk email that clogs your inbox</li>
<li><strong>Phishing:</strong> Fraudulent emails designed to steal credentials</li>
<li><strong>Spear Phishing:</strong> Targeted attacks against specific individuals</li>
<li><strong>Whaling:</strong> Attacks targeting high-profile executives</li>
<li><strong>Ransomware:</strong> Malicious software that encrypts your files</li>
</ol>
<p>Take your time to read through this material carefully. Understanding these concepts is crucial for the rest of the course.</p>
</div>', 
true, 120, 0, datetime('now'), datetime('now'));

-- Module 2: Video Module (must watch completely)
INSERT INTO course_modules (id, course_id, name, description, module_type, video_url, video_duration, must_complete, min_time_spent, order_index, created_date, modified_date) VALUES 
(2, 1, 'Real Phishing Examples Video', 'Watch real examples of phishing emails and learn how to spot them', 'video', 
'https://www.youtube.com/embed/dQw4w9WgXcQ', 300, true, 300, 1, datetime('now'), datetime('now'));

-- Module 3: Another HTML Module
INSERT INTO course_modules (id, course_id, name, description, module_type, content, must_complete, min_time_spent, order_index, created_date, modified_date) VALUES 
(3, 1, 'Email Security Best Practices', 'Learn practical steps to protect yourself', 'html',
'<div class="module-content">
<h2>Email Security Best Practices</h2>
<p>Follow these essential guidelines to protect yourself from email threats:</p>

<h3>🔒 Authentication and Verification</h3>
<ul>
<li>Always verify sender identity through alternative communication channels</li>
<li>Check email addresses carefully for typos or suspicious domains</li>
<li>Be wary of emails from unknown senders</li>
<li>Look for official company signatures and contact information</li>
</ul>

<h3>🔗 Link and Attachment Safety</h3>
<ul>
<li><strong>Never click suspicious links</strong> - hover to preview URLs first</li>
<li>Scan all attachments with antivirus software</li>
<li>Avoid downloading files from untrusted sources</li>
<li>Be extra cautious with .exe, .zip, and .pdf files</li>
</ul>

<h3>🛡️ Technical Safeguards</h3>
<ul>
<li>Enable two-factor authentication on all accounts</li>
<li>Keep your software and browser updated</li>
<li>Use reputable antivirus software</li>
<li>Enable spam filters</li>
</ul>

<div class="alert alert-warning">
<i class="fa fa-exclamation-triangle"></i> <strong>Remember:</strong> When in doubt, dont click! Its better to be safe than sorry.
</div>

<h3>🚨 What to Do If You Suspect Phishing</h3>
<ol>
<li>Do not click any links or download attachments</li>
<li>Report the email to your IT security team</li>
<li>Delete the email after reporting</li>
<li>If you accidentally clicked, immediately change your passwords</li>
<li>Run a full antivirus scan on your computer</li>
</ol>
</div>',
true, 180, 2, datetime('now'), datetime('now'));

-- Module 4: Quiz Module (standalone quiz as a module)
INSERT INTO course_modules (id, course_id, name, description, module_type, content, must_complete, order_index, created_date, modified_date) VALUES 
(4, 1, 'Email Security Knowledge Assessment', 'Test your understanding with this comprehensive quiz', 'quiz', 
'{"quiz_instructions": "This quiz will test your knowledge of email security concepts. You need 80% to pass. You have 3 attempts and 15 minutes per attempt."}',
true, 3, datetime('now'), datetime('now'));

-- Create the quiz linked to the quiz module
INSERT INTO course_quizzes (id, course_id, module_id, name, description, order_index, passing_score, time_limit, max_attempts, shuffle_questions, show_results, created_date) VALUES 
(1, 1, 4, 'Email Security Final Assessment', 'Comprehensive test of email security knowledge', 0, 80, 15, 3, 1, 1, datetime('now'));

-- Quiz Questions
INSERT INTO course_questions (id, quiz_id, question, question_type, points, explanation, order_index) VALUES 
(1, 1, 'Which of the following is the MOST reliable way to verify a suspicious email from your bank?', 'multiple_choice', 2, 
'Calling the bank using the official number on your card or statement is the most secure way to verify suspicious communications. Never use contact information from the suspicious email itself.', 0),

(2, 1, 'What should you do immediately if you accidentally click a suspicious link?', 'multiple_choice', 2,
'Acting quickly is essential. Change passwords, run antivirus scans, and inform IT security to prevent potential damage from the incident.', 1),

(3, 1, 'True or False: Its safe to click links in emails as long as they come from people you know.', 'true_false', 1,
'False. Email accounts can be compromised, and attackers often use hijacked accounts to send malicious emails to contacts. Always verify suspicious requests even from known senders.', 2),

(4, 1, 'Which email characteristic is LEAST likely to indicate a phishing attempt?', 'multiple_choice', 2,
'Professional formatting alone is not a reliable indicator. Modern phishing emails often mimic legitimate company designs very closely.', 3),

(5, 1, 'What is spear phishing?', 'multiple_choice', 1,
'Spear phishing uses personal information to create highly targeted and convincing attacks against specific individuals, making them much more dangerous than generic phishing attempts.', 4);

-- Question Options for Question 1
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(1, 1, 'Reply to the email asking for confirmation', 0, 0),
(2, 1, 'Click the link to verify it leads to the real bank website', 0, 1),
(3, 1, 'Call the bank using the phone number on your card or statement', 1, 2),
(4, 1, 'Forward the email to friends to get their opinion', 0, 3);

-- Question Options for Question 2
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(5, 2, 'Wait to see if anything bad happens', 0, 0),
(6, 2, 'Change your passwords, run antivirus scan, and notify IT', 1, 1),
(7, 2, 'Just close the browser and continue working', 0, 2),
(8, 2, 'Click the link again to see where it leads', 0, 3);

-- Question Options for Question 3 (True/False)
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(9, 3, 'True', 0, 0),
(10, 3, 'False', 1, 1);

-- Question Options for Question 4
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(11, 4, 'Urgent language demanding immediate action', 0, 0),
(12, 4, 'Requests for sensitive personal information', 0, 1),
(13, 4, 'Professional formatting and company logos', 1, 2),
(14, 4, 'Generic greetings like "Dear Customer"', 0, 3);

-- Question Options for Question 5
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(15, 5, 'Mass emails sent to many recipients simultaneously', 0, 0),
(16, 5, 'Targeted attacks using personal information about the victim', 1, 1),
(17, 5, 'Emails that contain malware attachments', 0, 2),
(18, 5, 'Fake emails pretending to be from social media sites', 0, 3);

-- Course 2: Video-Heavy Course
INSERT INTO courses (id, user_id, name, description, created_date, modified_date) VALUES 
(2, 1, 'Cybersecurity Fundamentals Video Course', 'Learn cybersecurity basics through engaging video content', datetime('now'), datetime('now'));

-- Video modules for course 2
INSERT INTO course_modules (id, course_id, name, description, module_type, video_url, video_duration, must_complete, min_time_spent, order_index, created_date, modified_date) VALUES 
(5, 2, 'Introduction to Cybersecurity', 'Overview of cybersecurity landscape', 'video', 'https://www.youtube.com/embed/inWWhr5tnEA', 480, true, 480, 0, datetime('now'), datetime('now')),
(6, 2, 'Password Security Explained', 'Deep dive into password security', 'video', 'https://www.youtube.com/embed/3NjQ9b3pgIg', 360, true, 360, 1, datetime('now'), datetime('now')),
(7, 2, 'Social Engineering Tactics', 'Understanding social engineering', 'video', 'https://www.youtube.com/embed/lc7scxvKQOo', 420, true, 420, 2, datetime('now'), datetime('now'));

-- HTML summary module
INSERT INTO course_modules (id, course_id, name, description, module_type, content, must_complete, min_time_spent, order_index, created_date, modified_date) VALUES 
(8, 2, 'Course Summary and Key Takeaways', 'Summary of all video content', 'html',
'<div class="module-content">
<h2>🎯 Key Takeaways</h2>
<p>Congratulations on completing the video modules! Here are the key points to remember:</p>

<h3>📚 What Youve Learned</h3>
<ul>
<li><strong>Cybersecurity Fundamentals:</strong> Understanding the threat landscape and basic security principles</li>
<li><strong>Password Security:</strong> Creating strong, unique passwords and using password managers</li>
<li><strong>Social Engineering:</strong> Recognizing psychological manipulation tactics used by attackers</li>
</ul>

<h3>🔄 Apply What You Learned</h3>
<ol>
<li>Implement strong password practices immediately</li>
<li>Be more aware of social engineering attempts</li>
<li>Share this knowledge with colleagues and family</li>
<li>Stay updated on new security threats</li>
</ol>

<div class="alert alert-success">
<i class="fa fa-check-circle"></i> <strong>Well Done!</strong> You have completed all video modules. Proceed to the final assessment.
</div>
</div>',
true, 60, 3, datetime('now'), datetime('now'));

-- Final quiz for video course
INSERT INTO course_modules (id, course_id, name, description, module_type, content, must_complete, order_index, created_date, modified_date) VALUES 
(9, 2, 'Final Assessment', 'Test your knowledge from the video course', 'quiz', 
'{"quiz_instructions": "Final test based on the video content. 75% required to pass."}',
true, 4, datetime('now'), datetime('now'));

-- Quiz for video course
INSERT INTO course_quizzes (id, course_id, module_id, name, description, order_index, passing_score, time_limit, max_attempts, shuffle_questions, show_results, created_date) VALUES 
(2, 2, 9, 'Video Course Assessment', 'Test knowledge from video content', 0, 75, 10, 3, 0, 1, datetime('now'));

-- Simple questions for video course quiz
INSERT INTO course_questions (id, quiz_id, question, question_type, points, explanation, order_index) VALUES 
(6, 2, 'What is the primary goal of cybersecurity?', 'multiple_choice', 1, 'Cybersecurity aims to protect digital assets from unauthorized access, use, disclosure, disruption, modification, or destruction.', 0),
(7, 2, 'Which password practice is most effective?', 'multiple_choice', 1, 'Using unique, complex passwords for each account provides the best security against credential-based attacks.', 1);

-- Options for video course quiz
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(19, 6, 'To make computers faster', 0, 0),
(20, 6, 'To protect digital assets from threats', 1, 1),
(21, 6, 'To reduce software costs', 0, 2),
(22, 6, 'To improve user experience', 0, 3),
(23, 7, 'Using the same strong password everywhere', 0, 0),
(24, 7, 'Using unique, complex passwords for each account', 1, 1),
(25, 7, 'Changing passwords daily', 0, 2),
(26, 7, 'Using only numbers in passwords', 0, 3);

-- Course 3: Mixed Content Course
INSERT INTO courses (id, user_id, name, description, created_date, modified_date) VALUES 
(3, 1, 'Complete Security Awareness Program', 'Comprehensive security training with mixed content types', datetime('now'), datetime('now'));

-- Mixed modules for course 3
INSERT INTO course_modules (id, course_id, name, description, module_type, content, must_complete, min_time_spent, order_index, created_date, modified_date) VALUES 
(10, 3, 'Security Policy Overview', 'Understanding company security policies', 'html',
'<div class="module-content">
<h2>📋 Company Security Policies</h2>
<p>Every organization needs clear security policies to protect its assets and data. This module covers:</p>
<ul>
<li>Acceptable Use Policies</li>
<li>Data Classification Standards</li>
<li>Incident Response Procedures</li>
<li>Remote Work Security Guidelines</li>
</ul>
<p><strong>Remember:</strong> Security is everyones responsibility!</p>
</div>',
true, 90, 0, datetime('now'), datetime('now'));

INSERT INTO course_modules (id, course_id, name, description, module_type, video_url, video_duration, must_complete, min_time_spent, order_index, created_date, modified_date) VALUES 
(11, 3, 'Data Protection Video', 'Learn about data protection principles', 'video', 'https://www.youtube.com/embed/hMtzBXWUEkA', 240, true, 240, 1, datetime('now'), datetime('now'));

INSERT INTO course_modules (id, course_id, name, description, module_type, content, must_complete, order_index, created_date, modified_date) VALUES 
(12, 3, 'Quick Knowledge Check', 'Mini quiz to check understanding', 'quiz', 
'{"quiz_instructions": "Quick 3-question check. 70% to pass."}',
true, 2, datetime('now'), datetime('now'));

-- Mini quiz for mixed course
INSERT INTO course_quizzes (id, course_id, module_id, name, description, order_index, passing_score, time_limit, max_attempts, shuffle_questions, show_results, created_date) VALUES 
(3, 3, 12, 'Quick Knowledge Check', 'Short quiz on policies and data protection', 0, 70, 5, 2, 0, 1, datetime('now'));

-- Quick quiz questions
INSERT INTO course_questions (id, quiz_id, question, question_type, points, order_index) VALUES 
(8, 3, 'Who is responsible for security in an organization?', 'multiple_choice', 1, 0),
(9, 3, 'Data protection is important because:', 'multiple_choice', 1, 1),
(10, 3, 'True or False: You can share your login credentials with trusted colleagues.', 'true_false', 1, 2);

-- Options for quick quiz
INSERT INTO question_options (id, question_id, option, is_correct, order_index) VALUES 
(27, 8, 'Only the IT department', 0, 0),
(28, 8, 'Everyone in the organization', 1, 1),
(29, 8, 'Only managers and executives', 0, 2),
(30, 9, 'It prevents legal issues and protects privacy', 1, 0),
(31, 9, 'It makes computers run faster', 0, 1),
(32, 9, 'It reduces software costs', 0, 2),
(33, 10, 'True', 0, 0),
(34, 10, 'False', 1, 1);