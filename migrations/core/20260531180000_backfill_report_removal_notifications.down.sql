-- Irreversible backfill; remove only rows created by this migration's kinds for affected users.
DELETE FROM user_notifications
WHERE kind IN ('listing_removed_seller', 'listing_removed_bidder', 'report_reward');
