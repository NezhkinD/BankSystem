-- password = password
INSERT INTO users (username, email, password)
VALUES ('testuser', 'test@example.com', '$2a$10$Vexamplehashedpassword1234567890');

INSERT INTO accounts (user_id, balance)
VALUES (1, 1000.00);