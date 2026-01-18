
export const API_CONFIG = {
    PRODUCT_SERVICE: import.meta.env.VITE_PRODUCT_SERVICE_URL || 'http://localhost:8001',
    PAYMENT_SERVICE: import.meta.env.VITE_PAYMENT_SERVICE_URL || 'http://localhost:8003',
    CART_SERVICE: import.meta.env.VITE_CART_SERVICE_URL || 'http://localhost:8002',
    // Auth endpoints use cart-order-service but at the API root, not under /api/carts
    AUTH_SERVICE: import.meta.env.VITE_AUTH_SERVICE_URL || 'http://localhost:8002',
};
