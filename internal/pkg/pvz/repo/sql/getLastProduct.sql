SELECT p.id, p.reception_time, p.category, p.reception_id
        FROM product p
        JOIN reception r ON p.reception_id = r.id
        WHERE r.pvz_id = $1 AND r.status = 'in_progress'
        ORDER BY p.reception_time DESC
        LIMIT 1