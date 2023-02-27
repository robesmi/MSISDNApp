DROP DATABASE IF EXISTS `msisdn`;
CREATE DATABASE `msisdn`;
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
    ("^48[0-9]{9}$",48,"pl",2),
    ("^971[0-9]{10}$",971,"ae",3),
    ("^850[0-9]{10}$",850,"kp",3),
    ("^43[0-9]{6,13}$",43,"at",2),
    ("^351[0-9]{9}$",351,"pt",3),
    ("^1246[0-9]{10}$",1246,"bb",3),
    ("^212[0-9]{9}$",212,"ma",3);

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
    ("pl","^5366[0-9]{6}$", "Polskie Sieci Cyfrowe Sp. z o.o.",2),
    ("ae","^(50|56)[0-9]{7}$", "Etisalat",2),
    ("ae","^(52|55)[0-9]{7}$", "Du",2),
    ("kp","^(191|192)[0-9]{7}$", "Koryolink",3),
    ("kp","^195[0-9]{7}$", "KangsongNET",3),
    ("at","^(660|699)[0-9]{3,10}$", "Hutchinson Drei",3),
    ("pt","^91[0-9]{7}$", "Vodafone Portugal",2),
    ("pt","^921[0-9]{6}$", "Vodafone Portugal",3),
    ("pt","^922[0-2][0-9]{5}$", "CTT CORREIOS DE PORTUGAL, S.A.",3),
    ("pt","^924[0-4][0-9]{5}$", "TMN - TELECOMUNICAÇÕES MÓVEIS NACIONAIS, SA",3),
    ("bb","^(23\d|24\d|25[0-4])[0-9]{7}$", "Liberty Latin America",3),
    ("bb","^(45[0-9])[0-9]{7}$", "Sunbeach",3),
    ("ma","^(?!^61(2|4|7|9)[0-9]{5}$)611[0-9]{6}$", "Maroc Telecom",3),
    ("ma","^61(2|4|7|9)[0-9]{6}$", "Orange Maroc",3);