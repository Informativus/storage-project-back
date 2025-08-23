-- Заполнение таблицы ролей валидными данными.
INSERT INTO roles (name, note) VALUES ('ADMIN', 'Owner of the system');
INSERT INTO roles (name, note) VALUES ('OWNER', 'Owner of the folder');
INSERT INTO roles (name, note) VALUES ('READER', 'Reader of the folder');
INSERT INTO roles (name, note) VALUES ('EDITOR', 'All CRUD operations of the folder');
