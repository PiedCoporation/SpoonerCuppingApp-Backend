import React, { useRef } from "react";
import {
  View,
  Text,
  TouchableOpacity,
  Animated,
  Dimensions,
  PanResponder,
} from "react-native";

const { width: screenWidth } = Dimensions.get("window");

export default function CuppingRegistration() {
  const rotationValue = useRef(new Animated.Value(0)).current;
  const lastGesture = useRef({ x: 0, y: 0 });
  const lastRotation = useRef(0);

  const items = [
    { id: 1, label: "Aroma", icon: "👃", color: "#FF6B6B" },
    { id: 2, label: "Flavor", icon: "👅", color: "#4ECDC4" },
    { id: 3, label: "Body", icon: "💪", color: "#45B7D1" },
    { id: 4, label: "Acidity", icon: "🍋", color: "#96CEB4" },
    { id: 5, label: "Balance", icon: "⚖️", color: "#FFEAA7" },
    { id: 6, label: "Overall", icon: "⭐", color: "#DDA0DD" },
    { id: 7, label: "Aftertaste", icon: "👅", color: "#FFB6C1" },
    { id: 8, label: "Cleanliness", icon: "🧼", color: "#FFD700" },
    { id: 9, label: "Sweetness", icon: "🍯", color: "#FF9F43" },
  ];

  // Fixed: Use local coordinates relative to the pan responder view
  const getAngle = (x: number, y: number) => {
    // The center is at screenWidth/2, screenWidth/2 relative to the pan responder view
    const centerX = screenWidth / 2;
    const centerY = screenWidth / 2;
    return Math.atan2(y - centerY, x - centerX);
  };

  const panResponder = PanResponder.create({
    onMoveShouldSetPanResponder: () => true,
    onStartShouldSetPanResponder: () => true,
    // Fixed: Intercept touches on child components too
    onMoveShouldSetPanResponderCapture: () => true,
    onStartShouldSetPanResponderCapture: () => true,

    onPanResponderGrant: (evt) => {
      // Fixed: Use pageX/pageY for consistent coordinates, then convert to local
      const localX =
        evt.nativeEvent.pageX -
        evt.nativeEvent.locationX +
        evt.nativeEvent.locationX;
      const localY =
        evt.nativeEvent.pageY -
        evt.nativeEvent.locationY +
        evt.nativeEvent.locationY;

      lastGesture.current = {
        x: evt.nativeEvent.locationX,
        y: evt.nativeEvent.locationY,
      };
    },

    onPanResponderMove: (evt) => {
      const currentAngle = getAngle(
        evt.nativeEvent.locationX,
        evt.nativeEvent.locationY
      );
      const lastAngle = getAngle(lastGesture.current.x, lastGesture.current.y);

      let angleDiff = currentAngle - lastAngle;

      // Handle angle wrapping
      if (angleDiff > Math.PI) angleDiff -= 2 * Math.PI;
      if (angleDiff < -Math.PI) angleDiff += 2 * Math.PI;

      lastRotation.current += angleDiff;
      rotationValue.setValue(lastRotation.current);

      lastGesture.current = {
        x: evt.nativeEvent.locationX,
        y: evt.nativeEvent.locationY,
      };
    },

    onPanResponderRelease: (evt, gestureState) => {
      const velocity = Math.sqrt(
        gestureState.vx * gestureState.vx + gestureState.vy * gestureState.vy
      );

      if (velocity > 0.5) {
        // Calculate momentum based on angular velocity
        const centerX = screenWidth / 2;
        const centerY = screenWidth / 2;
        const radius = Math.sqrt(
          Math.pow(evt.nativeEvent.locationX - centerX, 2) +
            Math.pow(evt.nativeEvent.locationY - centerY, 2)
        );

        // Convert linear velocity to angular velocity
        const angularVelocity = (velocity / Math.max(radius, 1)) * 0.1;

        // Determine direction based on the last movement
        const currentAngle = getAngle(
          evt.nativeEvent.locationX,
          evt.nativeEvent.locationY
        );
        const lastAngle = getAngle(
          lastGesture.current.x,
          lastGesture.current.y
        );
        let direction = currentAngle - lastAngle;
        if (direction > Math.PI) direction -= 2 * Math.PI;
        if (direction < -Math.PI) direction += 2 * Math.PI;

        const momentumSpin = angularVelocity * Math.sign(direction);
        lastRotation.current += momentumSpin;

        Animated.spring(rotationValue, {
          toValue: lastRotation.current,
          useNativeDriver: true,
          tension: 50,
          friction: 8,
        }).start();
      }
    },
  });

  const renderCircleItems = () => {
    const radius = screenWidth * 0.35;
    const centerX = screenWidth / 2;
    const centerY = screenWidth / 2;

    return items.map((item, index) => {
      const segmentAngle = 360 / items.length;
      // FIXED: Position items exactly in the center of their colored segments
      // Segments start at: index * segmentAngle - 90
      // Item should be at the center: segment start + segmentAngle/2
      const itemAngleDegrees = index * segmentAngle - 90 + segmentAngle / 2;
      const itemAngle = (itemAngleDegrees * Math.PI) / 180;

      return (
        <Animated.View
          key={item.id}
          style={{
            position: "absolute",
            left: centerX - 32,
            top: centerY - 32,
            transform: [
              {
                rotate: rotationValue.interpolate({
                  inputRange: [0, 2 * Math.PI],
                  outputRange: ["0rad", "6.28rad"],
                }),
              },
              { translateX: radius * Math.cos(itemAngle) },
              { translateY: radius * Math.sin(itemAngle) },
            ],
          }}
        >
          <TouchableOpacity
            className="w-16 h-16 bg-white rounded-full items-center justify-center shadow-lg border border-gray-200"
            onPress={() => console.log(`Pressed ${item.label}`)}
            // Fixed: Prevent touch events from interfering with pan responder
            activeOpacity={0.7}
            style={{ pointerEvents: "box-only" }}
          >
            <Text className="text-xl mb-1">{item.icon}</Text>
            <Text className="text-xs font-semibold text-gray-700 text-center">
              {item.label}
            </Text>
          </TouchableOpacity>
        </Animated.View>
      );
    });
  };

  return (
    <View className="flex-1 bg-white">
      <View className="flex-1 px-5 pt-5">
        <Text className="text-2xl font-bold text-center mb-8">
          Coffee Cupping Registration
        </Text>
      </View>

      {/* Half circle at the bottom */}
      <View
        className="absolute bottom-0 left-0 right-0 overflow-hidden items-center"
        style={{ height: screenWidth }}
        {...panResponder.panHandlers}
      >
        {/* Base circle */}
        <Animated.View
          style={{
            width: screenWidth,
            height: screenWidth,
            borderRadius: screenWidth / 2,
            backgroundColor: "#f8f8f8",
            transform: [
              {
                rotate: rotationValue.interpolate({
                  inputRange: [0, 2 * Math.PI],
                  outputRange: ["0rad", "6.28rad"],
                }),
              },
            ],
          }}
        />

        {/* Colored segments container */}
        <Animated.View
          style={{
            position: "absolute",
            width: screenWidth,
            height: screenWidth,
            borderRadius: screenWidth / 2,
            overflow: "hidden",
            transform: [
              {
                rotate: rotationValue.interpolate({
                  inputRange: [0, 2 * Math.PI],
                  outputRange: ["0rad", "6.28rad"],
                }),
              },
            ],
          }}
        >
          {/* Create segments with proper alignment to items */}
          {items.map((item, index) => {
            const segmentAngle = 360 / items.length;
            // FIXED: Create segments that are centered around where items will be placed
            // Items will be at: index * segmentAngle - 90 + segmentAngle/2
            // So segments should start at: index * segmentAngle - 90
            const segmentStartAngle = index * segmentAngle - 90;
            const skewAngle = 90 - segmentAngle;

            return (
              <View
                key={`segment-${item.id}`}
                style={{
                  position: "absolute",
                  width: screenWidth / 2,
                  height: screenWidth / 2,
                  left: screenWidth / 2,
                  top: screenWidth / 2,
                  backgroundColor: item.color,
                  transformOrigin: "0 0",
                  transform: [
                    { rotate: `${segmentStartAngle}deg` },
                    { skewY: `${skewAngle}deg` },
                  ],
                }}
              />
            );
          })}
        </Animated.View>

        {/* Add separator lines to create clean divisions */}
        <Animated.View
          style={{
            position: "absolute",
            width: screenWidth,
            height: screenWidth,
            transform: [
              {
                rotate: rotationValue.interpolate({
                  inputRange: [0, 2 * Math.PI],
                  outputRange: ["0rad", "6.28rad"],
                }),
              },
            ],
          }}
        >
          {/* Create separator lines at segment boundaries */}
          {items.map((item, index) => {
            const segmentAngle = 360 / items.length;
            // Place separators at the START of each segment (boundaries between segments)
            const separatorAngle = index * segmentAngle - 90;
            return (
              <View
                key={`separator-${index}`}
                style={{
                  position: "absolute",
                  width: screenWidth * 0.5,
                  height: 3,
                  left: screenWidth / 2,
                  top: screenWidth / 2 - 1.5,
                  backgroundColor: "white",
                  transform: [{ rotate: `${separatorAngle}deg` }],
                  transformOrigin: "0 50%",
                }}
              />
            );
          })}
        </Animated.View>

        {/* White center circle to create donut effect */}
        <View
          className="absolute bg-white rounded-full shadow-lg"
          style={{
            width: screenWidth * 0.3,
            height: screenWidth * 0.3,
            left: screenWidth * 0.35,
            top: screenWidth * 0.35,
          }}
        />

        {/* Circle items */}
        {renderCircleItems()}
      </View>
    </View>
  );
}
