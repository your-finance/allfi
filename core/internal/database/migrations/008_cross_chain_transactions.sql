-- 跨链交易表
CREATE TABLE IF NOT EXISTS cross_chain_transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL,

    -- 交易信息
    tx_hash VARCHAR(100) NOT NULL,
    bridge_protocol VARCHAR(50) NOT NULL,

    -- 源链信息
    source_chain VARCHAR(20) NOT NULL,
    source_token VARCHAR(20) NOT NULL,
    source_amount REAL NOT NULL,

    -- 目标链信息
    dest_chain VARCHAR(20) NOT NULL,
    dest_token VARCHAR(20) NOT NULL,
    dest_amount REAL NOT NULL,

    -- 费用
    bridge_fee REAL,
    gas_fee REAL,
    total_fee_usd REAL,

    -- 状态
    status VARCHAR(20) NOT NULL,
    initiated_at DATETIME NOT NULL,
    completed_at DATETIME,

    UNIQUE(tx_hash)
);

CREATE INDEX IF NOT EXISTS idx_cross_chain_tx_user ON cross_chain_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_cross_chain_tx_status ON cross_chain_transactions(status);
CREATE INDEX IF NOT EXISTS idx_cross_chain_tx_initiated ON cross_chain_transactions(initiated_at);
