import React, { useState, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { login, signup } from '../api/client';
import heroImage from '../assets/versace_hero.png';
import versaceLogo from '../assets/versace_logo.png';

// Floating Gold Particle Component
const GoldParticle: React.FC<{ delay: number; duration: number; left: number; size: number }> = ({
    delay, duration, left, size
}) => (
    <div
        className="absolute rounded-full pointer-events-none"
        style={{
            width: size,
            height: size,
            left: `${left}%`,
            bottom: '-10px',
            background: 'radial-gradient(circle, rgba(212,175,55,0.8) 0%, rgba(212,175,55,0) 70%)',
            animation: `floatUp ${duration}s ease-in-out ${delay}s infinite`,
            boxShadow: '0 0 10px rgba(212,175,55,0.5)',
        }}
    />
);

// Animated Greek Key Pattern
const GreekKeyPattern: React.FC<{ position: 'left' | 'right' }> = ({ position }) => (
    <div
        className={`absolute top-1/2 ${position === 'left' ? 'left-4' : 'right-4'} -translate-y-1/2 opacity-20`}
        style={{
            animation: `pulse 4s ease-in-out infinite ${position === 'right' ? '2s' : '0s'}`,
        }}
    >
        <svg width="24" height="120" viewBox="0 0 24 120" fill="none" className="text-versace-gold">
            {[0, 40, 80].map((y, i) => (
                <g key={i} transform={`translate(0, ${y})`}>
                    <path
                        d="M4 4h16v8h-8v8h8v8H4v-8h8v-8H4V4z"
                        stroke="currentColor"
                        strokeWidth="1"
                        fill="none"
                        style={{
                            animation: `drawLine 3s ease-in-out ${i * 0.5}s infinite`,
                        }}
                    />
                </g>
            ))}
        </svg>
    </div>
);

export const VersaceLogin: React.FC = () => {
    const [isSignUp, setIsSignUp] = useState(false);
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [name, setName] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');
    const [isPageLoaded, setIsPageLoaded] = useState(false);
    const [focusedField, setFocusedField] = useState<string | null>(null);

    const navigate = useNavigate();

    // Trigger entrance animations on mount
    useEffect(() => {
        const timer = setTimeout(() => setIsPageLoaded(true), 100);
        return () => clearTimeout(timer);
    }, []);

    // Generate particles with stable values using useMemo
    const particles = useMemo(() =>
        Array.from({ length: 15 }, (_, i) => ({
            id: i,
            delay: (i * 0.8) % 8,
            duration: 6 + (i % 4),
            left: 5 + ((i * 17) % 90),
            size: 4 + (i % 4) * 2,
        })),
        []);

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setErrorMessage('');

        try {
            const user = await login(email, password);
            localStorage.setItem('user', JSON.stringify(user));
            navigate('/home');
        } catch (err: any) {
            setErrorMessage(err.message || 'Login failed. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    const handleSignup = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setErrorMessage('');

        try {
            const user = await signup(name, email, password);
            localStorage.setItem('user', JSON.stringify(user));
            navigate('/home');
        } catch (err: any) {
            setErrorMessage(err.message || 'Sign up failed. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="h-screen w-full flex overflow-hidden">
            {/* CSS Animations */}
            <style>{`
                @keyframes floatUp {
                    0% {
                        transform: translateY(0) scale(1);
                        opacity: 0;
                    }
                    10% {
                        opacity: 1;
                    }
                    90% {
                        opacity: 1;
                    }
                    100% {
                        transform: translateY(-100vh) scale(0.5);
                        opacity: 0;
                    }
                }
                
                @keyframes pulse {
                    0%, 100% {
                        opacity: 0.1;
                        transform: translateY(-50%) scale(1);
                    }
                    50% {
                        opacity: 0.3;
                        transform: translateY(-50%) scale(1.05);
                    }
                }
                
                @keyframes shimmer {
                    0% {
                        transform: translateX(-100%);
                    }
                    100% {
                        transform: translateX(100%);
                    }
                }
                
                @keyframes glowPulse {
                    0%, 100% {
                        filter: drop-shadow(0 0 8px rgba(212,175,55,0.4)) drop-shadow(0 0 20px rgba(212,175,55,0.2));
                    }
                    50% {
                        filter: drop-shadow(0 0 15px rgba(212,175,55,0.7)) drop-shadow(0 0 35px rgba(212,175,55,0.4));
                    }
                }
                
                @keyframes slideInUp {
                    from {
                        opacity: 0;
                        transform: translateY(30px);
                    }
                    to {
                        opacity: 1;
                        transform: translateY(0);
                    }
                }
                
                @keyframes fadeIn {
                    from {
                        opacity: 0;
                    }
                    to {
                        opacity: 1;
                    }
                }
                
                @keyframes drawLine {
                    0%, 100% {
                        stroke-dasharray: 200;
                        stroke-dashoffset: 200;
                    }
                    50% {
                        stroke-dashoffset: 0;
                    }
                }
                
                @keyframes heroZoom {
                    0%, 100% {
                        transform: scale(1);
                    }
                    50% {
                        transform: scale(1.05);
                    }
                }
                
                @keyframes goldLineFlow {
                    0% {
                        background-position: -200% center;
                    }
                    100% {
                        background-position: 200% center;
                    }
                }
                
                @keyframes borderGlow {
                    0%, 100% {
                        box-shadow: 0 1px 0 0 rgba(212,175,55,0.5);
                    }
                    50% {
                        box-shadow: 0 1px 10px 0 rgba(212,175,55,0.8), 0 1px 20px 0 rgba(212,175,55,0.4);
                    }
                }
                
                @keyframes staggerFade {
                    from { opacity: 0; transform: translateY(10px); }
                    to { opacity: 1; transform: translateY(0); }
                }

                .animate-stagger-1 { animation: staggerFade 0.8s ease-out 0.1s both; }
                .animate-stagger-2 { animation: staggerFade 0.8s ease-out 0.2s both; }
                .animate-stagger-3 { animation: staggerFade 0.8s ease-out 0.3s both; }
                .animate-stagger-4 { animation: staggerFade 0.8s ease-out 0.4s both; }
                .animate-stagger-5 {
                    animation: slideInUp 0.8s ease-out 0.5s both;
                }
                .animate-stagger-6 {
                    animation: slideInUp 0.8s ease-out 0.6s both;
                }
                
                @keyframes slideDown {
                    from {
                        opacity: 0;
                        max-height: 0;
                        transform: translateY(-20px);
                    }
                    to {
                        opacity: 1;
                        max-height: 100px;
                        transform: translateY(0);
                    }
                }
                
                @keyframes slideUp {
                    from {
                        opacity: 1;
                        max-height: 100px;
                        transform: translateY(0);
                    }
                    to {
                        opacity: 0;
                        max-height: 0;
                        transform: translateY(-20px);
                    }
                }
                
                .name-field-enter {
                    animation: slideDown 0.4s ease-out forwards;
                    overflow: hidden;
                }
                
                .btn-shimmer::after {
                    content: '';
                    position: absolute;
                    top: 0;
                    left: 0;
                    width: 100%;
                    height: 100%;
                    background: linear-gradient(
                        90deg,
                        transparent,
                        rgba(255,255,255,0.4),
                        transparent
                    );
                    animation: shimmer 3s ease-in-out infinite;
                }
                
                .input-glow:focus {
                    animation: borderGlow 2s ease-in-out infinite;
                }
            `}</style>

            {/* Left Side - Fashion Image with Parallax Effect */}
            <div
                className="hidden lg:flex lg:w-1/2 relative overflow-hidden"
            >
                {/* Animated Background */}
                <div
                    className="absolute inset-0"
                    style={{
                        backgroundImage: `url(${heroImage})`,
                        backgroundSize: 'cover',
                        backgroundPosition: 'center',
                        animation: 'heroZoom 20s ease-in-out infinite',
                    }}
                />

                {/* Gradient Overlay */}
                <div
                    className="absolute inset-0"
                    style={{
                        background: 'linear-gradient(to bottom, rgba(0,0,0,0.2), rgba(0,0,0,0.5))',
                    }}
                />

                {/* Content with Fade In */}
                <div
                    className={`absolute inset-0 flex flex-col justify-end p-12 transition-all duration-1000 ${isPageLoaded ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'
                        }`}
                >
                    <h2
                        className="font-serif text-5xl text-white mb-4 tracking-wider"
                        style={{
                            textShadow: '0 0 30px rgba(212,175,55,0.3)',
                        }}
                    >
                        LA GRECA
                    </h2>
                    <p className="text-versace-gold text-lg tracking-[0.3em] uppercase">
                        Fall Winter 2026
                    </p>
                </div>

                {/* Floating particles on left side */}
                {particles.slice(0, 5).map((p) => (
                    <GoldParticle key={p.id} {...p} />
                ))}
            </div>

            {/* Right Side - Login Form - Centered and Fixed Height */}
            <div className="w-full lg:w-1/2 bg-versace-black flex flex-col justify-center items-center p-8 relative overflow-hidden">
                {/* Animated Gold Line - Top */}
                <div
                    className="absolute top-0 left-0 w-full h-1"
                    style={{
                        background: 'linear-gradient(90deg, transparent, transparent, rgba(212,175,55,0.3), #D4AF37, rgba(212,175,55,0.3), transparent, transparent)',
                        backgroundSize: '200% 100%',
                        animation: 'goldLineFlow 4s linear infinite',
                    }}
                />

                {/* Floating Gold Particles */}
                {particles.map((p) => (
                    <GoldParticle key={p.id} {...p} />
                ))}

                {/* Decorative Greek Key Patterns */}
                <GreekKeyPattern position="left" />
                <GreekKeyPattern position="right" />

                {/* Versace Logo with Glow Animation */}
                <div
                    className={`mb-6 flex-shrink-0 transition-all duration-1000 ${isPageLoaded ? 'opacity-100 scale-100' : 'opacity-0 scale-90'
                        }`}
                >
                    <img
                        src={versaceLogo}
                        alt="Versace"
                        className="w-36 h-auto mx-auto"
                        style={{
                            animation: 'glowPulse 3s ease-in-out infinite',
                        }}
                    />
                </div>

                {/* Form Container */}
                <div className="w-full max-w-md flex-shrink-0 relative z-10">
                    {/* Toggle Tabs with Staggered Animation */}
                    <div className={`flex mb-10 border-b border-gray-800 ${isPageLoaded ? 'animate-stagger-1' : 'opacity-0'}`}>
                        <button
                            onClick={() => { setIsSignUp(false); setErrorMessage(''); }}
                            className={`flex-1 pb-4 text-sm tracking-[0.2em] uppercase transition-all duration-500 relative ${!isSignUp
                                ? 'text-versace-gold'
                                : 'text-gray-500 hover:text-gray-300'
                                }`}
                        >
                            Sign In
                            {!isSignUp && (
                                <span
                                    className="absolute bottom-0 left-0 w-full h-0.5 bg-versace-gold"
                                    style={{
                                        boxShadow: '0 0 10px rgba(212,175,55,0.8), 0 0 20px rgba(212,175,55,0.4)',
                                    }}
                                />
                            )}
                        </button>
                        <button
                            onClick={() => { setIsSignUp(true); setErrorMessage(''); }}
                            className={`flex-1 pb-4 text-sm tracking-[0.2em] uppercase transition-all duration-500 relative ${isSignUp
                                ? 'text-versace-gold'
                                : 'text-gray-500 hover:text-gray-300'
                                }`}
                        >
                            Create Account
                            {isSignUp && (
                                <span
                                    className="absolute bottom-0 left-0 w-full h-0.5 bg-versace-gold"
                                    style={{
                                        boxShadow: '0 0 10px rgba(212,175,55,0.8), 0 0 20px rgba(212,175,55,0.4)',
                                    }}
                                />
                            )}
                        </button>
                    </div>

                    {/* Form */}
                    <form onSubmit={isSignUp ? handleSignup : handleLogin} className="space-y-6">
                        {/* Name Field (Sign Up only) - with smooth slide animation */}
                        <div
                            className={`group overflow-hidden transition-all duration-400 ease-out ${isSignUp
                                ? 'max-h-24 opacity-100 mb-0'
                                : 'max-h-0 opacity-0 -mb-6'
                                }`}
                            style={{
                                transitionProperty: 'max-height, opacity, margin-bottom',
                            }}
                        >
                            <label
                                className={`block text-xs tracking-[0.2em] uppercase mb-2 transition-colors duration-300 ${focusedField === 'name' ? 'text-versace-gold' : 'text-gray-400'
                                    }`}
                            >
                                Full Name
                            </label>
                            <input
                                type="text"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                onFocus={() => setFocusedField('name')}
                                onBlur={() => setFocusedField(null)}
                                className="w-full bg-transparent border-b border-gray-700 text-white py-3 px-1 focus:outline-none focus:border-versace-gold transition-all duration-300 font-sans input-glow"
                                placeholder="Enter your name"
                                required={isSignUp}
                                tabIndex={isSignUp ? 0 : -1}
                            />
                        </div>

                        {/* Email Field */}
                        <div className={`group ${isPageLoaded ? 'animate-stagger-2' : 'opacity-0'}`}>
                            <label
                                className={`block text-xs tracking-[0.2em] uppercase mb-2 transition-colors duration-300 ${focusedField === 'email' ? 'text-versace-gold' : 'text-gray-400'
                                    }`}
                            >
                                Email Address
                            </label>
                            <input
                                type="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                onFocus={() => setFocusedField('email')}
                                onBlur={() => setFocusedField(null)}
                                className="w-full bg-transparent border-b border-gray-700 text-white py-3 px-1 focus:outline-none focus:border-versace-gold transition-all duration-300 font-sans input-glow"
                                placeholder="Enter your email"
                                required
                            />
                        </div>

                        {/* Password Field */}
                        <div className={`group ${isPageLoaded ? 'animate-stagger-3' : 'opacity-0'}`}>
                            <label
                                className={`block text-xs tracking-[0.2em] uppercase mb-2 transition-colors duration-300 ${focusedField === 'password' ? 'text-versace-gold' : 'text-gray-400'
                                    }`}
                            >
                                Password
                            </label>
                            <input
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                onFocus={() => setFocusedField('password')}
                                onBlur={() => setFocusedField(null)}
                                className="w-full bg-transparent border-b border-gray-700 text-white py-3 px-1 focus:outline-none focus:border-versace-gold transition-all duration-300 font-sans input-glow"
                                placeholder="Enter your password"
                                required
                            />
                        </div>

                        {/* Error Message */}
                        {errorMessage && (
                            <p
                                className="text-red-400 text-sm text-center"
                                style={{ animation: 'slideInUp 0.3s ease-out' }}
                            >
                                {errorMessage}
                            </p>
                        )}

                        {/* Forgot Password (Sign In only) */}
                        {!isSignUp && (
                            <div className={`text-right ${isPageLoaded ? 'animate-stagger-4' : 'opacity-0'}`}>
                                <a
                                    href="#"
                                    className="text-gray-500 text-xs tracking-wider hover:text-versace-gold transition-all duration-300 hover:tracking-widest"
                                >
                                    Forgot Password?
                                </a>
                            </div>
                        )}

                        {/* Submit Button with Shimmer Effect */}
                        <button
                            type="submit"
                            disabled={isLoading}
                            className={`w-full mt-8 py-4 bg-versace-gold text-versace-black font-sans text-sm tracking-[0.2em] uppercase transition-all duration-500 disabled:opacity-50 disabled:cursor-not-allowed relative overflow-hidden group btn-shimmer hover:tracking-[0.3em] hover:shadow-[0_0_30px_rgba(212,175,55,0.5)] ${isPageLoaded ? 'animate-stagger-5' : 'opacity-0'
                                }`}
                            style={{
                                transform: isLoading ? 'scale(0.98)' : 'scale(1)',
                            }}
                        >
                            <span className={`relative z-10 transition-opacity ${isLoading ? 'opacity-0' : 'opacity-100'}`}>
                                {isSignUp ? 'Create Account' : 'Sign In'}
                            </span>
                            {isLoading && (
                                <span className="absolute inset-0 flex items-center justify-center">
                                    <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                                    </svg>
                                </span>
                            )}
                        </button>
                    </form>

                    {/* Premium Branding Section */}
                    <div className={`mt-12 text-center ${isPageLoaded ? 'animate-stagger-6' : 'opacity-0'}`}>
                        {/* Decorative Greek Key Pattern */}
                        <div className="flex items-center justify-center gap-2 mb-6">
                            <div
                                className="w-12 h-px bg-gradient-to-r from-transparent to-versace-gold"
                                style={{ animation: 'fadeIn 1s ease-out 0.8s both' }}
                            />
                            <svg
                                className="w-6 h-6 text-versace-gold"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                strokeWidth="1"
                                style={{ animation: 'glowPulse 3s ease-in-out infinite' }}
                            >
                                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
                            </svg>
                            <div
                                className="w-12 h-px bg-gradient-to-l from-transparent to-versace-gold"
                                style={{ animation: 'fadeIn 1s ease-out 0.8s both' }}
                            />
                        </div>
                        <p className="text-gray-500 text-xs tracking-[0.3em] uppercase hover:text-versace-gold transition-colors duration-500 cursor-default">
                            Exclusive Member Access
                        </p>
                    </div>
                </div>

                {/* Footer */}
                <div className="mt-auto pt-8 pb-4 text-center z-10">
                    <p
                        className="text-gray-600 text-xs tracking-wider transition-all duration-500 hover:text-gray-400 hover:tracking-widest cursor-default"
                        style={{ animation: 'fadeIn 1s ease-out 1s both' }}
                    >
                        © 2026 VERSACE. All Rights Reserved.
                    </p>
                </div>

                {/* Animated Gold Line - Bottom */}
                <div
                    className="absolute bottom-0 left-0 w-full h-1"
                    style={{
                        background: 'linear-gradient(90deg, transparent, transparent, rgba(212,175,55,0.3), #D4AF37, rgba(212,175,55,0.3), transparent, transparent)',
                        backgroundSize: '200% 100%',
                        animation: 'goldLineFlow 4s linear infinite reverse',
                    }}
                />
            </div>
        </div>
    );
};
