-- 文档中心表
CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    category VARCHAR(64) NOT NULL DEFAULT 'general',
    file_name VARCHAR(256) NOT NULL,
    file_path VARCHAR(512) NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    content_type VARCHAR(128) NOT NULL DEFAULT '',
    storage_backend VARCHAR(64) NOT NULL DEFAULT '',
    uploaded_by INT REFERENCES customers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_documents_category ON documents(category);
CREATE INDEX idx_documents_uploaded_by ON documents(uploaded_by);
