import { itemsV2 } from "@/data/dataV2";
import { Sample } from "@/types/sample";
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
  console.log(rotationValue);
  let normalizedRotation = rotationValue % 360;
  if (normalizedRotation < 0) normalizedRotation += 360;

  const currentSliceAngle = 360 / itemsArray.length;
  let sliceIndex = Math.floor((360 - normalizedRotation) / currentSliceAngle);

  sliceIndex = sliceIndex % itemsArray.length;
  if (sliceIndex < 0) sliceIndex += itemsArray.length;

  return itemsArray[sliceIndex];
}

// Function to calculate initial rotation to position FRUITY under the pointer
function calculateInitialRotation() {
  // FRUITY is the first item (index 0) in itemsV2
  // Each slice is 360/9 = 40 degrees
  // The wheel has a -90 degree rotation applied in SVG
  // To center FRUITY under the pointer, we need to account for this offset
  // Since FRUITY starts at 0 degrees and we want it centered under the pointer (top),
  // we need to rotate by negative half slice angle to move it left to center
  const halfSliceAngle = sliceAngle / 2;
  return -halfSliceAngle;
}

// Function to calculate initial rotation for second circle to position BERRY under the pointer
function calculateInitialRotation2(selectedParent: Sample | null) {
  if (!selectedParent || !(selectedParent as any).childrens) return 0;

  // BERRY is the first child (index 0) of FRUITY
  const childrenCount = (selectedParent as any).childrens.length;
  const childSliceAngle = 360 / childrenCount;
  const halfChildSliceAngle = childSliceAngle / 2;

  // Similar to first circle, rotate by negative half slice angle to center BERRY
  return -halfChildSliceAngle;
}

// Function to calculate initial rotation for third circle to position RASPBERRY under the pointer
function calculateInitialRotation3(selectedChild: Sample | null) {
  if (!selectedChild || !(selectedChild as any).childrens) return 0;

  // RASPBERRY is the second child (index 1) of BERRY
  const grandchildrenCount = (selectedChild as any).childrens.length;
  const grandchildSliceAngle = 360 / grandchildrenCount;

  // Position RASPBERRY (index 1) under the pointer
  // First slice starts at 0, second slice starts at grandchildSliceAngle
  // To center RASPBERRY, we need to rotate by -(grandchildSliceAngle + halfSliceAngle)
  const raspberryStartAngle = 1 * grandchildSliceAngle;
  const halfGrandchildSliceAngle = grandchildSliceAngle / 2;
  const raspberryCenterAngle = raspberryStartAngle + halfGrandchildSliceAngle;

  return -raspberryCenterAngle;
}

export default function CuppingRegistrationMinimalist() {
  // Calculate initial rotation to position FRUITY under the pointer
  const initialRotationValue = calculateInitialRotation();

  const rotation = useRef(new Animated.Value(initialRotationValue)).current;
  const lastRotation = useRef(initialRotationValue);
  const rotation2 = useRef(new Animated.Value(0)).current;
  const lastRotation2 = useRef(0);
  const rotation3 = useRef(new Animated.Value(0)).current;
  const lastRotation3 = useRef(0);

  // State for tracking spinning
  const [isSpinning, setIsSpinning] = useState(false);
  const [isSpinning2, setIsSpinning2] = useState(false);
  const [isSpinning3, setIsSpinning3] = useState(false);

  // Selection states - Initialize with FRUITY, BERRY, and RASPBERRY
  const [selectedParent, setSelectedParentState] = useState<Sample | null>(
    itemsV2[0] // FRUITY is at index 0
  );
  const [selectedChild, setSelectedChildState] = useState<Sample | null>(
    itemsV2[0] && (itemsV2[0] as any).childrens
      ? (itemsV2[0] as any).childrens[0]
      : null // BERRY is at index 0 of FRUITY's children
  );
  const [selectedChild2, setSelectedChild2] = useState<Sample | null>(
    itemsV2[0] &&
      (itemsV2[0] as any).childrens &&
      (itemsV2[0] as any).childrens[0] &&
      ((itemsV2[0] as any).childrens[0] as any).childrens
      ? ((itemsV2[0] as any).childrens[0] as any).childrens[1] // RASPBERRY is at index 1 of BERRY's children
      : null
  );

  // Saved selections for tags
  const [savedSelections, setSavedSelections] = useState<
    Array<{
      parent: Sample | null;
      child1: Sample | null;
      child2: Sample | null;
      id: string;
    }>
  >([]);

  // Refs to store current values for pan responders - Initialize with FRUITY, BERRY, and RASPBERRY
  const selectedParentRef = useRef<Sample | null>(itemsV2[0]);
  const selectedChildRef = useRef<Sample | null>(
    itemsV2[0] && (itemsV2[0] as any).childrens
      ? (itemsV2[0] as any).childrens[0]
      : null
  );

  // Wrapper function for selectedParent
  const setSelectedParent = (value: Sample | null) => {
    selectedParentRef.current = value;
    setSelectedParentState(value);
  };

  // Wrapper function for selectedChild
  const setSelectedChild = (value: Sample | null) => {
    selectedChildRef.current = value;
    setSelectedChildState(value);
    setSelectedChild2(null); // Reset third level selection when second level changes
  };

  // Set initial selection and rotation for second and third circles on component mount
  useEffect(() => {
    console.log("Initial selection set to:", itemsV2[0].label);
    setSelectedParent(itemsV2[0]); // FRUITY

    // Set initial selection for second circle (BERRY)
    if (itemsV2[0] && (itemsV2[0] as any).childrens) {
      const berry = (itemsV2[0] as any).childrens[0]; // BERRY
      setSelectedChild(berry);
      console.log("Initial child selection set to:", berry.label);

      // Set initial rotation for second circle to position BERRY under the pointer
      const initialRotation2Value = calculateInitialRotation2(itemsV2[0]);
      rotation2.setValue(initialRotation2Value);
      lastRotation2.current = initialRotation2Value;

      // Set initial selection for third circle (RASPBERRY)
      if (
        berry &&
        (berry as any).childrens &&
        (berry as any).childrens.length > 1
      ) {
        const raspberry = (berry as any).childrens[1]; // RASPBERRY is at index 1
        setSelectedChild2(raspberry);
        console.log("Initial grandchild selection set to:", raspberry.label);

        // Set initial rotation for third circle to position RASPBERRY under the pointer
        const initialRotation3Value = calculateInitialRotation3(berry);
        rotation3.setValue(initialRotation3Value);
        lastRotation3.current = initialRotation3Value;
      }
    }
  }, []);

  // Function to save current selections
  const saveCurrentSelection = () => {
    if (selectedParent) {
      const newSelection = {
        parent: selectedParent,
        child1: selectedChild,
        child2: selectedChild2,
        id: Date.now().toString(), // Unique ID for each selection
      };

      setSavedSelections((prev) => [...prev, newSelection]);

      console.log("Saved selection:", {
        parent: selectedParent?.label,
        child1: selectedChild?.label,
        child2: selectedChild2?.label,
      });
    }
  };

  // Function to remove a saved selection
  const removeSavedSelection = (id: string) => {
    setSavedSelections((prev) =>
      prev.filter((selection) => selection.id !== id)
    );
  };

  const panResponder = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => !isSpinning2 && !isSpinning3,
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
          const selectedSlice = getSelectedSlice(finalRotation, itemsV2);
          if (selectedSlice && selectedSlice.id) {
            console.log("Selected parent ID:", selectedSlice.id);
            console.log("Selected parent label:", selectedSlice.label);
            setSelectedParent(selectedSlice);
            setSelectedChild(null);

            // Reset second and third circle rotations when parent changes
            rotation2.setValue(0);
            lastRotation2.current = 0;
            rotation3.setValue(0);
            lastRotation3.current = 0;

            // If new parent has children, set first child as default and position it
            if (
              (selectedSlice as any).childrens &&
              (selectedSlice as any).childrens.length > 0
            ) {
              const firstChild = (selectedSlice as any).childrens[0];
              setSelectedChild(firstChild);

              // Set rotation to position first child under pointer
              const childInitialRotation =
                calculateInitialRotation2(selectedSlice);
              rotation2.setValue(childInitialRotation);
              lastRotation2.current = childInitialRotation;

              // If first child has grandchildren, set appropriate default and position it
              if (
                (firstChild as any).childrens &&
                (firstChild as any).childrens.length > 0
              ) {
                // For BERRY, use RASPBERRY (index 1), for others use first grandchild (index 0)
                const defaultGrandchildIndex =
                  firstChild.label === "BERRY" ? 1 : 0;
                const defaultGrandchild =
                  (firstChild as any).childrens[defaultGrandchildIndex] ||
                  (firstChild as any).childrens[0];
                setSelectedChild2(defaultGrandchild);

                // Set rotation to position default grandchild under pointer
                const grandchildInitialRotation =
                  firstChild.label === "BERRY"
                    ? calculateInitialRotation3(firstChild)
                    : calculateInitialRotation2(firstChild); // Use same logic for other categories
                rotation3.setValue(grandchildInitialRotation);
                lastRotation3.current = grandchildInitialRotation;
              }
            }
          }
        } catch (error) {
          console.log("Error getting selected slice:", error);
        }
      },
    })
  ).current;

  const panResponder2 = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => !isSpinning && !isSpinning3,
      onPanResponderGrant: () => setIsSpinning2(true),
      onPanResponderMove: (_, gesture) => {
        const newRotation = lastRotation2.current + gesture.dx / 2;
        rotation2.setValue(newRotation);
      },
      onPanResponderRelease: (_, gesture) => {
        setIsSpinning2(false);
        const finalRotation = lastRotation2.current + gesture.dx / 2;
        lastRotation2.current = finalRotation;

        if (
          selectedParentRef.current &&
          (selectedParentRef.current as any).childrens
        ) {
          try {
            const selectedSlice = getSelectedSlice(
              finalRotation,
              (selectedParentRef.current as any).childrens
            );

            if (selectedSlice && selectedSlice.id) {
              console.log("Selected child ID:", selectedSlice.id);
              console.log("Selected child label:", selectedSlice.label);
              setSelectedChild(selectedSlice);

              // Reset third circle rotation when child changes
              rotation3.setValue(0);
              lastRotation3.current = 0;

              // If new child has grandchildren, set appropriate default and position it
              if (
                (selectedSlice as any).childrens &&
                (selectedSlice as any).childrens.length > 0
              ) {
                // For BERRY, use RASPBERRY (index 1), for others use first grandchild (index 0)
                const defaultGrandchildIndex =
                  selectedSlice.label === "BERRY" ? 1 : 0;
                const defaultGrandchild =
                  (selectedSlice as any).childrens[defaultGrandchildIndex] ||
                  (selectedSlice as any).childrens[0];
                setSelectedChild2(defaultGrandchild);

                // Set rotation to position default grandchild under pointer
                const grandchildInitialRotation =
                  selectedSlice.label === "BERRY"
                    ? calculateInitialRotation3(selectedSlice)
                    : calculateInitialRotation2(selectedSlice); // Use same logic for other categories
                rotation3.setValue(grandchildInitialRotation);
                lastRotation3.current = grandchildInitialRotation;
              }
            }
          } catch (error) {
            console.log("Error getting selected child slice:", error);
          }
        }
      },
    })
  ).current;

  const panResponder3 = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => !isSpinning && !isSpinning2,
      onPanResponderGrant: () => setIsSpinning3(true),
      onPanResponderMove: (_, gesture) => {
        const newRotation = lastRotation3.current + gesture.dx / 2;
        rotation3.setValue(newRotation);
      },
      onPanResponderRelease: (_, gesture) => {
        setIsSpinning3(false);
        const finalRotation = lastRotation3.current + gesture.dx / 2;
        lastRotation3.current = finalRotation;

        if (
          selectedChildRef.current &&
          (selectedChildRef.current as any).childrens
        ) {
          try {
            const selectedSlice = getSelectedSlice(
              finalRotation,
              (selectedChildRef.current as any).childrens
            );

            if (selectedSlice && selectedSlice.id) {
              console.log("Selected child2 ID:", selectedSlice.id);
              console.log("Selected child2 label:", selectedSlice.label);
              setSelectedChild2(selectedSlice);
            }
          } catch (error) {
            console.log("Error getting selected child2 slice:", error);
          }
        }
      },
    })
  ).current;

  return (
    <View style={styles.container}>
      {/* Saved selections tags at the top */}
      <View style={styles.tagsContainer}>
        {savedSelections.map((selection) => (
          <View key={selection.id} style={styles.tag}>
            <Text style={styles.tagText}>
              {selection.parent?.label}
              {selection.child1 && ` → ${selection.child1.label}`}
              {selection.child2 && ` → ${selection.child2.label}`}
            </Text>
            <TouchableOpacity
              onPress={() => removeSavedSelection(selection.id)}
              style={styles.removeButton}
            >
              <Text style={styles.removeButtonText}>×</Text>
            </TouchableOpacity>
          </View>
        ))}
      </View>

      {/* Display current selection for debugging */}
      <View style={styles.selectionDisplay}>
        <Text style={styles.selectionText}>
          Selected: {selectedParent?.label || "None"}
        </Text>
      </View>

      {/* Save button in the center */}
      <TouchableOpacity
        style={[
          styles.saveButton,
          !selectedParent && styles.saveButtonDisabled,
        ]}
        onPress={saveCurrentSelection}
        disabled={!selectedParent}
      >
        <Text
          style={[
            styles.saveButtonText,
            styles.saveButtonPosition,
            !selectedParent && styles.saveButtonTextDisabled,
          ]}
        >
          SAVE
        </Text>
      </TouchableOpacity>

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
                {itemsV2.map((item, index) => {
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

            {itemsV2.map((item, index) => {
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
      {selectedParent && (
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
                  {(selectedParent as any).childrens.map(
                    (child: any, index: number) => {
                      const childSliceAngle =
                        360 / (selectedParent as any).childrens.length;
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

              {(selectedParent as any).childrens.map(
                (child: any, index: number) => {
                  const childSliceAngle =
                    360 / (selectedParent as any).childrens.length;
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
    fontSize: 16,
    fontWeight: "bold",
    textAlign: "center",
  },
  tagsContainer: {
    position: "absolute",
    top: 100,
    left: 20,
    right: 20,
    zIndex: 50,
    flexDirection: "row",
    flexWrap: "wrap",
    gap: 10,
  },
  tag: {
    backgroundColor: "#f0f0f0",
    borderRadius: 20,
    paddingHorizontal: 15,
    paddingVertical: 8,
    flexDirection: "row",
    alignItems: "center",
    shadowColor: "#000",
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.1,
    shadowRadius: 3,
    elevation: 3,
  },
  tagText: {
    fontSize: 12,
    color: "#333",
    fontWeight: "500",
    marginRight: 8,
  },
  removeButton: {
    width: 18,
    height: 18,
    backgroundColor: "#ff4444",
    borderRadius: 9,
    alignItems: "center",
    justifyContent: "center",
  },
  removeButtonText: {
    color: "white",
    fontSize: 12,
    fontWeight: "bold",
    lineHeight: 12,
  },
  saveButton: {
    position: "absolute",
    bottom: -50, // Fixed position at bottom of screen
    left: "50%",
    marginLeft: -75, // Half of button width to center it
    width: 150,
    height: 150,
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
    top: 55,
  },
  saveButtonDisabled: {
    backgroundColor: "#cccccc",
    shadowOpacity: 0.1,
  },
  saveButtonText: {
    color: "black",
    fontSize: 18,
    fontWeight: "bold",
  },
  saveButtonTextDisabled: {
    color: "#999999",
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
