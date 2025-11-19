import { X, Star, ShoppingBag, Minus, Plus } from 'lucide-react';
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
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div className="sticky top-0 bg-white border-b flex justify-between items-center p-6">
          <h2 className="text-2xl font-bold">{product.name}</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X size={24} />
          </button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 p-6">
          <div>
            <img
              src={product.image}
              alt={product.name}
              className="w-full h-full object-cover rounded-lg"
            />
          </div>

          <div className="space-y-6">
            <div>
              <div className="flex items-center gap-2 mb-2">
                {[...Array(5)].map((_, i) => (
                  <Star
                    key={i}
                    size={18}
                    className={i < Math.floor(product.rating) ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300'}
                  />
                ))}
              </div>
              <p className="text-gray-600">
                {product.reviews} reviews
              </p>
            </div>

            <div>
              <p className="text-4xl font-bold text-gray-900">
                ${product.price.toLocaleString()}
              </p>
              <p className={`mt-2 ${product.inStock ? 'text-green-600' : 'text-red-600'} font-semibold`}>
                {product.inStock ? 'In Stock' : 'Out of Stock'}
              </p>
            </div>

            <p className="text-gray-700 leading-relaxed">
              {product.description}
            </p>

            <div>
              <label className="block text-sm font-semibold text-gray-900 mb-3">
                Color
              </label>
              <div className="flex gap-3">
                {product.colors.map((color) => (
                  <button
                    key={color}
                    onClick={() => setSelectedColor(color)}
                    className={`px-4 py-2 rounded-lg border-2 font-medium transition-all ${
                      selectedColor === color
                        ? 'border-black bg-black text-white'
                        : 'border-gray-300 bg-white text-gray-900 hover:border-black'
                    }`}
                  >
                    {color}
                  </button>
                ))}
              </div>
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-900 mb-3">
                Size
              </label>
              <div className="grid grid-cols-4 gap-2">
                {product.sizes.map((size) => (
                  <button
                    key={size}
                    onClick={() => setSelectedSize(size)}
                    className={`py-2 rounded-lg border-2 font-medium transition-all ${
                      selectedSize === size
                        ? 'border-black bg-black text-white'
                        : 'border-gray-300 bg-white text-gray-900 hover:border-black'
                    }`}
                  >
                    {size}
                  </button>
                ))}
              </div>
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-900 mb-3">
                Quantity
              </label>
              <div className="flex items-center gap-4">
                <button
                  onClick={() => setQuantity(Math.max(1, quantity - 1))}
                  className="p-2 border border-gray-300 rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <Minus size={20} />
                </button>
                <span className="text-xl font-bold w-8 text-center">{quantity}</span>
                <button
                  onClick={() => setQuantity(quantity + 1)}
                  className="p-2 border border-gray-300 rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <Plus size={20} />
                </button>
              </div>
            </div>

            <button
              onClick={handleAddToCart}
              disabled={!product.inStock}
              className="w-full bg-black text-white py-4 rounded-lg font-bold text-lg hover:bg-yellow-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              <ShoppingBag size={24} />
              Add to Cart
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
