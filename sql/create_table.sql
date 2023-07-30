CREATE TABLE
    engine_log1(
        id UUID NOT NULL COMMENT 'Primary KEY',
        time DateTime COMMENT 'time',
      
        action String COMMENT 'action',

        uid Int32 COMMENT 'uid',
        market  String COMMENT 'market',

    
    ) ENGINE = MergeTree()
ORDER BY (id) PRIMARY KEY(id) COMMENT '';
