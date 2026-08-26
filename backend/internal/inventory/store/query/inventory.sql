-- name: GetLowStock :many
SELECT i.product_id, p.name AS product_name, i.quantity, i.restock_threshold 
FROM inventory i
JOIN products p ON p.id = i.product_id
WHERE i.quantity < i.restock_threshold
ORDER BY quantity ASC;

