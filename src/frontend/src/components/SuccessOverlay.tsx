
interface SuccessOverlayProps {
    isOpen: boolean;
    onClose: () => void;
}

export const SuccessOverlay = ({ isOpen, onClose }: SuccessOverlayProps) => {
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-[60] flex items-center justify-center animate-fade-in">
            <div className="absolute inset-0 bg-white/90 backdrop-blur-md" />

            <div className="relative bg-white p-12 rounded-2xl shadow-2xl max-w-lg w-full text-center space-y-8 animate-scale-up border border-gray-100">
                <div className="flex justify-center">
                    <div className="relative h-24 w-24 flex items-center justify-center">
                        {/* Pulsing background circle */}
                        <div className="absolute inset-0 bg-green-100 rounded-full animate-ping opacity-25" />
                        <div className="absolute inset-0 bg-green-50 rounded-full" />

                        <svg
                            className="w-16 h-16 text-green-500 relative z-10"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            strokeWidth="3"
                        >
                            <path
                                className="animate-draw"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M5 13l4 4L19 7"
                            />
                        </svg>
                    </div>
                </div>

                <div className="space-y-4">
                    <h2 className="text-4xl font-serif font-bold text-gray-900 tracking-tight">
                        Order Accomplished
                    </h2>
                    <p className="text-gray-500 text-lg">
                        Thank you for your purchase. We've received your order and are preparing it for delivery.
                    </p>
                </div>

                <div className="pt-6">
                    <button
                        onClick={onClose}
                        className="w-full bg-black text-white py-4 rounded-full font-bold text-lg hover:bg-gray-800 transition-all duration-300 transform hover:scale-[1.02] active:scale-95 shadow-lg"
                    >
                        Continue Shopping
                    </button>
                </div>

                <p className="text-sm text-gray-400 italic">
                    A confirmation email has been sent to your inbox.
                </p>
            </div>
        </div>
    );
};
