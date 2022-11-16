
-- +migrate Up
CREATE TABLE IF NOT EXISTS `be-wallet`.`transaction_record` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
    `member_id` BIGINT UNSIGNED NOT NULL COMMENT '交易者id',
    `wallet_id` BIGINT UNSIGNED NOT NULL COMMENT '錢包id',
    `action` TINYINT(4) UNSIGNED NOT NULL COMMENT '操作動作 1:入金 2:出金 3:贈送 4:利息 5:開倉 6:平倉 7:人工更改',
    `amount` DECIMAL(19,4) NOT NULL COMMENT '金額',
    `before_amount` DECIMAL(19,4) NULL DEFAULT NULL COMMENT '交易前金額',
    `after_amount` DECIMAL(19,4) NULL DEFAULT NULL COMMENT '交易後金額',
    `currency` VARCHAR(36) NOT NULL COMMENT '幣別',
    `committer_id` BIGINT UNSIGNED NOT NULL COMMENT '執行者id',
    `status` TINYINT(4) UNSIGNED NOT NULL COMMENT '交易狀態 1:待辦 2:成功 3:失敗 4:回滾',
    `remark` VARCHAR(128) COMMENT '註記',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '創建時間',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新時間',
    `rollback_before_amount` DECIMAL(19,4) NULL DEFAULT NULL COMMENT '回滾前金額',
    `rollback_after_amount` DECIMAL(19,4) NULL DEFAULT NULL COMMENT '回滾後金額',
    `rollbacker_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '回滾執行者id',
    
    INDEX (`created_at` DESC),
    INDEX (`member_id`, `created_at` DESC),
    FOREIGN KEY (`wallet_id`) REFERENCES wallet(`id`) ON DELETE CASCADE,
    PRIMARY KEY (`id`)
) AUTO_INCREMENT=1 CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易紀錄';

-- +migrate Down
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `transaction_record`;

