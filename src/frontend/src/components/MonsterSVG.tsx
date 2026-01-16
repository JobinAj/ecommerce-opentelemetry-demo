import React, { useEffect, useState, useRef } from 'react';

interface MonsterSVGProps {
    isPasswordFocused: boolean;
    isError: boolean;
    isSuccess: boolean;
    cursorPos: { x: number; y: number };
}

export const MonsterSVG: React.FC<MonsterSVGProps> = ({
    isPasswordFocused,
    isError,
    isSuccess,
    cursorPos,
}) => {
    const svgRef = useRef<SVGSVGElement>(null);
    const [pupilPos, setPupilPos] = useState({ x: 0, y: 0 });

    useEffect(() => {
        if (isPasswordFocused || !svgRef.current) return;

        // Calculate eye tracking
        const svgRect = svgRef.current.getBoundingClientRect();
        const svgCenterX = svgRect.left + svgRect.width / 2;
        const svgCenterY = svgRect.top + svgRect.height / 2;

        const angle = Math.atan2(cursorPos.y - svgCenterY, cursorPos.x - svgCenterX);
        const distance = Math.min(
            10,
            Math.hypot(cursorPos.x - svgCenterX, cursorPos.y - svgCenterY) / 10
        );

        setPupilPos({
            x: Math.cos(angle) * distance,
            y: Math.sin(angle) * distance,
        });
    }, [cursorPos, isPasswordFocused]);

    // Animation States
    const handsY = isPasswordFocused ? -45 : 0;

    // Mouth Path
    let mouthPath = "M 75,130 Q 100,140 125,130"; // Neutral
    if (isSuccess) mouthPath = "M 75,130 Q 100,160 125,130"; // Happy
    if (isError) mouthPath = "M 75,140 Q 100,120 125,140"; // Sad

    return (
        <svg
            ref={svgRef}
            viewBox="0 0 200 200"
            className={`w-48 h-48 transition-transform duration-300 ${isError ? 'animate-shake' : ''}`}
        >
            <defs>
                <mask id="face-mask">
                    <circle cx="100" cy="100" r="100" fill="white" />
                </mask>
            </defs>

            {/* Body/Face */}
            <circle cx="100" cy="100" r="90" fill="#4B5563" /> {/* Grey-700 */}

            {/* Eyes Container */}
            <g className="transition-all duration-300" style={{ opacity: isPasswordFocused ? 0 : 1 }}>
                {/* Left Eye */}
                <circle cx="70" cy="85" r="25" fill="white" />
                <circle cx={70 + pupilPos.x} cy={85 + pupilPos.y} r="10" fill="black" />

                {/* Right Eye */}
                <circle cx="130" cy="85" r="25" fill="white" />
                <circle cx={130 + pupilPos.x} cy={85 + pupilPos.y} r="10" fill="black" />
            </g>

            {/* Mouth */}
            <path
                d={mouthPath}
                fill="none"
                stroke="white"
                strokeWidth="5"
                strokeLinecap="round"
                className="transition-all duration-300"
            />

            {/* Hands (Initially hidden below, translate up) */}
            <g
                className="transition-transform duration-500 ease-in-out"
                style={{ transform: `translateY(${handsY}px)` }}
            >
                {/* Left Hand */}
                <path
                    d="M 20,160 Q 50,130 80,160 L 80,200 L 20,200 Z"
                    fill="#374151" // Grey-800
                    mask="url(#face-mask)"
                />
                {/* Fingers Left */}
                <path d="M 30,160 L 30,150 M 45,155 L 45,140 M 60,155 L 60,140" stroke="#374151" strokeWidth="10" strokeLinecap="round" />

                {/* Right Hand */}
                <path
                    d="M 120,160 Q 150,130 180,160 L 180,200 L 120,200 Z"
                    fill="#374151"
                    mask="url(#face-mask)"
                />
                {/* Fingers Right */}
                <path d="M 130,155 L 130,140 M 145,155 L 145,140 M 160,160 L 160,150" stroke="#374151" strokeWidth="10" strokeLinecap="round" />
            </g>
        </svg>
    );
};
