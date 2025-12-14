import { X } from 'lucide-react';
import { useState } from 'react';
import { CartItem, PaymentDetails } from '../types';

interface CheckoutModalProps {
  isOpen: boolean;
  cartItems: CartItem[];
  totalPrice: number;
  onClose: () => void;
  onSubmit: (paymentDetails: PaymentDetails) => void;
}

export const CheckoutModal = ({
  isOpen,
  cartItems,
  totalPrice,
  onClose,
  onSubmit,
}: CheckoutModalProps) => {
  const [formData, setFormData] = useState<PaymentDetails>({
    cardNumber: '',
    cardHolder: '',
    expiryDate: '',
    cvv: '',
  });

  const [email, setEmail] = useState('');
  const [address, setAddress] = useState('');
  const [processing, setProcessing] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setProcessing(true);

    try {
      await new Promise((resolve) => setTimeout(resolve, 1500));
      onSubmit(formData);
      setFormData({ cardNumber: '', cardHolder: '', expiryDate: '', cvv: '' });
      setEmail('');
      setAddress('');
    } finally {
      setProcessing(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="sticky top-0 bg-white border-b flex justify-between items-center p-6">
          <h2 className="text-2xl font-bold">Secure Checkout</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X size={24} />
          </button>
        </div>

        <div className="p-6 space-y-8">
          <div>
            <h3 className="text-lg font-bold mb-4">Order Summary</h3>
            <div className="space-y-3 max-h-40 overflow-y-auto">
              {cartItems.map((item) => (
                <div
                  key={`${item.product.id}-${item.selectedSize}-${item.selectedColor}`}
                  className="flex justify-between items-center text-sm"
                >
                  <span>
                    {item.product.name} ({item.selectedColor}, {item.selectedSize})
                  </span>
                  <span className="font-bold">
                    ${(item.product.price * item.quantity).toLocaleString()}
                  </span>
                </div>
              ))}
            </div>
            <div className="border-t mt-4 pt-4 flex justify-between items-center text-lg font-bold">
              <span>Total:</span>
              <span>${totalPrice.toLocaleString()}</span>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <h3 className="text-lg font-bold mb-4">Shipping Address</h3>
              <div className="space-y-3">
                <input
                  type="email"
                  placeholder="Email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black"
                />
                <textarea
                  placeholder="Shipping Address"
                  value={address}
                  onChange={(e) => setAddress(e.target.value)}
                  required
                  rows={3}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black resize-none"
                />
              </div>
            </div>

            <div>
              <h3 className="text-lg font-bold mb-4">Payment Information</h3>
              <div className="space-y-3">
                <input
                  type="text"
                  name="cardNumber"
                  placeholder="Card Number"
                  value={formData.cardNumber}
                  onChange={handleChange}
                  required
                  maxLength={19}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black"
                />

                <input
                  type="text"
                  name="cardHolder"
                  placeholder="Cardholder Name"
                  value={formData.cardHolder}
                  onChange={handleChange}
                  required
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black"
                />

                <div className="grid grid-cols-2 gap-3">
                  <input
                    type="text"
                    name="expiryDate"
                    placeholder="MM/YY"
                    value={formData.expiryDate}
                    onChange={handleChange}
                    required
                    maxLength={5}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black"
                  />
                  <input
                    type="text"
                    name="cvv"
                    placeholder="CVV"
                    value={formData.cvv}
                    onChange={handleChange}
                    required
                    maxLength={4}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black"
                  />
                </div>
              </div>
            </div>

            <button
              type="submit"
              disabled={processing}
              className="w-full bg-black text-white py-3 rounded-lg font-bold text-lg hover:bg-yellow-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {processing ? 'Processing...' : 'Complete Purchase'}
            </button>
          </form>

          <p className="text-xs text-gray-600 text-center">
            Your payment information is secure and encrypted. We never store your full card details.
          </p>
        </div>
      </div>
    </div>
  );
};
