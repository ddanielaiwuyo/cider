create table IF NOT EXISTS users (
	user_id SERIAL PRIMARY KEY,
	username VARCHAR(50) UNIQUE  NOT NULL,
	email VARCHAR(250) UNIQUE  NOT NULL,
	password VARCHAR(50) NOT NULL,
	created_at TIMESTAMP NOT NULL
);

create table IF NOT EXISTS messages (
	message_id SERIAL PRIMARY KEY,
	sender_id INT REFERENCES users(user_id),
	recipient_id INT REFERENCES users(user_id),
	content TEXT,
	timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
	 conversation_id INT REFERENCES conversations(conversation_id)
);

create table IF NOT EXISTS conversations (
	conversation_id SERIAL PRIMARY KEY
);
