-- 风险指标表
-- 存储每日计算的风险管理指标
CREATE TABLE IF NOT EXISTS risk_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    user_id INTEGER NOT NULL DEFAULT 1,
    metric_date DATE NOT NULL,                    -- 指标计算日期
    portfolio_value REAL NOT NULL DEFAULT 0,      -- 组合总价值（USD）
    var_95 REAL NOT NULL DEFAULT 0,               -- 95% 置信度 VaR
    var_99 REAL NOT NULL DEFAULT 0,               -- 99% 置信度 VaR
    sharpe_ratio REAL NOT NULL DEFAULT 0,         -- 夏普比率
    sortino_ratio REAL NOT NULL DEFAULT 0,        -- 索提诺比率
    max_drawdown REAL NOT NULL DEFAULT 0,         -- 最大回撤（百分比）
    max_drawdown_duration INTEGER NOT NULL DEFAULT 0, -- 最大回撤持续天数
    beta REAL NOT NULL DEFAULT 0,                 -- Beta 系数（相对 BTC）
    volatility REAL NOT NULL DEFAULT 0,           -- 波动率（年化）
    downside_deviation REAL NOT NULL DEFAULT 0,   -- 下行偏差
    calculation_period INTEGER NOT NULL DEFAULT 30, -- 计算周期（天数）
    UNIQUE(user_id, metric_date)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_risk_metrics_user_date ON risk_metrics(user_id, metric_date DESC);
CREATE INDEX IF NOT EXISTS idx_risk_metrics_date ON risk_metrics(metric_date DESC);
