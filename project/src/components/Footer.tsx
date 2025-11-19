import { Mail, Phone, MapPin } from 'lucide-react';

export const Footer = () => {
  return (
    <footer className="bg-black text-white py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
          <div>
            <h3 className="text-2xl font-bold mb-4 tracking-widest">VERSACE</h3>
            <p className="text-gray-400">
              Experience timeless luxury and modern elegance in every collection.
            </p>
          </div>

          <div>
            <h4 className="font-bold text-lg mb-4 uppercase tracking-wide">Shop</h4>
            <ul className="space-y-2 text-gray-400">
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  New Arrivals
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  Collections
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  Men
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  Women
                </a>
              </li>
            </ul>
          </div>

          <div>
            <h4 className="font-bold text-lg mb-4 uppercase tracking-wide">Help</h4>
            <ul className="space-y-2 text-gray-400">
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  Contact Us
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  Shipping
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  Returns
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-white transition-colors">
                  FAQs
                </a>
              </li>
            </ul>
          </div>

          <div>
            <h4 className="font-bold text-lg mb-4 uppercase tracking-wide">Contact</h4>
            <ul className="space-y-3 text-gray-400">
              <li className="flex items-center gap-2">
                <Mail size={18} />
                <span>support@versace.com</span>
              </li>
              <li className="flex items-center gap-2">
                <Phone size={18} />
                <span>+1 (800) 800-8000</span>
              </li>
              <li className="flex items-start gap-2">
                <MapPin size={18} className="flex-shrink-0 mt-1" />
                <span>Milan, Italy</span>
              </li>
            </ul>
          </div>
        </div>

        <div className="border-t border-gray-800 pt-8">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-gray-400 text-sm">
            <p>&copy; 2024 Versace. All rights reserved.</p>
            <div className="flex gap-6">
              <a href="#" className="hover:text-white transition-colors">
                Privacy Policy
              </a>
              <a href="#" className="hover:text-white transition-colors">
                Terms of Service
              </a>
              <a href="#" className="hover:text-white transition-colors">
                Cookie Settings
              </a>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
};
