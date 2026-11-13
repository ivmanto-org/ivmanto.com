/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './composables/**/*.ts',
    './plugins/**/*.ts',
    './app.vue',
  ],
  darkMode: ['selector', '[data-theme="dark"]'],
  // '.overline' is the design system's label class, not Tailwind's text-decoration utility.
  blocklist: ['overline'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['var(--font-sans)'],
        display: ['var(--font-display)'],
        mono: ['var(--font-mono)'],
      },
      colors: {
        // Semantic roles from docs/design-system/tokens (light/dark via [data-theme]).
        bg: 'var(--bg)',
        surface: 'var(--surface)',
        'surface-sunken': 'var(--surface-sunken)',
        text: 'var(--text)',
        muted: 'var(--text-muted)',
        border: 'var(--border)',
        'border-input': 'var(--border-input)',
        brand: 'var(--brand)',
        signal: 'var(--signal)',
        success: 'var(--success)',
        warning: 'var(--warning)',
        danger: 'var(--danger)',
        // Legacy names used by pages not yet redesigned, mapped onto the roles.
        primary: 'var(--accent)',
        'primary-light': 'var(--accent)',
        'primary-dark': 'var(--accent-hover)',
        'dark-slate': 'var(--text)',
        accent: 'var(--accent)',
        'accent-light': 'var(--accent)',
        secondary: 'var(--accent-hover)',
        'light-gray': 'var(--surface-sunken)',
        amber: 'var(--signal)',
        'amber-light': 'var(--signal)',
        'amber-dark': 'var(--signal)',
      },
      animation: {
        'fade-in-up': 'fadeInUp 0.8s ease-out forwards',
        'fade-in-up-delay-1': 'fadeInUp 0.8s ease-out 0.15s forwards',
        'fade-in-up-delay-2': 'fadeInUp 0.8s ease-out 0.3s forwards',
        'fade-in-up-delay-3': 'fadeInUp 0.8s ease-out 0.45s forwards',
        'fade-in': 'fadeIn 0.6s ease-out forwards',
        'blob-float': 'blobFloat 8s ease-in-out infinite',
        'blob-float-reverse': 'blobFloat 10s ease-in-out infinite reverse',
        'pulse-glow': 'pulseGlow 3s ease-in-out infinite',
        'text-reveal': 'textReveal 0.8s cubic-bezier(0.25, 0.46, 0.45, 0.94) forwards',
        'gradient-shift': 'gradientShift 6s ease-in-out infinite',
        'float': 'float 6s ease-in-out infinite',
      },
      keyframes: {
        fadeInUp: {
          '0%': { opacity: '0', transform: 'translateY(30px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        blobFloat: {
          '0%, 100%': { transform: 'translate(0, 0) scale(1)' },
          '25%': { transform: 'translate(20px, -30px) scale(1.05)' },
          '50%': { transform: 'translate(-10px, 15px) scale(0.95)' },
          '75%': { transform: 'translate(15px, 10px) scale(1.02)' },
        },
        pulseGlow: {
          '0%, 100%': { boxShadow: '0 0 20px rgba(0, 168, 150, 0.3)' },
          '50%': { boxShadow: '0 0 40px rgba(0, 168, 150, 0.6)' },
        },
        textReveal: {
          '0%': { opacity: '0', transform: 'translateY(20px) rotateX(-10deg)' },
          '100%': { opacity: '1', transform: 'translateY(0) rotateX(0)' },
        },
        gradientShift: {
          '0%, 100%': { backgroundPosition: '0% 50%' },
          '50%': { backgroundPosition: '100% 50%' },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-20px)' },
        },
      },
      backgroundSize: {
        '200%': '200% 200%',
      },
    },
  },
  plugins: [require('@tailwindcss/typography')],
}
