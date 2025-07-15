CREATE TABLE IF NOT EXISTS reviews
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    photo       VARCHAR(255) NOT NULL,
    pdf_file    VARCHAR(255),
    industry    VARCHAR(255),
    service     VARCHAR(255),
    description TEXT,
    rating      INT,
    created_at  DATETIME     NOT NULL,
    updated_at  DATETIME
);
use reviews;