/* change the data type to text to accommodate larger data */
-- ALTER TABLE properties ALTER COLUMN v TYPE text;
ALTER TABLE properties CHANGE COLUMN v TYPE text;