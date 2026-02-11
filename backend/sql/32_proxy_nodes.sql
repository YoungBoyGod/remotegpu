CREATE TABLE IF NOT EXISTS proxy_nodes (
    id              VARCHAR(64) PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    host            VARCHAR(256) NOT NULL,
    api_port        INT NOT NULL DEFAULT 9090,
    http_port       INT NOT NULL DEFAULT 9091,
    range_start     INT NOT NULL DEFAULT 20000,
    range_end       INT NOT NULL DEFAULT 60000,
    version         VARCHAR(32),
    status          VARCHAR(20) DEFAULT 'offline',
    active_mappings INT DEFAULT 0,
    used_ports      INT DEFAULT 0,
    last_heartbeat  TIMESTAMP,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
