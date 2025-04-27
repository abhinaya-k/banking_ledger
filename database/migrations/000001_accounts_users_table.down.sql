BEGIN;

  DROP TRIGGER IF EXISTS set_timestamp ON accounts;
  DROP TRIGGER IF EXISTS set_timestamp ON users;

  DROP index if exists "idx_user_email";
  DROP index if exists "idx_user_id";

  DROP TABLE IF EXISTS cisco_meraki_errors;
  DROP TABLE IF EXISTS kafka_topic_dropped_messages;
  DROP TABLE IF EXISTS users;
  DROP TABLE IF EXISTS accounts;

  DROP FUNCTION IF EXISTS TRIGGER_SET_TIMESTAMP;
  
COMMIT;
