import React from "react";
import { View, Text, ScrollView, TouchableOpacity, Image } from "react-native";
import { useNavigation } from "@react-navigation/native";

const HomeScreen: React.FC = () => {
  const navigation = useNavigation();

  const handleEventPress = (eventData: any) => {
    (navigation as any).navigate("CuppingEventDetail", eventData);
  };

  const handleCuppingRegistrationPress = () => {
    (navigation as any).navigate("CuppingRegistration");
  };

  return (
    <ScrollView className="flex-1 bg-gray-50">
      {/* Header Section */}
      <View className="bg-amber-900 px-6 py-8 rounded-b-3xl">
        <Text className="text-white text-2xl font-bold mb-2">
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

        <View className="flex-row justify-between mb-4">
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

          <TouchableOpacity className="bg-white p-4 rounded-xl flex-1 ml-2 shadow-sm border border-gray-100">
            <View className="items-center">
              <Text className="text-2xl mb-2">📊</Text>
              <Text className="text-gray-700 font-medium">My Ratings</Text>
              <Text className="text-gray-500 text-sm">View history</Text>
            </View>
          </TouchableOpacity>
        </View>

        <View className="flex-row justify-between">
          <TouchableOpacity className="bg-white p-4 rounded-xl flex-1 mr-2 shadow-sm border border-gray-100">
            <View className="items-center">
              <Text className="text-2xl mb-2">🛒</Text>
              <Text className="text-gray-700 font-medium">Coffee Shop</Text>
              <Text className="text-gray-500 text-sm">Browse beans</Text>
            </View>
          </TouchableOpacity>

          <TouchableOpacity className="bg-white p-4 rounded-xl flex-1 ml-2 shadow-sm border border-gray-100">
            <View className="items-center">
              <Text className="text-2xl mb-2">👥</Text>
              <Text className="text-gray-700 font-medium">Community</Text>
              <Text className="text-gray-500 text-sm">Connect & share</Text>
            </View>
          </TouchableOpacity>
        </View>
      </View>

      {/* Cupping Events Section */}
      <View className="px-6 py-4">
        <View className="flex-row justify-between items-center mb-4">
          <Text className="text-gray-800 text-xl font-semibold">
            Upcoming Cupping Events
          </Text>
          <TouchableOpacity>
            <Text className="text-amber-600 font-medium">View All</Text>
          </TouchableOpacity>
        </View>

        <View className="space-y-3">
          {/* Event Card 1 */}
          <TouchableOpacity
            className="bg-white rounded-xl p-4 shadow-sm border border-gray-100 mb-3"
            onPress={() =>
              handleEventPress({
                eventId: "1",
                eventTitle: "Ethiopian Highlands Tasting",
                eventDescription:
                  "Explore single-origin beans from Yirgacheffe and Sidamo regions. This comprehensive tasting session will guide you through the unique characteristics of Ethiopian coffee, known for its floral notes and bright acidity.",
                eventDate: "Today, 3:00 PM",
                participants: "12/15 participants",
                status: "Open",
                eventIcon: "☕",
              })
            }
          >
            <View className="flex-row">
              <View className="w-16 h-16 bg-amber-100 rounded-lg items-center justify-center mr-4">
                <Text className="text-2xl">☕</Text>
              </View>
              <View className="flex-1">
                <View className="flex-row justify-between items-start mb-1">
                  <Text
                    className="text-gray-800 font-semibold text-base flex-1 mr-2"
                    numberOfLines={2}
                  >
                    Ethiopian Highlands Tasting
                  </Text>
                  <View className="bg-green-100 px-2 py-1 rounded-full">
                    <Text className="text-green-700 text-xs font-medium">
                      Open
                    </Text>
                  </View>
                </View>
                <Text className="text-gray-600 text-sm mb-2" numberOfLines={2}>
                  Explore single-origin beans from Yirgacheffe and Sidamo
                  regions
                </Text>
                <View className="flex-row items-center justify-between">
                  <View className="flex-row items-center flex-1 mr-2">
                    <Text className="text-gray-500 text-sm" numberOfLines={1}>
                      📅 Today, 3:00 PM
                    </Text>
                  </View>
                  <View className="flex-row items-center">
                    <Text className="text-gray-500 text-sm" numberOfLines={1}>
                      👥 12/15 participants
                    </Text>
                  </View>
                </View>
              </View>
            </View>
          </TouchableOpacity>

          {/* Event Card 2 */}
          <TouchableOpacity
            className="bg-white rounded-xl p-4 shadow-sm border border-gray-100 mb-3"
            onPress={() =>
              handleEventPress({
                eventId: "2",
                eventTitle: "Organic Fair Trade Session",
                eventDescription:
                  "Compare sustainable coffees from Central America. Learn about ethical sourcing practices while enjoying coffees from Guatemala, Costa Rica, and Honduras that support local farming communities.",
                eventDate: "Tomorrow, 10:00 AM",
                participants: "8/12 participants",
                status: "Tomorrow",
                eventIcon: "🌱",
              })
            }
          >
            <View className="flex-row">
              <View className="w-16 h-16 bg-amber-100 rounded-lg items-center justify-center mr-4">
                <Text className="text-2xl">🌱</Text>
              </View>
              <View className="flex-1">
                <View className="flex-row justify-between items-start mb-1">
                  <Text
                    className="text-gray-800 font-semibold text-base flex-1 mr-2"
                    numberOfLines={2}
                  >
                    Organic Fair Trade Session
                  </Text>
                  <View className="bg-blue-100 px-2 py-1 rounded-full">
                    <Text className="text-blue-700 text-xs font-medium">
                      Tomorrow
                    </Text>
                  </View>
                </View>
                <Text className="text-gray-600 text-sm mb-2" numberOfLines={2}>
                  Compare sustainable coffees from Central America
                </Text>
                <View className="flex-row items-center justify-between">
                  <View className="flex-row items-center flex-1 mr-2">
                    <Text className="text-gray-500 text-sm" numberOfLines={1}>
                      📅 Tomorrow, 10:00 AM
                    </Text>
                  </View>
                  <View className="flex-row items-center">
                    <Text className="text-gray-500 text-sm" numberOfLines={1}>
                      👥 8/12 participants
                    </Text>
                  </View>
                </View>
              </View>
            </View>
          </TouchableOpacity>

          {/* Event Card 3 */}
          <TouchableOpacity
            className="bg-white rounded-xl p-4 shadow-sm border border-gray-100 mb-3"
            onPress={() =>
              handleEventPress({
                eventId: "3",
                eventTitle: "Championship Roast Finals",
                eventDescription:
                  "Professional cupping competition with award-winning roasters. Witness master roasters compete in this prestigious event and learn from the best in the industry.",
                eventDate: "Friday, 2:00 PM",
                participants: "5/8 participants",
                status: "Premium",
                eventIcon: "🏆",
              })
            }
          >
            <View className="flex-row">
              <View className="w-16 h-16 bg-amber-100 rounded-lg items-center justify-center mr-4">
                <Text className="text-2xl">🏆</Text>
              </View>
              <View className="flex-1">
                <View className="flex-row justify-between items-start mb-1">
                  <Text
                    className="text-gray-800 font-semibold text-base flex-1 mr-2"
                    numberOfLines={2}
                  >
                    Championship Roast Finals
                  </Text>
                  <View className="bg-purple-100 px-2 py-1 rounded-full">
                    <Text className="text-purple-700 text-xs font-medium">
                      Premium
                    </Text>
                  </View>
                </View>
                <Text className="text-gray-600 text-sm mb-2" numberOfLines={2}>
                  Professional cupping competition with award-winning roasters
                </Text>
                <View className="flex-row items-center justify-between">
                  <View className="flex-row items-center flex-1 mr-2">
                    <Text className="text-gray-500 text-sm" numberOfLines={1}>
                      📅 Friday, 2:00 PM
                    </Text>
                  </View>
                  <View className="flex-row items-center">
                    <Text className="text-gray-500 text-sm" numberOfLines={1}>
                      👥 5/8 participants
                    </Text>
                  </View>
                </View>
              </View>
            </View>
          </TouchableOpacity>

          {/* Event Card 4 */}
          <TouchableOpacity
            className="bg-white rounded-xl p-4 shadow-sm border border-gray-100"
            onPress={() =>
              handleEventPress({
                eventId: "4",
                eventTitle: "Beginner's Cupping Workshop",
                eventDescription:
                  "Learn the basics of coffee cupping and flavor identification. Perfect for newcomers to coffee tasting, this workshop covers fundamental techniques and vocabulary.",
                eventDate: "Saturday, 11:00 AM",
                participants: "18/20 participants",
                status: "Free",
                eventIcon: "🎓",
              })
            }
          >
            <View className="flex-row">
              <View className="w-16 h-16 bg-amber-100 rounded-lg items-center justify-center mr-4">
                <Text className="text-2xl">🎓</Text>
              </View>
              <View className="flex-1">
                <View className="flex-row justify-between items-start mb-1">
                  <Text
                    className="text-gray-800 font-semibold text-base flex-1 mr-2"
                    numberOfLines={2}
                  >
                    Beginner's Cupping Workshop
                  </Text>
                  <View className="bg-yellow-100 px-2 py-1 rounded-full">
                    <Text className="text-yellow-700 text-xs font-medium">
                      Free
                    </Text>
                  </View>
                </View>
                <Text className="text-gray-600 text-sm mb-2" numberOfLines={2}>
                  Learn the basics of coffee cupping and flavor identification
                </Text>
                <View className="flex-row items-center justify-between">
                  <View className="flex-row items-center flex-1 mr-2">
                    <Text className="text-gray-500 text-sm" numberOfLines={1}>
                      📅 Saturday, 11:00 AM
                    </Text>
                  </View>
                  <View className="flex-row items-center">
                    <Text className="text-gray-500 text-sm" numberOfLines={1}>
                      👥 18/20 participants
                    </Text>
                  </View>
                </View>
              </View>
            </View>
          </TouchableOpacity>
        </View>
      </View>

      {/* Recent Activity */}
      <View className="px-6 py-4 mb-6">
        <Text className="text-gray-800 text-xl font-semibold mb-4">
          Recent Activity
        </Text>

        <View className="bg-white rounded-xl p-4 shadow-sm border border-gray-100">
          <View className="flex-row items-center mb-3">
            <View className="w-10 h-10 bg-green-100 rounded-full items-center justify-center mr-3">
              <Text className="text-green-600 font-semibold">✓</Text>
            </View>
            <View className="flex-1">
              <Text className="text-gray-800 font-medium">
                Completed cupping session
              </Text>
              <Text className="text-gray-500 text-sm">
                Ethiopian Sidamo - Rated 4.5/5
              </Text>
            </View>
            <Text className="text-gray-400 text-sm">2h ago</Text>
          </View>

          <View className="flex-row items-center mb-3">
            <View className="w-10 h-10 bg-blue-100 rounded-full items-center justify-center mr-3">
              <Text className="text-blue-600 font-semibold">+</Text>
            </View>
            <View className="flex-1">
              <Text className="text-gray-800 font-medium">
                Added new coffee to wishlist
              </Text>
              <Text className="text-gray-500 text-sm">Costa Rican Tarrazú</Text>
            </View>
            <Text className="text-gray-400 text-sm">1d ago</Text>
          </View>

          <View className="flex-row items-center">
            <View className="w-10 h-10 bg-purple-100 rounded-full items-center justify-center mr-3">
              <Text className="text-purple-600 font-semibold">🏆</Text>
            </View>
            <View className="flex-1">
              <Text className="text-gray-800 font-medium">
                Achievement unlocked
              </Text>
              <Text className="text-gray-500 text-sm">
                Coffee Connoisseur - 50 tastings completed
              </Text>
            </View>
            <Text className="text-gray-400 text-sm">3d ago</Text>
          </View>
        </View>
      </View>
    </ScrollView>
  );
};

export default HomeScreen;
