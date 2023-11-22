-- +migrate Up
-- +migrate StatementBegin

CREATE TYPE task_status AS ENUM ('New', 'In Progress', 'Pending', 'Done');

CREATE SEQUENCE IF NOT EXISTS role_pkey_seq;
CREATE TABLE "role"
(
    id             BIGINT NOT NULL             DEFAULT nextval('role_pkey_seq'::regclass),
    role_name      VARCHAR(50) NOT NULL,
    description    VARCHAR(255),
    permission     VARCHAR(255),
    created_by     BIGINT,
    created_at     TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by     BIGINT,
    updated_at     TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted        BOOLEAN DEFAULT FALSE,
    CONSTRAINT pk_role_id PRIMARY KEY (id),
    CONSTRAINT uq_role_role_name UNIQUE (role_name)
);

INSERT INTO "role"(role_name, description, permission, created_by, updated_by) VALUES
('Admin', 'All Access for All Menu', '{"user":["insert","update","delete","view"],"task":["insert","update","delete","view"]}', 1, 1),
('Owner', 'Own Access for Task Menu', '{"task":["insert","update-own","delete-own","view-own"]}', 1, 1);

CREATE SEQUENCE IF NOT EXISTS user_pkey_seq;
CREATE TABLE "user"
(
    id             BIGINT NOT NULL             DEFAULT nextval('user_pkey_seq'::regclass),
    first_name     VARCHAR(50) NOT NULL,
    last_name      VARCHAR(100),
    username       VARCHAR(20) NOT NULL,
    password       VARCHAR(20) NOT NULL,
    email          VARCHAR(255) NOT NULL,
    role_id        BIGINT,
    created_by     BIGINT,
    created_at     TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by     BIGINT,
    updated_at     TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted        BOOLEAN DEFAULT FALSE,
    CONSTRAINT pk_user_id PRIMARY KEY (id),
    CONSTRAINT fk_user_role_id_role_id FOREIGN KEY (role_id) REFERENCES role(id),
    CONSTRAINT uq_user_name UNIQUE (first_name, last_name)
);

INSERT INTO "user"(first_name, username, password, email) VALUES
('System', 'System', 'System', 'System');

CREATE SEQUENCE IF NOT EXISTS task_pkey_seq;
CREATE TABLE "task"
(
    id             BIGINT NOT NULL             DEFAULT nextval('task_pkey_seq'::regclass),
    title          VARCHAR(255) NOT NULL,
    description    TEXT,
    due_date       TIMESTAMP WITHOUT TIME ZONE,
    status         task_status DEFAULT 'New',
    user_id        BIGINT,
    created_by     BIGINT,
    created_at     TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by     BIGINT,
    updated_at     TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted        BOOLEAN DEFAULT FALSE,
    CONSTRAINT pk_task_id PRIMARY KEY (id),
    CONSTRAINT fk_task_user_id_user_id FOREIGN KEY (user_id) REFERENCES "user"(id)
);

-- +migrate StatementEnd