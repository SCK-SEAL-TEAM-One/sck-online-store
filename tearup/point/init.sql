DROP DATABASE IF EXISTS point;
CREATE DATABASE IF NOT EXISTS point CHARACTER SET utf8 COLLATE utf8_general_ci;
USE point;

CREATE TABLE organizations (
    id BIGINT AUTO_INCREMENT,
    name varchar(255),
    created timestamp DEFAULT current_timestamp,
    updated timestamp DEFAULT current_timestamp ON UPDATE current_timestamp,
    PRIMARY KEY (id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;

INSERT INTO organizations (name) VALUES ("sck-online-store");

CREATE TABLE points (
    id BIGINT AUTO_INCREMENT,
    org_id BIGINT,
    user_id int,
    amount int,
    created timestamp DEFAULT current_timestamp,
    updated timestamp DEFAULT current_timestamp ON UPDATE current_timestamp,
    PRIMARY KEY (id),
    FOREIGN KEY (org_id) REFERENCES organizations(id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;