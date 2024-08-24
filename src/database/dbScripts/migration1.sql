-- Defining the UserIdType here, so we don't have to keep repeating 
-- the same thing over and over again
CREATE DOMAIN UserIdType AS VARCHAR(16) NOT NULL;

CREATE TABLE IF NOT EXISTS "user_info" (
    user_id UserIdType PRIMARY KEY,
    full_name VARCHAR(127) NOT NULL,
    email VARCHAR(127) UNIQUE NOT NULL CHECK (email ~* '^[A-Za-z0-9._+%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    auth_hash VARCHAR(15) NOT NULL,
    password VARCHAR(511) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_banned BOOLEAN DEFAULT FALSE,
    ban_reason TEXT,
    role VARCHAR(10) DEFAULT 'student' CHECK (role IN ('admin', 'student', 'teacher'))
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
CREATE OR REPLACE FUNCTION create_user_info(
    p_user_id UserIdType,
    p_full_name VARCHAR(127),
    p_email VARCHAR(127),
    p_auth_hash VARCHAR(15),
    p_password VARCHAR(511),
    p_role VARCHAR(10) DEFAULT 'student'
) RETURNS UserIdType AS $$
DECLARE
    new_user_id UserIdType := '0'; -- UserIdType does not allow NULL values
BEGIN
    INSERT INTO "user_info" (user_id, full_name, email, auth_hash, password, role)
    VALUES (p_user_id, p_full_name, p_email, p_auth_hash, p_password, p_role)
    RETURNING user_id INTO new_user_id;
    
    RETURN new_user_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_user_password(
    p_user_id UserIdType,
    p_auth_hash VARCHAR(15),
    p_new_password VARCHAR(511)
) RETURNS VOID AS $$
BEGIN
    IF p_new_password = p_auth_hash THEN
        RAISE EXCEPTION 'Password and auth_hash cannot be the same';
    END IF;

    UPDATE "user_info"
    SET password = p_new_password,
        auth_hash = p_auth_hash
    WHERE user_id = p_user_id;
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

-- Checks if the password and auth hash are the same
CREATE OR REPLACE FUNCTION check_password_auth_hash()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.password = NEW.auth_hash THEN
        RAISE EXCEPTION 'Password and auth_hash cannot be the same';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Prevents from the password and auth hash being the same
CREATE  OR REPLACE TRIGGER prevent_same_password_auth_hash
BEFORE INSERT ON user_info
FOR EACH ROW
EXECUTE FUNCTION check_password_auth_hash();


---------------------------------------------------------------


CREATE TABLE IF NOT EXISTS "topic_info" (
    topic_id SERIAL PRIMARY KEY,
    topic_name VARCHAR(127) NOT NULL
);

COMMENT ON COLUMN "topic_info".topic_id IS 'Unique identifier for the topic';
COMMENT ON COLUMN "topic_info".topic_name IS 'Name of the topic';

-- Create a function to create a new topic
-- Example of calling this function:
-- SELECT create_topic_info('Mathematics');
CREATE OR REPLACE FUNCTION create_topic_info(
    p_topic_name VARCHAR(127)
) RETURNS INTEGER AS $$
DECLARE
    new_topic_id INTEGER;
BEGIN
    INSERT INTO "topic_info" (topic_name)
    VALUES (p_topic_name)
    RETURNING topic_id INTO new_topic_id;
    
    RETURN new_topic_id;
END;
$$ LANGUAGE plpgsql;

---------------------------------------------------------------

-- User topic stat table to store user's progress in each topic
-- This table will be updated automatically using triggers defined later
-- in SQL code, each time the backend should only query this table.
-- This table should not be cached on backend.
CREATE TABLE IF NOT EXISTS "user_topic_stat" (
    user_id UserIdType,
    topic_id INTEGER NOT NULL,
    current_exp INTEGER DEFAULT 0,
    total_exp INTEGER DEFAULT 0,
    current_level INTEGER DEFAULT 1,
    last_visited TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, topic_id),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES "user_info"(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_topic_id FOREIGN KEY (topic_id) REFERENCES "topic_info"(topic_id) ON DELETE CASCADE ON UPDATE CASCADE
);

COMMENT ON TABLE "user_topic_stat" IS 'Stores user progress in each topic';
COMMENT ON COLUMN "user_topic_stat".user_id IS 'ID of the user';
COMMENT ON COLUMN "user_topic_stat".topic_id IS 'ID of the topic';
COMMENT ON COLUMN "user_topic_stat".current_exp IS 'Current experience points in the topic';
COMMENT ON COLUMN "user_topic_stat".total_exp IS 'Total experience points earned in the topic';
COMMENT ON COLUMN "user_topic_stat".current_level IS 'Current level in the topic';
COMMENT ON COLUMN "user_topic_stat".last_visited IS 'Timestamp when the user last visited the topic (earned exp)';

---------------------------------------------------------------

CREATE TABLE IF NOT EXISTS "course_info" (
    course_id SERIAL PRIMARY KEY,
    topic_id INTEGER NOT NULL,
    course_name VARCHAR(127) NOT NULL,
    course_description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    added_by UserIdType,

    CONSTRAINT fk_added_by FOREIGN KEY (added_by) REFERENCES "user_info"(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_topic_id FOREIGN KEY (topic_id) REFERENCES "topic_info"(topic_id) ON DELETE CASCADE ON UPDATE CASCADE
);

COMMENT ON COLUMN "course_info".course_id IS 'Unique identifier for the course';
COMMENT ON COLUMN "course_info".course_name IS 'Name of the course';
COMMENT ON COLUMN "course_info".course_description IS 'Description of the course';
COMMENT ON COLUMN "course_info".created_at IS 'Timestamp when the course was created';
COMMENT ON COLUMN "course_info".added_by IS 'ID of the user who added the course';

-- Creates a new course info in the database.
-- Example of calling this function:
-- SELECT create_course_info(
--    p_course_name := 'Mathematics 101',
--    p_topic_id := 1,
--    p_course_description := 'This is an introductory course to mathematics.',
--    p_added_by := '101'
-- );
CREATE OR REPLACE FUNCTION create_course_info(
    p_course_name VARCHAR(127),
    p_topic_id INTEGER,
    p_course_description TEXT,
    p_added_by UserIdType
) RETURNS INTEGER AS $$
DECLARE
    new_course_id INTEGER;
BEGIN
    INSERT INTO "course_info" (course_name, topic_id, course_description, added_by)
    VALUES (p_course_name, p_topic_id, p_course_description, p_added_by)
    RETURNING course_id INTO new_course_id;
    
    RETURN new_course_id;
END;
$$ LANGUAGE plpgsql;

---------------------------------------------------------------

CREATE TABLE IF NOT EXISTS "exam_info" (
    exam_id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL,
    exam_title VARCHAR(63) NOT NULL,
    exam_description VARCHAR(63),
    price VARCHAR(16) DEFAULT '0T',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    exam_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    duration INTEGER NOT NULL DEFAULT 60, -- Duration in minutes (e.g. 120)
    created_by UserIdType,
    is_public BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES "user_info"(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_course_id FOREIGN KEY (course_id) REFERENCES "course_info"(course_id) ON DELETE CASCADE ON UPDATE CASCADE
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
--     p_exam_title := 'Math Midterm Exam 1403',
--     p_exam_description := 'This is a midterm exam for the Math course.',
--     p_price := 149.99,
--     p_created_by := 101,
--     p_is_public := TRUE,
--     p_duration := 120,
--     p_exam_date := '2023-12-31 14:00:00+00'
-- );
CREATE OR REPLACE FUNCTION create_exam_info(
    p_course_id INTEGER,
    p_exam_title VARCHAR(63),
    p_exam_description VARCHAR(63),
    p_created_by UserIdType,
    p_price VARCHAR(16) DEFAULT '0T',
    p_is_public BOOLEAN DEFAULT FALSE,
    p_duration INTEGER DEFAULT 60,
    p_exam_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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
        duration
    )
    VALUES (
        p_course_id,
        p_exam_title,
        p_exam_description,
        p_price,
        p_exam_date, 
        p_created_by, 
        p_is_public, 
        p_duration
    )
    RETURNING exam_id INTO new_exam_id;
    
    RETURN new_exam_id;
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
    ei.exam_date >= CURRENT_TIMESTAMP AND ei.is_public = TRUE
ORDER BY 
    ei.exam_date DESC;

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
