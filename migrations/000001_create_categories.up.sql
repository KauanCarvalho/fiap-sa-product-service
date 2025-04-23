CREATE TABLE IF NOT EXISTS `categories` (
    `id` SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_name` (`name`),
    INDEX `idx_deleted_at` (`deleted_at`)
);

INSERT INTO `categories` (`name`) VALUES ('lanche');
INSERT INTO `categories` (`name`) VALUES ('bebida');
INSERT INTO `categories` (`name`) VALUES ('sobremesa');
INSERT INTO `categories` (`name`) VALUES ('acompanhamento');
