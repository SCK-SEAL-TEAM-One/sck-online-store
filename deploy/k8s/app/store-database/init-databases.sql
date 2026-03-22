-- Create databases required by the application
CREATE DATABASE IF NOT EXISTS store CHARACTER SET utf8 COLLATE utf8_general_ci;
CREATE DATABASE IF NOT EXISTS point CHARACTER SET utf8 COLLATE utf8_general_ci;

-- Grant all privileges to app user
GRANT ALL ON *.* TO 'user'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
