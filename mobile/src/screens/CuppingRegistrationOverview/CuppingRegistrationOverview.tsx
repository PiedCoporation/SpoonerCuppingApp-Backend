import React, { useRef, useState, useEffect, useCallback } from "react";
import {
  View,
  Dimensions,
  PanResponder,
  Animated,
  StyleSheet,
  Text,
  TouchableOpacity,
  Easing,
} from "react-native";
import Svg, { Path, Polygon, G } from "react-native-svg";
import { firstCircle, secondCircle } from "@/data/dataV1";
import { Text as SvgText } from "react-native-svg";
import { FlavorItem } from "@/types/Flavor";

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

function darkenColor(color: string, factor: number = 0.4): string {
  const hex = color.replace("#", "");
  const r = parseInt(hex.substr(0, 2), 16);
  const g = parseInt(hex.substr(2, 2), 16);
  const b = parseInt(hex.substr(4, 2), 16);
  const newR = Math.floor(r * factor);
  const newG = Math.floor(g * factor);
  const newB = Math.floor(b * factor);
  const newHex =
    "#" +
    newR.toString(16).padStart(2, "0") +
    newG.toString(16).padStart(2, "0") +
    newB.toString(16).padStart(2, "0");
  return newHex;
}

function calculateInitialRotation(
  data: Array<{ id: string; label: string; color: string; numbers: number }>,
  total: number
) {
  const fruityIndex = 0;
  const fruitySliceAngle = (data[fruityIndex].numbers / total) * 360;
  const fruityCenterAngle = fruitySliceAngle / 2;
  return fruityCenterAngle;
}

function calculateInitialRotation2(
  data: Array<{ id: string; label: string; color: string; numbers: number }>,
  total: number
) {
  const berryIndex = 0;
  const berrySliceAngle = (data[berryIndex].numbers / total) * 360;
  const berryCenterAngle = berrySliceAngle / 2;
  return berryCenterAngle;
}

function calculateRotationForParent(
  parentId: string,
  data: Array<{ id: string; label: string; color: string; numbers: number }>,
  total: number
) {
  const parentIndex = data.findIndex((item) => item.id === parentId);
  if (parentIndex === -1) return 0;

  let cumulativeAngle = 0;
  for (let i = 0; i < parentIndex; i++) {
    cumulativeAngle += (data[i].numbers / total) * 360;
  }

  const parentSliceAngle = (data[parentIndex].numbers / total) * 360;
  const parentCenterAngle = cumulativeAngle + parentSliceAngle / 2;

  let targetRotation = parentCenterAngle - 90;
  if (targetRotation < 0) targetRotation += 360;
  targetRotation = (360 - targetRotation) % 360;

  return targetRotation;
}

function calculateRotationForChild(
  parentId: string,
  data: Array<{
    id: string;
    label: string;
    color: string;
    numbers: number;
    parentId?: string;
  }>,
  total: number
) {
  const firstChildIndex = data.findIndex((item) => item.parentId === parentId);
  if (firstChildIndex === -1) return 0;

  let cumulativeAngle = 0;
  for (let i = 0; i < firstChildIndex; i++) {
    cumulativeAngle += (data[i].numbers / total) * 360;
  }

  const childSliceAngle = (data[firstChildIndex].numbers / total) * 360;
  const childCenterAngle = cumulativeAngle + childSliceAngle / 2;

  let targetRotation = childCenterAngle - 90;
  if (targetRotation < 0) targetRotation += 360;
  targetRotation = (360 - targetRotation) % 360;

  return targetRotation;
}

const { width: screenWidth } = Dimensions.get("window");
const radius = screenWidth * 1.6;
const center = radius / 2;
const data = firstCircle.data;
const total = firstCircle.numbers;
const radius2 = screenWidth * 2.5;
const center2 = radius2 / 2;

export default function CuppingRegistrationOverview() {
  const [data2, setData2] = useState<FlavorItem[]>(secondCircle.data);
  const [total2, setTotal2] = useState(secondCircle.numbers);

  const stateRef = useRef({
    data2: secondCircle.data,
    total2: secondCircle.numbers,
    selectedParent: data[0],
    selectedChild: secondCircle.data[0],
    isUpdating: false,
  });

  useEffect(() => {
    stateRef.current.data2 = data2;
    stateRef.current.total2 = total2;
  }, [data2, total2]);

  const initialRotationValue = calculateInitialRotation(data, total);
  const initialRotationValue2 = calculateInitialRotation2(data2, total2);

  const rotation = useRef(new Animated.Value(initialRotationValue)).current;
  const lastRotation = useRef(initialRotationValue);

  const rotation2 = useRef(new Animated.Value(initialRotationValue2)).current;
  const lastRotation2 = useRef(initialRotationValue2);

  const [isSpinning, setIsSpinning] = useState(false);
  const [isSpinning2, setIsSpinning2] = useState(false);

  const [selectedParent, setSelectedParentState] = useState<any | null>(
    data[0]
  );
  const [selectedChild, setSelectedChildState] = useState<any | null>(data2[0]);

  const setSelectedParent = useCallback((value: any | null) => {
    stateRef.current.selectedParent = value;
    setSelectedParentState(value);
  }, []);

  const setSelectedChild = useCallback((value: any | null) => {
    stateRef.current.selectedChild = value;
    setSelectedChildState(value);
  }, []);

  const blindItemTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // FIX: Simplified BLIND item management without rotation adjustment
  const addBlindItem = useCallback((parentId: string) => {
    if (stateRef.current.isUpdating) return;

    if (blindItemTimeoutRef.current) {
      clearTimeout(blindItemTimeoutRef.current);
    }

    blindItemTimeoutRef.current = setTimeout(() => {
      stateRef.current.isUpdating = true;

      const blindItem: FlavorItem = {
        id: "0_0_0",
        label: "BLIND",
        color: "#cbcfce",
        numbers: 1,
        parentId: parentId,
      };

      setData2((prevData) => {
        const dataWithoutBlind = prevData.filter((item) => item.id !== "0_0_0");
        const insertionIndex = dataWithoutBlind.findIndex(
          (item) => item.parentId === parentId
        );
        const finalInsertionIndex =
          insertionIndex === -1 ? dataWithoutBlind.length : insertionIndex;

        const newData = [
          ...dataWithoutBlind.slice(0, finalInsertionIndex),
          blindItem,
          ...dataWithoutBlind.slice(finalInsertionIndex),
        ];

        // REMOVED: No automatic rotation adjustment
        // The circle will stay exactly where the user positioned it

        return newData;
      });

      setTotal2((prevTotal) => {
        const hasExistingBlind = stateRef.current.data2.some(
          (item) => item.id === "0_0_0"
        );
        return hasExistingBlind ? prevTotal : prevTotal + 1;
      });

      setTimeout(() => {
        stateRef.current.isUpdating = false;
      }, 100);
    }, 50);
  }, []);

  const removeBlindItem = useCallback(() => {
    const blindItem = stateRef.current.data2.find(
      (item) => item.id === "0_0_0"
    );
    if (blindItem) {
      setData2((prevData) => prevData.filter((item) => item.id !== "0_0_0"));
      setTotal2((prevTotal) => prevTotal - 1);
    }
  }, []);

  useEffect(() => {
    setSelectedParent(data[0]);
    setSelectedChild(data2[0]);
  }, [setSelectedParent, setSelectedChild]);

  useEffect(() => {
    if (selectedParent && selectedParent.id && !stateRef.current.isUpdating) {
      addBlindItem(selectedParent.id);
    }
  }, [selectedParent, addBlindItem]);

  // FIX: Added flag to prevent unwanted rotation during manual pan
  const isManuallyPanning2 = useRef(false);

  const panResponder2 = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => !isSpinning,
      onPanResponderGrant: () => {
        setIsSpinning2(true);
        isManuallyPanning2.current = true; // Mark as manual panning
      },
      onPanResponderMove: (_, gesture) => {
        const newRotation = lastRotation2.current + gesture.dx / 2;
        rotation2.setValue(newRotation);
      },
      onPanResponderRelease: (_, gesture) => {
        setIsSpinning2(false);
        const finalRotation = lastRotation2.current + gesture.dx / 2;
        lastRotation2.current = finalRotation;

        try {
          const selectedSlice = getSelectedSlice(
            finalRotation,
            stateRef.current.data2,
            stateRef.current.total2
          );

          if (selectedSlice && selectedSlice.id) {
            const currentParentId = stateRef.current.selectedParent?.id;
            const childParentId = selectedSlice.parentId;

            if (childParentId && currentParentId !== childParentId) {
              const newParent = data.find((item) => item.id === childParentId);

              if (newParent) {
                const targetRotation = calculateRotationForParent(
                  childParentId,
                  data,
                  total
                );

                Animated.timing(rotation, {
                  toValue: targetRotation,
                  duration: 800,
                  useNativeDriver: false,
                  easing: Easing.out(Easing.quad),
                }).start();

                lastRotation.current = targetRotation;
                setSelectedParent(newParent);
              }
            }

            setSelectedChild(selectedSlice);
          }
        } catch (error) {
          console.log("Error getting selected child slice:", error);
        }

        // Reset manual panning flag after a delay
        setTimeout(() => {
          isManuallyPanning2.current = false;
        }, 200);
      },
    })
  ).current;

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
            const currentParentId = stateRef.current.selectedParent?.id;
            const newParentId = selectedSlice.id;

            if (currentParentId !== newParentId) {
              const firstChild = stateRef.current.data2.find(
                (item) => item.parentId === newParentId
              );

              if (firstChild) {
                // FIX: Only auto-rotate second circle if it's not being manually panned
                if (!isManuallyPanning2.current) {
                  const targetRotation2 = calculateRotationForChild(
                    newParentId,
                    stateRef.current.data2,
                    stateRef.current.total2
                  );

                  Animated.timing(rotation2, {
                    toValue: targetRotation2,
                    duration: 800,
                    useNativeDriver: false,
                    easing: Easing.out(Easing.quad),
                  }).start();

                  lastRotation2.current = targetRotation2;
                }
                setSelectedChild(firstChild);
              }
            }

            setSelectedParent(selectedSlice);
          }
        } catch (error) {
          console.log("Error getting selected slice:", error);
        }
      },
    })
  ).current;

  useEffect(() => {
    return () => {
      if (blindItemTimeoutRef.current) {
        clearTimeout(blindItemTimeoutRef.current);
      }
    };
  }, []);

  return (
    <View style={styles.container}>
      <TouchableOpacity style={[styles.saveButton]}>
        <Text style={[styles.saveButtonText, styles.saveButtonPosition]}>
          SAVE
        </Text>
      </TouchableOpacity>

      <View style={styles.selectionDisplay}>
        <Text style={styles.selectionText}>
          Parent: {selectedParent?.label || "None"} ({selectedParent?.id})
        </Text>
        <Text style={styles.selectionText}>
          Child: {selectedChild?.label || "None"} ({selectedChild?.id})
        </Text>
        <Text style={styles.selectionText}>
          Child Parent: {selectedChild?.parentId || "None"}
        </Text>
        <Text style={styles.selectionText}>
          Match: {selectedChild?.parentId === selectedParent?.id ? "✅" : "❌"}
        </Text>
        <Text style={styles.selectionText}>Total Items: {total2}</Text>
        <Text style={styles.selectionText}>
          BLIND Items: {data2.filter((item) => item.id === "0_0_0").length}
        </Text>
      </View>

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

                    const innerRadius = radius / 2 + 145;
                    const labelRadius = innerRadius + 15;

                    const rad = ((midAngle - 90) * Math.PI) / 180;
                    const x = center2 + labelRadius * Math.cos(rad);
                    const y = center2 + labelRadius * Math.sin(rad);

                    const textRotation = midAngle + 90;

                    const isChildOfSelectedParent =
                      selectedParent && item.parentId === selectedParent.id;

                    const fillColor = isChildOfSelectedParent
                      ? item.color
                      : darkenColor(item.color, 0.15);
                    const fillOpacity = isChildOfSelectedParent ? 1 : 0.25;
                    const textOpacity = isChildOfSelectedParent ? 1 : 0.3;

                    const isSelected =
                      selectedChild && selectedChild.id === item.id;
                    const strokeWidth = isSelected ? 4 : 2;
                    const strokeColor = isSelected ? "#FFD700" : "#fff";

                    const arc = (
                      <React.Fragment key={item.id}>
                        <Path
                          d={path}
                          fill={fillColor}
                          fillOpacity={fillOpacity}
                          stroke={strokeColor}
                          strokeWidth={strokeWidth}
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

                    const isSelected =
                      selectedParent && selectedParent.id === item.id;
                    const strokeWidth = isSelected ? 4 : 2;
                    const strokeColor = isSelected ? "#FFD700" : "#fff";

                    const arc = (
                      <Path
                        key={item.id}
                        d={path}
                        fill={item.color}
                        stroke={strokeColor}
                        strokeWidth={strokeWidth}
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
    fontSize: 11,
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
    zIndex: 30,
  },
  secondCircleWrapper: {
    position: "absolute",
    bottom: -radius2 * 0.095,
    left: (screenWidth - radius2) / 2,
    width: radius2,
    height: radius2 / 2,
    overflow: "hidden",
    alignItems: "center",
    zIndex: 20,
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
