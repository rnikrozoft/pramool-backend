SET statement_timeout = 0;

--bun:split

ALTER TABLE auctions
    ADD COLUMN IF NOT EXISTS carrier_code VARCHAR(50),
    ADD COLUMN IF NOT EXISTS carrier_name TEXT,
    ADD COLUMN IF NOT EXISTS tracking_number TEXT,
    ADD COLUMN IF NOT EXISTS shipment_status VARCHAR(50) NOT NULL DEFAULT 'pending';

--bun:split

CREATE TABLE IF NOT EXISTS auction_shipment_events (
    event_id BIGSERIAL PRIMARY KEY,
    auction_id VARCHAR(40) NOT NULL REFERENCES auctions (auction_id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    location TEXT,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_auction_shipment_events_auction ON auction_shipment_events (auction_id, created_at DESC);
