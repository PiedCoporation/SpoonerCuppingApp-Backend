import React from "react";
import { View, Text, ScrollView, TouchableOpacity } from "react-native";
import { MainTabProps } from "../../types/navigation";

export default function SessionsScreen({
  navigation,
}: MainTabProps<"Sessions">) {
  return (
    <View className="flex-1 bg-cream-50">
      <View className="px-6 pt-12 pb-6">
        <Text className="text-3xl font-bold text-coffee-800 mb-2">
          Cupping Sessions
        </Text>
        <Text className="text-coffee-600">
          Track and review your coffee evaluations
        </Text>
      </View>

      <ScrollView className="flex-1 px-6">
        <TouchableOpacity className="bg-coffee-600 rounded-lg p-4 mb-6">
          <Text className="text-white font-semibold text-center text-lg">
            + Start New Session
          </Text>
        </TouchableOpacity>

        <Text className="text-lg font-semibold text-coffee-800 mb-4">
          All Sessions
        </Text>

        {/* Session List */}
        <View className="space-y-4">
          {[1, 2, 3, 4, 5].map((item) => (
            <TouchableOpacity
              key={item}
              className="bg-white rounded-lg p-4 border border-cream-200"
            >
              <View className="flex-row justify-between items-start mb-2">
                <Text className="font-semibold text-coffee-800 flex-1">
                  Brazilian Santos
                </Text>
                <Text className="text-coffee-500 text-sm">Mar {item}</Text>
              </View>
              <Text className="text-coffee-600 text-sm mb-2">
                Medium roast • Washed process
              </Text>
              <View className="flex-row justify-between">
                <Text className="text-coffee-500 text-sm">Overall Score</Text>
                <Text className="font-semibold text-coffee-700">8{item}.2</Text>
              </View>
            </TouchableOpacity>
          ))}
        </View>
      </ScrollView>
    </View>
  );
}
