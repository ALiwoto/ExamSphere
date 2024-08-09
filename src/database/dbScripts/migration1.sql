CREATE TABLE IF NOT EXISTS "user_info" (
    user_id INTEGER PRIMARY KEY,
    full_name VARCHAR(127) NOT NULL,
    email VARCHAR(127) UNIQUE NOT NULL CHECK (email ~* '^[A-Za-z0-9._+%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    password VARCHAR(127) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_banned BOOLEAN DEFAULT FALSE,
    ban_reason TEXT,
    role VARCHAR(10) DEFAULT 'user'
);

-- Comments to describe the fields, these can be easily viewed from PgAdmin
COMMENT ON COLUMN "user_info".user_id IS 'Unique identifier for the user';
COMMENT ON COLUMN "user_info".full_name IS 'User''s full name';
COMMENT ON COLUMN "user_info".email IS 'Email address of the user (unique and has validation)';
COMMENT ON COLUMN "user_info".password IS 'Hashed password';
COMMENT ON COLUMN "user_info".created_at IS 'Timestamp when the user account was created';
COMMENT ON COLUMN "user_info".is_banned IS 'Flag indicating if the user is banned';
COMMENT ON COLUMN "user_info".ban_reason IS 'Reason for banning the user, if applicable';
COMMENT ON COLUMN "user_info".role IS 'Role of the user in the system (admin, student, teacher)';

-- Create user function
CREATE OR REPLACE FUNCTION create_user(
    p_user_id INTEGER,
    p_full_name VARCHAR(127),
    p_email VARCHAR(127),
    p_password VARCHAR(127),
    p_role VARCHAR(10) DEFAULT 'student'
) RETURNS INTEGER AS $$
DECLARE
    new_user_id INTEGER;
BEGIN
    INSERT INTO "user_info" (user_id, full_name, email, password, role)
    VALUES (p_user_id, p_full_name, p_email, p_password, p_role)
    RETURNING user_id INTO new_user_id;
    
    RETURN new_user_id;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION update_created_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER set_created_at
BEFORE INSERT ON "user_info"
FOR EACH ROW
EXECUTE FUNCTION update_created_at();


---------------------------------------------------------------

CREATE TABLE IF NOT EXISTS "course_info" (
    course_id SERIAL PRIMARY KEY,
    course_name VARCHAR(127) NOT NULL,
    course_description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN "course_info".course_id IS 'Unique identifier for the course';
COMMENT ON COLUMN "course_info".course_name IS 'Name of the course';
COMMENT ON COLUMN "course_info".course_description IS 'Description of the course';
COMMENT ON COLUMN "course_info".created_at IS 'Timestamp when the course was created';

CREATE OR REPLACE FUNCTION create_course_info(
    p_course_name VARCHAR(127),
    p_course_description TEXT DEFAULT NULL
) RETURNS INTEGER AS $$
DECLARE
    new_course_id INTEGER;
BEGIN
    INSERT INTO "course_info" (course_name, course_description)
    VALUES (p_course_name, p_course_description)
    RETURNING course_id INTO new_course_id;
    
    RETURN new_course_id;
END;
$$ LANGUAGE plpgsql;

---------------------------------------------------------------

CREATE TABLE IF NOT EXISTS "exam_info" (
    exam_id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL,
    price NUMERIC(10, 2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    exam_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    duration INTEGER NOT NULL DEFAULT 60, -- Duration in minutes (e.g. 120)
    created_by INTEGER NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES "user_info"(user_id),
    CONSTRAINT fk_course_id FOREIGN KEY (course_id) REFERENCES "course_info"(course_id)
);

-- Add comments to describe the fields
COMMENT ON TABLE "exam_info" IS 'Stores information about exams';
COMMENT ON COLUMN "exam_info".exam_id IS 'Unique identifier for the exam';
COMMENT ON COLUMN "exam_info".price IS 'Price of the exam';
COMMENT ON COLUMN "exam_info".duration IS 'Duration of the exam in minutes';
COMMENT ON COLUMN "exam_info".created_at IS 'Timestamp when the exam info was created';
COMMENT ON COLUMN "exam_info".exam_date IS 'Date when the exam is scheduled';
COMMENT ON COLUMN "exam_info".created_by IS 'ID of the user who created this exam info';
COMMENT ON COLUMN "exam_info".is_public IS 'Flag indicating if the exam is public';

-- functions for creating a single exam_info
-- examples for calling this function:
-- SELECT create_exam_info(
--     p_course_id := 2,
--     p_price := 149.99,
--     p_created_by := 101,
--     p_is_public := TRUE,
--     p_duration := 120,
--     p_exam_date := '2023-12-31 14:00:00+00'
-- );
CREATE OR REPLACE FUNCTION create_exam_info(
    p_course_id INTEGER,
    p_price NUMERIC(10, 2),
    p_created_by INTEGER,
    p_is_public BOOLEAN DEFAULT FALSE,
    p_duration INTEGER DEFAULT 60,
    p_exam_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
) RETURNS INTEGER AS $$
DECLARE
    new_exam_id INTEGER;
BEGIN
    INSERT INTO "exam_info" (course_id, price, exam_date, created_by, is_public, duration)
    VALUES (p_course_id, p_price, p_exam_date, p_created_by, p_is_public, p_duration)
    RETURNING exam_id INTO new_exam_id;
    
    RETURN new_exam_id;
END;
$$ LANGUAGE plpgsql;

-- is_started returns true if the exam is started and false otherwise
CREATE OR REPLACE FUNCTION has_exam_started(p_exam_id INTEGER)
RETURNS BOOLEAN AS $$
DECLARE
    exam_start_time TIMESTAMP WITH TIME ZONE;
BEGIN
    SELECT exam_date INTO exam_start_time
    FROM exam_info
    WHERE exam_id = p_exam_id;

    IF exam_start_time IS NULL THEN
        RAISE EXCEPTION 'Exam with ID % not found', p_exam_id;
    END IF;

    RETURN CURRENT_TIMESTAMP >= exam_start_time;
END;
$$ LANGUAGE plpgsql;


-- Returns true if the exam is finished, false otherwise
CREATE OR REPLACE FUNCTION has_exam_finished(p_exam_id INTEGER)
RETURNS BOOLEAN AS $$
DECLARE
    exam_end_time TIMESTAMP WITH TIME ZONE;
BEGIN
    SELECT exam_date + (duration || ' minutes')::INTERVAL INTO exam_end_time
    FROM exam_info
    WHERE exam_id = p_exam_id;

    IF exam_end_time IS NULL THEN
        RAISE EXCEPTION 'Exam with ID % not found', p_exam_id;
    END IF;

    RETURN CURRENT_TIMESTAMP > exam_end_time;
END;
$$ LANGUAGE plpgsql;

-- Returns the number of minutes until the exam starts.
-- Returns -1 if the exam has already started.
CREATE OR REPLACE FUNCTION get_exam_starts_in(p_exam_id INTEGER)
RETURNS INTEGER AS $$
DECLARE
    exam_start_time TIMESTAMP WITH TIME ZONE;
BEGIN
    SELECT exam_date INTO exam_start_time
    FROM exam_info
    WHERE exam_id = p_exam_id;

    IF exam_start_time IS NULL THEN
        RAISE EXCEPTION 'Exam with ID % not found', p_exam_id;
    END IF;

    RETURN EXTRACT(EPOCH FROM (exam_start_time - CURRENT_TIMESTAMP)) / 60;
END;
$$ LANGUAGE plpgsql;

-- Returns the number of minutes until the exam finishes
-- If the exam has already finished, it returns 0
-- If the exam is ongoing, it returns the remaining minutes as integer
-- If the exam is not yet started, this function will return all the minutes
-- remaining until the exam starts + duration to finish.
-- Example of using it: SELECT get_exam_finishes_in(3)
CREATE OR REPLACE FUNCTION get_exam_finishes_in(p_exam_id INTEGER)
RETURNS INTEGER AS $$
DECLARE
    exam_end_time TIMESTAMP WITH TIME ZONE;
BEGIN
    SELECT exam_date + (duration || ' minutes')::INTERVAL INTO exam_end_time
    FROM exam_info
    WHERE exam_id = p_exam_id;

    IF exam_end_time IS NULL THEN
        RAISE EXCEPTION 'Exam with ID % not found', p_exam_id;
    END IF;

    RETURN GREATEST(0, EXTRACT(EPOCH FROM (exam_end_time - CURRENT_TIMESTAMP)) / 60)::INTEGER;
END;
$$ LANGUAGE plpgsql;

---------------------------------------------------------------

CREATE TABLE IF NOT EXISTS exam_question (
    question_id SERIAL PRIMARY KEY,
    exam_id INTEGER NOT NULL,
    question_title VARCHAR(2048) NOT NULL,
    description TEXT,
    option1 TEXT,
    option2 TEXT,
    option3 TEXT,
    option4 TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_exam_info FOREIGN KEY (exam_id) REFERENCES exam_info(exam_id)
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

CREATE TABLE IF NOT EXISTS "given_answer" (
    exam_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    answered_by INTEGER NOT NULL,
    chosen_option TEXT DEFAULT NULL,
    answer_text TEXT,
    answered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (exam_id, question_id, answered_by),

    CONSTRAINT fk_exam FOREIGN KEY (exam_id) REFERENCES "exam_info"(exam_id),
    CONSTRAINT fk_question FOREIGN KEY (question_id) REFERENCES "exam_question"(question_id),
    CONSTRAINT fk_user FOREIGN KEY (answered_by) REFERENCES "user_info"(user_id)
);

COMMENT ON TABLE given_answer IS 'Stores answers given by users to exam questions';
COMMENT ON COLUMN given_answer.exam_id IS 'ID of the exam';
COMMENT ON COLUMN given_answer.question_id IS 'ID of the question';
COMMENT ON COLUMN given_answer.answered_by IS 'ID of the user who answered';
COMMENT ON COLUMN given_answer.chosen_option IS 'The title of option chosen by the user';
COMMENT ON COLUMN given_answer.answer_text IS 'Text answer provided by the user, if applicable';
COMMENT ON COLUMN given_answer.answered_at IS 'Timestamp when the answer was submitted';

CREATE OR REPLACE FUNCTION give_answer_to_exam(
    p_exam_id INTEGER,
    p_question_id INTEGER,
    p_answered_by INTEGER,
    p_chosen_option TEXT DEFAULT NULL,
    p_answer_text TEXT DEFAULT NULL
) RETURNS VOID AS $$
BEGIN
    -- Check if the exam has started
    IF NOT has_exam_started(p_exam_id) THEN
        RAISE EXCEPTION 'Exam % has not started yet', p_exam_id;
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

