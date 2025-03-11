CREATE TABLE IF NOT EXISTS `bpm_event_log_queue` (
    `ID` int unsigned NOT NULL AUTO_INCREMENT,
    `DATA` text COLLATE utf8mb4_general_ci DEFAULT NULL,
    `CREATED_AT` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `ATTEMPTS` INT NOT NULL DEFAULT 0,
    `LAST_ATTEMPT_AT` datetime
PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TRIGGER after_insert_b_event_log
    AFTER INSERT ON b_event_log
    FOR EACH ROW
    INSERT INTO bpm_event_log_queue (DATA)
    VALUES (
               CONCAT_WS("",
                       '{"ID":"', NEW.ID,
                       '","SEVERITY":"', NEW.SEVERITY,
                       '","MODULE_ID":"', NEW.MODULE_ID,
                       '","ITEM_ID":"', NEW.ITEM_ID,
                       '","REQUEST_URI":"', NEW.REQUEST_URI,
                       '","AUDIT_TYPE_ID":"', NEW.AUDIT_TYPE_ID,
                       '","DESCRIPTION":"', NEW.DESCRIPTION,
                       '","TIMESTAMP_X":"', NEW.TIMESTAMP_X,
                       '","DATE":"', NOW(), '"}\n'
               )
    );
