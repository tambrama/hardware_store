-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS images (
    image_id UUID PRIMARY KEY,
    image BYTEA NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS address (
    address_id UUID PRIMARY KEY,
    country TEXT NOT NULL,
    city TEXT,
    street TEXT
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS category (
    category_id UUID PRIMARY KEY,
    category TEXT NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS client (
    client_id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    birthday DATE CHECK (
        birthday <= CURRENT_DATE
        AND birthday >= '1900-01-01'
    ),
    gender TEXT CHECK (gender IN ('male', 'female')),
    registration_date TIMESTAMPTZ DEFAULT NOW(),
    address_id UUID,
    FOREIGN KEY (address_id) REFERENCES address(address_id) ON DELETE RESTRICT -- Запрещает удаление, если есть ссылки
    ON UPDATE CASCADE
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS supplier (
    supplier_id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    address_id UUID,
    phone_number VARCHAR(20) CHECK (phone_number ~ '^\+?[0-9\s\-\(\)]+$'),
    FOREIGN KEY (address_id) REFERENCES address(address_id) ON DELETE RESTRICT ON UPDATE CASCADE
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product (
    product_id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    category_id UUID,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    available_stock INTEGER NOT NULL DEFAULT 0 CHECK (available_stock >= 0),
    last_update_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
    CHECK (last_update_date >= '2000-01-01'),
    supplier_id UUID,
    image_id UUID,
    FOREIGN KEY (supplier_id) REFERENCES supplier(supplier_id) ON DELETE
    SET NULL ON UPDATE CASCADE,
        FOREIGN KEY (image_id) REFERENCES images(image_id) ON DELETE
    SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (category_id) REFERENCES category(category_id) ON DELETE
    SET NULL ON UPDATE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS supplier;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS client;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS category;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS address;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd