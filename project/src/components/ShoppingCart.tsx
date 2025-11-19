import { X, Trash2, Minus, Plus } from 'lucide-react';
import { CartItem } from '../types';

interface ShoppingCartProps {
  cartItems: CartItem[];
  isOpen: boolean;
  onClose: () => void;
  onUpdateQuantity: (productId: string, size: string, color: string, quantity: number) => void;
  onRemoveItem: (productId: string, size: string, color: string) => void;
  onCheckout: () => void;
}

export const ShoppingCart = ({
  cartItems,
  isOpen,
  onClose,
  onUpdateQuantity,
  onRemoveItem,
  onCheckout,
}: ShoppingCartProps) => {
  const totalPrice = cartItems.reduce(
    (total, item) => total + item.product.price * item.quantity,
    0
  );

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-40">
      <div className="absolute right-0 top-0 h-full w-full max-w-md bg-white shadow-xl flex flex-col">
        <div className="sticky top-0 bg-white border-b flex justify-between items-center p-6">
          <h2 className="text-2xl font-bold">Shopping Cart</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X size={24} />
          </button>
        </div>

        {cartItems.length === 0 ? (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center">
              <p className="text-xl text-gray-500 mb-4">Your cart is empty</p>
              <button
                onClick={onClose}
                className="text-black hover:text-yellow-600 transition-colors font-semibold"
              >
                Continue Shopping
              </button>
            </div>
          </div>
        ) : (
          <>
            <div className="flex-1 overflow-y-auto p-6 space-y-6">
              {cartItems.map((item) => (
                <div key={`${item.product.id}-${item.selectedSize}-${item.selectedColor}`} className="border-b pb-6">
                  <div className="flex gap-4">
                    <img
                      src={item.product.image}
                      alt={item.product.name}
                      className="w-20 h-20 object-cover rounded-lg"
                    />
                    <div className="flex-1">
                      <h3 className="font-bold text-gray-900">{item.product.name}</h3>
                      <p className="text-sm text-gray-600">
                        {item.selectedColor} - {item.selectedSize}
                      </p>
                      <p className="font-semibold text-gray-900 mt-2">
                        ${item.product.price.toLocaleString()}
                      </p>
                    </div>
                    <button
                      onClick={() =>
                        onRemoveItem(
                          item.product.id,
                          item.selectedSize,
                          item.selectedColor
                        )
                      }
                      className="text-red-600 hover:bg-red-50 p-1 rounded transition-colors"
                    >
                      <Trash2 size={18} />
                    </button>
                  </div>

                  <div className="flex items-center gap-3 mt-4">
                    <button
                      onClick={() =>
                        onUpdateQuantity(
                          item.product.id,
                          item.selectedSize,
                          item.selectedColor,
                          item.quantity - 1
                        )
                      }
                      className="p-1 border border-gray-300 rounded hover:bg-gray-100 transition-colors"
                    >
                      <Minus size={16} />
                    </button>
                    <span className="font-bold text-center w-8">{item.quantity}</span>
                    <button
                      onClick={() =>
                        onUpdateQuantity(
                          item.product.id,
                          item.selectedSize,
                          item.selectedColor,
                          item.quantity + 1
                        )
                      }
                      className="p-1 border border-gray-300 rounded hover:bg-gray-100 transition-colors"
                    >
                      <Plus size={16} />
                    </button>
                    <span className="ml-auto font-bold">
                      ${(item.product.price * item.quantity).toLocaleString()}
                    </span>
                  </div>
                </div>
              ))}
            </div>

            <div className="border-t p-6 space-y-4">
              <div className="flex justify-between items-center text-lg font-bold">
                <span>Total:</span>
                <span>${totalPrice.toLocaleString()}</span>
              </div>
              <button
                onClick={onCheckout}
                className="w-full bg-black text-white py-3 rounded-lg font-bold hover:bg-yellow-600 transition-colors"
              >
                Proceed to Checkout
              </button>
              <button
                onClick={onClose}
                className="w-full border border-gray-300 text-gray-900 py-3 rounded-lg font-bold hover:bg-gray-50 transition-colors"
              >
                Continue Shopping
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
};
