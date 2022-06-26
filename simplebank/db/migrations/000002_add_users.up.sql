CREATE TABLE `users` (
  `username` varchar(255) PRIMARY KEY,
  `hashed_password` varchar(255) NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `password_changed_at` datetime NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE `accounts` ADD FOREIGN KEY (`owner`) REFERENCES `users` (`username`);

CREATE UNIQUE INDEX `idx_owner_currency` ON `accounts` (`owner`, `currency`);
