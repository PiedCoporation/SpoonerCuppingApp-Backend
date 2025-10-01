import React from "react";
import { View, Text, ScrollView, TouchableOpacity, Image } from "react-native";
import { Event } from "@/types/event";
import { mockEvents } from "@/data/events.mock";
import EventPostCard from "@/components/ui/EventPostCard";
import { useNavigation } from "@react-navigation/native";
import { getMockSamplesForEvent } from "@/data/samples.mock";

const HomeScreen: React.FC = () => {
  const navigation = useNavigation();

  const handleEventPress = (eventData: any) => {
    (navigation as any).navigate("CuppingEventDetail", eventData);
  };

  const handleCuppingRegistrationPress = () => {
    (navigation as any).navigate("CuppingRegistrationMinimalist");
  };

  const handleCuppingRegistration1Press = () => {
    (navigation as any).navigate("CuppingRegistrationOverview");
  };

  const handlePressEvent = (event: Event) => {
    (navigation as any).navigate("CuppingEventDetail", {
      event,
      samples: getMockSamplesForEvent(event.id),
    });
  };

  return (
    <ScrollView className="flex-1 bg-gray-50">
      {/* Header Section */}
      <View className="bg-amber-900 px-6 py-8 rounded-b-3xl">
        <Text className="text-white text-2xl mt-4 font-bold mb-2">
          Coffee Cupping Mall
        </Text>
        <Text className="text-amber-100 text-base">
          Discover, taste, and rate premium coffee beans
        </Text>
      </View>

      {/* Quick Actions */}
      <View className="px-6 py-6">
        <Text className="text-gray-800 text-xl font-semibold mb-4">
          Quick Actions
        </Text>

        <View className="flex-row justify-between">
          <TouchableOpacity
            onPress={handleCuppingRegistrationPress}
            className="bg-white p-4 rounded-xl flex-1 mr-2 shadow-sm border border-gray-100"
          >
            <View className="items-center">
              <Text className="text-2xl mb-2">☕</Text>
              <Text className="text-gray-700 font-medium">New Cupping</Text>
              <Text className="text-gray-500 text-sm">Start tasting</Text>
            </View>
          </TouchableOpacity>

          <TouchableOpacity
            onPress={handleCuppingRegistration1Press}
            className="bg-white p-4 rounded-xl flex-1 ml-2 shadow-sm border border-gray-100"
          >
            <View className="items-center">
              <Text className="text-2xl mb-2">📊</Text>
              <Text className="text-gray-700 font-medium">My Ratings</Text>
              <Text className="text-gray-500 text-sm">View history</Text>
            </View>
          </TouchableOpacity>
        </View>
      </View>

      {/* Cupping Events Section */}
      <View className="px-6 py-4">
        <View className="flex-row justify-between items-center mb-4">
          <Text className="text-gray-800 text-xl font-semibold">
            Coffee News Feed
          </Text>
        </View>

        <View className="space-y-3">
          {mockEvents.map((event) => (
            <EventPostCard
              key={event.id}
              event={event}
              onPress={handlePressEvent}
              commentCount={10}
              likeCount={10}
            />
          ))}
        </View>
      </View>
    </ScrollView>
  );
};

export default HomeScreen;
