CREATE INDEX idx_receptions_pvz_id ON receptions(pvz_id);
CREATE INDEX idx_products_reception_id ON products(reception_id);
CREATE INDEX idx_pvz_city ON pvz(city);
CREATE INDEX idx_receptions_created_at ON receptions(created_at);
CREATE INDEX idx_products_reception_id_created ON products(reception_id, created_at);

CREATE UNIQUE INDEX idx_unique_active_reception
    ON receptions (pvz_id)
    WHERE status = 'in_progress';