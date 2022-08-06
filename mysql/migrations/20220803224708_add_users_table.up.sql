CREATE TABLE `users` (
    `id` VARCHAR(100) NOT NULL,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100) NOT NULL,
    `nickname` VARCHAR(100) NOT NULL,
    `email` VARCHAR(100) NOT NULL,
    `encoded_password` VARCHAR(200) NOT NULL,
    `country` VARCHAR(8) NOT NULL ,
    `created_at` TIMESTAMP(6) NOT NULL DEFAULT current_timestamp(6),
    `updated_at` TIMESTAMP(6) NOT NULL DEFAULT current_timestamp(6) ON UPDATE current_timestamp(6),
    `deleted_at` TIMESTAMP(6) NULL,
    PRIMARY KEY (`id`),
    INDEX (`email`)
) ENGINE=InnoDB CHARSET=utf8 COLLATE utf8_general_ci;
