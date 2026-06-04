-- +goose Up
-- +goose StatementBegin
-- Демо-админ для защиты проекта: admin@gmail.com / 12341234
-- Хэш = bcrypt(пароль + salt), как в internal/domain/user/password.go.
INSERT INTO users (id, email, username, password_hash, salt, role, subscription_type) VALUES
    ('00000000-0000-0000-0000-0000000000ad', 'admin@gmail.com', 'admin',
     '$2a$10$Mq1zPIlIX5sWVSa0nmh6KuAD60ZMCeZHIMrRXMd/HDXURCPy51G4.',
     'c0ffee1234567890abcdef9876543210',
     'ADMIN', 'PREMIUM')
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id = '00000000-0000-0000-0000-0000000000ad';
-- +goose StatementEnd
