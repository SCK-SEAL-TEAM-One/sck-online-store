import { faker } from "@faker-js/faker";
import fs from "fs";
import path from "path";

const minimumUser = 1;
const minimumProducts = 2;
const minimumShippingMethod = 1;
const minimumPaymentMethod = 1;

const seedUsers = parseInt(process.env.SEED_USERS) - minimumUser;
const seedProducts = parseInt(process.env.SEED_PRODUCTS) - minimumProducts;
const seedShippingMethods =
  parseInt(process.env.SEED_SHIPPING_METHODS) - minimumShippingMethod;
const seedPaymentMethods =
  parseInt(process.env.SEED_PAYMENT_METHODS) - minimumPaymentMethod;

if (
  seedUsers < 0 ||
  seedProducts < 0 ||
  seedShippingMethods < 0 ||
  seedPaymentMethods < 0
) {
  throw new Error(`Please set number of seed data more than minimum data:
    No. of minimum users: ${minimumUser}
    No. of minimum products: ${minimumProducts}
    No. of minimum shipping methods: ${minimumShippingMethod}
    No. of minimum payment methods: ${minimumPaymentMethod}`);
}

const file = "/app/output/tearup/store/init.sql";
const dir = path.dirname(file);

// Ensure directory exists
fs.mkdirSync(dir, { recursive: true });

// Create database schema
let sql = `DROP DATABASE IF EXISTS store;
CREATE DATABASE IF NOT EXISTS store CHARACTER SET utf8 COLLATE utf8_general_ci;
USE store;

CREATE TABLE users (
  id BIGINT AUTO_INCREMENT,
  first_name varchar(255),
  last_name varchar(255),
  created timestamp DEFAULT current_timestamp,
  updated timestamp DEFAULT current_timestamp ON UPDATE current_timestamp,
  PRIMARY KEY (id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE products (
  id BIGINT AUTO_INCREMENT,
  product_name varchar(255),
  product_brand varchar(255),
  stock int,
  product_price double,
  image_url varchar(255),
  created timestamp DEFAULT current_timestamp,
  updated timestamp DEFAULT current_timestamp ON UPDATE current_timestamp,
  PRIMARY KEY (id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE shipping_methods (
  id BIGINT AUTO_INCREMENT,
  name varchar(255),
  description varchar(255),
  fee double,
  created timestamp DEFAULT current_timestamp,
  updated timestamp DEFAULT current_timestamp ON UPDATE current_timestamp,
  PRIMARY KEY (id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE payment_methods (
  id BIGINT AUTO_INCREMENT,
  name varchar(255),
  description varchar(255),
  created timestamp DEFAULT current_timestamp,
  updated timestamp DEFAULT current_timestamp ON UPDATE current_timestamp,
  PRIMARY KEY (id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE carts (
  id BIGINT AUTO_INCREMENT,
  user_id BIGINT,
  product_id BIGINT,
  quantity int,
  created timestamp DEFAULT current_timestamp,
  updated timestamp DEFAULT current_timestamp ON UPDATE current_timestamp,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (product_id) REFERENCES products(id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE orders (
  id BIGINT AUTO_INCREMENT,
  user_id BIGINT,
  shipping_method_id BIGINT,
  payment_method_id BIGINT,
  sub_total_price double,
  discount_price double,
  total_price double,
  burn_point int,
  earn_point int,
  shipping_fee double,
  transaction_id varchar(255) DEFAULT '',
  status ENUM('created', 'paid', 'completed','cancel') DEFAULT 'created',
  authorized timestamp DEFAULT current_timestamp,
  created timestamp DEFAULT current_timestamp,
  updated timestamp DEFAULT current_timestamp ON UPDATE current_timestamp,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (shipping_method_id) REFERENCES shipping_methods(id),
  FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE order_product (
  order_id BIGINT,
  product_id BIGINT,
  quantity int,
  product_price double
) CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE order_shipping (
  id int AUTO_INCREMENT,
  order_id BIGINT,
  user_id BIGINT,
  method_id BIGINT,
  address varchar(255),
  sub_district varchar(255),
  district varchar(255),
  province varchar(255),
  zip_code varchar(5),
  recipient_first_name varchar(255),
  recipient_last_name varchar(255),
  phone_number varchar(13),
  created timestamp DEFAULT current_timestamp,
  updated timestamp DEFAULT current_timestamp ON UPDATE current_timestamp,
  PRIMARY KEY (id),
  FOREIGN KEY (order_id) REFERENCES orders(id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (method_id) REFERENCES shipping_methods(id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;
`;

// Insert users

// Default user
sql += `
INSERT INTO users (first_name, last_name) VALUES
("ponsakorn", "rungruangsap")`;

// Seed users
for (let i = 0; i < seedUsers; i++) {
  sql += `,
("${faker.person.firstName()}", "${faker.person.lastName()}")`;
}

sql += ";";

// Insert products

// Default products
sql += `
INSERT INTO products (product_name, product_brand, stock, product_price, image_url) VALUE 
("Balance Training Bicycle", "SportsFun", 100, 119.95, "/Balance_Training_Bicycle.png"),
("43 Piece dinner Set", "CoolKidz", 200, 12.95, "/43_Piece_dinner_Set.png")`;

// Seed products
for (let i = 0; i < seedProducts; i++) {
  sql += `,
("${faker.commerce.productName()}", "${faker.company.name()}", ${faker.number.int(
    { min: 0, max: 1000 }
  )}, ${faker.number.float({
    min: 0,
    max: 1000,
    fractionDigits: 2,
  })}, "/product_${minimumProducts + i + 1}.png")`;
}

sql += ";";

// Insert shipping methods

// Default shipping method
sql += `
INSERT INTO shipping_methods (name, description, fee) VALUE 
("Kerry", "4-5 business days", 50)`;

// Seed shipping methods
const shippingMethods = [
  "Thai Post",
  "Lineman",
  "Kerry",
  "Flash",
  "J&T",
  "SCG",
  "BEST",
  "Nim",
  "Grab",
  "Lalamove",
  "DHL",
  "FedEx",
  "UPS",
];
const shippingDurations = [
  "1-2 business days",
  "1-3 business days",
  "1-5 business days",
  "2-3 business days",
  "2-4 business days",
  "2-5 business days",
  "3-4 business days",
  "3-5 business days",
  "4-5 business days",
  "7-14 business days",
  "10-15 business days",
];

const totalShippingMethods = minimumShippingMethod + shippingMethods.length;

if (seedShippingMethods > totalShippingMethods) {
  throw new Error(
    `Shipping methods cannot be more than ${totalShippingMethods} methods`
  );
}

const shuffledShippingMethods = faker.helpers.shuffle(shippingMethods);

const selectedShippingMethods = shuffledShippingMethods.slice(
  0,
  seedShippingMethods
);

for (let i = 0; i < seedShippingMethods; i++) {
  sql += `,
("${selectedShippingMethods[i]}", "${faker.helpers.arrayElement(
    shippingDurations
  )}", ${faker.number.int({ min: 20, max: 150 })})`;
}

sql += ";";

// Insert payment methods

// Default payment method
sql += `
INSERT INTO payment_methods (name, description) VALUE 
("Credit Card / Debit Card", "")`;

// Seed payment methods
const paymentMethods = [
  "QR Code (PromptPay)",
  "TrueMoney Wallet",
  "Rabbit LINE Pay",
  "ShopeePay",
  "Mobile Banking",
  "Cash on Delivery",
  "Buy Now, Pay Later",
  "Over-the-Counter Payments",
];

const totalPaymentMethods = minimumPaymentMethod + paymentMethods.length;

if (seedPaymentMethods > totalPaymentMethods) {
  throw new Error(
    `Payment methods cannot be more than ${totalPaymentMethods} methods`
  );
}

const shuffledPaymentMethods = faker.helpers.shuffle(paymentMethods);

const selectedPaymentMethods = shuffledPaymentMethods.slice(
  0,
  seedPaymentMethods
);

for (let i = 0; i < seedPaymentMethods; i++) {
  sql += `,
("${selectedPaymentMethods[i]}", "")`;
}

sql += ";";

try {
  fs.writeFileSync(file, sql, "utf8");
} catch (err) {
  throw new Error("Error generate store data SQL file:", err);
}
