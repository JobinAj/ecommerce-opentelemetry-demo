import React, { useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { MonsterSVG } from './MonsterSVG';
import { login, signup } from '../api/client';

export const MonsterLogin: React.FC = () => {
    const [isSignUp, setIsSignUp] = useState(false);
    const [cursorPos, setCursorPos] = useState({ x: 0, y: 0 });
    const [isPasswordFocused, setIsPasswordFocused] = useState(false);
    const [isError, setIsError] = useState(false);
    const [isSuccess, setIsSuccess] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [name, setName] = useState('');

    const navigate = useNavigate();
    const containerRef = useRef<HTMLDivElement>(null);

    const handleMouseMove = (e: React.MouseEvent) => {
        if (containerRef.current) {
            const rect = containerRef.current.getBoundingClientRect();
            setCursorPos({
                x: e.clientX - rect.left,
                y: e.clientY - rect.top
            });
        }
    };

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsError(false);
        setErrorMessage('');

        try {
            const user = await login(email, password);
            setIsSuccess(true);
            localStorage.setItem('user', JSON.stringify(user));
            setTimeout(() => navigate('/'), 1500); // Wait for success animation
        } catch (err: any) {
            setIsError(true);
            setErrorMessage(err.message);
            setTimeout(() => setIsError(false), 500); // Reset shake
        }
    };

    const handleSignup = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsError(false);
        setErrorMessage('');

        try {
            const user = await signup(name, email, password);
            setIsSuccess(true);
            // Auto login after signup
            localStorage.setItem('user', JSON.stringify(user));
            setTimeout(() => navigate('/'), 1500);
        } catch (err: any) {
            setIsError(true);
            setErrorMessage(err.message);
            setTimeout(() => setIsError(false), 500);
        }
    };

    return (
        <div
            className="min-h-screen bg-gray-100 flex items-center justify-center p-4"
            onMouseMove={handleMouseMove}
        >
            <div
                ref={containerRef}
                className="bg-white rounded-2xl shadow-2xl w-full max-w-4xl min-h-[600px] relative overflow-hidden flex"
            >
                {/* Monster Container - Absolute positioned to track cursor locally */}
                <div className="absolute top-4 left-1/2 transform -translate-x-1/2 z-20 pointer-events-none">
                    <MonsterSVG
                        cursorPos={cursorPos}
                        isPasswordFocused={isPasswordFocused}
                        isError={isError}
                        isSuccess={isSuccess}
                    />
                </div>

                {/* Sign Up Form Container */}
                <div className={`w-1/2 p-12 flex flex-col justify-center items-center transition-all duration-700 ease-in-out absolute top-0 h-full left-0 ${isSignUp ? 'translate-x-[100%] opacity-0 z-0' : 'opacity-100 z-10'}`}>
                    {/* This is actually Sign In View (Logic inverted in standard templates, let's keep it simple) */}
                    {/* Wait, if isSignUp is false, we see Sign In form on the left? No, usually overlays slide. */}
                    {/* Let's build 2 panels side-by-side and slide the OVERLAY. */}
                </div>

                {/* Re-thinking layout for "Exact" sliding overlay style */}

                {/* Sign Up Form (Hidden by default, slides in) */}
                <div className={`absolute top-0 h-full transition-all duration-700 ease-in-out left-0 w-1/2 flex flex-col justify-center items-center p-12 bg-white ${isSignUp ? 'translate-x-[100%] opacity-100 z-50' : 'opacity-0 z-0'}`}>
                    <h1 className="text-3xl font-bold mb-6 pt-20">Create Account</h1>
                    <form onSubmit={handleSignup} className="w-full flex flex-col gap-4">
                        <input
                            type="text"
                            placeholder="Name"
                            className="bg-gray-100 border-none p-3 rounded-lg w-full focus:ring-2 focus:ring-purple-500 outline-none"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                        />
                        <input
                            type="email"
                            placeholder="Email"
                            className="bg-gray-100 border-none p-3 rounded-lg w-full focus:ring-2 focus:ring-purple-500 outline-none"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                        />
                        <input
                            type="password"
                            placeholder="Password"
                            className="bg-gray-100 border-none p-3 rounded-lg w-full focus:ring-2 focus:ring-purple-500 outline-none"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            onFocus={() => setIsPasswordFocused(true)}
                            onBlur={() => setIsPasswordFocused(false)}
                        />
                        {errorMessage && <p className="text-red-500 text-sm">{errorMessage}</p>}
                        <button className="bg-purple-600 text-white font-bold py-3 px-12 rounded-full uppercase tracking-wider hover:bg-purple-700 transition transform hover:scale-105 mt-4">
                            Sign Up
                        </button>
                    </form>
                </div>

                {/* Sign In Form (Visible by default) */}
                <div className={`absolute top-0 h-full transition-all duration-700 ease-in-out left-0 w-1/2 flex flex-col justify-center items-center p-12 bg-white ${isSignUp ? 'translate-x-[100%] opacity-0 z-0' : 'opacity-100 z-10'}`}>
                    <h1 className="text-3xl font-bold mb-6 pt-20">Sign in</h1>
                    <form onSubmit={handleLogin} className="w-full flex flex-col gap-4">
                        <input
                            type="email"
                            placeholder="Email"
                            className="bg-gray-100 border-none p-3 rounded-lg w-full focus:ring-2 focus:ring-purple-500 outline-none"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                        />
                        <input
                            type="password"
                            placeholder="Password"
                            className="bg-gray-100 border-none p-3 rounded-lg w-full focus:ring-2 focus:ring-purple-500 outline-none"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            onFocus={() => setIsPasswordFocused(true)}
                            onBlur={() => setIsPasswordFocused(false)}
                        />
                        {errorMessage && <p className="text-red-500 text-sm">{errorMessage}</p>}
                        <a href="#" className="text-sm text-gray-500 hover:text-gray-800 mb-4">Forgot your password?</a>
                        <button className="bg-purple-600 text-white font-bold py-3 px-12 rounded-full uppercase tracking-wider hover:bg-purple-700 transition transform hover:scale-105">
                            Sign In
                        </button>
                    </form>
                </div>

                {/* Overlay Container */}
                <div className={`absolute top-0 left-1/2 w-1/2 h-full overflow-hidden transition-transform duration-700 ease-in-out z-100 ${isSignUp ? '-translate-x-full' : ''}`}>
                    <div className={`bg-gradient-to-r from-purple-600 to-indigo-600 text-white relative -left-full h-full w-[200%] transform transition-transform duration-700 ease-in-out flex flex-row ${isSignUp ? 'translate-x-1/2' : 'translate-x-0'}`}>

                        {/* Left Panel (Shown when IsSignUp is true, effectively overlay is on left) */}
                        <div className={`w-1/2 flex flex-col justify-center items-center p-12 text-center transform transition-transform duration-700 ease-in-out ${isSignUp ? 'translate-x-0' : 'translate-x-20'}`}>
                            <h1 className="text-3xl font-bold mb-4">Welcome Back!</h1>
                            <p className="mb-8">To keep connected with us please login with your personal info</p>
                            <button
                                className="bg-transparent border-2 border-white text-white font-bold py-3 px-12 rounded-full uppercase tracking-wider hover:bg-white hover:text-purple-600 transition"
                                onClick={() => setIsSignUp(false)}
                            >
                                Sign In
                            </button>
                        </div>

                        {/* Right Panel (Shown when IsSignUp is false, overlay on right) */}
                        <div className={`w-1/2 flex flex-col justify-center items-center p-12 text-center transform transition-transform duration-700 ease-in-out ${isSignUp ? 'translate-x-20' : 'translate-x-0'}`}>
                            <h1 className="text-3xl font-bold mb-4">Hello, Friend!</h1>
                            <p className="mb-8">Enter your personal details and start journey with us</p>
                            <button
                                className="bg-transparent border-2 border-white text-white font-bold py-3 px-12 rounded-full uppercase tracking-wider hover:bg-white hover:text-purple-600 transition"
                                onClick={() => setIsSignUp(true)}
                            >
                                Sign Up
                            </button>
                        </div>
                    </div>
                </div>

            </div>
        </div>
    );
};
