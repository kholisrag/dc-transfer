INSERT INTO `test_snapshot_and_increment` (`id`, `value`, `value_10`)
VALUES
    (10,  0,                                                                   0),
    (11,  9999999999999999999999999999999999,                                  99999),
    (12,  99999999999999999999999999999999999,                                 9999999),
    (13,  9999999999999999999999999999999999999999999999999999999999999999,    9999999999),
    (14,  99999999999999999999999999999999999999999999999999999999999999999,   9999999999),
    (15,  999999999999999999999999999999999999.99999999999999999999999999999,  9999999999),
    (16,  99999999999999999999999999999999999.999999999999999999999999999999,  9999999999),
    (17,  1.000000000000000000000000000001,                                    1),
    (18,  NULL,                                                                9999999999),
    (19,  99999999999999999999999999999999999.999999999999999999999999999999,  NULL)
;