SELECT 
    pvz.id, pvz.city, pvz.registration_date,
    reception.id, reception.reception_time, reception.status,
    product.id, product.reception_time, product.category
FROM pvz
LEFT JOIN reception 
    ON pvz.id = reception.pvz_id 
    AND ($1::timestamptz IS NULL OR reception.reception_time >= $1)
    AND ($2::timestamptz IS NULL OR reception.reception_time <= $2)
LEFT JOIN product 
    ON product.reception_id = reception.id
ORDER BY pvz.id, reception.reception_time NULLS LAST
LIMIT $3 OFFSET $4;
