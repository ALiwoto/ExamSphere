

ALTER TABLE topic_info ADD CONSTRAINT unique_topic_name UNIQUE (topic_name);
ALTER TABLE course_info ADD CONSTRAINT unique_course_name UNIQUE (course_name);
ALTER TABLE exam_info ADD CONSTRAINT unique_exam_title UNIQUE (exam_title);
