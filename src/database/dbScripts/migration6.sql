
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
