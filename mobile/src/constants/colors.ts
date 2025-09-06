export const colors = {
  // Coffee palette
  coffee: {
    50: "#FAF7F2",
    100: "#F5EFE6",
    200: "#E6D7C0",
    300: "#D7BF9A",
    400: "#C8A774",
    500: "#8B4513", // Primary
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
    500: "#F5E6D3", // Primary
    600: "#E6D0B8",
    700: "#D7BA9D",
    800: "#C8A482",
    900: "#B98E67",
  },

  // Semantic colors
  primary: "#8B4513",
  secondary: "#F5E6D3",
  accent: "#D2691E",
  background: "#FFFEF7",
  surface: "#FAF7F2",

  text: {
    primary: "#2D1810",
    secondary: "#5A453A",
    muted: "#8B7355",
  },

  status: {
    success: "#22C55E",
    warning: "#F59E0B",
    error: "#EF4444",
    info: "#3B82F6",
  },

  // Common usage
  border: "#E6D7C0",
  divider: "#D7BF9A",
  overlay: "rgba(45, 24, 16, 0.5)",
} as const;

export type ColorKey = keyof typeof colors;
