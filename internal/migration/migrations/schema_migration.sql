CREATE TABLE IF NOT EXISTS schema_migration (

schema TEXT PRIMARY KEY,

created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP

);