import { useState } from 'react';
import { Header } from './components/Header';
import { Hero } from './components/Hero';
import { ProductCatalog } from './components/ProductCatalog';
import { ProductDetail } from './components/ProductDetail';
import { ShoppingCart } from './components/ShoppingCart';
import { CheckoutModal } from './components/CheckoutModal';
import { Footer } from './components/Footer';
import { useCart } from './hooks/useCart';
import { products } from './data/products';
import { Product, PaymentDetails } from './types';

function App() {
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

  const handleCheckout = (paymentDetails: PaymentDetails) => {
    console.log('Processing payment with details:', paymentDetails);
    console.log('Cart items:', cartItems);
    setCheckoutOpen(false);
    setOrderComplete(true);
    clearCart();

    setTimeout(() => {
      setOrderComplete(false);
    }, 3000);
  };

  return (
    <div className="min-h-screen bg-white">
      <Header cartItemsCount={getTotalItems()} onCartClick={() => setCartOpen(true)} />

      {orderComplete && (
        <div className="fixed top-4 right-4 bg-green-500 text-white px-6 py-4 rounded-lg shadow-lg z-40 animate-pulse">
          <p className="font-bold">Order placed successfully!</p>
          <p className="text-sm">Thank you for your purchase.</p>
        </div>
      )}

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
    </div>
  );
}

export default App;
