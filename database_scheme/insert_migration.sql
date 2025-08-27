-- Заполнение таблицы ролей валидными данными.
INSERT INTO roles (id, name, note) VALUES (1, 'ADMIN', 'Owner of the system');
INSERT INTO roles (id, name, note) VALUES (2, 'USER', 'User of system');
INSERT INTO roles (id, name, note) VALUES (3, 'OWNER', 'Owner of the folder');
INSERT INTO roles (id, name, note) VALUES (4, 'READER', 'Reader of the folder');
INSERT INTO roles (id, name, note) VALUES (5, 'EDITOR', 'All CRUD operations of the folder');
