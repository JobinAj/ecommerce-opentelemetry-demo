

export const Footer = () => {
  return (
    <footer className="bg-white text-versace-black py-16 px-4 sm:px-6 lg:px-8 border-t border-gray-100">
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-12 mb-16">
          <div className="col-span-1 md:col-span-1">
            <h3 className="text-xl font-serif font-bold mb-6 tracking-widest">VERSACE</h3>
            <div className="flex gap-4">
              {/* Social placeholders */}
            </div>
          </div>

          <div>
            <h4 className="font-bold text-xs mb-6 uppercase tracking-widest">Client Services</h4>
            <ul className="space-y-4 text-sm font-light text-gray-600">
              <li><a href="#" className="hover:text-versace-gold transition-colors">Contact Us</a></li>
              <li><a href="#" className="hover:text-versace-gold transition-colors">Shipping & Returns</a></li>
              <li><a href="#" className="hover:text-versace-gold transition-colors">Track your Order</a></li>
              <li><a href="#" className="hover:text-versace-gold transition-colors">FAQ</a></li>
            </ul>
          </div>

          <div>
            <h4 className="font-bold text-xs mb-6 uppercase tracking-widest">The Company</h4>
            <ul className="space-y-4 text-sm font-light text-gray-600">
              <li><a href="#" className="hover:text-versace-gold transition-colors">About Versace</a></li>
              <li><a href="#" className="hover:text-versace-gold transition-colors">Careers</a></li>
              <li><a href="#" className="hover:text-versace-gold transition-colors">Sustainability</a></li>
              <li><a href="#" className="hover:text-versace-gold transition-colors">Privacy Policy</a></li>
            </ul>
          </div>

          <div>
            <h4 className="font-bold text-xs mb-6 uppercase tracking-widest">Store Locator</h4>
            <p className="text-sm font-light text-gray-600 mb-4">
              Find the nearest boutique and experience the collection.
            </p>
            <a href="#" className="text-xs font-bold border-b border-black pb-1 hover:text-versace-gold hover:border-versace-gold transition-colors">
              FIND A STORE
            </a>
          </div>
        </div>

        <div className="border-t border-gray-100 pt-8 flex flex-col md:flex-row justify-between items-center gap-4 text-xs font-light text-gray-500">
          <p>&copy; 2026 VERSACE. ALL RIGHTS RESERVED.</p>
          <div className="flex gap-6">
            <a href="#" className="hover:text-black transition-colors">Terms & Conditions</a>
            <a href="#" className="hover:text-black transition-colors">Privacy Policy</a>
            <a href="#" className="hover:text-black transition-colors">Cookie Policy</a>
          </div>
        </div>
      </div>
    </footer>
  );
};
