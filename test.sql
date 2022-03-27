-- POST data
-- [
--   {
--     "id": 4,
--     "customer_id": 1,
--     "purchase_orders": [{
--         "id": 4112347,
--         "order_id": 1,
--         "sells": [{
--             "purchase_order_id": 4112347,
--             "product_id": 14520,
--             "amount": 2
--         }],
--         "status": [{
--             "purchase_order_id": 4112347,
--             "user_id": 14510,
--             "description": "pending"
--         }]
--     }], 
--     "center_id": 114
--   }
-- ]
UPDATE users
SET active = True
WHERE id = 14510;

INSERT INTO customers
VALUES(1, current_timestamp, current_timestamp, null, 'Lima', 'lima@email.com');

INSERT INTO centers
VALUES(114, '', '');

INSERT INTO products
VALUES(14520, '', '', 125.0);