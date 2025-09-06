/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./App.{js,jsx,ts,tsx}", "./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        // Coffee-themed primary colors
        coffee: {
          50: "#FAF7F2",
          100: "#F5EFE6",
          200: "#E6D7C0",
          300: "#D7BF9A",
          400: "#C8A774",
          500: "#8B4513", // Primary coffee brown
          600: "#7A3D11",
          700: "#693610",
          800: "#582E0E",
          900: "#47260C",
        },
        cream: {
          50: "#FFFEF7",
          100: "#FFFDF0",
          200: "#FEF9E1",
          300: "#FDF5D2",
          400: "#FCF1C3",
          500: "#F5E6D3", // Primary cream
          600: "#E6D0B8",
          700: "#D7BA9D",
          800: "#C8A482",
          900: "#B98E67",
        },
        // UI colors
        primary: "#8B4513", // Coffee brown
        secondary: "#F5E6D3", // Cream
        accent: "#D2691E", // Chocolate
        background: "#FFFEF7", // Light cream
        surface: "#FAF7F2", // Off-white
        text: {
          primary: "#2D1810", // Dark coffee
          secondary: "#5A453A", // Medium coffee
          muted: "#8B7355", // Light coffee
        },
        success: "#22C55E",
        warning: "#F59E0B",
        error: "#EF4444",
        info: "#3B82F6",
      },
    },
  },
  plugins: [],
};
