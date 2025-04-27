BEGIN;

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW."updated_at" = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS service_errors (
        id SERIAL PRIMARY KEY,
        "priority" INT,
        "error_message" TEXT,
        "additional_info" TEXT,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );
    
CREATE TABLE IF NOT EXISTS kafka_topic_dropped_messages (
    id SERIAL PRIMARY KEY,
    "topic_name" TEXT,
    "error_type" TEXT,
    "kafka_message" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    "user_id" SERIAL PRIMARY KEY,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "password_hash" TEXT NOT NULL,
    "first_name" VARCHAR(100),                   
    "last_name" VARCHAR(100),   
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS accounts (
    "account_id" SERIAL PRIMARY KEY,                        
    "user_id" INT NOT NULL,                          -- Foreign key to users table (this links to a user)
    "balance" INT NOT NULL DEFAULT 0,                 -- The account balance in rupees (using INT to avoid float issues)
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    CONSTRAINT "fk_user" FOREIGN KEY("user_id") REFERENCES users(user_id) ON DELETE CASCADE, -- Foreign key constraint
    CHECK ("balance" >= 0)                           -- Ensure that balance cannot be negative
);

CREATE TRIGGER set_timestamp BEFORE
UPDATE ON accounts FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_user_email ON users("email");
CREATE INDEX idx_user_id ON accounts("user_id");

COMMIT;