CREATE TABLE users (
  id              BIGINT        NOT NULL AUTO_INCREMENT PRIMARY KEY,
  email           varchar(255)  NOT NULL,
  hashed_password varchar(255)  NOT NULL,
  confirmed_at    datetime      DEFAULT NULL,
  likes_count     int           NOT NULL DEFAULT 0
);

CREATE TABLE posts (
  id          BIGINT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id     BIGINT  NOT NULL,
  body        text    NOT NULL,
  likes_count int     NOT NULL,
  INDEX idx_user_id (user_id)
);
