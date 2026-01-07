
export const API_CONFIG = {
    PRODUCT_SERVICE: import.meta.env.VITE_PRODUCT_SERVICE_URL || 'http://localhost:8001',
    PAYMENT_SERVICE: import.meta.env.VITE_PAYMENT_SERVICE_URL || 'http://localhost:8003',
    CART_SERVICE: import.meta.env.VITE_CART_SERVICE_URL || 'http://localhost:8002',
};
