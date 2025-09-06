import React from "react";
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  TextInput,
} from "react-native";
import { MainTabProps } from "../../types/navigation";

export default function ExploreScreen({ navigation }: MainTabProps<"Explore">) {
  return (
    <ScrollView className="flex-1 bg-cream-50">
      <View className="px-6 pt-12 pb-6">
        <Text className="text-3xl font-bold text-coffee-800 mb-4">
          Explore Coffee
        </Text>

        {/* Search Bar */}
        <TextInput
          className="bg-white border border-cream-200 rounded-lg px-4 py-3"
          placeholder="Search coffee origins, farms..."
          placeholderTextColor="#8B7355"
        />
      </View>

      {/* Categories */}
      <View className="px-6 mb-6">
        <Text className="text-lg font-semibold text-coffee-800 mb-4">
          Browse by Origin
        </Text>
        <ScrollView horizontal showsHorizontalScrollIndicator={false}>
          <View className="flex-row space-x-3">
            {["Ethiopia", "Colombia", "Brazil", "Guatemala", "Kenya"].map(
              (origin) => (
                <TouchableOpacity
                  key={origin}
                  className="bg-coffee-100 rounded-full px-4 py-2"
                >
                  <Text className="text-coffee-700 font-medium">{origin}</Text>
                </TouchableOpacity>
              )
            )}
          </View>
        </ScrollView>
      </View>

      {/* Featured Coffees */}
      <View className="px-6">
        <Text className="text-lg font-semibold text-coffee-800 mb-4">
          Featured Coffees
        </Text>
        <View className="space-y-4">
          {[
            {
              name: "Ethiopian Yirgacheffe",
              farm: "Konga Cooperative",
              score: 88.5,
            },
            { name: "Colombian Huila", farm: "Finca El Paraiso", score: 86.2 },
            { name: "Kenyan AA", farm: "Kiambu Estate", score: 87.8 },
          ].map((coffee, index) => (
            <TouchableOpacity
              key={index}
              className="bg-white rounded-lg p-4 border border-cream-200"
            >
              <Text className="font-semibold text-coffee-800 text-lg">
                {coffee.name}
              </Text>
              <Text className="text-coffee-600 mt-1">{coffee.farm}</Text>
              <View className="flex-row justify-between mt-3">
                <Text className="text-coffee-500">Average Score</Text>
                <Text className="font-semibold text-coffee-700">
                  {coffee.score}
                </Text>
              </View>
            </TouchableOpacity>
          ))}
        </View>
      </View>
    </ScrollView>
  );
}
