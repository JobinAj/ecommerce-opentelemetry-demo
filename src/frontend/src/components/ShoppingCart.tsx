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
    <div className="fixed inset-0 bg-black/50 z-40 backdrop-blur-sm">
      <div className="absolute right-0 top-0 h-full w-full max-w-md bg-white shadow-2xl flex flex-col">
        <div className="sticky top-0 bg-white border-b border-gray-100 flex justify-between items-center p-6">
          <h2 className="text-xl font-serif font-bold tracking-wider">YOUR SHOPPING BAG</h2>
          <button
            onClick={onClose}
            className="hover:text-versace-gold transition-colors"
          >
            <X size={24} />
          </button>
        </div>

        {cartItems.length === 0 ? (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center">
              <p className="text-lg text-gray-500 mb-6 font-light">Your bag is empty</p>
              <button
                onClick={onClose}
                className="text-black border-b border-black hover:text-versace-gold hover:border-versace-gold transition-colors font-bold text-sm tracking-widest pb-1"
              >
                CONTINUE SHOPPING
              </button>
            </div>
          </div>
        ) : (
          <>
            <div className="flex-1 overflow-y-auto p-6 space-y-8">
              {cartItems.map((item) => (
                <div key={`${item.product.id}-${item.selectedSize}-${item.selectedColor}`} className="border-b border-gray-100 pb-8">
                  <div className="flex gap-4">
                    <img
                      src={item.product.image}
                      alt={item.product.name}
                      className="w-24 h-32 object-cover"
                    />
                    <div className="flex-1 flex flex-col justify-between">
                      <div>
                        <div className="flex justify-between items-start">
                          <h3 className="font-bold text-sm text-gray-900 tracking-wide">{item.product.name}</h3>
                          <button
                            onClick={() =>
                              onRemoveItem(
                                item.product.id,
                                item.selectedSize,
                                item.selectedColor
                              )
                            }
                            className="text-gray-400 hover:text-red-500 transition-colors"
                          >
                            <Trash2 size={16} />
                          </button>
                        </div>
                        <p className="text-xs text-gray-500 mt-2 uppercase">
                          {item.selectedColor} | Size: {item.selectedSize}
                        </p>
                      </div>

                      <div className="flex justify-between items-end mt-4">
                        <div className="flex items-center border border-gray-200">
                          <button
                            onClick={() =>
                              onUpdateQuantity(
                                item.product.id,
                                item.selectedSize,
                                item.selectedColor,
                                item.quantity - 1
                              )
                            }
                            className="p-1 hover:bg-gray-50 transition-colors"
                          >
                            <Minus size={14} />
                          </button>
                          <span className="font-bold text-xs w-8 text-center">{item.quantity}</span>
                          <button
                            onClick={() =>
                              onUpdateQuantity(
                                item.product.id,
                                item.selectedSize,
                                item.selectedColor,
                                item.quantity + 1
                              )
                            }
                            className="p-1 hover:bg-gray-50 transition-colors"
                          >
                            <Plus size={14} />
                          </button>
                        </div>
                        <p className="font-bold text-sm text-gray-900">
                          ${(item.product.price * item.quantity).toLocaleString()}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>

            <div className="border-t border-gray-100 p-6 space-y-4 bg-gray-50">
              <div className="flex justify-between items-center text-lg font-serif font-bold">
                <span>TOTAL</span>
                <span>${totalPrice.toLocaleString()}</span>
              </div>
              <p className="text-xs text-gray-500 text-center mb-4">
                Shipping and taxes calculated at checkout
              </p>
              <button
                onClick={onCheckout}
                className="w-full bg-black text-white py-4 font-bold text-sm tracking-[0.2em] hover:bg-versace-gold transition-colors"
              >
                CHECKOUT
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
};
