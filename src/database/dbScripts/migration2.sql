
CREATE TABLE IF NOT EXISTS "exam_question" (
    question_id SERIAL PRIMARY KEY,
    exam_id INTEGER NOT NULL,
    question_title VARCHAR(2048) NOT NULL,
    description TEXT,
    option1 TEXT,
    option2 TEXT,
    option3 TEXT,
    option4 TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_exam_info FOREIGN KEY (exam_id) REFERENCES exam_info(exam_id) ON DELETE CASCADE ON UPDATE CASCADE
);

COMMENT ON TABLE exam_question IS 'Stores information about exam questions';
COMMENT ON COLUMN exam_question.question_id IS 'Unique identifier for the question';
COMMENT ON COLUMN exam_question.exam_id IS 'ID of the exam this question belongs to';
COMMENT ON COLUMN exam_question.question_title IS 'Title or short description of the question';
COMMENT ON COLUMN exam_question.description IS 'Detailed description or full text of the question';
COMMENT ON COLUMN exam_question.option1 IS 'First answer option';
COMMENT ON COLUMN exam_question.option2 IS 'Second answer option';
COMMENT ON COLUMN exam_question.option3 IS 'Third answer option (optional)';
COMMENT ON COLUMN exam_question.option4 IS 'Fourth answer option (optional)';
COMMENT ON COLUMN exam_question.created_at IS 'Timestamp when the question was created';

-- Function to create a single exam question.
-- Returns the question_id of the newly created question.
-- Example usage:
--      SELECT create_exam_question(
--         p_exam_id := 1234,
--         p_question_title := 'What is the capital of France?',
--         p_description := 'Choose the correct option from the following.',
--         p_option1 := 'Paris',
--         p_option2 := 'London',
--         p_option3 := 'Berlin',
--         p_option4 := 'Madrid'
--      );
CREATE OR REPLACE FUNCTION create_exam_question(
    p_exam_id INTEGER,
    p_question_title VARCHAR(2048),
    p_description TEXT DEFAULT NULL,
    p_option1 TEXT DEFAULT NULL,
    p_option2 TEXT DEFAULT NULL,
    p_option3 TEXT DEFAULT NULL,
    p_option4 TEXT DEFAULT NULL
) RETURNS INTEGER AS $$
DECLARE
    new_question_id INTEGER;
BEGIN
    INSERT INTO exam_question (exam_id, question_title, description, option1, option2, option3, option4)
    VALUES (p_exam_id, p_question_title, p_description, p_option1, p_option2, p_option3, p_option4)
    RETURNING question_id INTO new_question_id;
    
    RETURN new_question_id;
END;
$$ LANGUAGE plpgsql;

---------------------------------------------------------------

-- Given exam holds information about exams taken by users.
CREATE TABLE IF NOT EXISTS "given_exam" (
    user_id UserIdType,
    exam_id INTEGER NOT NULL,
    price VARCHAR(16) DEFAULT '0T', -- the price paid by the user to take the exam
    added_by VARCHAR(16) DEFAULT NULL,
    scored_by VARCHAR(16) DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    final_score VARCHAR(63) DEFAULT NULL,
    PRIMARY KEY (user_id, exam_id),

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user_info"(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_exam FOREIGN KEY (exam_id) REFERENCES "exam_info"(exam_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_added_by FOREIGN KEY (added_by) REFERENCES "user_info"(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_scored_by FOREIGN KEY (scored_by) REFERENCES "user_info"(user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

COMMENT ON TABLE given_exam IS 'Stores information about exams taken by users';
COMMENT ON COLUMN given_exam.user_id IS 'ID of the user taking the exam';
COMMENT ON COLUMN given_exam.exam_id IS 'ID of the exam being taken';
COMMENT ON COLUMN given_exam.price IS 'Price paid for the exam';
COMMENT ON COLUMN given_exam.added_by IS 'ID of the user who added this exam entry (can be null)';
COMMENT ON COLUMN given_exam.scored_by IS 'ID of the user (teacher) who scored this exam entry (can be null)';
COMMENT ON COLUMN given_exam.created_at IS 'Timestamp when the exam entry was created';
COMMENT ON COLUMN given_exam.final_score IS 'Final score of the user in the exam; has to be decided by teacher';

-- Returns true if the user has participated in the exam, false otherwise
-- Please note that if the user has been forcefully added by someone else to the exam,
-- this function will still return true.
-- Example usage:
--      SELECT has_participated_in_exam(1234, '5678');
CREATE OR REPLACE FUNCTION has_participated_in_exam(p_exam_id INTEGER, p_user_id UserIdType)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1
        FROM given_exam
        WHERE exam_id = p_exam_id AND user_id = p_user_id
    );
END;
$$ LANGUAGE plpgsql;

-- Returns true if the user can participate in the exam, false otherwise.
-- Example usage:
--      SELECT can_participate_in_exam(1234, '5678');
CREATE OR REPLACE FUNCTION can_participate_in_exam(p_exam_id INTEGER, p_user_id UserIdType)
RETURNS BOOLEAN AS $$
DECLARE
    is_public BOOLEAN;
BEGIN
    -- Just return true if the user already participated inside of this exam
    IF has_participated_in(p_exam_id, p_user_id) THEN
        RETURN TRUE;
    END IF;

    -- Check if the exam is public (later on we can add more conditions here)
    SELECT is_public INTO is_public
    FROM exam_info
    WHERE exam_id = p_exam_id;

    IF is_public IS NULL THEN
        RAISE EXCEPTION 'Exam with ID % not found', p_exam_id;
    END IF;

    RETURN is_public;
END;
$$ LANGUAGE plpgsql;

-- Sets final_score and scored_by for a user in a given_exam.
-- Example usage:
--    CALL set_score_for_user_in_exam(
--        p_exam_id := 1001,           -- p_exam_id: The ID of the exam
--        p_user_id := 'user123',      -- p_user_id: The ID of the user (assuming UserIdType is a string)
--        p_final_score := '85/100',   -- p_final_score: The final score as a string
--        p_scored_by := 5             -- p_scored_by: The ID of the user who scored the exam
--    );
CREATE OR REPLACE PROCEDURE set_score_for_user_in_exam(
    p_exam_id INTEGER,
    p_user_id UserIdType,
    p_final_score VARCHAR(63),
    p_scored_by INTEGER
)
LANGUAGE plpgsql
AS $$
BEGIN
    BEGIN -- start a transaction
        IF NOT EXISTS (
            SELECT 1 FROM given_exam
            WHERE exam_id = p_exam_id AND user_id = p_user_id
        ) THEN
            RAISE EXCEPTION 'No exam entry found for user % in exam %', p_user_id, p_exam_id;
        END IF;

        UPDATE given_exam
        SET final_score = p_final_score,
            scored_by = p_scored_by
        WHERE exam_id = p_exam_id AND user_id = p_user_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Failed to update score for user % in exam %', p_user_id, p_exam_id;
        END IF;

        COMMIT; -- commit the transaction
    EXCEPTION
        WHEN OTHERS THEN
            -- just rollback the transaction if any error occurs
            ROLLBACK;
            RAISE;
    END;
END;
$$;

-- Procedure to add a user to an exam.
-- Example usage:
--    CALL add_user_in_exam(
--        p_user_id := 'user123',
--        p_exam_id := 1001,
--        p_price := '0T',
--        p_added_by := 'admin'
--    );
CREATE OR REPLACE PROCEDURE add_user_in_exam(
    p_user_id UserIdType,
    p_exam_id INTEGER,
    p_price VARCHAR(16) DEFAULT '0T',
    p_added_by VARCHAR(16) DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
BEGIN
    BEGIN
        -- Check if the user already exists in the exam
        IF EXISTS (
            SELECT 1 FROM given_exam
            WHERE user_id = p_user_id AND exam_id = p_exam_id
        ) THEN
            RAISE EXCEPTION 'User % is already registered for exam %', p_user_id, p_exam_id;
        END IF;

        -- Insert the new entry
        INSERT INTO "given_exam" (user_id, exam_id, price, added_by)
        VALUES (p_user_id, p_exam_id, p_price, p_added_by);

        -- If we get here, the insertion was successful, so commit the transaction
        COMMIT;
    EXCEPTION
        WHEN OTHERS THEN
            ROLLBACK; -- rollback all the changes
            -- Re-raise the exception
            RAISE;
    END;
END;
$$;


-- View to get all courses a user has ever enrolled in in their
-- lifetime.
-- An example of using this view would be:
--      SELECT course_id, course_name
--          FROM user_courses
--          WHERE user_id = '1234';
CREATE OR REPLACE VIEW user_courses AS
SELECT DISTINCT u.user_id, c.course_id, c.course_name
FROM "given_exam" g
JOIN "exam_info" e ON g.exam_id = e.exam_id
JOIN "course_info" c ON e.course_id = c.course_id
JOIN "user_info" u ON g.user_id = u.user_id;

COMMENT ON VIEW user_courses IS 'View to get all courses a user has ever enrolled in in their lifetime';


-- View to get user_id and full_name of all users who have ever
-- participated in any exam related to a certain course.
-- An example of using this view would be:
--      SELECT user_id, full_name
--          FROM course_participants
--          WHERE course_id = '1234';
CREATE OR REPLACE VIEW course_participants AS
SELECT DISTINCT
    u.user_id,
    u.full_name,
    c.course_id,
    c.course_name
FROM "course_info" c
JOIN "exam_info" e ON c.course_id = e.course_id
JOIN "given_exam" g ON e.exam_id = g.exam_id
JOIN "user_info" u ON g.user_id = u.user_id
ORDER BY c.course_id, u.user_id;

COMMENT ON VIEW course_participants IS 'View to get user_id and full_name of all users who have ever participated in any exam related to a certain course.';


-- View to get all exams (exam_id and exam_title and when it starts) that a user has ever participated in
-- and are not finished
-- An example of using this view would be:
--      SELECT exam_id, exam_title, start_time
--          FROM user_ongoing_exams
--          WHERE user_id = '1234';
CREATE OR REPLACE VIEW user_ongoing_exams AS
SELECT DISTINCT
    u.user_id,
    e.exam_id,
    e.exam_title,
    e.start_time
FROM "exam_info" e
JOIN "given_exam" g ON e.exam_id = g.exam_id
JOIN "user_info" u ON g.user_id = u.user_id
WHERE CURRENT_TIMESTAMP < (e.exam_date + (e.duration || ' minutes')::INTERVAL);

-- View to get all exams (exam_id and exam_title and when it starts) that a user
-- has participated in the past and now are finished.
-- An example of using this view would be:
--      SELECT exam_id, exam_title, start_time
--          FROM user_exams_history
--          WHERE user_id = '1234';
CREATE OR REPLACE VIEW user_exams_history AS
SELECT DISTINCT
    u.user_id,
    e.exam_id,
    e.exam_title,
    e.start_time
FROM "exam_info" e
JOIN "given_exam" g ON e.exam_id = g.exam_id
JOIN "user_info" u ON g.user_id = u.user_id
WHERE CURRENT_TIMESTAMP > (e.exam_date + (e.duration || ' minutes')::INTERVAL);

COMMENT ON VIEW user_ongoing_exams IS 'View to get all exams (exam_id and exam_title and when it starts) that a user has ever participated in and are not finished';

-- View to get all participants (user_id, full_name, etc) of an exam.
-- An example of using this view would be:
--      SELECT user_id, full_name
--          FROM exam_participants
--          WHERE exam_id = 1234;
CREATE OR REPLACE VIEW exam_participants AS
SELECT DISTINCT
    u.user_id,
    u.full_name,
    e.exam_id,
    e.exam_title
FROM "exam_info" e
JOIN "given_exam" g ON e.exam_id = g.exam_id
JOIN "user_info" u ON g.user_id = u.user_id
WHERE e.exam_id = g.exam_id;

---------------------------------------------------------------
