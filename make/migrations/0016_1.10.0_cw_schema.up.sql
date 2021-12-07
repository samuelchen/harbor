-- cw: hacked for unique with tag_filter in immutable_tag_rule
-- cw: CRC hash will have hash conflict on big table. (93k records causes 1% change)
-- cw: DO NOT use SHA / MD5. They use higher CPU & store.
DELIMITER $$

CREATE TRIGGER `tag_filter_hash_insert` BEFORE INSERT ON `immutable_tag_rule` FOR EACH ROW
BEGIN
  set NEW.tag_filter_hash=crc32(NEW.tag_filter);
END $$

CREATE TRIGGER `tag_filter_hash_update` BEFORE UPDATE ON `immutable_tag_rule` FOR EACH ROW
BEGIN
  set NEW.tag_filter_hash=crc32(NEW.tag_filter);
END $$

DELIMITER ;
