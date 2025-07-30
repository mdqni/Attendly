CREATE SCHEMA IF NOT EXISTS "lesson";

SET search_path TO "lesson";

CREATE TABLE IF NOT EXISTS lessons
(
    lesson_id  UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    teacher_id UUID                     NOT NULL, -- проверка роли через user-сервис
    subject    TEXT                     NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time   TIMESTAMP WITH TIME ZONE NOT NULL,

    latitude   NUMERIC(9, 6)            NOT NULL,
    longitude  NUMERIC(9, 6)            NOT NULL,
    room       TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_lessons_teacher_id ON lessons (teacher_id);
CREATE INDEX idx_lessons_start_time ON lessons (start_time);
