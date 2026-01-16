import { useState } from 'react';
import { Header } from './Header';
import { Hero } from './Hero';
import { ProductCatalog } from './ProductCatalog';
import { ProductDetail } from './ProductDetail';
import { ShoppingCart } from './ShoppingCart';
import { CheckoutModal } from './CheckoutModal';
import { Footer } from './Footer';
import { SuccessOverlay } from './SuccessOverlay';
import { useCart } from '../hooks/useCart';
import { products } from '../data/products';
import { Product, PaymentDetails } from '../types';
import { processPayment, createCart, addItemToCart, createOrder } from '../api/client';

export function Home() {
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
    const [cartOpen, setCartOpen] = useState(false);
    const [checkoutOpen, setCheckoutOpen] = useState(false);
    const [orderComplete, setOrderComplete] = useState(false);

    const {
        cartItems,
        addToCart,
        removeFromCart,
        updateQuantity,
        clearCart,
        getTotalPrice,
        getTotalItems,
    } = useCart();

    const handleAddToCart = (product: Product) => {
        addToCart(product, 1, product.sizes[0], product.colors[0]);
    };

    const handleProductDetailAddToCart = (
        product: Product,
        quantity: number,
        size: string,
        color: string
    ) => {
        addToCart(product, quantity, size, color);
        setSelectedProduct(null);
        setCartOpen(true);
    };

    const handleCheckout = async (paymentDetails: PaymentDetails) => {
        try {
            // 1. Create a cart in the backend
            const storedUser = localStorage.getItem('user');
            const user = storedUser ? JSON.parse(storedUser) : null;
            const userId = user ? user.id : 'guest';

            const cart = await createCart(userId);
            const cartId = cart.id;

            // 2. Add each item to the backend cart
            for (const item of cartItems) {
                await addItemToCart(cartId, item);
            }

            // 3. Create an order from the cart
            const order = await createOrder(cartId);
            const orderId = order.id;

            // 4. Process payment for the newly created order
            await processPayment(paymentDetails, order.total, orderId);

            console.log('Payment successful');
            setCheckoutOpen(false);
            setOrderComplete(true);
            clearCart();

            // We don't need setTimeout to hide it automatically, users can click "Continue Shopping"
        } catch (error: any) {
            console.error('Checkout failed:', error);
            alert(`Checkout failed: ${error.message}`);
            throw error;
        }
    };

    return (
        <div className="min-h-screen bg-white">
            <Header cartItemsCount={getTotalItems()} onCartClick={() => setCartOpen(true)} />

            <Hero onShopClick={() => window.scrollTo({ top: 500, behavior: 'smooth' })} />
            <ProductCatalog
                products={products}
                onViewDetails={setSelectedProduct}
                onAddToCart={handleAddToCart}
            />
            <Footer />

            {selectedProduct && (
                <ProductDetail
                    product={selectedProduct}
                    onClose={() => setSelectedProduct(null)}
                    onAddToCart={handleProductDetailAddToCart}
                />
            )}

            <ShoppingCart
                cartItems={cartItems}
                isOpen={cartOpen}
                onClose={() => setCartOpen(false)}
                onUpdateQuantity={updateQuantity}
                onRemoveItem={removeFromCart}
                onCheckout={() => {
                    setCartOpen(false);
                    setCheckoutOpen(true);
                }}
            />

            <CheckoutModal
                isOpen={checkoutOpen}
                cartItems={cartItems}
                totalPrice={getTotalPrice()}
                onClose={() => setCheckoutOpen(false)}
                onSubmit={handleCheckout}
            />

            <SuccessOverlay
                isOpen={orderComplete}
                onClose={() => setOrderComplete(false)}
            />
        </div>
    );
}
