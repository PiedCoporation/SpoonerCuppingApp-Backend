import { items } from "@/data/data";
import { Sample } from "@/types/sample";
import React, { useRef, useState, useEffect } from "react";
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
const radius = screenWidth * 1; // First circle - made bigger
const center = radius / 2;
const sliceAngle = 360 / 9;

// Second circle properties - larger and behind first circle
const radius2 = screenWidth * 1.5; // Much larger second circle
const center2 = radius2 / 2;

// Third circle properties - even larger and behind second circle
const radius3 = screenWidth * 2; // Even larger third circle
const center3 = radius3 / 2;

// Function to create a sector path
function createSector(
  startAngle: number,
  endAngle: number,
  outerRadius: number,
  centerPoint: number
) {
  const largeArc = endAngle - startAngle <= 180 ? 0 : 1;
  const x1 = centerPoint + outerRadius * Math.cos((Math.PI * startAngle) / 180);
  const y1 = centerPoint + outerRadius * Math.sin((Math.PI * startAngle) / 180);
  const x2 = centerPoint + outerRadius * Math.cos((Math.PI * endAngle) / 180);
  const y2 = centerPoint + outerRadius * Math.sin((Math.PI * endAngle) / 180);
  return `M${centerPoint},${centerPoint} L${x1},${y1} A${outerRadius},${outerRadius} 0 ${largeArc} 1 ${x2},${y2} Z`;
}

// Function to determine which slice is selected based on rotation
function getSelectedSlice(rotationValue: number, itemsArray: Sample[]) {
  let normalizedRotation = rotationValue % 360;
  if (normalizedRotation < 0) normalizedRotation += 360;

  const currentSliceAngle = 360 / itemsArray.length;
  let sliceIndex = Math.floor((360 - normalizedRotation) / currentSliceAngle);

  sliceIndex = sliceIndex % itemsArray.length;
  if (sliceIndex < 0) sliceIndex += itemsArray.length;

  return itemsArray[sliceIndex];
}

export default function CuppingRegistration() {
  // Add this debug log at the very beginning
  console.log("🔍 Component render start");

  const rotation = useRef(new Animated.Value(0)).current;
  const lastRotation = useRef(0);
  const rotation2 = useRef(new Animated.Value(0)).current;
  const lastRotation2 = useRef(0);
  const rotation3 = useRef(new Animated.Value(0)).current;
  const lastRotation3 = useRef(0);

  // State for tracking spinning
  const [isSpinning, setIsSpinning] = useState(false);
  const [isSpinning2, setIsSpinning2] = useState(false);
  const [isSpinning3, setIsSpinning3] = useState(false);

  // Debug version of selectedParent state
  const [selectedParentState, setSelectedParentState] = useState<Sample | null>(
    null
  );
  const [selectedChild, setSelectedChildState] = useState<Sample | null>(null);
  const [selectedChild2, setSelectedChild2State] = useState<Sample | null>(
    null
  );

  // Refs to store current values for pan responders
  const selectedParentRef = useRef<Sample | null>(null);
  const selectedChildRef = useRef<Sample | null>(null);

  // Wrapper function to log when selectedParent changes
  const setSelectedParent = (value: Sample | null) => {
    console.log("🔄 Setting selectedParent to:", value);
    console.trace("Called from:"); // This will show you the call stack
    selectedParentRef.current = value; // Update ref immediately
    setSelectedParentState(value);
  };

  // Wrapper function to log when selectedChild changes
  const setSelectedChild = (value: Sample | null) => {
    console.log("🔄 Setting selectedChild to:", value);
    selectedChildRef.current = value; // Update ref immediately
    setSelectedChildState(value);
    setSelectedChild2State(null); // Reset third level selection when second level changes
  };

  // Wrapper function for selectedChild2
  const setSelectedChild2 = (value: Sample | null) => {
    console.log("🔄 Setting selectedChild2 to:", value);
    setSelectedChild2State(value);
  };

  // Monitor selectedParent changes
  useEffect(() => {
    console.log("📍 selectedParent changed to:", selectedParentState);
  }, [selectedParentState]);

  // Monitor selectedChild changes
  useEffect(() => {
    console.log("📍 selectedChild changed to:", selectedChild);
  }, [selectedChild]);

  // Monitor selectedChild2 changes
  useEffect(() => {
    console.log("📍 selectedChild2 changed to:", selectedChild2);
  }, [selectedChild2]);

  // Debug log current state
  console.log("🔍 Current selectedParent:", selectedParentState);
  console.log("🔍 Current selectedChild:", selectedChild);
  console.log("🔍 Current selectedChild2:", selectedChild2);
  console.log("🔍 Current isSpinning:", isSpinning);
  console.log("🔍 Current isSpinning2:", isSpinning2);
  console.log("🔍 Current isSpinning3:", isSpinning3);

  const panResponder = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => {
        console.log(
          "🎯 First circle - onStartShouldSetPanResponder, isSpinning2:",
          isSpinning2,
          "isSpinning3:",
          isSpinning3
        );
        return !isSpinning2 && !isSpinning3; // Don't respond if second or third circle is spinning
      },
      onPanResponderGrant: () => {
        console.log("🎯 First circle - onPanResponderGrant");
        setIsSpinning(true);
      },
      onPanResponderMove: (_, gesture) => {
        const newRotation = lastRotation.current + gesture.dx / 2;
        rotation.setValue(newRotation);
      },
      onPanResponderRelease: (_, gesture) => {
        console.log("🎯 First circle - onPanResponderRelease");
        setIsSpinning(false);
        const finalRotation = lastRotation.current + gesture.dx / 2;
        lastRotation.current = finalRotation;

        try {
          const selectedSlice = getSelectedSlice(finalRotation, items);
          if (selectedSlice && selectedSlice.id) {
            console.log("Selected parent ID:", selectedSlice.id);
            console.log("Selected parent label:", selectedSlice.label);
            setSelectedParent(selectedSlice);
            setSelectedChild(null); // Reset child selection
          }
        } catch (error) {
          console.log("Error getting selected slice:", error);
        }
      },
    })
  ).current;

  const panResponder2 = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => {
        console.log(
          "🎯 Second circle - onStartShouldSetPanResponder, isSpinning:",
          isSpinning,
          "isSpinning3:",
          isSpinning3
        );
        return !isSpinning && !isSpinning3; // Don't respond if first or third circle is spinning
      },
      onPanResponderGrant: () => {
        console.log("🎯 Second circle - onPanResponderGrant");
        setIsSpinning2(true);
      },
      onPanResponderMove: (_, gesture) => {
        const newRotation = lastRotation2.current + gesture.dx / 2;
        rotation2.setValue(newRotation);
      },
      onPanResponderRelease: (_, gesture) => {
        console.log("🎯 Second circle - onPanResponderRelease");
        setIsSpinning2(false);
        console.log("=== Second circle spin ended ===");
        const finalRotation = lastRotation2.current + gesture.dx / 2;
        lastRotation2.current = finalRotation;

        console.log("Final rotation:", finalRotation);
        console.log("Selected parent from ref:", selectedParentRef.current);
        console.log("Selected parent from state:", selectedParentState);

        if (
          selectedParentRef.current &&
          (selectedParentRef.current as any).childrens
        ) {
          console.log(
            "Childrens array length:",
            (selectedParentRef.current as any).childrens.length
          );
          try {
            const selectedSlice = getSelectedSlice(
              finalRotation,
              (selectedParentRef.current as any).childrens
            );
            console.log("Selected slice result:", selectedSlice);

            if (selectedSlice && selectedSlice.id) {
              console.log("✅ Selected child ID:", selectedSlice.id);
              console.log("✅ Selected child label:", selectedSlice.label);
              setSelectedChild(selectedSlice);
            } else {
              console.log("❌ No valid selectedSlice or missing ID");
            }
          } catch (error) {
            console.log("❌ Error getting selected child slice:", error);
          }
        } else {
          console.log("❌ No selectedParent found or no childrens");
        }
      },
    })
  ).current;

  const panResponder3 = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => {
        console.log(
          "🎯 Third circle - onStartShouldSetPanResponder, isSpinning:",
          isSpinning,
          "isSpinning2:",
          isSpinning2
        );
        return !isSpinning && !isSpinning2; // Don't respond if first or second circle is spinning
      },
      onPanResponderGrant: () => {
        console.log("🎯 Third circle - onPanResponderGrant");
        setIsSpinning3(true);
      },
      onPanResponderMove: (_, gesture) => {
        const newRotation = lastRotation3.current + gesture.dx / 2;
        rotation3.setValue(newRotation);
      },
      onPanResponderRelease: (_, gesture) => {
        console.log("🎯 Third circle - onPanResponderRelease");
        setIsSpinning3(false);
        console.log("=== Third circle spin ended ===");
        const finalRotation = lastRotation3.current + gesture.dx / 2;
        lastRotation3.current = finalRotation;

        console.log("Final rotation:", finalRotation);
        console.log("Selected child from ref:", selectedChildRef.current);
        console.log("Selected child from state:", selectedChild);

        if (
          selectedChildRef.current &&
          (selectedChildRef.current as any).childrens
        ) {
          console.log(
            "Child's childrens array length:",
            (selectedChildRef.current as any).childrens.length
          );
          try {
            const selectedSlice = getSelectedSlice(
              finalRotation,
              (selectedChildRef.current as any).childrens
            );
            console.log("Selected slice result:", selectedSlice);

            if (selectedSlice && selectedSlice.id) {
              console.log("✅ Selected child2 ID:", selectedSlice.id);
              console.log("✅ Selected child2 label:", selectedSlice.label);
              setSelectedChild2(selectedSlice);
            } else {
              console.log("❌ No valid selectedSlice or missing ID");
            }
          } catch (error) {
            console.log("❌ Error getting selected child2 slice:", error);
          }
        } else {
          console.log("❌ No selectedChild found or no childrens");
        }
      },
    })
  ).current;

  return (
    <View style={styles.container}>
      {/* First circle - Parent categories */}
      <View style={styles.firstCircleWrapper}>
        <View style={styles.pointer}>
          <Svg width={20} height={20} style={styles.pointerSvg}>
            <Polygon
              points="10,18 6,10 14,10"
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
                      d={createSector(startAngle, endAngle, center, center)}
                      fill={item.color}
                      stroke="white"
                      strokeWidth="2"
                    />
                  );
                })}
              </G>
            </Svg>

            {items.map((item, index) => {
              const startAngle = index * sliceAngle;
              const endAngle = startAngle + sliceAngle;
              const midAngle = (startAngle + endAngle) / 2;
              const adjustedAngle = (midAngle - 90) * (Math.PI / 180);
              const distanceFromCenter = center * 0.7;
              const x = center + distanceFromCenter * Math.cos(adjustedAngle);
              const y = center + distanceFromCenter * Math.sin(adjustedAngle);

              return (
                <Animated.View
                  key={item.id + "-label"}
                  style={{
                    position: "absolute",
                    left: x - 40,
                    top: y - 15,
                    width: 80,
                    height: 30,
                    alignItems: "center",
                    justifyContent: "center",
                    transform: [
                      {
                        rotate: rotation.interpolate({
                          inputRange: [-360, 360],
                          outputRange: ["360deg", "-360deg"],
                        }),
                      },
                    ],
                  }}
                >
                  <Text
                    style={{
                      fontSize: 10,
                      color: "white",
                      textAlign: "center",
                      fontWeight: "bold",
                      backgroundColor: "rgba(0, 0, 0, 0.3)",
                      paddingHorizontal: 4,
                      paddingVertical: 2,
                      borderRadius: 3,
                      lineHeight: 12,
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

      {/* Second circle - Child categories */}
      {selectedParentState && (
        <View style={styles.secondCircleWrapper}>
          <View style={styles.pointer2}>
            <Svg width={25} height={25} style={styles.pointerSvg}>
              <Polygon
                points="12.5,22 8,12 17,12"
                fill="#000"
                stroke="#fff"
                strokeWidth="2"
              />
            </Svg>
          </View>

          <View style={styles.halfCircle2} {...panResponder2.panHandlers}>
            <Animated.View
              style={{
                transform: [
                  {
                    rotate: rotation2.interpolate({
                      inputRange: [-360, 360],
                      outputRange: ["-360deg", "360deg"],
                    }),
                  },
                ],
              }}
            >
              <Svg width={radius2} height={radius2}>
                <G rotation={-90} originX={center2} originY={center2}>
                  {(selectedParentState as any).childrens.map(
                    (child: any, index: number) => {
                      const childSliceAngle =
                        360 / (selectedParentState as any).childrens.length;
                      const startAngle = index * childSliceAngle;
                      const endAngle = startAngle + childSliceAngle;
                      return (
                        <Path
                          key={child.id}
                          d={createSector(
                            startAngle,
                            endAngle,
                            center2,
                            center2
                          )}
                          fill={child.color}
                          stroke="white"
                          strokeWidth="2"
                        />
                      );
                    }
                  )}
                </G>
              </Svg>

              {(selectedParentState as any).childrens.map(
                (child: any, index: number) => {
                  const childSliceAngle =
                    360 / (selectedParentState as any).childrens.length;
                  const startAngle = index * childSliceAngle;
                  const endAngle = startAngle + childSliceAngle;
                  const midAngle = (startAngle + endAngle) / 2;
                  const adjustedAngle = (midAngle - 90) * (Math.PI / 180);
                  const distanceFromCenter = center2 * 0.85;
                  const x =
                    center2 + distanceFromCenter * Math.cos(adjustedAngle);
                  const y =
                    center2 + distanceFromCenter * Math.sin(adjustedAngle);

                  return (
                    <Animated.View
                      key={child.id + "-label"}
                      style={{
                        position: "absolute",
                        left: x - 50,
                        top: y - 20,
                        width: 100,
                        height: 40,
                        alignItems: "center",
                        justifyContent: "center",
                        transform: [
                          {
                            rotate: rotation2.interpolate({
                              inputRange: [-360, 360],
                              outputRange: ["360deg", "-360deg"],
                            }),
                          },
                        ],
                      }}
                    >
                      <Text
                        style={{
                          fontSize: 11,
                          color: "white",
                          textAlign: "center",
                          fontWeight: "bold",
                          backgroundColor: "rgba(0, 0, 0, 0.4)",
                          paddingHorizontal: 5,
                          paddingVertical: 3,
                          borderRadius: 4,
                          lineHeight: 13,
                        }}
                      >
                        {child.label.replace("/", "/\n")}
                      </Text>
                    </Animated.View>
                  );
                }
              )}
            </Animated.View>
          </View>
        </View>
      )}
      {/* Third circle - Grandchild categories */}
      {selectedChild &&
        (selectedChild as any).childrens &&
        (selectedChild as any).childrens.length > 0 && (
          <View style={styles.thirdCircleWrapper}>
            <View style={styles.pointer3}>
              <Svg width={30} height={30} style={styles.pointerSvg}>
                <Polygon
                  points="15,27 10,15 20,15"
                  fill="#000"
                  stroke="#fff"
                  strokeWidth="2"
                />
              </Svg>
            </View>

            <View style={styles.halfCircle3} {...panResponder3.panHandlers}>
              <Animated.View
                style={{
                  transform: [
                    {
                      rotate: rotation3.interpolate({
                        inputRange: [-360, 360],
                        outputRange: ["-360deg", "360deg"],
                      }),
                    },
                  ],
                }}
              >
                <Svg width={radius3} height={radius3}>
                  <G rotation={-90} originX={center3} originY={center3}>
                    {(selectedChild as any).childrens.map(
                      (grandchild: any, index: number) => {
                        const grandchildSliceAngle =
                          360 / (selectedChild as any).childrens.length;
                        const startAngle = index * grandchildSliceAngle;
                        const endAngle = startAngle + grandchildSliceAngle;
                        return (
                          <Path
                            key={grandchild.id}
                            d={createSector(
                              startAngle,
                              endAngle,
                              center3,
                              center3
                            )}
                            fill={grandchild.color}
                            stroke="white"
                            strokeWidth="2"
                          />
                        );
                      }
                    )}
                  </G>
                </Svg>

                {(selectedChild as any).childrens.map(
                  (grandchild: any, index: number) => {
                    const grandchildSliceAngle =
                      360 / (selectedChild as any).childrens.length;
                    const startAngle = index * grandchildSliceAngle;
                    const endAngle = startAngle + grandchildSliceAngle;
                    const midAngle = (startAngle + endAngle) / 2;
                    const adjustedAngle = (midAngle - 90) * (Math.PI / 180);
                    const distanceFromCenter = center3 * 0.9;
                    const x =
                      center3 + distanceFromCenter * Math.cos(adjustedAngle);
                    const y =
                      center3 + distanceFromCenter * Math.sin(adjustedAngle);

                    return (
                      <Animated.View
                        key={grandchild.id + "-label"}
                        style={{
                          position: "absolute",
                          left: x - 60,
                          top: y - 25,
                          width: 120,
                          height: 50,
                          alignItems: "center",
                          justifyContent: "center",
                          transform: [
                            {
                              rotate: rotation3.interpolate({
                                inputRange: [-360, 360],
                                outputRange: ["360deg", "-360deg"],
                              }),
                            },
                          ],
                        }}
                      >
                        <Text
                          style={{
                            fontSize: 12,
                            color: "white",
                            textAlign: "center",
                            fontWeight: "bold",
                            backgroundColor: "rgba(0, 0, 0, 0.5)",
                            paddingHorizontal: 6,
                            paddingVertical: 4,
                            borderRadius: 5,
                            lineHeight: 14,
                          }}
                        >
                          {grandchild.label.replace("/", "/\n")}
                        </Text>
                      </Animated.View>
                    );
                  }
                )}
              </Animated.View>
            </View>
          </View>
        )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "white",
    position: "relative",
  },
  firstCircleWrapper: {
    position: "absolute",
    bottom: 0,
    left: (screenWidth - radius) / 2, // Center horizontally
    width: radius,
    height: radius / 2,
    overflow: "hidden",
    alignItems: "center",
    zIndex: 30, // Highest z-index to appear in front
  },
  secondCircleWrapper: {
    position: "absolute",
    bottom: 0, // Also anchored to bottom
    left: (screenWidth - radius2) / 2, // Center horizontally
    width: radius2,
    height: radius2 / 2,
    overflow: "hidden",
    alignItems: "center",
    zIndex: 20, // Middle z-index
  },
  thirdCircleWrapper: {
    position: "absolute",
    bottom: 0, // Also anchored to bottom
    left: (screenWidth - radius3) / 2, // Center horizontally
    width: radius3,
    height: radius3 / 2,
    overflow: "hidden",
    alignItems: "center",
    zIndex: 10, // Lowest z-index to appear behind
  },
  halfCircle: {
    width: radius,
    height: radius,
    borderRadius: radius / 2,
    overflow: "hidden",
    alignItems: "center",
    justifyContent: "center",
  },
  halfCircle2: {
    width: radius2,
    height: radius2,
    borderRadius: radius2 / 2,
    overflow: "hidden",
    alignItems: "center",
    justifyContent: "center",
  },
  halfCircle3: {
    width: radius3,
    height: radius3,
    borderRadius: radius3 / 2,
    overflow: "hidden",
    alignItems: "center",
    justifyContent: "center",
  },
  pointer: {
    position: "absolute",
    top: -5,
    zIndex: 40, // Highest z-index for pointer
    alignItems: "center",
  },
  pointer2: {
    position: "absolute",
    top: -8,
    zIndex: 30, // Higher than second circle but lower than first pointer
    alignItems: "center",
  },
  pointer3: {
    position: "absolute",
    top: -10,
    zIndex: 20, // Higher than third circle but lower than second pointer
    alignItems: "center",
  },
  pointerSvg: {
    zIndex: 10,
  },
});
