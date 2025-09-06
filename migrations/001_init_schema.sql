-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       phone_number VARCHAR(20) UNIQUE NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_phone_number ON users(phone_number);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Add a table to track OTP requests for rate limiting (optional)
CREATE TABLE otp_requests (
                              id SERIAL PRIMARY KEY,
                              phone_number VARCHAR(20) NOT NULL,
                              requested_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                              successful BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_otp_requests_phone_number ON otp_requests(phone_number);
CREATE INDEX idx_otp_requests_requested_at ON otp_requests(requested_at);

-- +goose Down
DROP TABLE IF EXISTS otp_requests;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";