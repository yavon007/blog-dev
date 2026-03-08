import {
  defineConfig,
  presetUno,
  presetIcons,
  presetTypography,
} from 'unocss'

export default defineConfig({
  presets: [
    presetUno(),
    presetTypography(),
    presetIcons({
      scale: 1.2,
      warn: true,
      extraProperties: {
        display: 'inline-block',
        'vertical-align': 'middle',
      },
    }),
  ],
  theme: {
    colors: {
      primary: {
        50:  '#f0f9ff',
        100: '#e0f2fe',
        200: '#bae6fd',
        300: '#7dd3fc',
        400: '#38bdf8',
        500: '#0ea5e9',
        600: '#0284c7',
        700: '#0369a1',
        800: '#075985',
        900: '#0c4a6e',
      },
    },
    breakpoints: {
      sm: '640px',
      md: '768px',
      lg: '1024px',
      xl: '1280px',
      '2xl': '1536px',
    },
  },
  shortcuts: [
    ['btn', 'px-4 py-2 rounded font-medium transition-colors duration-200 cursor-pointer'],
    ['btn-primary', 'btn bg-primary-500 text-white hover:bg-primary-600'],
    ['btn-secondary', 'btn bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600'],
    ['btn-danger', 'btn bg-red-500 text-white hover:bg-red-600'],
    ['card', 'bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-100 dark:border-gray-700'],
    ['input-base', 'w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-primary-500'],
  ],
})
