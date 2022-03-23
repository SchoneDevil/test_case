CREATE TABLE IF NOT EXISTS users (
    timestamp UInt64,
    userid      String
) ENGINE = Kafka SETTINGS
    kafka_broker_list = 'broker:9092',
    kafka_topic_list = 'users',
    kafka_group_name = 'statistics',
    kafka_format = 'JSONEachRow',
    kafka_num_consumers = 2;


CREATE TABLE IF NOT EXISTS users_stats (
    timestamp UInt64,
    userid      String
    ) ENGINE = MergeTree()
    ORDER BY timestamp;

CREATE MATERIALIZED VIEW users_queue_mv TO users_stats AS SELECT * FROM users;