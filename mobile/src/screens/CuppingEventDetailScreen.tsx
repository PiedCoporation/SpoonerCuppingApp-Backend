import React from "react";
import { View, Text, ScrollView, TouchableOpacity, Image } from "react-native";
import { useRoute, useNavigation } from "@react-navigation/native";
import { EventDetail } from "@/types/event";

const CuppingEventDetailScreen: React.FC = () => {
  const route = useRoute();
  const navigation = useNavigation();
  const { event, samples } = route.params as unknown as EventDetail;

  const getStatusColor = (status: string) => {
    switch (status) {
      case "Open":
        return "bg-green-100 text-green-700";
      case "Tomorrow":
        return "bg-blue-100 text-blue-700";
      case "Premium":
        return "bg-purple-100 text-purple-700";
      case "Free":
        return "bg-yellow-100 text-yellow-700";
      default:
        return "bg-gray-100 text-gray-700";
    }
  };

  return (
    <ScrollView className="flex-1 bg-gray-50">
      {/* Header Image Section */}
      <View className="h-48 bg-amber-100 items-center justify-center">
        {Boolean(event.event_image) ? (
          // @ts-ignore className support by NativeWind
          <Image
            source={{ uri: event.event_image as string }}
            className="w-full h-48"
            resizeMode="cover"
          />
        ) : (
          <Text className="text-6xl">☕</Text>
        )}
      </View>

      {/* Event Info Section */}
      <View className="px-6 py-6 bg-white">
        <View className="flex-row justify-between items-start mb-4">
          <Text className="text-2xl font-bold text-gray-800 flex-1 mr-4">
            {event.name}
          </Text>
          <View
            className={`px-3 py-1 rounded-full ${getStatusColor(
              event.register_status
            )}`}
          >
            <Text
              className={`text-sm font-medium ${
                getStatusColor(event.register_status).split(" ")[1]
              }`}
            >
              {event.register_status}
            </Text>
          </View>
        </View>

        <Text className="text-gray-600 text-base leading-6 mb-6">
          {event.name} • {event.number_samples} samples
        </Text>

        {/* Event Details */}
        <View className="space-y-4">
          <View className="flex-row items-center">
            <Text className="text-2xl mr-4">📅</Text>
            <View>
              <Text className="text-gray-800 font-semibold">Date & Time</Text>
              <Text className="text-gray-600">
                {event.date_of_event} {event.start_time} - {event.end_time}
              </Text>
            </View>
          </View>

          <View className="flex-row items-center">
            <Text className="text-2xl mr-4">👥</Text>
            <View>
              <Text className="text-gray-800 font-semibold">Participants</Text>
              <Text className="text-gray-600">Limit {event.limit}</Text>
            </View>
          </View>

          <View className="flex-row items-center">
            <Text className="text-2xl mr-4">📍</Text>
            <View>
              <Text className="text-gray-800 font-semibold">Location</Text>
              <Text className="text-gray-600">
                {event.event_address?.[0]?.street},{" "}
                {event.event_address?.[0]?.district},{" "}
                {event.event_address?.[0]?.province}
              </Text>
            </View>
          </View>

          <View className="flex-row items-center">
            <Text className="text-2xl mr-4">⏱️</Text>
            <View>
              <Text className="text-gray-800 font-semibold">Duration</Text>
              <Text className="text-gray-600">2 hours</Text>
            </View>
          </View>
        </View>
      </View>

      {/* What to Expect Section */}
      <View className="px-6 py-6 bg-white mt-4">
        <Text className="text-xl font-bold text-gray-800 mb-4">
          What to Expect
        </Text>

        <View className="space-y-3">
          <View className="flex-row items-start">
            <View className="w-6 h-6 bg-amber-500 rounded-full items-center justify-center mr-3 mt-1">
              <Text className="text-white text-xs font-bold">1</Text>
            </View>
            <View className="flex-1">
              <Text className="text-gray-800 font-medium">
                Introduction & Setup
              </Text>
              <Text className="text-gray-600 text-sm">
                Brief overview of cupping process and coffee origins
              </Text>
            </View>
          </View>

          <View className="flex-row items-start">
            <View className="w-6 h-6 bg-amber-500 rounded-full items-center justify-center mr-3 mt-1">
              <Text className="text-white text-xs font-bold">2</Text>
            </View>
            <View className="flex-1">
              <Text className="text-gray-800 font-medium">
                Aroma Evaluation
              </Text>
              <Text className="text-gray-600 text-sm">
                Assess dry and wet fragrance of coffee grounds
              </Text>
            </View>
          </View>

          <View className="flex-row items-start">
            <View className="w-6 h-6 bg-amber-500 rounded-full items-center justify-center mr-3 mt-1">
              <Text className="text-white text-xs font-bold">3</Text>
            </View>
            <View className="flex-1">
              <Text className="text-gray-800 font-medium">
                Tasting & Scoring
              </Text>
              <Text className="text-gray-600 text-sm">
                Systematic evaluation of flavor, acidity, body, and finish
              </Text>
            </View>
          </View>

          <View className="flex-row items-start">
            <View className="w-6 h-6 bg-amber-500 rounded-full items-center justify-center mr-3 mt-1">
              <Text className="text-white text-xs font-bold">4</Text>
            </View>
            <View className="flex-1">
              <Text className="text-gray-800 font-medium">
                Discussion & Comparison
              </Text>
              <Text className="text-gray-600 text-sm">
                Share findings and compare notes with other participants
              </Text>
            </View>
          </View>
        </View>
      </View>

      {/* Featured Coffees for this Event */}
      <View className="px-6 py-6 bg-white mt-4">
        <Text className="text-xl font-bold text-gray-800 mb-4">
          Featured Coffees
        </Text>

        <View className="space-y-3">
          {samples && samples.length > 0 ? (
            samples.map((s) => (
              <View key={s.id} className="bg-gray-50 rounded-xl p-4">
                <View className="flex-row items-center mb-2">
                  <Text className="text-lg font-semibold text-gray-800">
                    {s.name}
                  </Text>
                  <View className="ml-auto bg-amber-100 px-2 py-1 rounded">
                    <Text className="text-amber-700 text-xs font-medium">
                      {s.grow_nation}
                    </Text>
                  </View>
                </View>
                <Text className="text-gray-600 text-sm mb-2">
                  {s.pre_processing} • Roast: {s.roast_level}
                </Text>
                <View className="flex-row items-center">
                  <Text className="text-gray-500 text-sm">
                    Altitude: {s.altitude_grow}
                  </Text>
                  <Text className="text-gray-400 mx-2">•</Text>
                  <Text className="text-gray-500 text-sm">
                    Roastery: {s.roastery_name}
                  </Text>
                </View>
              </View>
            ))
          ) : (
            <Text className="text-gray-500">No samples provided.</Text>
          )}
        </View>
      </View>

      {/* Requirements Section */}
      <View className="px-6 py-6 bg-white mt-4">
        <Text className="text-xl font-bold text-gray-800 mb-4">
          Requirements
        </Text>

        <View className="space-y-2">
          <View className="flex-row items-center">
            <Text className="text-green-600 mr-2">✓</Text>
            <Text className="text-gray-600">
              No prior cupping experience required
            </Text>
          </View>
          <View className="flex-row items-center">
            <Text className="text-green-600 mr-2">✓</Text>
            <Text className="text-gray-600">
              All materials and equipment provided
            </Text>
          </View>
          <View className="flex-row items-center">
            <Text className="text-green-600 mr-2">✓</Text>
            <Text className="text-gray-600">Notebook and pen recommended</Text>
          </View>
          <View className="flex-row items-center">
            <Text className="text-amber-600 mr-2">!</Text>
            <Text className="text-gray-600">
              Please avoid strong perfumes or scents
            </Text>
          </View>
        </View>
      </View>

      {/* Action Buttons */}
      <View className="px-6 py-6 bg-white mt-4 mb-6">
        <TouchableOpacity className="bg-amber-600 rounded-xl py-4 items-center mb-3">
          <Text className="text-white font-semibold text-lg">
            Join This Event
          </Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
};

export default CuppingEventDetailScreen;
