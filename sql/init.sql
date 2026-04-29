DROP TABLE IF EXISTS tasks;

CREATE TABLE tasks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tasks (title, description, completed) VALUES
    ('Изучить гошку', 'Пройти базовый курс', false),
    ('Написать Апишку', '', false),
    ('Зарелизить приложку', 'Развернуть на серваке', false);