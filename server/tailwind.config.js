/** @type {import('tailwindcss').Config} */
const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
  content: ["./views/**/*.{html,js}", "./public/**/*.{html,js}"],
  theme: {
    extend: {
      colors: {
        lgray: "#f6f9fc", // backgrounds (light gray)
        bgray: "#e1e4e6", // borders (border gray)
        bgray2: "#d9d9d9",
        bgray3: "#a4a4a4",
        dgray: "#dfdfdf", // backgrounds (dark gray)
      },
      fontFamily: {
        sans: ['"Poppins"', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  plugins: [],
};
