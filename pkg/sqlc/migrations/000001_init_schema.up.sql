CREATE TABLE `accounts` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `owner` varchar(255) NOT NULL,
  `balance` bigint NOT NULL,
  `currency` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `entries` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `account_id` bigint NOT NULL,
  `amount` bigint NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `transfers` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `from_account_id` bigint NOT NULL,
  `to_account_id` bigint NOT NULL,
  `amount` bigint NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX `idx_owner` ON `accounts` (`owner`);
CREATE INDEX `idx_account_id` ON `entries` (`account_id`);
CREATE INDEX `idx_from_account_id` ON `transfers` (`from_account_id`);
CREATE INDEX `idx_to_account_id` ON `transfers` (`to_account_id`);
CREATE INDEX `idx_from_to_account_id` ON `transfers` (`from_account_id`, `to_account_id`);

ALTER TABLE `entries` ADD CONSTRAINT `fk_account_id` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`);
ALTER TABLE `transfers` ADD CONSTRAINT `fk_from_account_id` FOREIGN KEY (`from_account_id`) REFERENCES `accounts` (`id`);
ALTER TABLE `transfers` ADD CONSTRAINT `fk_to_account_id` FOREIGN KEY (`to_account_id`) REFERENCES `accounts` (`id`);
