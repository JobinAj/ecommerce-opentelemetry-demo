
import { API_CONFIG } from './config';
import { PaymentDetails } from '../types';

export interface PaymentResponse {
    success: boolean;
    message: string;
    payment?: any;
}


export interface CartResponse {
    id: string;
    userId: string;
    items: any[];
    total: number;
}

export interface OrderResponse {
    id: string;
    userId: string;
    items: any[];
    total: number;
    status: string;
}

export const createCart = async (userId: string = 'user123'): Promise<CartResponse> => {
    const response = await fetch(`${API_CONFIG.CART_SERVICE}/api/carts`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ userId }),
    });
    if (!response.ok) throw new Error('Failed to create cart');
    return response.json();
};

export const addItemToCart = async (cartId: string, item: any): Promise<any> => {
    const url = `${API_CONFIG.CART_SERVICE}/api/carts/${cartId}/items`;
    console.log(`Adding item to cart: ${url}`, item);

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                productId: item.product.id,
                productName: item.product.name,
                price: item.product.price,
                quantity: item.quantity,
                selectedSize: item.selectedSize,
                selectedColor: item.selectedColor,
            }),
        });

        if (!response.ok) {
            const text = await response.text();
            console.error(`Failed to add item to cart. Status: ${response.status}, URL: ${url}, Response: ${text}`);
            throw new Error(`Failed to add item to cart: ${response.status} ${text}`);
        }
        return response.json();
    } catch (error) {
        console.error("Network error in addItemToCart:", error);
        throw error;
    }
};

export const createOrder = async (cartId: string): Promise<OrderResponse> => {
    const response = await fetch(`${API_CONFIG.CART_SERVICE}/api/carts/${cartId}/orders`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
    });
    if (!response.ok) throw new Error('Failed to create order');
    return response.json();
};

export const processPayment = async (details: PaymentDetails, amount: number, orderId: string): Promise<PaymentResponse> => {
    const response = await fetch(`${API_CONFIG.PAYMENT_SERVICE}/api/payments`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            orderId,
            amount,
            currency: 'USD',
            cardNumber: details.cardNumber,
            cardHolder: details.cardHolder,
            expiryDate: details.expiryDate,
            cvv: details.cvv,
        }),
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Payment failed');
    }

    return response.json();
};

export const signup = async (name: string, email: string, password: string): Promise<any> => {
    const response = await fetch(`${API_CONFIG.AUTH_SERVICE}/api/signup`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, email, password }),
    });
    if (!response.ok) {
        let errorMessage = 'Signup failed';
        const text = await response.text();
        try {
            const errorData = JSON.parse(text);
            errorMessage = errorData.error || errorData.message || errorMessage;
        } catch {
            errorMessage = text || errorMessage;
        }
        throw new Error(errorMessage);
    }
    return response.json();
};

export const login = async (email: string, password: string): Promise<any> => {
    const response = await fetch(`${API_CONFIG.AUTH_SERVICE}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
    });
    if (!response.ok) {
        let errorMessage = 'Login failed';
        const text = await response.text();
        try {
            const errorData = JSON.parse(text);
            errorMessage = errorData.error || errorData.message || errorMessage;
        } catch {
            errorMessage = text || errorMessage;
        }
        throw new Error(errorMessage);
    }
    return response.json();
};

