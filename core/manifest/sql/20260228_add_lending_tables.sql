-- 借贷仓位表
CREATE TABLE IF NOT EXISTS lending_positions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL DEFAULT 1,
    protocol VARCHAR(50) NOT NULL,  -- aave, compound
    chain VARCHAR(50) NOT NULL,     -- ethereum, polygon, arbitrum
    wallet_address VARCHAR(100) NOT NULL,

    -- 存款信息
    supply_token VARCHAR(20),
    supply_amount DECIMAL(30, 8) DEFAULT 0,
    supply_value_usd DECIMAL(20, 2) DEFAULT 0,
    supply_apy DECIMAL(10, 4) DEFAULT 0,

    -- 借款信息
    borrow_token VARCHAR(20),
    borrow_amount DECIMAL(30, 8) DEFAULT 0,
    borrow_value_usd DECIMAL(20, 2) DEFAULT 0,
    borrow_apy DECIMAL(10, 4) DEFAULT 0,

    -- 健康因子
    health_factor DECIMAL(10, 4),
    liquidation_threshold DECIMAL(10, 4),
    ltv DECIMAL(10, 4),  -- Loan-to-Value 比率

    -- 净收益
    net_apy DECIMAL(10, 4),

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_lending_positions_user ON lending_positions(user_id);
CREATE INDEX IF NOT EXISTS idx_lending_positions_protocol ON lending_positions(protocol);
CREATE INDEX IF NOT EXISTS idx_lending_positions_health ON lending_positions(health_factor);

-- 借贷利率历史表
CREATE TABLE IF NOT EXISTS lending_rate_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    protocol VARCHAR(50) NOT NULL,
    chain VARCHAR(50) NOT NULL,
    token VARCHAR(20) NOT NULL,

    -- 存款利率
    supply_apy DECIMAL(10, 4),

    -- 借款利率
    borrow_apy_stable DECIMAL(10, 4),
    borrow_apy_variable DECIMAL(10, 4),

    -- 利用率
    utilization_rate DECIMAL(10, 4),

    recorded_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_lending_rate_protocol_token ON lending_rate_history(protocol, token, recorded_at);
