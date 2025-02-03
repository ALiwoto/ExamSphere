
SET TIME ZONE 'UTC';

-- update the database default timezone
DO $$
DECLARE
    db_name text;
BEGIN
    SELECT current_database() INTO db_name;
    EXECUTE format('ALTER DATABASE %I SET timezone TO ''UTC''', db_name);
END
$$;

-- I should have done this sooner, may God forgive me for my sins
UPDATE exam_info 
SET exam_date = exam_date AT TIME ZONE 'UTC',
    created_at = created_at AT TIME ZONE 'UTC';

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
            'can_participate_in_exam',
            'give_answer_to_exam_question'
        )
    LOOP
        EXECUTE 'DROP FUNCTION IF EXISTS public.' || quote_ident(func_record.proname) || '(' || func_record.args || ') CASCADE';
    END LOOP;
END $$;



-- Returns true if the user can participate in the exam, false otherwise.
-- Example usage:
--      SELECT can_participate_in_exam(1234, '5678');
CREATE OR REPLACE FUNCTION can_participate_in_exam(p_exam_id INTEGER, p_user_id UserIdType)
RETURNS BOOLEAN AS $$
DECLARE
    is_exam_public BOOLEAN;
BEGIN
    -- Just return true if the user already participated inside of this exam
    IF has_participated_in_exam(p_exam_id, p_user_id) THEN
        RETURN TRUE;
    END IF;

    -- Check if the exam is public (later on we can add more conditions here)
    SELECT "is_public" INTO is_exam_public
    FROM "exam_info"
    WHERE exam_id = p_exam_id;

    IF is_exam_public IS NULL THEN
        RAISE EXCEPTION 'Exam with ID % not found', p_exam_id;
    END IF;

    RETURN is_exam_public;
END;
$$ LANGUAGE plpgsql;


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
END;
$$;


-- give_answer_to_exam_question function is used to insert or update
-- an answer given by a user to an exam question.
-- Example usage:
--      SELECT give_answer_to_exam_question(
--          p_exam_id := 1,
--          p_question_id := 1,
--          p_answered_by := '1234',
--          p_chosen_option := 'A',
--          p_seconds_taken := 30,
--          p_answer_text := NULL
--      );
CREATE OR REPLACE FUNCTION give_answer_to_exam_question(
    p_exam_id INTEGER,
    p_question_id INTEGER,
    p_answered_by UserIdType,
    p_chosen_option TEXT DEFAULT NULL,
    p_seconds_taken INTEGER DEFAULT 0,
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
    INSERT INTO given_answer (
        exam_id,
        question_id,
        answered_by,
        chosen_option,
        seconds_taken,
        answer_text
    )
    VALUES (
        p_exam_id,
        p_question_id,
        p_answered_by,
        p_chosen_option,
        p_seconds_taken,
        p_answer_text
    )
    ON CONFLICT (exam_id, question_id, answered_by)
    DO UPDATE SET -- Just update the answer if it already exists
        chosen_option = EXCLUDED.chosen_option,
        answer_text = EXCLUDED.answer_text,
        answered_at = CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;


