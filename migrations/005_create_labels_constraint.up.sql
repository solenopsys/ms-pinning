-- remove duplicate labels
CREATE TEMP TABLE temp_labels AS SELECT DISTINCT ON (name, value, pin_id) * FROM labels;
TRUNCATE TABLE labels;
INSERT INTO labels SELECT * FROM temp_labels;
DROP TABLE temp_labels;

-- add primary key
ALTER TABLE labels ADD CONSTRAINT labels_pk PRIMARY KEY (name, value, pin_id);
