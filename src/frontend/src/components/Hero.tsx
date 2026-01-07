interface HeroProps {
  onShopClick: () => void;
}

export const Hero = ({ onShopClick }: HeroProps) => {
  return (
    <div className="relative w-full h-[90vh] bg-black overflow-hidden">
      <div
        className="absolute inset-0 opacity-60"
        style={{
          backgroundImage:
            'url("https://images.unsplash.com/photo-1549497538-303791108f95?q=80&w=1920&auto=format&fit=crop")',
          backgroundSize: 'cover',
          backgroundPosition: 'center',
        }}
      />

      <div className="relative h-full flex flex-col items-center justify-center text-center px-4">
        <h2 className="text-white text-sm md:text-base font-bold tracking-[0.3em] mb-4">
          FALL WINTER 2026
        </h2>
        <h1 className="text-5xl md:text-8xl font-serif font-bold text-white mb-8">
          LA GRECA
        </h1>
        <button
          onClick={onShopClick}
          className="px-10 py-4 bg-white text-black font-bold text-sm tracking-widest hover:bg-versace-gold hover:text-white transition-colors duration-300 rounded-none"
        >
          DISCOVER COLLECTION
        </button>
      </div>
    </div>
  );
};
