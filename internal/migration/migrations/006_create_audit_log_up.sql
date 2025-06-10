DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'action_type') THEN
        CREATE TYPE action_type AS ENUM ('CREATE', 'UPDATE', 'READ', 'DELETE');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS audit_log (
    id BIGSERIAL PRIMARY KEY,
    request_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    action_type action_type NOT NULL,
    affected_table VARCHAR(255) NOT NULL,
    affected_record_id BIGINT,
    created_by BIGINT NOT NULL,
    ip_address VARCHAR(255),
    old_data TEXT,
    new_data TEXT,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_audit_log_request_id ON audit_log (request_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_created_at ON audit_log (created_at);
CREATE INDEX IF NOT EXISTS idx_audit_log_event_type ON audit_log (event_type);
CREATE INDEX IF NOT EXISTS idx_audit_log_action_type ON audit_log (action_type);
CREATE INDEX IF NOT EXISTS idx_audit_log_affected_table ON audit_log (affected_table);
CREATE INDEX IF NOT EXISTS idx_audit_log_affected_record_id ON audit_log (affected_record_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_created_by ON audit_log (created_by);
