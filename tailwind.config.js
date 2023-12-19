/** @type {import('tailwindcss').Config} */
const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
  content: ["./server/views/**/*.{html,js}", "./server/public/**/*.{html,js}"],
  theme: {
    extend: {
      colors: {
        lgray: "#f6f9fc", // backgrounds (light gray)
        bgray: "#e1e4e6", // borders (border gray)
      },
      fontFamily: {
        sans: ['"Poppins"', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  plugins: [],
};
