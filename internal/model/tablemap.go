package model

type Table string

const (
	USERS         Table = "users"
	REIMBURSEMENT Table = "reimbursement"
	ATTENDANCE    Table = "attendance"
	OVERTIME      Table = "overtime"
	PAYROLL       Table = "payroll"
	AUDITLOG      Table = "audit_log"
)
