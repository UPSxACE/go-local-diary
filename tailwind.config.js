/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./server/views/**/*.{html,js}", "./server/public/**/*.{html,js}"],
  theme: {
    extend: {
      colors: {
        lgray: "#f6f9fc", // backgrounds (light gray)
        bgray: "#e1e4e6", // borders (border gray)
      },
    },
  },
  plugins: [],
};
