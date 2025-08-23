CREATE TABLE IF NOT EXISTS messages (
  id SERIAL PRIMARY KEY,
  phone_number VARCHAR(20) NOT NULL,
  content VARCHAR(160) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',

  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
  sent_at TIMESTAMPTZ,
  response_message_id VARCHAR(255),
  fail_reason TEXT,
  retry_count INT DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_messages_status_created ON messages(status, created_at);
CREATE INDEX IF NOT EXISTS idx_messages_phone_number ON messages(phone_number);
