ALTER TABLE `accounts` DROP INDEX `idx_owner_currency`;

ALTER TABLE `accounts` DROP FOREIGN KEY `accounts_ibfk_1`;

DROP TABLE IF EXISTS users;
