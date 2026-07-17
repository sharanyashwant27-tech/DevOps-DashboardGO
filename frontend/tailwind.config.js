/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      fontFamily: {
        display: ['"Sora"', 'system-ui', 'sans-serif'],
        body: ['"IBM Plex Sans"', 'system-ui', 'sans-serif'],
        mono: ['"IBM Plex Mono"', 'ui-monospace', 'monospace'],
      },
      colors: {
        ink: {
          950: '#07101f',
          900: '#0c1a2e',
          800: '#12243a',
        },
        accent: {
          DEFAULT: '#22d3ee',
          soft: '#67e8f9',
          sky: '#38bdf8',
          warm: '#fbbf24',
          coral: '#fb7185',
          mint: '#34d399',
          warn: '#f59e0b',
          danger: '#f43f5e',
        },
      },
      boxShadow: {
        glow: '0 0 40px rgba(34, 211, 238, 0.22)',
        warm: '0 0 36px rgba(251, 191, 36, 0.2)',
      },
      keyframes: {
        rise: {
          '0%': { opacity: 0, transform: 'translateY(12px)' },
          '100%': { opacity: 1, transform: 'translateY(0)' },
        },
        pulseSoft: {
          '0%, 100%': { opacity: 1 },
          '50%': { opacity: 0.7 },
        },
      },
      animation: {
        rise: 'rise 0.45s ease-out both',
        pulseSoft: 'pulseSoft 2.2s ease-in-out infinite',
      },
    },
  },
  plugins: [],
};
