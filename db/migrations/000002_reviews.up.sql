CREATE TABLE IF NOT EXISTS reviews
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    photo       VARCHAR(255) NOT NULL,
    description TEXT,
    rating      INT,
    created_at  DATETIME     NOT NULL,
    updated_at  DATETIME
);
use reviews;