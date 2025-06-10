CREATE TABLE IF NOT EXISTS reimbursement (
    id BIGSERIAL PRIMARY KEY,
    request_id UUID NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    role user_role NOT NULL,
    reimbursement_amount DECIMAL(20, 2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_reimbursement_user_id ON reimbursement (user_id);
CREATE INDEX IF NOT EXISTS idx_reimbursement_request_id ON reimbursement (request_id);
