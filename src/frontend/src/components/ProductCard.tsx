import { Product } from '../types';

interface ProductCardProps {
  product: Product;
  onViewDetails: (product: Product) => void;
}

export const ProductCard = ({
  product,
  onViewDetails,
}: ProductCardProps) => {
  return (
    <div className="group cursor-pointer" onClick={() => onViewDetails(product)}>
      <div className="relative w-full aspect-[3/4] overflow-hidden bg-gray-100">
        <img
          src={product.image}
          alt={product.name}
          className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
        />
        {!product.inStock && (
          <div className="absolute inset-0 bg-black/40 flex items-center justify-center">
            <span className="text-white text-sm font-bold tracking-widest border border-white px-4 py-2">
              SOLD OUT
            </span>
          </div>
        )}
      </div>

      <div className="pt-4 text-center">
        <h3 className="text-sm font-bold tracking-wide text-gray-900 group-hover:text-versace-gold transition-colors">
          {product.name}
        </h3>

        <p className="text-xs text-gray-500 mt-1 mb-2 capitalize">
          {product.category}
        </p>

        <span className="text-sm font-normal text-gray-900">
          ${product.price.toLocaleString()}
        </span>
      </div>
    </div>
  );
};
