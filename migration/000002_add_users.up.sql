CREATE TABLE `users` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(255) UNIQUE NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ALTER TABLE `users` ADD CONSTRAINT `fk_owner` FOREIGN KEY (`owner`) REFERENCES `users` (`username`);

-- CREATE UNIQUE INDEX `idx_account_id_currency` ON `wallets` (`account_id`, `currency`);
-- CREATE UNIQUE INDEX `idx_username` ON `users` (`username`);
