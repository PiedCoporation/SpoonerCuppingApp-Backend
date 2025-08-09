import { items } from "@/data/data";
import React, { useRef } from "react";
import {
  View,
  Dimensions,
  PanResponder,
  Animated,
  StyleSheet,
  Text,
} from "react-native";
import Svg, { G, Path, Polygon } from "react-native-svg";

const { width: screenWidth } = Dimensions.get("window");
const radius = screenWidth; // make wheel big enough
const center = radius / 2;
const sliceAngle = 360 / 9;

// Function to create a sector path
function createSector(startAngle: number, endAngle: number) {
  const largeArc = endAngle - startAngle <= 180 ? 0 : 1;
  const x1 = center + center * Math.cos((Math.PI * startAngle) / 180);
  const y1 = center + center * Math.sin((Math.PI * startAngle) / 180);
  const x2 = center + center * Math.cos((Math.PI * endAngle) / 180);
  const y2 = center + center * Math.sin((Math.PI * endAngle) / 180);
  return `M${center},${center} L${x1},${y1} A${center},${center} 0 ${largeArc} 1 ${x2},${y2} Z`;
}

// Function to determine which slice is selected based on rotation
function getSelectedSlice(rotationValue: number) {
  // Normalize rotation to 0-360 range
  let normalizedRotation = rotationValue % 360;
  if (normalizedRotation < 0) normalizedRotation += 360;

  // The slices are arranged starting from 0 degrees and going clockwise
  // But visually, we need to account for the SVG rotation and the fact that
  // the pointer is at the top center of the visible half-circle

  // Adjust for the visual layout - the mapping was exactly reversed
  let sliceIndex = Math.floor((360 - normalizedRotation) / sliceAngle);

  // Ensure we stay within bounds
  sliceIndex = sliceIndex % items.length;
  if (sliceIndex < 0) sliceIndex += items.length;

  console.log(
    "Rotation:",
    rotationValue,
    "Normalized:",
    normalizedRotation,
    "Index:",
    sliceIndex
  );

  return items[sliceIndex];
}

export default function CuppingRegistration() {
  const rotation = useRef(new Animated.Value(0)).current;
  const lastRotation = useRef(0);

  const panResponder = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => true,
      onPanResponderMove: (_, gesture) => {
        const newRotation = lastRotation.current + gesture.dx / 2; // adjust speed
        rotation.setValue(newRotation);
      },
      onPanResponderRelease: (_, gesture) => {
        const finalRotation = lastRotation.current + gesture.dx / 2;
        lastRotation.current = finalRotation;

        // Get the selected slice when user stops spinning
        try {
          const selectedSlice = getSelectedSlice(finalRotation);
          if (selectedSlice && selectedSlice.id) {
            console.log("Selected slice ID:", selectedSlice.id);
            console.log("Selected slice label:", selectedSlice.label);
          } else {
            console.log("No valid slice selected");
          }
        } catch (error) {
          console.log("Error getting selected slice:", error);
        }
      },
    })
  ).current;

  return (
    <View style={styles.container}>
      {/* Half circle container */}
      <View style={styles.halfCircleWrapper}>
        {/* Fixed pointer/indicator */}
        <View style={styles.pointer}>
          <Svg width={30} height={30} style={styles.pointerSvg}>
            <Polygon
              points="15,25 10,15 20,15"
              fill="#000"
              stroke="#fff"
              strokeWidth="2"
            />
          </Svg>
        </View>

        <View style={styles.halfCircle} {...panResponder.panHandlers}>
          <Animated.View
            style={{
              transform: [
                {
                  rotate: rotation.interpolate({
                    inputRange: [-360, 360],
                    outputRange: ["-360deg", "360deg"],
                  }),
                },
              ],
            }}
          >
            <Svg width={radius} height={radius}>
              <G rotation={-90} originX={center} originY={center}>
                {items.map((item, index) => {
                  const startAngle = index * sliceAngle;
                  const endAngle = startAngle + sliceAngle;
                  return (
                    <Path
                      key={item.id}
                      d={createSector(startAngle, endAngle)}
                      fill={item.color}
                      stroke="white"
                      strokeWidth="1"
                    />
                  );
                })}
              </G>
            </Svg>

            {/* Text labels positioned inside Animated.View but counter-rotated to stay horizontal */}
            {items.map((item, index) => {
              const startAngle = index * sliceAngle;
              const endAngle = startAngle + sliceAngle;
              const midAngle = (startAngle + endAngle) / 2;

              // Convert to radians and adjust for the -90 degree rotation of the SVG
              const adjustedAngle = (midAngle - 90) * (Math.PI / 180);

              // Position text closer to the middle of each slice so it's not cut off
              const distanceFromCenter = center * 0.75; // Move back from edge a bit
              const x = center + distanceFromCenter * Math.cos(adjustedAngle);
              const y = center + distanceFromCenter * Math.sin(adjustedAngle);

              return (
                <Animated.View
                  key={item.id + "-label"}
                  style={{
                    position: "absolute",
                    left: x - 60, // Wider container for longer text
                    top: y - 20, // Taller container for multi-line text
                    width: 120, // Increased width
                    height: 40, // Increased height for multi-line
                    alignItems: "center",
                    justifyContent: "center",
                    transform: [
                      {
                        rotate: rotation.interpolate({
                          inputRange: [-360, 360],
                          outputRange: ["360deg", "-360deg"], // Counter-rotate to keep text horizontal
                        }),
                      },
                    ],
                  }}
                >
                  <Text
                    style={{
                      fontSize: 12, // Slightly smaller font to fit better
                      color: "white",
                      textAlign: "center",
                      fontWeight: "bold",
                      backgroundColor: "rgba(0, 0, 0, 0.3)", // Slightly more opaque background
                      paddingHorizontal: 6,
                      paddingVertical: 3,
                      borderRadius: 4,
                      lineHeight: 16, // Better line spacing for multi-line text
                    }}
                  >
                    {item.label.replace("/", "/\n")}
                  </Text>
                </Animated.View>
              );
            })}
          </Animated.View>
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "white" },
  halfCircleWrapper: {
    position: "absolute",
    bottom: 0,
    width: screenWidth,
    height: screenWidth / 2,
    overflow: "hidden",
    alignItems: "center",
  },
  halfCircle: {
    width: radius,
    height: radius,
    borderRadius: radius / 2,
    overflow: "hidden",
    alignItems: "center",
    justifyContent: "center",
  },
  pointer: {
    position: "absolute",
    top: -5,
    zIndex: 10,
    alignItems: "center",
  },
  pointerSvg: {
    zIndex: 10,
  },
});
