-- Задание 2: Наполнение orders на основе существующих products и customers
INSERT INTO orders (id, create_time, status, price, customer_id, product_id, product_amount)
SELECT
    n.number + 1 AS id,
    now() - (rand() % 2592000) AS create_time,
    arrayElement(['created','paid','shipped','delivered','cancelled'], (n.number % 5) + 1) AS status,
    p.price * ((rand() % 5) + 1) AS price,
    c.id AS customer_id,
    p.id AS product_id,
    (rand() % 5) + 1 AS product_amount
FROM numbers(1000) AS n
JOIN (SELECT id FROM customers ORDER BY rand()) AS c ON n.number % (SELECT count() FROM customers) = c.id % (SELECT count() FROM customers)
JOIN (SELECT id, price FROM products ORDER BY rand()) AS p ON n.number % (SELECT count() FROM products) = p.id % (SELECT count() FROM products);