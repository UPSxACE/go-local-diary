/** @type {import('tailwindcss').Config} */
const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
  content: ["./views/**/*.{html,js}", "./assets/**/*.{html,js}"],
  theme: {
    extend: {
      screens: {
        xs: "450px",
      },
      colors: {
        lgray: "#f6f9fc", // backgrounds (light gray)
        bgray: "#e1e4e6", // borders (border gray)
        bgray2: "#d9d9d9",
        bgray3: "#a4a4a4",
        dgray: "#dfdfdf", // backgrounds (dark gray)
        bsep: "#1b1f232b", // border (separation)
      },
      fontFamily: {
        sans: ['"Poppins"', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  variants: {},
  plugins: [],
};
