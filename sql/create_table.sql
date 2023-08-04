CREATE TABLE
    record(
        id UUID NOT NULL COMMENT 'Primary KEY',
        time DateTime COMMENT 'time',
        time_utc DateTime COMMENT 'time_utc',
        name String COMMENT 'String',
        action String COMMENT 'action',
        duration Int32 COMMENT 'duration',
    ) ENGINE = MergeTree()
ORDER BY (id) PRIMARY KEY(id) COMMENT '';