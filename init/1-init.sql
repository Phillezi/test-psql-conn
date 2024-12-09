-- mock random database
-- Users table
CREATE TABLE
    users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) NOT NULL UNIQUE,
        email VARCHAR(100) NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Insert sample data into users
INSERT INTO
    users (username, email)
VALUES
    ('john_doe', 'john@example.com'),
    ('jane_doe', 'jane@example.com'),
    ('sam_smith', 'sam@example.com');

-- Products table
CREATE TABLE
    products (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        description TEXT,
        price DECIMAL(10, 2) NOT NULL,
        stock INT NOT NULL CHECK (stock >= 0),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Insert sample data into products
INSERT INTO
    products (name, description, price, stock)
VALUES
    ('Laptop', 'High-performance laptop', 1200.00, 10),
    (
        'Smartphone',
        'Latest model smartphone',
        800.00,
        25
    ),
    (
        'Headphones',
        'Noise-canceling headphones',
        150.00,
        50
    );

-- Orders table
CREATE TABLE
    orders (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Insert sample data into orders
INSERT INTO
    orders (user_id)
VALUES
    (1),
    (2),
    (1);

-- Order Items table
CREATE TABLE
    order_items (
        id SERIAL PRIMARY KEY,
        order_id INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
        product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
        quantity INT NOT NULL CHECK (quantity > 0),
        price DECIMAL(10, 2) NOT NULL
    );

-- Insert sample data into order_items
INSERT INTO
    order_items (order_id, product_id, quantity, price)
VALUES
    (1, 1, 1, 1200.00),
    (1, 2, 2, 800.00),
    (2, 3, 1, 150.00),
    (3, 1, 1, 1200.00);