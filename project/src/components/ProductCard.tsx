import { Star, ShoppingBag } from 'lucide-react';
import { Product } from '../types';

interface ProductCardProps {
  product: Product;
  onViewDetails: (product: Product) => void;
  onAddToCart: (product: Product) => void;
}

export const ProductCard = ({
  product,
  onViewDetails,
  onAddToCart,
}: ProductCardProps) => {
  return (
    <div className="bg-white rounded-lg overflow-hidden shadow-md hover:shadow-xl transition-all duration-300 group cursor-pointer">
      <div
        className="relative w-full h-64 overflow-hidden bg-gray-200 cursor-pointer"
        onClick={() => onViewDetails(product)}
      >
        <img
          src={product.image}
          alt={product.name}
          className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
        />
        {!product.inStock && (
          <div className="absolute inset-0 bg-black bg-opacity-40 flex items-center justify-center">
            <span className="text-white text-lg font-bold">OUT OF STOCK</span>
          </div>
        )}
      </div>

      <div className="p-4">
        <p className="text-xs text-gray-500 font-semibold uppercase tracking-wide mb-2">
          {product.category}
        </p>

        <h3
          className="text-lg font-bold text-gray-900 mb-2 hover:text-yellow-600 transition-colors cursor-pointer line-clamp-2"
          onClick={() => onViewDetails(product)}
        >
          {product.name}
        </h3>

        <div className="flex items-center gap-2 mb-3">
          <div className="flex items-center gap-1">
            {[...Array(5)].map((_, i) => (
              <Star
                key={i}
                size={14}
                className={i < Math.floor(product.rating) ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300'}
              />
            ))}
          </div>
          <span className="text-xs text-gray-600">({product.reviews})</span>
        </div>

        <div className="flex items-center justify-between">
          <span className="text-2xl font-bold text-gray-900">
            ${product.price.toLocaleString()}
          </span>
          <button
            onClick={() => onAddToCart(product)}
            disabled={!product.inStock}
            className="p-2 bg-black text-white rounded-lg hover:bg-yellow-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <ShoppingBag size={20} />
          </button>
        </div>
      </div>
    </div>
  );
};
