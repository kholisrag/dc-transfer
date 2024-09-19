CREATE TABLE `time_test` (
    `id` integer NOT NULL PRIMARY KEY,
    `col_d` date,
    `col_dt` datetime,
    `col_ts` timestamp
) engine=innodb default charset=utf8;

INSERT INTO `time_test` VALUES
    (1, '2020-12-23', '2020-12-23 14:15:16', '2020-12-23 14:15:16'),
    (2, '2020-12-24', '2020-12-24 14:15:16', '2020-12-24 14:15:16'),
    (3, '1970-01-01', '1970-01-01 00:00:00', '1970-01-01 00:00:00'), -- yt has minimal allowed value for 1970-01-01
    (4, NULL, NULL, NULL),
    (5, '1989-11-09', '1989-11-09 19:02:03.456789', '1989-11-09 19:02:03.456789'),
    (6, '1970-01-01', '1970-01-01 00:00:00', '1970-01-01 00:00:00'),
    (7, '2025-05-25', '2025-05-25 00:05:25.555', '2025-05-25 00:05:25.555555');