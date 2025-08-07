import React, { useState, useRef } from "react";
import {
  View,
  ScrollView,
  Text,
  TouchableOpacity,
  Animated,
} from "react-native";

export default function CuppingRegistration() {
  const [selectedItems, setSelectedItems] = useState<string[]>([]);
  const [currentRotation, setCurrentRotation] = useState(0);
  const rotationValue = useRef(new Animated.Value(0)).current;

  // Items around the circle
  const circleItems = [
    { id: 1, label: "Espresso", emoji: "☕", color: "#8B4513" },
    { id: 2, label: "Americano", emoji: "🫖", color: "#CD853F" },
    { id: 3, label: "Cappuccino", emoji: "🥛", color: "#DEB887" },
    { id: 4, label: "Latte", emoji: "🍼", color: "#F5DEB3" },
    { id: 5, label: "Mocha", emoji: "🍫", color: "#A0522D" },
    { id: 6, label: "Macchiato", emoji: "☕", color: "#D2B48C" },
    { id: 7, label: "Frappé", emoji: "🧊", color: "#87CEEB" },
    { id: 8, label: "Turkish", emoji: "🫖", color: "#8B4513" },
  ];

  const itemAngle = 360 / circleItems.length;

  const getCurrentSelectedItem = () => {
    const normalizedRotation = ((currentRotation % 360) + 360) % 360;
    const selectedIndex =
      Math.floor((360 - normalizedRotation + itemAngle / 2) / itemAngle) %
      circleItems.length;
    return circleItems[selectedIndex];
  };

  const addSelectedItem = () => {
    const selectedItem = getCurrentSelectedItem();
    if (!selectedItems.includes(selectedItem.label)) {
      setSelectedItems([...selectedItems, selectedItem.label]);
    }
  };

  const removeItem = (item: string) => {
    setSelectedItems(selectedItems.filter((i) => i !== item));
  };

  const rotateCircle = () => {
    const newRotation = currentRotation + 45; // Rotate by 45 degrees each click
    setCurrentRotation(newRotation);

    Animated.timing(rotationValue, {
      toValue: newRotation,
      duration: 300,
      useNativeDriver: true,
    }).start();
  };

  const renderCircleItem = (item: any, index: number) => {
    const angle = (index * itemAngle - 90) * (Math.PI / 180); // -90 to start from top
    const radius = 80;
    const x = Math.cos(angle) * radius;
    const y = Math.sin(angle) * radius;

    const isSelected = getCurrentSelectedItem().id === item.id;

    return (
      <View
        key={item.id}
        style={{
          position: "absolute",
          left: x,
          top: y,
          transform: [{ translateX: -20 }, { translateY: -20 }],
          backgroundColor: item.color,
          borderWidth: isSelected ? 3 : 0,
          borderColor: "#FFF",
        }}
        className="w-10 h-10 rounded-full items-center justify-center"
      >
        <Text className="text-white text-lg">{item.emoji}</Text>
      </View>
    );
  };
  return (
    <View className="flex-1 bg-gray-50">
      <ScrollView className="flex-1">
        <Text>CuppingRegistration</Text>
      </ScrollView>

      <View className="absolute bottom-0 left-0 right-0 p-6 bg-white border-t border-gray-200">
        <View className="items-center relative">
          {/* Bigger circle (background/under) */}
          <TouchableOpacity className="bg-amber-600 w-20 h-20 rounded-full items-center justify-center">
            <Text className="text-white font-semibold text-2xl">-</Text>
          </TouchableOpacity>

          {/* Smaller circle (foreground/over) */}
          <TouchableOpacity
            className="absolute bg-red-600 w-16 h-16 rounded-full items-center justify-center"
            onPress={addSelectedItem}
          >
            <Text className="text-white font-semibold text-2xl">+</Text>
          </TouchableOpacity>
        </View>
      </View>
    </View>
  );
}
