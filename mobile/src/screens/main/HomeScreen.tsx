import React from "react";
import { View, Text, ScrollView, TouchableOpacity } from "react-native";
import { MainTabProps } from "../../types/navigation";

export default function HomeScreen({ navigation }: MainTabProps<"Home">) {
  return (
    <ScrollView className="flex-1 bg-cream-50">
      <View className="px-6 pt-12 pb-6">
        <Text className="text-3xl font-bold text-coffee-800 mb-2">
          Good Morning ☕
        </Text>
        <Text className="text-coffee-600 text-lg">
          Ready for today's cupping session?
        </Text>
      </View>

      {/* Quick Actions */}
      <View className="px-6 mb-6">
        <Text className="text-xl font-semibold text-coffee-800 mb-4">
          Quick Actions
        </Text>
        <View className="flex-row space-x-4">
          <TouchableOpacity className="flex-1 bg-coffee-600 rounded-lg p-4">
            <Text className="text-white font-semibold text-center">
              New Session
            </Text>
          </TouchableOpacity>
          <TouchableOpacity className="flex-1 bg-cream-300 rounded-lg p-4">
            <Text className="text-coffee-800 font-semibold text-center">
              View Results
            </Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* Recent Sessions */}
      <View className="px-6">
        <Text className="text-xl font-semibold text-coffee-800 mb-4">
          Recent Sessions
        </Text>
        <View className="space-y-3">
          {[1, 2, 3].map((item) => (
            <TouchableOpacity
              key={item}
              className="bg-white rounded-lg p-4 border border-cream-200"
            >
              <Text className="font-semibold text-coffee-800">
                Ethiopian Yirgacheffe
              </Text>
              <Text className="text-coffee-600 text-sm mt-1">
                March {item}, 2024 • Score: 85.5
              </Text>
            </TouchableOpacity>
          ))}
        </View>
      </View>
    </ScrollView>
  );
}
