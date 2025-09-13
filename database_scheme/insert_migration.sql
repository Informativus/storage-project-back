-- Заполнение таблицы ролей валидными данными.
INSERT INTO roles (id, name, note) VALUES (1, 'ADMIN', 'Owner of the system');
INSERT INTO roles (id, name, note) VALUES (2, 'USER', 'User of system');
INSERT INTO roles (id, name, note) VALUES (3, 'OWNER', 'Owner of the folder');

-- Заполнение таблицы доступных действий для хранилища
INSERT INTO actions (id, name, note) VALUES 
    (1, 'READ_FILE', 'Чтение файлов'),
    (2, 'UPLOAD_FILE', 'Загрузка файлов'),
    (3, 'DELETE_FILE', 'Удаление файлов'),
    (4, 'EDIT_FILE', 'Редактирование файлов');

-- Protection Groups
INSERT INTO protection_groups (id, name, description) VALUES
  ('5f9352a3-4d78-409c-b42b-eeba7fbd779c', 'Reader', 'Базовый доступ только к чтению'),
  ('4d81a112-169b-4659-9d1f-d6745aa96f1c', 'Editor', 'Доступ к чтению, редактированию и добавлению файлов'),
  ('87d8ac99-6755-44e6-b4a9-cbd69ab31a4a', 'Owner', 'Полный контроль, включая управление доступом (скрытая группа)');

-- Reader → только READ_FILE
INSERT INTO protection_group_actions (group_id, action_id) VALUES
  ('5f9352a3-4d78-409c-b42b-eeba7fbd779c', 1);

-- Editor → READ_FILE, UPLOAD_FILE, EDIT_FILE
INSERT INTO protection_group_actions (group_id, action_id) VALUES
  ('4d81a112-169b-4659-9d1f-d6745aa96f1c', 1),
  ('4d81a112-169b-4659-9d1f-d6745aa96f1c', 2),
  ('4d81a112-169b-4659-9d1f-d6745aa96f1c', 4);

-- Owner → все действия
INSERT INTO protection_group_actions (group_id, action_id) VALUES
  ('87d8ac99-6755-44e6-b4a9-cbd69ab31a4a', 1),
  ('87d8ac99-6755-44e6-b4a9-cbd69ab31a4a', 2),
  ('87d8ac99-6755-44e6-b4a9-cbd69ab31a4a', 3),
  ('87d8ac99-6755-44e6-b4a9-cbd69ab31a4a', 4);
