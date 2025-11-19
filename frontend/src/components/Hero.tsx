interface HeroProps {
  onShopClick: () => void;
}

export const Hero = ({ onShopClick }: HeroProps) => {
  return (
    <div className="relative w-full h-96 md:h-[500px] bg-gradient-to-r from-gray-900 via-black to-gray-900 overflow-hidden">
      <div
        className="absolute inset-0 opacity-20"
        style={{
          backgroundImage:
            'url("https://images.pexels.com/photos/3622613/pexels-photo-3622613.jpeg?auto=compress&cs=tinysrgb&w=1200")',
          backgroundSize: 'cover',
          backgroundPosition: 'center',
        }}
      />

      <div className="relative h-full flex items-center justify-center">
        <div className="text-center space-y-6 px-4">
          <h1 className="text-5xl md:text-7xl font-bold text-white tracking-tighter">
            VERSACE LUXURY
          </h1>
          <p className="text-xl md:text-2xl text-gray-300 max-w-2xl mx-auto">
            Experience the epitome of Italian elegance and sophisticated design
          </p>
          <button
            onClick={onShopClick}
            className="inline-block px-8 py-4 bg-yellow-500 hover:bg-yellow-600 text-black font-bold text-lg rounded-lg transition-colors duration-300 mt-6"
          >
            SHOP NOW
          </button>
        </div>
      </div>

      <div className="absolute bottom-0 left-0 right-0 h-24 bg-gradient-to-t from-gray-50 to-transparent" />
    </div>
  );
};
