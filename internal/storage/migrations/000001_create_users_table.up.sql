CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		passport_number VARCHAR(20) NOT NULL UNIQUE
	);

CREATE INDEX IF NOT EXISTS users_passport_number_idx ON users (passport_number);

CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		UserID INT NOT NULL REFERENCES users (id),
		Description TEXT NOT NULL,
		StartTime TIMESTAMP NOT NULL,
		EndTime TIMESTAMP NOT NULL
	);