-- cw: hacked for unique with tag_filter in immutable_tag_rule
--     unique_together supports max bytes 3072. tried that only support varchar around 1000 
-- cw: CRC hash will have hash conflict on big table. (93k records causes 1% change)
-- cw: Use MD5 instead. Although it use higher CPU, but **its NO PROBLEM because its UI ACTION** .

-- REMOVED. Use code caclulate instead.

-- CREATE TRIGGER `immutable_tag_rule_insert` BEFORE INSERT ON `immutable_tag_rule` FOR EACH ROW
-- BEGIN
--   set NEW.tag_filter_hash=md5(NEW.tag_filter);
-- END;

-- CREATE TRIGGER `immutable_tag_rule_update` BEFORE UPDATE ON `immutable_tag_rule` FOR EACH ROW
-- BEGIN
--   set NEW.tag_filter_hash=md5(NEW.tag_filter);
-- END;
