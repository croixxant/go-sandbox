CREATE TABLE `wallets` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `balance` bigint NOT NULL,
  `currency` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `entries` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `wallet_id` bigint NOT NULL,
  `amount` bigint NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `transfers` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `from_wallet_id` bigint NOT NULL,
  `to_wallet_id` bigint NOT NULL,
  `amount` bigint NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CREATE INDEX `idx_user_id` ON `users` (`user_id`);
-- CREATE INDEX `idx_wallet_id` ON `entries` (`wallet_id`);
-- CREATE INDEX `idx_from_wallet_id` ON `transfers` (`from_wallet_id`);
-- CREATE INDEX `idx_to_wallet_id` ON `transfers` (`to_wallet_id`);
-- CREATE INDEX `idx_from_to_wallet_id` ON `transfers` (`from_wallet_id`, `to_wallet_id`);

-- ALTER TABLE `entries` ADD CONSTRAINT `fk_wallet_id` FOREIGN KEY (`wallet_id`) REFERENCES `wallets` (`id`);
-- ALTER TABLE `transfers` ADD CONSTRAINT `fk_wallet_id` FOREIGN KEY (`from_wallet_id`) REFERENCES `wallets` (`id`);
-- ALTER TABLE `transfers` ADD CONSTRAINT `fk_to_wallet_id` FOREIGN KEY (`to_wallet_id`) REFERENCES `wallets` (`id`);
