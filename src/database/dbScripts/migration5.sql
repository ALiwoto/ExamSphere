

-- since we are going to re-declare some functions in this migration which might have
-- been declared previously, but with different arguments, we need to first remove them
-- and then re-declare them with the new arguments. This is necessary because function
-- overloading is a huge pain in PostgreSQL.
DO $$ 
DECLARE 
    func_record RECORD;
BEGIN
    FOR func_record IN 
        SELECT p.proname, pg_get_function_identity_arguments(p.oid) as args
        FROM pg_proc p
        JOIN pg_namespace n ON p.pronamespace = n.oid
        WHERE p.proname IN (
            'create_exam_info',
            'create_exam_question'
        )
    LOOP
        EXECUTE 'DROP FUNCTION IF EXISTS public.' || quote_ident(func_record.proname) || '(' || func_record.args || ') CASCADE';
    END LOOP;
END $$;



ALTER TABLE exam_info ADD COLUMN is_sample_exam BOOLEAN DEFAULT FALSE;
ALTER TABLE exam_info ADD COLUMN is_strict BOOLEAN DEFAULT FALSE;
ALTER TABLE exam_info ADD COLUMN max_questions_seconds INTEGER DEFAULT 0;
ALTER TABLE exam_info ADD COLUMN needs_video_call BOOLEAN DEFAULT FALSE;
ALTER TABLE exam_info ADD COLUMN needs_voice_call BOOLEAN DEFAULT FALSE;

COMMENT ON COLUMN exam_info.is_sample_exam IS 'If the exam is a sample exam';
COMMENT ON COLUMN exam_info.is_strict IS 'If strict rules apply to this exam';
COMMENT ON COLUMN exam_info.max_questions_seconds IS 'Maximum time allowed for each question';
COMMENT ON COLUMN exam_info.needs_video_call IS 'If the exam requires a video call';
COMMENT ON COLUMN exam_info.needs_voice_call IS 'If the exam requires a voice call';

-- functions for creating a single exam_info
-- examples for calling this function:
-- SELECT create_exam_info(
--     p_course_id := 2,
--     p_exam_title := 'Math Midterm Exam 1403',
--     p_exam_description := 'This is a midterm exam for the Math course.',
--     p_created_by := 101,
--     p_price := '149.99',
--     p_is_public := TRUE,
--     p_duration := 120,
--     p_exam_date := '2023-12-31 14:00:00+00',
--     p_is_strict := TRUE,
--     p_is_sample_exam := FALSE,
--     p_max_questions_seconds := 60,
--     p_needs_video_call := TRUE,
--     p_needs_voice_call := TRUE
-- );
CREATE OR REPLACE FUNCTION create_exam_info(
    p_course_id INTEGER,
    p_exam_title VARCHAR(63),
    p_exam_description VARCHAR(63),
    p_created_by UserIdType,
    p_price VARCHAR(16) DEFAULT '0T',
    p_is_public BOOLEAN DEFAULT FALSE,
    p_duration INTEGER DEFAULT 60,
    p_exam_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    p_is_strict BOOLEAN DEFAULT FALSE,
    p_is_sample_exam BOOLEAN DEFAULT FALSE,
    p_max_questions_seconds INTEGER DEFAULT 0,
    p_needs_video_call BOOLEAN DEFAULT FALSE,
    p_needs_voice_call BOOLEAN DEFAULT FALSE
) RETURNS INTEGER AS $$
DECLARE
    new_exam_id INTEGER;
BEGIN
    INSERT INTO "exam_info" (
        course_id,
        exam_title,
        exam_description,
        price,
        exam_date,
        created_by,
        is_public,
        duration,
        is_strict,
        is_sample_exam,
        max_questions_seconds,
        needs_video_call,
        needs_voice_call
    )
    VALUES (
        p_course_id,
        p_exam_title,
        p_exam_description,
        p_price,
        p_exam_date, 
        p_created_by, 
        p_is_public, 
        p_duration,
        p_is_strict,
        p_is_sample_exam,
        p_max_questions_seconds,
        p_needs_video_call,
        p_needs_voice_call
    )
    RETURNING exam_id INTO new_exam_id;
    
    RETURN new_exam_id;
END;
$$ LANGUAGE plpgsql;


-- question pointer is when a question is pointing to some other questions
-- with a specific pointer_count, the question is not a real question but a pointer
-- to other questions. This is used to create a question that is a combination of
-- multiple questions; this has to be handled on the backend.
ALTER TABLE exam_question ADD COLUMN is_pointer BOOLEAN DEFAULT FALSE;
ALTER TABLE exam_question ADD COLUMN pointer_count INTEGER DEFAULT 0;

ALTER TABLE exam_question 
ADD COLUMN pointer_to_exam_id INTEGER DEFAULT NULL 
REFERENCES exam_info(exam_id) 
ON DELETE CASCADE 
ON UPDATE CASCADE;

-- Function to create a single exam question.
-- Returns the question_id of the newly created question.
-- Example usage:
-- SELECT create_exam_question(
--     p_exam_id := 1234,
--     p_question_title := 'What is the capital of France?',
--     p_description := 'Choose the correct option from the following.',
--     p_option1 := 'Paris',
--     p_option2 := 'London',
--     p_option3 := 'Berlin',
--     p_option4 := 'Madrid',
--     p_is_pointer := FALSE,
--     p_pointer_count := 0,
--     p_pointer_to_exam_id := NULL
-- );
CREATE OR REPLACE FUNCTION create_exam_question(
    p_exam_id INTEGER,
    p_question_title VARCHAR(2048),
    p_description TEXT DEFAULT NULL,
    p_option1 TEXT DEFAULT NULL,
    p_option2 TEXT DEFAULT NULL,
    p_option3 TEXT DEFAULT NULL,
    p_option4 TEXT DEFAULT NULL,
    p_is_pointer BOOLEAN DEFAULT FALSE,
    p_pointer_count INTEGER DEFAULT 0,
    p_pointer_to_exam_id INTEGER DEFAULT NULL
) RETURNS INTEGER AS $$
DECLARE
    new_question_id INTEGER;
BEGIN
    -- Check if p_is_pointer is FALSE and p_pointer_to_exam_id is NOT NULL
    IF p_is_pointer = FALSE AND p_pointer_to_exam_id IS NOT NULL THEN
        RAISE EXCEPTION 'p_pointer_to_exam_id must be NULL when p_is_pointer is FALSE';
    END IF;

    -- Check if p_is_pointer is TRUE and p_pointer_to_exam_id is NULL
    IF p_is_pointer = TRUE AND p_pointer_to_exam_id IS NULL THEN
        RAISE EXCEPTION 'p_pointer_to_exam_id must NOT be NULL when p_is_pointer is TRUE';
    END IF;

    INSERT INTO exam_question (
        exam_id,
        question_title,
        description,
        option1,
        option2,
        option3,
        option4,
        is_pointer,
        pointer_count,
        pointer_to_exam_id
    )
    VALUES (
        p_exam_id,
        p_question_title,
        p_description,
        p_option1,
        p_option2,
        p_option3,
        p_option4,
        p_is_pointer,
        p_pointer_count,
        p_pointer_to_exam_id
    )
    RETURNING question_id INTO new_question_id;
    
    RETURN new_question_id;
END;
$$ LANGUAGE plpgsql;


-- View to get the most recent exams.
-- The results this view returns are ordered by exam_date in ascending order,
-- meaning the exams that are going to happen soon will be shown first.
-- Example usage:
--   SELECT * FROM most_recent_exams_view LIMIT 10 OFFSET 0;
-- It is strongly recommended that you use pagination when querying this view.
CREATE OR REPLACE VIEW most_recent_exams_view AS
SELECT 
    ei.exam_id,
    ei.course_id,
    ei.exam_title,
    ei.exam_description,
    ei.price,
    ei.created_at,
    ei.exam_date,
    ei.duration,
    ei.created_by,
    ei.is_public
FROM 
    exam_info ei
WHERE 
    ei.exam_date >= CURRENT_TIMESTAMP AND
    ei.is_public = TRUE AND
    ei.is_sample_exam = FALSE
ORDER BY 
    ei.exam_date DESC;



ALTER TABLE given_answer ADD COLUMN seen_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;

COMMENT ON COLUMN given_answer.seen_at IS 'Timestamp when the answer was seen by the student';

-- Function to mark an answer as seen by the student.
-- Example usage:
--      SELECT mark_given_answer_as_seen(
--          p_exam_id := 1,
--          p_question_id := 1,
--          p_answered_by := '1234'
--      );
CREATE OR REPLACE FUNCTION mark_given_answer_as_seen(
    p_exam_id INTEGER,
    p_question_id INTEGER,
    p_answered_by UserIdType
) RETURNS VOID AS $$
BEGIN
    UPDATE given_answer
    SET seen_at = CURRENT_TIMESTAMP
    WHERE exam_id = p_exam_id AND question_id = p_question_id AND answered_by = p_answered_by;
END;
$$ LANGUAGE plpgsql;

-- Function to mark some answers as seen by the student.
-- NOTE: the main differences
-- SELECT mark_given_answers_as_seen(
--     p_exam_id := 1,
--     p_question_ids := ARRAY[1, 2, 3, 4, 5],
--     p_answered_by := 'user123'
-- );
CREATE OR REPLACE FUNCTION mark_given_answers_as_seen(
    p_exam_id INTEGER,
    p_question_ids INTEGER[],
    p_answered_by UserIdType
) RETURNS VOID AS $$
BEGIN
    UPDATE given_answer
    SET seen_at = CURRENT_TIMESTAMP
    WHERE exam_id = p_exam_id 
    AND question_id = ANY(p_question_ids) 
    AND answered_by = p_answered_by
    AND seen_at IS NULL;  -- Only update if seen_at is null
END;
$$ LANGUAGE plpgsql;

-- Function to choose questions for a participant in an exam.
-- This function is used to choose questions for a participant in an exam.
-- It loops through all questions for the given exam and inserts them into the given_answer table.
-- If a question is a pointer, it randomly selects pointer_count questions from the pointer_to_exam_id exam.
-- Example usage:
--      SELECT choose_questions_for_participant(
--          p_exam_id := 1,
--          p_participant_id := '1234'
--      );
CREATE OR REPLACE FUNCTION choose_questions_for_participant(
    p_exam_id INTEGER,
    p_participant_id UserIdType
) RETURNS VOID AS $$
DECLARE
    r RECORD;
    pointer_question RECORD;
    exam_info_record RECORD;
BEGIN
    -- Check if the exam exists and if it is a sample exam
    SELECT * INTO exam_info_record FROM exam_info WHERE exam_id = p_exam_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'No exam found with id %', p_exam_id;
    END IF;
    
    IF exam_info_record.is_sample_exam = TRUE THEN
        RAISE EXCEPTION 'Cannot participate in a sample exam with id %', p_exam_id;
    END IF;

    -- Loop through all questions for the given exam
    FOR r IN SELECT * FROM exam_question WHERE exam_id = p_exam_id
    LOOP
        IF r.is_pointer = false THEN
            -- Insert the question directly into given_answer with empty answers
            INSERT INTO given_answer (exam_id, question_id, answered_by, chosen_option, seconds_taken, answer_text, answered_at)
            VALUES (p_exam_id, r.question_id, p_participant_id, NULL, 0, NULL, CURRENT_TIMESTAMP);
        ELSE
            -- Randomly select pointer_count questions from the pointer_to_exam_id exam
            FOR pointer_question IN
                SELECT question_id FROM exam_question
                WHERE exam_id = r.pointer_to_exam_id
                ORDER BY RANDOM()
                LIMIT r.pointer_count
            LOOP
                -- Insert the randomly selected questions into given_answer with empty answers
                INSERT INTO given_answer (exam_id, question_id, answered_by, chosen_option, seconds_taken, answer_text, answered_at)
                VALUES (p_exam_id, pointer_question.question_id, p_participant_id, NULL, 0, NULL, CURRENT_TIMESTAMP);
            END LOOP;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;


-- Function to get exam questions for a participant in an exam.
-- The target exam cannot be a sample exam.
-- If no questions are found for the participant, it calls choose_questions_for_participant.
-- Example usage:
--      SELECT * FROM get_exam_questions_for_participant(
--          p_exam_id := 2,
--          p_participant_id := '1234'
--      );
CREATE OR REPLACE FUNCTION get_exam_questions_for_participant(
    p_exam_id INTEGER,
    p_participant_id UserIdType
) RETURNS SETOF exam_question AS $$
DECLARE
    r RECORD;
    exam_info_record RECORD;
BEGIN
    -- Check if the exam exists and if it is a sample exam
    SELECT * INTO exam_info_record FROM exam_info WHERE exam_id = p_exam_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'No exam found with id %', p_exam_id;
    END IF;
    
    IF exam_info_record.is_sample_exam = TRUE THEN
        RAISE EXCEPTION 'Cannot participate in a sample exam with id %', p_exam_id;
    END IF;

    -- Try to select from given_answer table
    FOR r IN
        SELECT eq.*
        FROM exam_question eq
        JOIN given_answer ga ON eq.question_id = ga.question_id
        WHERE ga.exam_id = p_exam_id AND ga.answered_by = p_participant_id
    LOOP
        RETURN NEXT r;
    END LOOP;

    -- If no questions found, call choose_questions_for_participant
    IF NOT FOUND THEN
        PERFORM choose_questions_for_participant(p_exam_id, p_participant_id);

        -- Try to select from given_answer table again
        FOR r IN
            SELECT eq.*
            FROM exam_question eq
            JOIN given_answer ga ON eq.question_id = ga.question_id
            WHERE ga.exam_id = p_exam_id AND ga.answered_by = p_participant_id
        LOOP
            RETURN NEXT r;
        END LOOP;

        -- If still no questions found, raise an exception
        IF NOT FOUND THEN
            RAISE EXCEPTION 'No questions found for exam_id % and participant_id %', p_exam_id, p_participant_id;
        END IF;
    END IF;
END;
$$ LANGUAGE plpgsql;

