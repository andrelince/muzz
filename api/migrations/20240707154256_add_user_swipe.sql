-- +goose Up
CREATE TABLE user_swipes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    swiped_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    swipe_status BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, swiped_user_id)
);

CREATE INDEX idx_user_swipes_user_id ON user_swipes(user_id);

CREATE INDEX idx_user_swipes_swiped_user_id ON user_swipes(swiped_user_id);

-- +goose Down
DROP TABLE user_swipes;