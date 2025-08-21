import React, { useRef, useState, useEffect } from "react";
import {
  View,
  Dimensions,
  PanResponder,
  Animated,
  StyleSheet,
  Text,
  TouchableOpacity,
} from "react-native";
import Svg, { Path, Polygon, G } from "react-native-svg";
import { firstCircle, secondCircle } from "@/data/dataV1";
import { Text as SvgText } from "react-native-svg";

function polarToCartesian(
  centerX: number,
  centerY: number,
  radius: number,
  angleInDegrees: number
) {
  const angleInRadians = ((angleInDegrees - 90) * Math.PI) / 180.0;
  return {
    x: centerX + radius * Math.cos(angleInRadians),
    y: centerY + radius * Math.sin(angleInRadians),
  };
}

function getSelectedSlice(
  rotationValue: number,
  data: Array<{
    id: string;
    label: string;
    color: string;
    numbers: number;
    parentId?: string;
  }>,
  total: number
) {
  let normalizedRotation = rotationValue % 360;
  if (normalizedRotation < 0) normalizedRotation += 360;

  normalizedRotation = (360 - normalizedRotation) % 360;
  normalizedRotation = (normalizedRotation + 90) % 360;

  let cumulativeAngle = 0;

  for (let i = 0; i < data.length; i++) {
    const item = data[i];
    const sliceAngle = (item.numbers / total) * 360;
    const endAngle = cumulativeAngle + sliceAngle;

    if (
      normalizedRotation >= cumulativeAngle &&
      normalizedRotation < endAngle
    ) {
      console.log(`Selected slice: ${item.label}`);
      return item;
    }

    cumulativeAngle = endAngle;
  }

  return data[0];
}

function describeArc(
  x: number,
  y: number,
  radius: number,
  startAngle: number,
  endAngle: number
) {
  const start = polarToCartesian(x, y, radius, endAngle);
  const end = polarToCartesian(x, y, radius, startAngle);
  const largeArcFlag = endAngle - startAngle <= 180 ? "0" : "1";
  const d = [
    `M ${x} ${y}`,
    `L ${start.x} ${start.y}`,
    `A ${radius} ${radius} 0 ${largeArcFlag} 0 ${end.x} ${end.y}`,
    "Z",
  ].join(" ");
  return d;
}

// Function to darken a color by reducing brightness
function darkenColor(color: string, factor: number = 0.4): string {
  // Remove # if present
  const hex = color.replace("#", "");

  // Parse RGB values
  const r = parseInt(hex.substr(0, 2), 16);
  const g = parseInt(hex.substr(2, 2), 16);
  const b = parseInt(hex.substr(4, 2), 16);

  // Darken by reducing each component
  const newR = Math.floor(r * factor);
  const newG = Math.floor(g * factor);
  const newB = Math.floor(b * factor);

  // Convert back to hex
  const newHex =
    "#" +
    newR.toString(16).padStart(2, "0") +
    newG.toString(16).padStart(2, "0") +
    newB.toString(16).padStart(2, "0");

  return newHex;
}

// Function to calculate initial rotation to position FRUITY under the pointer
function calculateInitialRotation(
  data: Array<{ id: string; label: string; color: string; numbers: number }>,
  total: number
) {
  const fruityIndex = 0;
  const fruitySliceAngle = (data[fruityIndex].numbers / total) * 360;
  const fruityCenterAngle = fruitySliceAngle / 2;
  return fruityCenterAngle;
}

// Function to calculate initial rotation for second circle
function calculateInitialRotation2(
  data: Array<{ id: string; label: string; color: string; numbers: number }>,
  total: number
) {
  // BERRY is the first item in secondCircle data
  const berryIndex = 0;
  const berrySliceAngle = (data[berryIndex].numbers / total) * 360;
  const berryCenterAngle = berrySliceAngle / 2;
  return berryCenterAngle;
}

// FIXED FUNCTION: Calculate rotation needed to center a specific parent item
function calculateRotationForParent(
  parentId: string,
  data: Array<{ id: string; label: string; color: string; numbers: number }>,
  total: number
) {
  const parentIndex = data.findIndex((item) => item.id === parentId);
  if (parentIndex === -1) return 0;

  // Calculate cumulative angle up to this parent
  let cumulativeAngle = 0;
  for (let i = 0; i < parentIndex; i++) {
    cumulativeAngle += (data[i].numbers / total) * 360;
  }

  // Add half of this parent's slice angle to get center
  const parentSliceAngle = (data[parentIndex].numbers / total) * 360;
  const parentCenterAngle = cumulativeAngle + parentSliceAngle / 2;

  // Apply the same transformations as in getSelectedSlice but in reverse
  // We want the normalized rotation after transformations to equal parentCenterAngle
  // Working backwards from getSelectedSlice transformations:
  // 1. normalizedRotation = (normalizedRotation + 90) % 360 -> subtract 90
  // 2. normalizedRotation = (360 - normalizedRotation) % 360 -> apply 360 - x

  // Step 1: Subtract 90 degrees
  let targetRotation = parentCenterAngle - 90;
  if (targetRotation < 0) targetRotation += 360;

  // Step 2: Reverse the direction flip
  targetRotation = (360 - targetRotation) % 360;

  return targetRotation;
}

const { width: screenWidth } = Dimensions.get("window");

// First circle (parent categories)
const radius = screenWidth * 1.6;
const center = radius / 2;
const data = firstCircle.data;
const total = firstCircle.numbers;

// Second circle (child categories) - larger and behind first circle
const radius2 = screenWidth * 2.5;
const center2 = radius2 / 2;
const data2 = secondCircle.data;
const total2 = secondCircle.numbers;

export default function CuppingRegistrationOverview() {
  // Calculate initial rotations
  const initialRotationValue = calculateInitialRotation(data, total);
  const initialRotationValue2 = calculateInitialRotation2(data2, total2);

  const rotation = useRef(new Animated.Value(initialRotationValue)).current;
  const lastRotation = useRef(initialRotationValue);

  const rotation2 = useRef(new Animated.Value(initialRotationValue2)).current;
  const lastRotation2 = useRef(initialRotationValue2);

  const [isSpinning, setIsSpinning] = useState(false);
  const [isSpinning2, setIsSpinning2] = useState(false);

  // Initialize with FRUITY and BERRY selected
  const [selectedParent, setSelectedParentState] = useState<any | null>(
    data[0]
  );
  const [selectedChild, setSelectedChildState] = useState<any | null>(data2[0]);

  const selectedParentRef = useRef<any | null>(data[0]);
  const selectedChildRef = useRef<any | null>(data2[0]);

  const setSelectedParent = (value: any | null) => {
    selectedParentRef.current = value;
    setSelectedParentState(value);
  };

  const setSelectedChild = (value: any | null) => {
    selectedChildRef.current = value;
    setSelectedChildState(value);
  };

  // Set initial selection on component mount
  useEffect(() => {
    console.log("Initial parent selection set to:", data[0].label);
    console.log("Initial child selection set to:", data2[0].label);
    setSelectedParent(data[0]); // FRUITY
    setSelectedChild(data2[0]); // BERRY
  }, []);

  const panResponder = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => !isSpinning2,
      onPanResponderGrant: () => setIsSpinning(true),
      onPanResponderMove: (_, gesture) => {
        const newRotation = lastRotation.current + gesture.dx / 2;
        rotation.setValue(newRotation);
      },
      onPanResponderRelease: (_, gesture) => {
        setIsSpinning(false);
        const finalRotation = lastRotation.current + gesture.dx / 2;
        lastRotation.current = finalRotation;

        try {
          const selectedSlice = getSelectedSlice(finalRotation, data, total);
          if (selectedSlice && selectedSlice.id) {
            console.log("Selected parent ID:", selectedSlice.id);
            console.log("Selected parent label:", selectedSlice.label);
            setSelectedParent(selectedSlice);
          }
        } catch (error) {
          console.log("Error getting selected slice:", error);
        }
      },
    })
  ).current;

  const panResponder2 = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => !isSpinning,
      onPanResponderGrant: () => setIsSpinning2(true),
      onPanResponderMove: (_, gesture) => {
        const newRotation = lastRotation2.current + gesture.dx / 2;
        rotation2.setValue(newRotation);
      },
      onPanResponderRelease: (_, gesture) => {
        setIsSpinning2(false);
        const finalRotation = lastRotation2.current + gesture.dx / 2;
        lastRotation2.current = finalRotation;

        try {
          const selectedSlice = getSelectedSlice(finalRotation, data2, total2);
          if (selectedSlice && selectedSlice.id) {
            console.log("Selected child ID:", selectedSlice.id);
            console.log("Selected child label:", selectedSlice.label);

            // Check if the selected child belongs to a different parent
            const currentParentId = selectedParentRef.current?.id;
            const childParentId = selectedSlice.parentId;

            if (childParentId && currentParentId !== childParentId) {
              // Find the new parent in the first circle data
              const newParent = data.find((item) => item.id === childParentId);

              if (newParent) {
                console.log("Switching to new parent:", newParent.label);

                // Calculate the rotation needed to center the new parent (FIXED)
                const targetRotation = calculateRotationForParent(
                  childParentId,
                  data,
                  total
                );

                console.log(
                  `Rotating to ${targetRotation} degrees for parent ${newParent.label}`
                );

                // Animate the first circle to the new position
                Animated.timing(rotation, {
                  toValue: targetRotation,
                  duration: 500, // 500ms animation
                  useNativeDriver: false,
                }).start(() => {
                  // Verify the rotation worked correctly
                  console.log("Animation completed, verifying selection...");
                  const verifySlice = getSelectedSlice(
                    targetRotation,
                    data,
                    total
                  );
                  console.log("Verification result:", verifySlice?.label);
                });

                // Update the lastRotation ref and selectedParent
                lastRotation.current = targetRotation;
                setSelectedParent(newParent);
              }
            }

            setSelectedChild(selectedSlice);
          }
        } catch (error) {
          console.log("Error getting selected child slice:", error);
        }
      },
    })
  ).current;

  return (
    <View style={styles.container}>
      <TouchableOpacity style={[styles.saveButton]}>
        <Text style={[styles.saveButtonText, styles.saveButtonPosition]}>
          SAVE
        </Text>
      </TouchableOpacity>

      {/* Display current selection for debugging */}
      <View style={styles.selectionDisplay}>
        <Text style={styles.selectionText}>
          Parent: {selectedParent?.label || "None"}
        </Text>
        <Text style={styles.selectionText}>
          Child: {selectedChild?.label || "None"}
        </Text>
      </View>

      {/* Second circle - Child categories (larger, behind first circle) */}
      <View style={styles.secondCircleWrapper}>
        <View style={styles.pointer2}>
          <Svg width={25} height={25} style={styles.pointerSvg}>
            <Polygon
              points="12.5,22 8,12 17,12"
              fill="#333"
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
                {(() => {
                  let cumulativeAngle = 0;
                  return data2.map((item, idx) => {
                    const sliceAngle = (item.numbers / total2) * 360;
                    const endAngle = cumulativeAngle + sliceAngle;
                    const path = describeArc(
                      center2,
                      center2,
                      radius2 / 2 - 10,
                      cumulativeAngle,
                      endAngle
                    );
                    const midAngle = (cumulativeAngle + endAngle) / 2;

                    // Position text closer to inner radius (top of slice)
                    const innerRadius = radius / 2 + 145; // Start from first circle's outer edge plus padding
                    const labelRadius = innerRadius + 15; // Small offset from inner edge

                    const rad = ((midAngle - 90) * Math.PI) / 180;
                    const x = center2 + labelRadius * Math.cos(rad);
                    const y = center2 + labelRadius * Math.sin(rad);

                    const textRotation = midAngle + 90;

                    // Check if this item belongs to the selected parent
                    const isChildOfSelectedParent =
                      selectedParent && item.parentId === selectedParent.id;

                    // Determine the fill color and opacity based on parent relationship
                    const fillColor = isChildOfSelectedParent
                      ? item.color
                      : darkenColor(item.color, 0.3);
                    const fillOpacity = isChildOfSelectedParent ? 1 : 0.4;
                    const textOpacity = isChildOfSelectedParent ? 1 : 0.5;

                    const arc = (
                      <React.Fragment key={item.id}>
                        <Path
                          d={path}
                          fill={fillColor}
                          fillOpacity={fillOpacity}
                          stroke="#fff"
                          strokeWidth={2}
                          strokeOpacity={isChildOfSelectedParent ? 1 : 0.5}
                        />
                        <SvgText
                          x={x}
                          y={y}
                          fill="#fff"
                          fillOpacity={textOpacity}
                          fontSize="14"
                          fontWeight={
                            isChildOfSelectedParent ? "bold" : "normal"
                          }
                          textAnchor="start"
                          alignmentBaseline="middle"
                          transform={`rotate(${textRotation}, ${x}, ${y})`}
                        >
                          {item.label.replace("/", "\n")}
                        </SvgText>
                      </React.Fragment>
                    );
                    cumulativeAngle = endAngle;
                    return arc;
                  });
                })()}
              </G>
            </Svg>
          </Animated.View>
        </View>
      </View>

      {/* First circle - Parent categories (in front) */}
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
                {(() => {
                  let cumulativeAngle = 0;
                  return data.map((item, idx) => {
                    const sliceAngle = (item.numbers / total) * 360;
                    const endAngle = cumulativeAngle + sliceAngle;
                    const path = describeArc(
                      center,
                      center,
                      radius / 2 - 10,
                      cumulativeAngle,
                      endAngle
                    );

                    const arc = (
                      <Path
                        key={item.id}
                        d={path}
                        fill={item.color}
                        stroke="#fff"
                        strokeWidth={2}
                      />
                    );
                    cumulativeAngle = endAngle;
                    return arc;
                  });
                })()}
              </G>
            </Svg>

            {(() => {
              let angleOffset = 0;
              return data.map((item, idx) => {
                const itemAngle = (item.numbers / total) * 360;
                const centerAngle = angleOffset + itemAngle / 2;

                const radians = ((centerAngle - 90 - 90) * Math.PI) / 180;
                const distance = (radius / 2 - 10) * 0.85;
                const x = center + distance * Math.cos(radians);
                const y = center + distance * Math.sin(radians);

                const label = (
                  <Animated.View
                    key={item.id + "-label"}
                    style={{
                      position: "absolute",
                      left: x - 60,
                      top: y - 20,
                      width: 120,
                      height: 40,
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
                        fontSize: 12,
                        color: "white",
                        textAlign: "center",
                        fontWeight: "bold",
                        backgroundColor: "rgba(0, 0, 0, 0.25)",
                        paddingHorizontal: 6,
                        paddingVertical: 3,
                        borderRadius: 4,
                        lineHeight: 16,
                      }}
                    >
                      {item.label.replace("/", "/\n")}
                    </Text>
                  </Animated.View>
                );

                angleOffset += itemAngle;
                return label;
              });
            })()}
          </Animated.View>
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#cbcfce",
    position: "relative",
  },
  selectionDisplay: {
    position: "absolute",
    top: 50,
    left: 20,
    right: 20,
    backgroundColor: "rgba(0, 0, 0, 0.8)",
    padding: 10,
    borderRadius: 8,
    zIndex: 50,
  },
  selectionText: {
    color: "white",
    fontSize: 14,
    fontWeight: "bold",
    textAlign: "center",
  },
  firstCircleWrapper: {
    position: "absolute",
    bottom: -radius * 0.15,
    left: (screenWidth - radius) / 2,
    width: radius,
    height: radius / 2,
    overflow: "hidden",
    alignItems: "center",
    zIndex: 30, // Higher z-index to appear in front
  },
  secondCircleWrapper: {
    position: "absolute",
    bottom: -radius2 * 0.095,
    left: (screenWidth - radius2) / 2,
    width: radius2,
    height: radius2 / 2,
    overflow: "hidden",
    alignItems: "center",
    zIndex: 20, // Lower z-index to appear behind first circle
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
  pointer: {
    position: "absolute",
    top: -5,
    zIndex: 40,
    alignItems: "center",
  },
  pointer2: {
    position: "absolute",
    top: -8,
    zIndex: 30,
    alignItems: "center",
  },
  pointerSvg: {
    zIndex: 10,
  },
  saveButton: {
    position: "absolute",
    bottom: 18,
    left: "50%",
    marginLeft: -85,
    width: 170,
    height: 80,
    backgroundColor: "#ffffff",
    borderRadius: 100,
    alignItems: "center",
    justifyContent: "center",
    zIndex: 40,
    shadowColor: "#000",
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.25,
    shadowRadius: 4,
    elevation: 5,
  },
  saveButtonPosition: {
    position: "absolute",
    top: 30,
  },
  saveButtonText: {
    color: "black",
    fontSize: 18,
    fontWeight: "bold",
  },
});
