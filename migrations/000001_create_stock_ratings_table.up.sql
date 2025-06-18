CREATE TABLE IF NOT EXISTS stock_ratings (
    id BIGSERIAL PRIMARY KEY,
    ticker VARCHAR(10) NOT NULL,
    target_from VARCHAR(50),
    target_to VARCHAR(50),
    company VARCHAR(255),
    action VARCHAR(50),
    brokerage VARCHAR(255),
    rating_from VARCHAR(50),
    rating_to VARCHAR(50),
    time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_stock_ratings_ticker ON stock_ratings(ticker);
CREATE INDEX IF NOT EXISTS idx_stock_ratings_time ON stock_ratings(time);
CREATE INDEX IF NOT EXISTS idx_stock_ratings_ticker_time ON stock_ratings(ticker, time DESC); 