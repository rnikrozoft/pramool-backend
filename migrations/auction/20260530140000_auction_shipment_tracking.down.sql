SET statement_timeout = 0;

--bun:split

DROP INDEX IF EXISTS idx_auction_shipment_events_auction;

--bun:split

DROP TABLE IF EXISTS auction_shipment_events;

--bun:split

ALTER TABLE auctions
    DROP COLUMN IF EXISTS carrier_code,
    DROP COLUMN IF EXISTS carrier_name,
    DROP COLUMN IF EXISTS tracking_number,
    DROP COLUMN IF EXISTS shipment_status;
