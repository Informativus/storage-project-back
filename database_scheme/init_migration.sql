CREATE TABLE ROLES (ID SERIAL PRIMARY KEY, NAME VARCHAR(32) NOT NULL, NOTE VARCHAR(512));

create table users (
    id uuid primary key,
    name varchar(255) not null,
    blocked boolean default false,
    role_id smallint not null references roles(id) on delete cascade,
    created_at timestamp default now(),
    updated_at timestamp default now()
);

CREATE TABLE user_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(512) UNIQUE NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
);

CREATE TABLE folders (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id UUID REFERENCES folders(id) ON DELETE SET NULL,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    main_folder_id UUID REFERENCES folders(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE files (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    folder_id UUID NOT NULL REFERENCES folders(id) ON DELETE CASCADE,
    size BIGINT,
    mime_type VARCHAR(128),
    storage_key VARCHAR(512),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE folder_access (
    folder_id UUID NOT NULL REFERENCES folders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id SMALLINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (folder_id, user_id)
);

CREATE VIEW MAIN_USER_FOLDER AS SELECT FA.USER_ID, F.ID AS FOLDER_ID, F.NAME FROM FOLDERS F JOIN FOLDER_ACCESS FA ON F.ID = FA.FOLDER_ID WHERE F.MAIN_FOLDER_ID IS NULL;
