import React, {
  useRef,
  useState,
  useEffect,
  useCallback,
  useMemo,
} from "react";
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
import { firstCircle, secondCircle, thirdCircle } from "@/data/dataV1";
import { Text as SvgText } from "react-native-svg";
import { FlavorItem } from "@/types/Flavor";

const polarToCartesian = (
  centerX: number,
  centerY: number,
  radius: number,
  angleInDegrees: number
) => {
  const angleInRadians = ((angleInDegrees - 90) * Math.PI) / 180.0;
  return {
    x: centerX + radius * Math.cos(angleInRadians),
    y: centerY + radius * Math.sin(angleInRadians),
  };
};

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

const describeArc = (
  x: number,
  y: number,
  radius: number,
  startAngle: number,
  endAngle: number
) => {
  const start = polarToCartesian(x, y, radius, endAngle);
  const end = polarToCartesian(x, y, radius, startAngle);
  const largeArcFlag = endAngle - startAngle <= 180 ? "0" : "1";
  return [
    `M ${x} ${y}`,
    `L ${start.x} ${start.y}`,
    `A ${radius} ${radius} 0 ${largeArcFlag} 0 ${end.x} ${end.y}`,
    "Z",
  ].join(" ");
};

// First circle color darkening function
function darkenColor(
  color: string,
  factor: number = 0.4,
  fallback: string = "#666666"
): string {
  if (!color || color === "") return fallback;
  const hex = color.replace("#", "");
  const r = parseInt(hex.slice(0, 2), 16);
  const g = parseInt(hex.slice(2, 4), 16);
  const b = parseInt(hex.slice(4, 6), 16);
  const newR = Math.floor(r * factor);
  const newG = Math.floor(g * factor);
  const newB = Math.floor(b * factor);
  return `#${newR.toString(16).padStart(2, "0")}${newG
    .toString(16)
    .padStart(2, "0")}${newB.toString(16).padStart(2, "0")}`;
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

function calculateInitialRotation3(
  data: Array<{
    id: string;
    label: string;
    color: string;
    numbers: number;
    parentId?: string;
  }>,
  total: number,
  parentId: string
) {
  const firstChildIndex = data.findIndex((item) => item.parentId === parentId);
  if (firstChildIndex === -1) return 0;

  const firstChildSliceAngle = (data[firstChildIndex].numbers / total) * 360;
  const firstChildCenterAngle = firstChildSliceAngle / 2;
  return firstChildCenterAngle;
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

function calculateRotationForSpecificGrandchild(
  grandchildId: string,
  data: Array<{
    id: string;
    label: string;
    color: string;
    numbers: number;
    parentId?: string;
  }>,
  total: number
) {
  const grandchildIndex = data.findIndex((item) => item.id === grandchildId);
  if (grandchildIndex === -1) return 0;

  let cumulativeAngle = 0;
  for (let i = 0; i < grandchildIndex; i++) {
    cumulativeAngle += (data[i].numbers / total) * 360;
  }

  const grandchildSliceAngle = (data[grandchildIndex].numbers / total) * 360;
  const grandchildCenterAngle = cumulativeAngle + grandchildSliceAngle / 2;

  let targetRotation = grandchildCenterAngle - 90;
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
const radius3 = screenWidth * 3.2;
const center3 = radius3 / 2;

export default function CuppingRegistrationOverview() {
  const [data2, setData2] = useState<FlavorItem[]>(secondCircle.data);
  const [total2, setTotal2] = useState(secondCircle.numbers);
  const [data3, setData3] = useState<FlavorItem[]>(thirdCircle.data);
  const [total3, setTotal3] = useState(thirdCircle.numbers);

  const stateRef = useRef({
    data2: secondCircle.data,
    total2: secondCircle.numbers,
    data3: thirdCircle.data,
    total3: thirdCircle.numbers,
    selectedParent: data[0],
    selectedChild: secondCircle.data[0],
    selectedGrandchild: thirdCircle.data[0],
    isUpdating: false,
    isUpdating3: false,
  });

  useEffect(() => {
    stateRef.current.data2 = data2;
    stateRef.current.total2 = total2;
    stateRef.current.data3 = data3;
    stateRef.current.total3 = total3;
  }, [data2, total2, data3, total3]);

  const initialRotationValue = calculateInitialRotation(data, total);
  const initialRotationValue2 = calculateInitialRotation2(data2, total2);
  const initialRotationValue3 = calculateInitialRotation3(
    data3,
    total3,
    secondCircle.data[0].id
  );

  const rotation = useRef(new Animated.Value(initialRotationValue)).current;
  const lastRotation = useRef(initialRotationValue);

  const rotation2 = useRef(new Animated.Value(initialRotationValue2)).current;
  const lastRotation2 = useRef(initialRotationValue2);

  const rotation3 = useRef(new Animated.Value(initialRotationValue3)).current;
  const lastRotation3 = useRef(initialRotationValue3);

  const [isSpinning, setIsSpinning] = useState(false);
  const [isSpinning2, setIsSpinning2] = useState(false);
  const [isSpinning3, setIsSpinning3] = useState(false);

  const [selectedParent, setSelectedParentState] = useState<any | null>(
    data[0]
  );
  const [selectedChild, setSelectedChildState] = useState<any | null>(data2[0]);
  const [selectedGrandchild, setSelectedGrandchildState] = useState<any | null>(
    data3[0]
  );

  const setSelectedParent = useCallback((value: any | null) => {
    stateRef.current.selectedParent = value;
    setSelectedParentState(value);
  }, []);

  const setSelectedChild = useCallback((value: any | null) => {
    stateRef.current.selectedChild = value;
    setSelectedChildState(value);
  }, []);

  const setSelectedGrandchild = useCallback((value: any | null) => {
    stateRef.current.selectedGrandchild = value;
    setSelectedGrandchildState(value);
  }, []);

  const blindItemTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const blindItemTimeoutRef3 = useRef<NodeJS.Timeout | null>(null);

  // Add helper function to check if third circle should be disabled
  const isThirdCircleDisabled = useCallback(() => {
    return (
      !selectedChild ||
      selectedChild.id === "0_0_0" ||
      selectedChild.label === "BLIND"
    );
  }, [selectedChild]);

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
        hasChild: false,
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

  const addBlindItem3 = useCallback((parentId: string) => {
    if (stateRef.current.isUpdating3) return;

    if (blindItemTimeoutRef3.current) {
      clearTimeout(blindItemTimeoutRef3.current);
    }

    blindItemTimeoutRef3.current = setTimeout(() => {
      stateRef.current.isUpdating3 = true;

      const blindItem: FlavorItem = {
        id: "0_0_0_0",
        label: "BLIND",
        color: "#cbcfce",
        numbers: 1,
        parentId: parentId,
        hasChild: false, // NEW: Added hasChild property
      };

      setData3((prevData) => {
        const dataWithoutBlind = prevData.filter(
          (item) => item.id !== "0_0_0_0"
        );
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

        return newData;
      });

      setTotal3((prevTotal) => {
        const hasExistingBlind = stateRef.current.data3.some(
          (item) => item.id === "0_0_0_0"
        );
        return hasExistingBlind ? prevTotal : prevTotal + 1;
      });

      setTimeout(() => {
        stateRef.current.isUpdating3 = false;
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

  const removeBlindItem3 = useCallback(() => {
    const blindItem = stateRef.current.data3.find(
      (item) => item.id === "0_0_0_0"
    );
    if (blindItem) {
      setData3((prevData) => prevData.filter((item) => item.id !== "0_0_0_0"));
      setTotal3((prevTotal) => prevTotal - 1);
    }
  }, []);

  useEffect(() => {
    setSelectedParent(data[0]);
    setSelectedChild(data2[0]);
    setSelectedGrandchild(data3[0]);
  }, [setSelectedParent, setSelectedChild, setSelectedGrandchild]);

  useEffect(() => {
    if (selectedParent && selectedParent.id && !stateRef.current.isUpdating) {
      addBlindItem(selectedParent.id);
    }
  }, [selectedParent, addBlindItem]);

  // Modified useEffect for handling BLIND items in second circle
  useEffect(() => {
    if (selectedChild && selectedChild.id && !stateRef.current.isUpdating3) {
      // Check if the selected child is a BLIND item
      if (selectedChild.id === "0_0_0" || selectedChild.label === "BLIND") {
        // If BLIND item is selected, don't add BLIND to third circle
        // and clear any existing grandchild selection
        setSelectedGrandchild(null);
        // Remove any existing BLIND items from third circle
        removeBlindItem3();
      } else if (selectedChild.hasChild === true) {
        // NEW: Only add BLIND item to third circle if selectedChild has children
        addBlindItem3(selectedChild.id);
      } else {
        // NEW: If hasChild is false, just remove any existing BLIND but don't add new one
        removeBlindItem3();
      }
    }
  }, [selectedChild, addBlindItem3, removeBlindItem3]);

  const isManuallyPanning2 = useRef(false);
  const isManuallyPanning3 = useRef(false);

  // Modified panResponder3 to prevent interaction when disabled
  const panResponder3 = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () =>
        !isSpinning && !isSpinning2 && !isThirdCircleDisabled(),
      onPanResponderGrant: () => {
        if (isThirdCircleDisabled()) return;
        setIsSpinning3(true);
        isManuallyPanning3.current = true;
      },
      onPanResponderMove: (_, gesture) => {
        if (isThirdCircleDisabled()) return;
        const newRotation = lastRotation3.current + gesture.dx / 2;
        rotation3.setValue(newRotation);
      },
      onPanResponderRelease: (_, gesture) => {
        if (isThirdCircleDisabled()) return;
        setIsSpinning3(false);
        const finalRotation = lastRotation3.current + gesture.dx / 2;
        lastRotation3.current = finalRotation;

        try {
          const selectedSlice = getSelectedSlice(
            finalRotation,
            stateRef.current.data3,
            stateRef.current.total3
          );

          if (selectedSlice && selectedSlice.id) {
            const currentChildId = stateRef.current.selectedChild?.id;
            const grandchildParentId = selectedSlice.parentId;

            if (grandchildParentId && currentChildId !== grandchildParentId) {
              const newChild = stateRef.current.data2.find(
                (item) => item.id === grandchildParentId
              );

              if (newChild) {
                const targetRotation2 = calculateRotationForParent(
                  grandchildParentId,
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
                setSelectedChild(newChild);

                const newChildParentId = newChild.parentId;
                const currentParentId = stateRef.current.selectedParent?.id;

                if (newChildParentId && currentParentId !== newChildParentId) {
                  const newParent = data.find(
                    (item) => item.id === newChildParentId
                  );
                  if (newParent) {
                    const targetRotation = calculateRotationForParent(
                      newChildParentId,
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
              }
            }

            setSelectedGrandchild(selectedSlice);
          }
        } catch (error) {
          console.log("Error getting selected grandchild slice:", error);
        }

        setTimeout(() => {
          isManuallyPanning3.current = false;
        }, 200);
      },
    })
  ).current;

  // Modified panResponder2 to handle BLIND selection
  // Modified panResponder2 to point to first real child, not BLIND
  const panResponder2 = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => !isSpinning && !isSpinning3,
      onPanResponderGrant: () => {
        setIsSpinning2(true);
        isManuallyPanning2.current = true;
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

            // Only rotate third circle if the selected item is not BLIND
            if (
              !isManuallyPanning3.current &&
              selectedSlice.id !== "0_0_0" &&
              selectedSlice.label !== "BLIND"
            ) {
              // Find the first REAL grandchild (not BLIND) in the current data
              const realGrandchild = stateRef.current.data3.find(
                (item) =>
                  item.parentId === selectedSlice.id &&
                  item.id !== "0_0_0_0" &&
                  item.label !== "BLIND"
              );

              if (realGrandchild) {
                // Create the BLIND item that WILL be added
                const futureBlindItem = {
                  id: "0_0_0_0",
                  label: "BLIND",
                  color: "#cbcfce",
                  numbers: 1,
                  parentId: selectedSlice.id,
                };

                // Simulate the data structure AFTER BLIND is added
                const dataWithoutBlind = stateRef.current.data3.filter(
                  (item) => item.id !== "0_0_0_0"
                );
                const insertionIndex = dataWithoutBlind.findIndex(
                  (item) => item.parentId === selectedSlice.id
                );
                const finalInsertionIndex =
                  insertionIndex === -1
                    ? dataWithoutBlind.length
                    : insertionIndex;

                const simulatedData3 = [
                  ...dataWithoutBlind.slice(0, finalInsertionIndex),
                  futureBlindItem,
                  ...dataWithoutBlind.slice(finalInsertionIndex),
                ];

                // Find the real grandchild's NEW position in the simulated array
                const realGrandchildInSimulated = simulatedData3.find(
                  (item) => item.id === realGrandchild.id
                );

                if (realGrandchildInSimulated) {
                  // Calculate rotation to point to the REAL grandchild's position
                  const targetRotation3 =
                    calculateRotationForSpecificGrandchild(
                      realGrandchild.id,
                      simulatedData3,
                      stateRef.current.total3 + 1
                    );

                  Animated.timing(rotation3, {
                    toValue: targetRotation3,
                    duration: 800,
                    useNativeDriver: false,
                    easing: Easing.out(Easing.quad),
                  }).start();

                  lastRotation3.current = targetRotation3;

                  // Set the grandchild to the REAL child, not BLIND
                  setSelectedGrandchild(realGrandchild);
                }
              }
            } else if (
              selectedSlice.id === "0_0_0" ||
              selectedSlice.label === "BLIND"
            ) {
              // Clear grandchild selection when BLIND is selected
              setSelectedGrandchild(null);
            }

            setSelectedChild(selectedSlice);
          }
        } catch (error) {
          console.log("Error getting selected child slice:", error);
        }

        setTimeout(() => {
          isManuallyPanning2.current = false;
        }, 200);
      },
    })
  ).current;

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
          const selectedSlice = getSelectedSlice(finalRotation, data, total);
          if (selectedSlice && selectedSlice.id) {
            const currentParentId = stateRef.current.selectedParent?.id;
            const newParentId = selectedSlice.id;

            if (currentParentId !== newParentId) {
              const firstChild = stateRef.current.data2.find(
                (item) => item.parentId === newParentId
              );

              if (firstChild) {
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

                if (!isManuallyPanning3.current) {
                  const firstGrandchild = stateRef.current.data3.find(
                    (item) => item.parentId === firstChild.id
                  );

                  if (firstGrandchild) {
                    const targetRotation3 =
                      calculateRotationForSpecificGrandchild(
                        firstGrandchild.id,
                        stateRef.current.data3,
                        stateRef.current.total3
                      );

                    Animated.timing(rotation3, {
                      toValue: targetRotation3,
                      duration: 800,
                      useNativeDriver: false,
                      easing: Easing.out(Easing.quad),
                    }).start();

                    lastRotation3.current = targetRotation3;
                    setSelectedGrandchild(firstGrandchild);
                  }
                }
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
      if (blindItemTimeoutRef3.current) {
        clearTimeout(blindItemTimeoutRef3.current);
      }
    };
  }, []);

  const renderedThirdCircle = useMemo(() => {
    let cumulativeAngle = 0;
    return data3.map((item) => {
      const sliceAngle = (item.numbers / total3) * 360;
      const endAngle = cumulativeAngle + sliceAngle;
      const path = describeArc(
        center3,
        center3,
        radius3 / 2 - 10,
        cumulativeAngle,
        endAngle
      );
      const midAngle = (cumulativeAngle + endAngle) / 2;

      // FIXED: Use radius3 for consistent positioning
      const innerRadius = radius3 / 2 - 50;
      const labelRadius = innerRadius - 30;

      const rad = ((midAngle - 90) * Math.PI) / 180;
      const x = center3 + labelRadius * Math.cos(rad);
      const y = center3 + labelRadius * Math.sin(rad);

      const textRotation = midAngle + 90;

      const isGrandchildOfSelectedChild =
        selectedChild && item.parentId === selectedChild.id;

      const isEmptyLabel = item.label === "" || !item.label;

      const fillColor = isEmptyLabel
        ? "transparent"
        : isGrandchildOfSelectedChild
        ? item.color
        : darkenColor(item.color, 0.15);
      const fillOpacity = isEmptyLabel
        ? 0
        : isGrandchildOfSelectedChild
        ? 1
        : 0.25;

      const textOpacity = isEmptyLabel
        ? 0
        : isGrandchildOfSelectedChild
        ? 1
        : 0.6;

      const isSelected =
        selectedGrandchild && selectedGrandchild.id === item.id;
      const strokeWidth = isSelected ? 4 : 2;
      const strokeColor = isSelected ? "#FFD700" : "#fff";
      const strokeOpacity = isEmptyLabel ? 0 : 1;

      const arc = (
        <React.Fragment key={item.id}>
          <Path
            d={path}
            fill={fillColor}
            fillOpacity={fillOpacity}
            stroke={strokeColor}
            strokeWidth={strokeWidth}
            strokeOpacity={strokeOpacity}
          />
          {!isEmptyLabel && (
            <SvgText
              x={x}
              y={y}
              fill="#fff"
              fillOpacity={textOpacity}
              fontSize="11"
              fontWeight={isGrandchildOfSelectedChild ? "bold" : "normal"}
              textAnchor="middle"
              alignmentBaseline="middle"
              transform={`rotate(${textRotation}, ${x}, ${y})`}
            >
              {item.label.replace("/", "\n")}
            </SvgText>
          )}
        </React.Fragment>
      );
      cumulativeAngle = endAngle;
      return arc;
    });
  }, [data3, total3, selectedChild, selectedGrandchild]);

  const renderedSecondCircle = useMemo(() => {
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
        : darkenColor(item.color, 0.5);
      const fillOpacity = 1;
      const textOpacity = isChildOfSelectedParent ? 1 : 0.6;

      const isSelected = selectedChild && selectedChild.id === item.id;
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
            strokeOpacity={1}
          />
          <SvgText
            x={x}
            y={y}
            fill="#fff"
            fillOpacity={textOpacity}
            fontSize="14"
            fontWeight={isChildOfSelectedParent ? "bold" : "normal"}
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
  }, [data2, total2, selectedParent, selectedChild]);

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
          Grandchild: {selectedGrandchild?.label || "None"} (
          {selectedGrandchild?.id})
        </Text>
        <Text style={styles.selectionText}>
          Third Circle: {isThirdCircleDisabled() ? "DISABLED" : "ENABLED"}
        </Text>
        <Text style={styles.selectionText}>
          Child Match:{" "}
          {selectedChild?.parentId === selectedParent?.id ? "✅" : "❌"}
        </Text>
        <Text style={styles.selectionText}>
          Grandchild Match:{" "}
          {selectedGrandchild?.parentId === selectedChild?.id ? "✅" : "❌"}
        </Text>
        <Text style={styles.selectionText}>
          Total2: {total2} | Total3: {total3}
        </Text>
      </View>

      {/* Third circle - largest and behind - only render when not disabled */}
      {!isThirdCircleDisabled() && (
        <View style={styles.thirdCircleWrapper}>
          <View style={styles.pointer3}>
            <Svg width={30} height={30} style={styles.pointerSvg}>
              <Polygon
                points="15,27 10,15 20,15"
                fill="#666"
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
                  {renderedThirdCircle}
                </G>
              </Svg>
            </Animated.View>
          </View>
        </View>
      )}

      {/* Second circle - middle */}
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
                {renderedSecondCircle}
              </G>
            </Svg>
          </Animated.View>
        </View>
      </View>

      {/* First circle - front */}
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

                    const fillColor = isSelected
                      ? item.color
                      : darkenColor(item.color, 0.5);

                    const strokeWidth = isSelected ? 4 : 2;
                    const strokeColor = isSelected ? "#FFD700" : "#fff";

                    const arc = (
                      <Path
                        key={item.id}
                        d={path}
                        fill={fillColor}
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
    fontSize: 10,
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
  thirdCircleWrapper: {
    position: "absolute",
    bottom: -radius3 * 0.075,
    left: (screenWidth - radius3) / 2,
    width: radius3,
    height: radius3 / 2,
    overflow: "hidden",
    alignItems: "center",
    zIndex: 10,
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
    zIndex: 40,
    alignItems: "center",
  },
  pointer2: {
    position: "absolute",
    top: -8,
    zIndex: 30,
    alignItems: "center",
  },
  pointer3: {
    position: "absolute",
    top: -10,
    zIndex: 20,
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
