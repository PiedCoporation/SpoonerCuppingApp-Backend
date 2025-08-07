import React from "react";
import { View, Text, ScrollView, TouchableOpacity, Image } from "react-native";
import { useRoute, useNavigation } from "@react-navigation/native";

interface CuppingEventDetailProps {
  route: {
    params: {
      eventId: string;
      eventTitle: string;
      eventDescription: string;
      eventDate: string;
      participants: string;
      status: string;
      eventIcon: string;
    };
  };
}

const CuppingEventDetailScreen: React.FC = () => {
  const route = useRoute();
  const navigation = useNavigation();
  const {
    eventTitle,
    eventDescription,
    eventDate,
    participants,
    status,
    eventIcon,
  } = route.params as any;

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
        <Text className="text-6xl">{eventIcon}</Text>
      </View>

      {/* Event Info Section */}
      <View className="px-6 py-6 bg-white">
        <View className="flex-row justify-between items-start mb-4">
          <Text className="text-2xl font-bold text-gray-800 flex-1 mr-4">
            {eventTitle}
          </Text>
          <View className={`px-3 py-1 rounded-full ${getStatusColor(status)}`}>
            <Text
              className={`text-sm font-medium ${
                getStatusColor(status).split(" ")[1]
              }`}
            >
              {status}
            </Text>
          </View>
        </View>

        <Text className="text-gray-600 text-base leading-6 mb-6">
          {eventDescription}
        </Text>

        {/* Event Details */}
        <View className="space-y-4">
          <View className="flex-row items-center">
            <Text className="text-2xl mr-4">📅</Text>
            <View>
              <Text className="text-gray-800 font-semibold">Date & Time</Text>
              <Text className="text-gray-600">{eventDate}</Text>
            </View>
          </View>

          <View className="flex-row items-center">
            <Text className="text-2xl mr-4">👥</Text>
            <View>
              <Text className="text-gray-800 font-semibold">Participants</Text>
              <Text className="text-gray-600">{participants}</Text>
            </View>
          </View>

          <View className="flex-row items-center">
            <Text className="text-2xl mr-4">📍</Text>
            <View>
              <Text className="text-gray-800 font-semibold">Location</Text>
              <Text className="text-gray-600">
                Coffee Cupping Lab, 2nd Floor
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
          <View className="bg-gray-50 rounded-xl p-4">
            <View className="flex-row items-center mb-2">
              <Text className="text-lg font-semibold text-gray-800">
                Ethiopian Yirgacheffe
              </Text>
              <View className="ml-auto bg-amber-100 px-2 py-1 rounded">
                <Text className="text-amber-700 text-xs font-medium">
                  Single Origin
                </Text>
              </View>
            </View>
            <Text className="text-gray-600 text-sm mb-2">
              Floral, citrusy notes with bright acidity and tea-like body
            </Text>
            <View className="flex-row items-center">
              <Text className="text-gray-500 text-sm">
                Altitude: 1,800-2,200m
              </Text>
              <Text className="text-gray-400 mx-2">•</Text>
              <Text className="text-gray-500 text-sm">Process: Washed</Text>
            </View>
          </View>

          <View className="bg-gray-50 rounded-xl p-4">
            <View className="flex-row items-center mb-2">
              <Text className="text-lg font-semibold text-gray-800">
                Ethiopian Sidamo
              </Text>
              <View className="ml-auto bg-amber-100 px-2 py-1 rounded">
                <Text className="text-amber-700 text-xs font-medium">
                  Single Origin
                </Text>
              </View>
            </View>
            <Text className="text-gray-600 text-sm mb-2">
              Wine-like characteristics with complex fruit flavors and medium
              body
            </Text>
            <View className="flex-row items-center">
              <Text className="text-gray-500 text-sm">
                Altitude: 1,400-2,200m
              </Text>
              <Text className="text-gray-400 mx-2">•</Text>
              <Text className="text-gray-500 text-sm">Process: Natural</Text>
            </View>
          </View>
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

        <View className="flex-row space-x-3">
          <TouchableOpacity className="flex-1 bg-gray-100 rounded-xl py-3 items-center">
            <Text className="text-gray-700 font-medium">Add to Calendar</Text>
          </TouchableOpacity>

          <TouchableOpacity className="flex-1 bg-gray-100 rounded-xl py-3 items-center">
            <Text className="text-gray-700 font-medium">Share Event</Text>
          </TouchableOpacity>
        </View>
      </View>
    </ScrollView>
  );
};

export default CuppingEventDetailScreen;
