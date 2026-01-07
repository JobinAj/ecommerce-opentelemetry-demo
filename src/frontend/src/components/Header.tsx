import { ShoppingBag, Menu, X, Search } from 'lucide-react';
import { useState } from 'react';

interface HeaderProps {
  cartItemsCount: number;
  onCartClick: () => void;
}

export const Header = ({ cartItemsCount, onCartClick }: HeaderProps) => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <header className="bg-white text-versace-black sticky top-0 z-50 border-b border-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-24">

          {/* Mobile Menu Button */}
          <button
            className="md:hidden p-2"
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          >
            {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
          </button>

          {/* Left Navigation (Desktop) */}
          <nav className="hidden md:flex items-center gap-8 flex-1">
            <a href="#" className="text-xs font-bold tracking-widest hover:text-versace-gold transition-colors">WOMEN</a>
            <a href="#" className="text-xs font-bold tracking-widest hover:text-versace-gold transition-colors">MEN</a>
            <a href="#" className="text-xs font-bold tracking-widest hover:text-versace-gold transition-colors">CHILDREN</a>
          </nav>

          {/* Logo (Centered) */}
          <div className="flex-shrink-0 flex justify-center flex-1">
            <div className="text-3xl font-serif font-bold tracking-widest">VERSACE</div>
          </div>

          {/* Right Icons */}
          <div className="flex items-center justify-end gap-6 flex-1">
            <button className="hidden md:block hover:text-versace-gold transition-colors">
              <Search size={20} />
            </button>
            <button
              onClick={onCartClick}
              className="relative hover:text-versace-gold transition-colors"
            >
              <ShoppingBag size={20} />
              {cartItemsCount > 0 && (
                <span className="absolute -top-2 -right-2 bg-versace-black text-white text-[10px] font-bold rounded-full w-4 h-4 flex items-center justify-center">
                  {cartItemsCount}
                </span>
              )}
            </button>
          </div>
        </div>

        {/* Mobile Menu */}
        {mobileMenuOpen && (
          <nav className="md:hidden pb-6 space-y-4 pt-2 border-t border-gray-100">
            <a href="#" className="block text-sm font-bold tracking-widest hover:text-versace-gold transition-colors">WOMEN</a>
            <a href="#" className="block text-sm font-bold tracking-widest hover:text-versace-gold transition-colors">MEN</a>
            <a href="#" className="block text-sm font-bold tracking-widest hover:text-versace-gold transition-colors">CHILDREN</a>
            <div className="pt-4 border-t border-gray-100">
              <a href="#" className="block text-sm text-gray-600 hover:text-versace-gold transition-colors">Sign In</a>
            </div>
          </nav>
        )}
      </div>
    </header>
  );
};
