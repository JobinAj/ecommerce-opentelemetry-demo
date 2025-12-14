import { ShoppingBag, Menu, X } from 'lucide-react';
import { useState } from 'react';

interface HeaderProps {
  cartItemsCount: number;
  onCartClick: () => void;
}

export const Header = ({ cartItemsCount, onCartClick }: HeaderProps) => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <header className="bg-black text-white sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-20">
          <div className="flex items-center gap-2">
            <div className="text-3xl font-bold tracking-widest">VERSACE</div>
          </div>

          <nav className="hidden md:flex items-center gap-8">
            <a href="#" className="hover:text-gold transition-colors text-sm font-semibold">
              SHOP
            </a>
            <a href="#" className="hover:text-gold transition-colors text-sm font-semibold">
              COLLECTION
            </a>
            <a href="#" className="hover:text-gold transition-colors text-sm font-semibold">
              ABOUT
            </a>
            <a href="#" className="hover:text-gold transition-colors text-sm font-semibold">
              CONTACT
            </a>
          </nav>

          <div className="flex items-center gap-6">
            <button
              onClick={onCartClick}
              className="relative p-2 hover:bg-gray-900 rounded-lg transition-colors"
            >
              <ShoppingBag size={24} />
              {cartItemsCount > 0 && (
                <span className="absolute top-1 right-1 bg-yellow-500 text-black text-xs font-bold rounded-full w-5 h-5 flex items-center justify-center">
                  {cartItemsCount}
                </span>
              )}
            </button>

            <button
              className="md:hidden p-2"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
            </button>
          </div>
        </div>

        {mobileMenuOpen && (
          <nav className="md:hidden pb-4 space-y-3">
            <a href="#" className="block hover:text-gold transition-colors text-sm font-semibold">
              SHOP
            </a>
            <a href="#" className="block hover:text-gold transition-colors text-sm font-semibold">
              COLLECTION
            </a>
            <a href="#" className="block hover:text-gold transition-colors text-sm font-semibold">
              ABOUT
            </a>
            <a href="#" className="block hover:text-gold transition-colors text-sm font-semibold">
              CONTACT
            </a>
          </nav>
        )}
      </div>
    </header>
  );
};
