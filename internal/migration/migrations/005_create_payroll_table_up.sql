CREATE TABLE IF NOT EXISTS payroll (
    id BIGSERIAL PRIMARY KEY,
    start_period TIMESTAMP WITH TIME ZONE NOT NULL,
    end_period TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_payroll_start_period ON payroll (start_period);
CREATE INDEX IF NOT EXISTS idx_payroll_end_period ON payroll (end_period);
CREATE INDEX IF NOT EXISTS idx_payroll_created_by ON payroll (created_by);
