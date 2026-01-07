import { X, Minus, Plus } from 'lucide-react';
import { useState } from 'react';
import { Product } from '../types';

interface ProductDetailProps {
  product: Product;
  onClose: () => void;
  onAddToCart: (product: Product, quantity: number, size: string, color: string) => void;
}

export const ProductDetail = ({
  product,
  onClose,
  onAddToCart,
}: ProductDetailProps) => {
  const [quantity, setQuantity] = useState(1);
  const [selectedSize, setSelectedSize] = useState(product.sizes[0]);
  const [selectedColor, setSelectedColor] = useState(product.colors[0]);

  const handleAddToCart = () => {
    onAddToCart(product, quantity, selectedSize, selectedColor);
    setQuantity(1);
  };

  return (
    <div className="fixed inset-0 bg-black/60 z-50 flex items-center justify-center p-0 md:p-4">
      <div className="bg-white max-w-6xl w-full h-full md:h-[90vh] overflow-y-auto flex flex-col md:flex-row shadow-2xl">

        {/* Close Button (Mobile sticky) */}
        <div className="md:hidden sticky top-0 bg-white p-4 flex justify-end border-b border-gray-100 z-10">
          <button onClick={onClose}>
            <X size={24} />
          </button>
        </div>

        {/* Image Section */}
        <div className="w-full md:w-1/2 bg-gray-50">
          <img
            src={product.image}
            alt={product.name}
            className="w-full h-full object-cover"
          />
        </div>

        {/* Details Section */}
        <div className="w-full md:w-1/2 p-8 md:p-16 flex flex-col relative">

          {/* Close Button (Desktop) */}
          <button
            onClick={onClose}
            className="hidden md:block absolute top-8 right-8 hover:text-versace-gold transition-colors"
          >
            <X size={24} />
          </button>

          <div className="mb-8">
            <h2 className="text-3xl md:text-4xl font-serif font-bold mb-2">{product.name}</h2>
            <p className="text-sm text-gray-500 uppercase tracking-widest mb-6">{product.category}</p>
            <p className="text-2xl font-normal text-gray-900 mb-4">
              ${product.price.toLocaleString()}
            </p>
          </div>

          <p className="text-gray-600 leading-relaxed mb-10 font-light">
            {product.description}
          </p>

          <div className="space-y-8 flex-1">
            {/* Color Selection */}
            <div>
              <label className="block text-xs font-bold uppercase tracking-widest mb-4">
                Color: <span className="text-gray-500 font-normal">{selectedColor}</span>
              </label>
              <div className="flex gap-3">
                {product.colors.map((color) => (
                  <button
                    key={color}
                    onClick={() => setSelectedColor(color)}
                    className={`w-10 h-10 rounded-full border border-gray-200 transition-all ${selectedColor === color ? 'ring-2 ring-offset-2 ring-black' : 'hover:scale-110'
                      }`}
                    style={{ backgroundColor: color.toLowerCase() }}
                    title={color}
                  />
                ))}
              </div>
            </div>

            {/* Size Selection */}
            <div>
              <label className="block text-xs font-bold uppercase tracking-widest mb-4">
                Size: <span className="text-gray-500 font-normal">{selectedSize}</span>
              </label>
              <div className="grid grid-cols-4 gap-2">
                {product.sizes.map((size) => (
                  <button
                    key={size}
                    onClick={() => setSelectedSize(size)}
                    className={`py-3 text-sm font-bold border transition-all ${selectedSize === size
                      ? 'border-black bg-black text-white'
                      : 'border-gray-200 hover:border-black text-gray-900'
                      }`}
                  >
                    {size}
                  </button>
                ))}
              </div>
            </div>

            {/* Quantity */}
            <div>
              <label className="block text-xs font-bold uppercase tracking-widest mb-4">
                Quantity
              </label>
              <div className="flex items-center border border-gray-200 w-32">
                <button
                  onClick={() => setQuantity(Math.max(1, quantity - 1))}
                  className="p-3 hover:bg-gray-50 transition-colors"
                >
                  <Minus size={16} />
                </button>
                <span className="flex-1 text-center font-bold">{quantity}</span>
                <button
                  onClick={() => setQuantity(quantity + 1)}
                  className="p-3 hover:bg-gray-50 transition-colors"
                >
                  <Plus size={16} />
                </button>
              </div>
            </div>
          </div>

          {/* Add to Cart */}
          <button
            onClick={handleAddToCart}
            disabled={!product.inStock}
            className="w-full bg-black text-white py-5 font-bold text-sm tracking-[0.2em] hover:bg-versace-gold transition-colors disabled:opacity-50 disabled:cursor-not-allowed mt-8 uppercase"
          >
            {product.inStock ? 'Add to Shopping Bag' : 'Out of Stock'}
          </button>

          <div className="mt-6 space-y-2 text-xs text-gray-500">
            <p>COMPLIMENTARY SHIPPING AND RETURNS</p>
            <p>SECURE PAYMENTS</p>
          </div>
        </div>
      </div>
    </div>
  );
};
