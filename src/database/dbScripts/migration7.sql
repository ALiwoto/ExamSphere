-- View to get all exams (exam_id and exam_title and when it starts) that a user has ever participated in
-- and are not finished
-- An example of using this view would be:
--      SELECT exam_id, exam_title, exam_date
--          FROM user_ongoing_exams
--          WHERE user_id = '1234';
CREATE OR REPLACE VIEW user_ongoing_exams AS
SELECT DISTINCT
    u.user_id,
    e.exam_id,
    e.exam_title,
    e.exam_date
FROM "exam_info" e
JOIN "given_exam" g ON e.exam_id = g.exam_id
JOIN "user_info" u ON g.user_id = u.user_id
WHERE 
    CURRENT_TIMESTAMP BETWEEN e.exam_date AND (e.exam_date + (e.duration || ' minutes')::INTERVAL)
    AND e.is_sample_exam = FALSE;

COMMENT ON VIEW user_ongoing_exams IS 'View to get all exams (exam_id and exam_title and when it starts) that a user has ever participated in and are not finished';

-- View to get all exams (exam_id and exam_title and when it starts) that a user
-- has participated in the past and now are finished.
-- An example of using this view would be:
--      SELECT exam_id, exam_title, exam_date
--          FROM user_exams_history
--          WHERE user_id = '1234';
CREATE OR REPLACE VIEW user_exams_history AS
SELECT DISTINCT
    u.user_id,
    e.exam_id,
    e.exam_title,
    e.exam_date
FROM "exam_info" e
JOIN "given_exam" g ON e.exam_id = g.exam_id
JOIN "user_info" u ON g.user_id = u.user_id
WHERE
    CURRENT_TIMESTAMP > (e.exam_date + (e.duration || ' minutes')::INTERVAL)
    AND e.is_sample_exam = FALSE;

COMMENT ON VIEW user_ongoing_exams IS 'View to get all exams (exam_id and exam_title and when it starts) that a user has ever participated in and are not finished';


CREATE OR REPLACE VIEW user_future_exams AS
SELECT DISTINCT
    u.user_id,
    e.exam_id,
    e.exam_title,
    e.exam_date
FROM "exam_info" e
JOIN "given_exam" g ON e.exam_id = g.exam_id
JOIN "user_info" u ON g.user_id = u.user_id
WHERE 
    e.exam_date > CURRENT_TIMESTAMP
    AND e.is_sample_exam = FALSE;

COMMENT ON VIEW user_future_exams IS 'View to get all exams (exam_id and exam_title and when it starts) that a user has ever participated in and are not finished';
