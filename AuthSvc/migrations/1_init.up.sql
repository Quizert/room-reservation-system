CREATE TABLE users (
    ID SERIAL PRIMARY KEY,
    Username TEXT NOT NULL,
    ChatID TEXT NOT NULL,
    Password TEXT NOT NULL,
    IsHotelier BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_chat_id on users (ChatID);