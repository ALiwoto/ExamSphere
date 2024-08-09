
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

-- Function to create a single exam question
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
CREATE TABLE IF NOT EXISTS "given_exam" (
    user_id VARCHAR(16) NOT NULL,
    exam_id INTEGER NOT NULL,
    price NUMERIC(10, 2),
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
CREATE OR REPLACE FUNCTION has_participated_in_exam(p_exam_id INTEGER, p_user_id VARCHAR(16))
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1
        FROM given_exam
        WHERE exam_id = p_exam_id AND user_id = p_user_id
    );
END;
$$ LANGUAGE plpgsql;

-- Returns true if the user can participate in the exam, false otherwise
CREATE OR REPLACE FUNCTION can_participate_in_exam(p_exam_id INTEGER, p_user_id VARCHAR(16))
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
CREATE OR REPLACE FUNCTION set_score_for_user_in_exam(
    p_exam_id INTEGER,
    p_user_id VARCHAR(16),
    p_final_score VARCHAR(63),
    p_scored_by INTEGER
) RETURNS VOID AS $$
BEGIN
    -- Check if the exam entry exists
    IF NOT EXISTS (
        SELECT 1 FROM given_exam
        WHERE exam_id = p_exam_id AND user_id = p_user_id
    ) THEN
        RAISE EXCEPTION 'No exam entry found for user % in exam %', p_user_id, p_exam_id;
    END IF;

    -- Update the final_score and added_by
    UPDATE given_exam
    SET final_score = p_final_score,
        scored_by = p_added_by
    WHERE exam_id = p_exam_id AND user_id = p_user_id;

    -- Check if any rows were affected
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Failed to update score for user % in exam %', p_user_id, p_exam_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

---------------------------------------------------------------
