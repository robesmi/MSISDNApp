DROP TABLE IF EXISTS `countries`;
CREATE TABLE `countries` (
    `country_number_format` varchar(100) NOT NULL,
    `country_code` varchar(6) NOT NULL,
    `country_identifier` varchar(20) NOT NULL,
    `country_code_length` int NOT NULL,
    PRIMARY KEY (`country_number_format`)
);
INSERT INTO `countries` VALUES
    ("^389[0-9]{8}$",389,"mk",3),
    ("^350[0-9]{5}$",350,"gi",3),
    ("^242[0-9]{9}$",242,"cg",3);

DROP TABLE IF EXISTS `mobile_operators`;
CREATE TABLE `mobile_operators` (
    `country_identifier` varchar(20) NOT NULL,
    `prefix_format` varchar(100) NOT NULL,
    `mno`           varchar(100) NOT NULL,
    `prefix_length` int NOT NULL,
    PRIMARY KEY (`country_identifier`, `prefix_format`)
);

INSERT INTO `mobile_operators` VALUES
    ("mk","^77[0-9]{6}$","A1",2),
    ("mk","^71[0-9]{6}$", "Telekom",2);
