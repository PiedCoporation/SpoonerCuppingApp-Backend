import React from "react";
import { View, Text } from "react-native";

export default function NotificationScreen() {
  return (
    <View className="flex-1 bg-white justify-center items-center">
      <Text className="text-2xl font-bold text-amber-900">Notifications</Text>
      <Text className="text-amber-700 mt-2">Stay updated with latest news</Text>
    </View>
  );
}
