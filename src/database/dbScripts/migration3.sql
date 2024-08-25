
ALTER TABLE "user_info" ADD COLUMN IF NOT EXISTS user_address VARCHAR(255);
ALTER TABLE "user_info" ADD COLUMN IF NOT EXISTS phone_number VARCHAR(16);
ALTER TABLE "user_info" ADD COLUMN IF NOT EXISTS setup_completed BOOLEAN DEFAULT FALSE;

COMMENT ON COLUMN "user_info".user_address IS 'Address of the user';
COMMENT ON COLUMN "user_info".phone_number IS 'Phone number of the user';

DO
$$
DECLARE
    r RECORD;
BEGIN
    FOR r IN SELECT proname, pg_get_function_identity_arguments(p.oid) AS args
             FROM pg_proc p
             JOIN pg_namespace n ON p.pronamespace = n.oid
             WHERE n.nspname NOT IN ('pg_catalog', 'information_schema')
             AND pg_function_is_visible(p.oid)
             AND proname = 'create_user_info'
    LOOP
        EXECUTE format('DROP FUNCTION IF EXISTS %I(%s);', r.proname, r.args);
    END LOOP;
END
$$;

-- Example usage:
--      SELECT create_user_info(
--         p_user_id := '1234',
--         p_full_name := 'John Doe',
--         p_email := 'email@gmail.com',
--         p_auth_hash := '1234567890',
--         p_password := 'password',
--         p_role := 'student',
--         p_user_address := '123 Main St, City, Country',
--         p_phone_number := '123-456-7890',
--         p_setup_completed := TRUE
--      );
CREATE OR REPLACE FUNCTION create_user_info(
    p_user_id UserIdType,
    p_full_name VARCHAR(127),
    p_email VARCHAR(127),
    p_auth_hash VARCHAR(15),
    p_password VARCHAR(511),
    p_role VARCHAR(10) DEFAULT 'student',
    p_user_address VARCHAR(255) DEFAULT NULL,
    p_phone_number VARCHAR(16) DEFAULT NULL,
    p_setup_completed BOOLEAN DEFAULT FALSE
) RETURNS UserIdType AS $$
DECLARE
    new_user_id UserIdType := '0'; -- UserIdType does not allow NULL values
BEGIN
    INSERT INTO "user_info" (
        user_id, 
        full_name, 
        email, 
        auth_hash, 
        password, 
        role, 
        user_address, 
        phone_number, 
        setup_completed
    )
    VALUES (
        p_user_id, 
        p_full_name, 
        p_email, 
        p_auth_hash, 
        p_password, 
        p_role, 
        p_user_address, 
        p_phone_number, 
        p_setup_completed
    )
    RETURNING user_id INTO new_user_id;
    
    RETURN new_user_id;
END;
$$ LANGUAGE plpgsql;

-----------------------------------------------------------

CREATE TABLE IF NOT EXISTS "given_answer" (
    exam_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    answered_by UserIdType,
    chosen_option TEXT DEFAULT NULL,
    seconds_taken INTEGER DEFAULT 0,
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
COMMENT ON COLUMN given_answer.seconds_taken IS 'Time taken (seconds) by the user to answer the question';
COMMENT ON COLUMN given_answer.answer_text IS 'Text answer provided by the user, if applicable';
COMMENT ON COLUMN given_answer.answered_at IS 'Timestamp when the answer was submitted';

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


CREATE OR REPLACE FUNCTION update_user_topic_stat()
RETURNS TRIGGER AS $$
DECLARE
    topic_id_var INTEGER;
    has_hit_limit BOOLEAN;
BEGIN
    SELECT topic_id INTO topic_id_var
    FROM exam_info ei
    JOIN topic_info ti ON ei.course_id = ti.topic_id
    WHERE ei.exam_id = NEW.exam_id;

    IF NOT FOUND THEN
        RETURN NEW;
    END IF;

    INSERT INTO user_topic_stat (user_id, topic_id)
    VALUES (NEW.answered_by, topic_id_var)
    ON CONFLICT (user_id, topic_id) DO NOTHING;

    -- Check if user has hit exp limit
    SELECT (recent_exp >= 10 AND last_visited > NOW() - INTERVAL '30 minutes')
    INTO has_hit_limit
    FROM user_topic_stat
    WHERE user_id = NEW.answered_by AND topic_id = topic_id_var;

    IF has_hit_limit THEN
        RETURN NEW;
    END IF;

    UPDATE user_topic_stat
    SET 
        current_exp = CASE 
            WHEN current_exp + 1 >= current_level * 10 THEN 0
            ELSE current_exp + 1
        END,
        total_exp = total_exp + 1,
        recent_exp = CASE 
            WHEN last_visited < NOW() - INTERVAL '30 minutes' THEN 1
            ELSE recent_exp + 1
        END,
        last_visited = NOW(),
        current_level = CASE 
            WHEN current_exp + 1 >= current_level * 10 THEN current_level + 1
            ELSE current_level
        END
    WHERE user_id = NEW.answered_by AND topic_id = topic_id_var;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_topic_stat_trigger
AFTER INSERT ON given_answer
FOR EACH ROW
EXECUTE FUNCTION update_user_topic_stat();

---------------------------------------------------------------
