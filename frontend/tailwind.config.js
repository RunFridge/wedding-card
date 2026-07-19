/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: 'rgb(var(--c-primary) / <alpha-value>)',
        secondary: 'rgb(var(--c-secondary) / <alpha-value>)',
        accent: 'rgb(var(--c-accent) / <alpha-value>)',
        grass: 'rgb(var(--c-grass) / <alpha-value>)',
        sky: 'rgb(var(--c-sky) / <alpha-value>)',
        wood: 'rgb(var(--c-wood) / <alpha-value>)',
        'wood-dark': 'rgb(var(--c-wood-dark) / <alpha-value>)',
        parchment: 'rgb(var(--c-parchment) / <alpha-value>)',
        'parchment-dark': 'rgb(var(--c-parchment-dark) / <alpha-value>)',
      },
      fontFamily: {
        pixel: ['NeoDunggeunmo', 'monospace'],
      },
    },
  },
  plugins: [],
}
