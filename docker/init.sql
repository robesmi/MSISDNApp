DROP DATABASE IF EXISTS `msisdn`;
CREATE DATABASE IF NOT EXISTS `msisdn`;
USE `msisdn`;
DROP TABLE IF EXISTS `countries`;
CREATE TABLE `countries` (
    `country_number_format` varchar(20) NOT NULL,
    `country_code` varchar(6) NOT NULL,
    `country_identifier` varchar(3) NOT NULL,
    `country_code_length` int NOT NULL,
    PRIMARY KEY (`country_number_format`)
);
INSERT INTO `countries` VALUES
    ("^389[0-9]{8}$",389,"mk",3),
    ("^350[0-9]{5}$",350,"gi",3),
    ("^242[0-9]{9}$",242,"cg",3),
    ("^423[0-9]{8}$",423,"li",3),
    ("^48[0-9]{9}$",48,"pl",2);

DROP TABLE IF EXISTS `mobile_operators`;
CREATE TABLE `mobile_operators` (
    `country_identifier` varchar(3) NOT NULL,
    `prefix_format` varchar(60) NOT NULL,
    `mno`           varchar(100) NOT NULL,
    `prefix_length` int NOT NULL,
    PRIMARY KEY (`country_identifier`, `prefix_format`)
);

INSERT INTO `mobile_operators` VALUES
    ("mk","^77[0-9]{6}$","A1",2),
    ("mk","^71[0-9]{6}$", "Telekom",2),
    ("li","^6[0-9]{7}$", "Lietuvos",1),
    ("pl","^510[0-9]{6}$", "Mobile telephoOrange",2),
    ("pl","(?!^5329[0-9]{5}$)(?!^5366[0-9]{5}$)^53[0-7][0-9]{6}$", "Orange Polska S.A",2),
    ("pl","^53(2|8|9)[0-9]{6}$", "T-MOBILE POLSKA S.A.",2),
    ("pl","^5366[0-9]{6}$", "Polskie Sieci Cyfrowe Sp. z o.o.",2);