// Cupping wheel constants
export const CUPING_CONSTANTS = {
  // Screen dimensions
  SCREEN_WIDTH_MULTIPLIER: {
    FIRST_CIRCLE: 1.0,
    SECOND_CIRCLE: 1.5,
    THIRD_CIRCLE: 2.0,
  },

  // Circle positioning
  CIRCLE_OFFSETS: {
    FIRST: -0.15,
    SECOND: -0.095,
    THIRD: -0.075,
  },

  // Animation settings
  ANIMATION: {
    DURATION: 800,
    EASING: "quad",
    GESTURE_SENSITIVITY: 2,
  },

  // Blind item IDs
  BLIND_ITEM_IDS: {
    SECOND_CIRCLE: "0_0_0",
    THIRD_CIRCLE: "0_0_0_0",
  },

  // Default selections
  DEFAULT_SELECTIONS: {
    FRUITY_INDEX: 0,
    BERRY_INDEX: 0,
    RASPBERRY_INDEX: 1,
  },

  // Z-index values
  Z_INDEX: {
    FIRST_CIRCLE: 30,
    SECOND_CIRCLE: 20,
    THIRD_CIRCLE: 10,
    POINTER: 40,
    POINTER_2: 30,
    POINTER_3: 20,
    SAVE_BUTTON: 40,
    SELECTION_DISPLAY: 50,
  },

  // Colors
  COLORS: {
    BLIND: "#cbcfce",
    SELECTED_STROKE: "#FFD700",
    DEFAULT_STROKE: "#fff",
    BACKGROUND: "#cbcfce",
    SAVE_BUTTON: "#ffffff",
  },

  // Dimensions
  DIMENSIONS: {
    POINTER: {
      FIRST: { width: 20, height: 20 },
      SECOND: { width: 25, height: 25 },
      THIRD: { width: 30, height: 30 },
    },
    SAVE_BUTTON: {
      WIDTH: 170,
      HEIGHT: 80,
      RADIUS: 100,
    },
  },
} as const;

// Animation easing types
export type EasingType = "quad" | "cubic" | "sin" | "exp";

// Circle types
export type CircleType = "first" | "second" | "third";

// Selection state interface
export interface SelectionState {
  parent: any | null;
  child: any | null;
  grandchild: any | null;
}
