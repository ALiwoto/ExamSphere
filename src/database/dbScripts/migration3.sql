
CREATE TABLE IF NOT EXISTS "given_answer" (
    exam_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    answered_by VARCHAR(16) NOT NULL,
    chosen_option TEXT DEFAULT NULL,
    answer_text TEXT,
    answered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (exam_id, question_id, answered_by),

    CONSTRAINT fk_exam FOREIGN KEY (exam_id) REFERENCES "exam_info"(exam_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_question FOREIGN KEY (question_id) REFERENCES "exam_question"(question_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (answered_by) REFERENCES "user_info"(user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

COMMENT ON TABLE given_answer IS 'Stores answers given by users to exam questions';
COMMENT ON COLUMN given_answer.exam_id IS 'ID of the exam';
COMMENT ON COLUMN given_answer.question_id IS 'ID of the question';
COMMENT ON COLUMN given_answer.answered_by IS 'ID of the user who answered';
COMMENT ON COLUMN given_answer.chosen_option IS 'The title of option chosen by the user';
COMMENT ON COLUMN given_answer.answer_text IS 'Text answer provided by the user, if applicable';
COMMENT ON COLUMN given_answer.answered_at IS 'Timestamp when the answer was submitted';

CREATE OR REPLACE FUNCTION give_answer_to_exam_question(
    p_exam_id INTEGER,
    p_question_id INTEGER,
    p_answered_by INTEGER,
    p_chosen_option TEXT DEFAULT NULL,
    p_answer_text TEXT DEFAULT NULL
) RETURNS VOID AS $$
BEGIN
    -- Check if the user has participated in the exam
    IF NOT has_participated_in_exam(p_exam_id, p_answered_by) THEN
        RAISE EXCEPTION 'User has not participated in exam % yet', p_exam_id;
    END IF;

    -- Check if the exam has finished
    IF has_exam_finished(p_exam_id) THEN
        RAISE EXCEPTION 'Exam % has already finished', p_exam_id;
    END IF;

    -- If the exam is ongoing, proceed with inserting or updating the answer
    INSERT INTO given_answer (exam_id, question_id, answered_by, chosen_option, answer_text)
    VALUES (p_exam_id, p_question_id, p_answered_by, p_chosen_option, p_answer_text)
    ON CONFLICT (exam_id, question_id, answered_by)
    DO UPDATE SET -- Just update the answer if it already exists
        chosen_option = EXCLUDED.chosen_option,
        answer_text = EXCLUDED.answer_text,
        answered_at = CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;


---------------------------------------------------------------
