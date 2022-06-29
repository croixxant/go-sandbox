ALTER TABLE `entries` DROP FOREIGN KEY `fk_account_id`;
ALTER TABLE `transfers` DROP FOREIGN KEY `fk_from_account_id`;
ALTER TABLE `transfers` DROP FOREIGN KEY `fk_to_account_id`;

DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS transfers;