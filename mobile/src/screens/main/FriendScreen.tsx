import React from "react";
import { View, Text } from "react-native";

export default function FriendScreen() {
  return (
    <View className="flex-1 bg-white justify-center items-center">
      <Text className="text-2xl font-bold text-amber-900">Friends</Text>
      <Text className="text-amber-700 mt-2">
        Connect with coffee enthusiasts
      </Text>
    </View>
  );
}
