import React from "react";
import { View, Text, TouchableOpacity, Image, ScrollView } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";
import { RootStackParamList } from "../../types/navigation";
import { useAuth } from "../../contexts/AuthContext";
import { Event } from "@/types/event";
import { mockEvents } from "@/data/events.mock";
import EventPostCard from "@/components/ui/EventPostCard";
import { getMockSamplesForEvent } from "@/data/samples.mock";

type ProfileNavigationProp = NativeStackNavigationProp<RootStackParamList>;

export default function ProfileScreen() {
  const navigation = useNavigation<ProfileNavigationProp>();
  const { logout, user } = useAuth();

  const handleAccountSecurity = () => {
    navigation.navigate("AccountSecurity");
  };

  const handleCuppingStyle = () => {
    navigation.navigate("CuppingStyle");
  };

  const handleCreateEvent = () => {
    navigation.navigate("CreateCuppingEvent");
  };

  const handleLogout = () => {
    logout();
    // navigation.dispatch(
    //   CommonActions.reset({
    //     index: 0,
    //     routes: [{ name: "Auth" }],
    //   })
    // );
  };

  const handlePressEvent = (event: Event) => {
    (navigation as any).navigate("CuppingEventDetail", {
      event,
      samples: getMockSamplesForEvent(event.id),
    });
  };

  const organizedHostId =
    user?.id && user.id.startsWith("u_") ? user.id : "u_1";
  const organizedEvents = mockEvents.filter(
    (e) => e.host_by.id === organizedHostId
  );

  return (
    <SafeAreaView
      className="flex-1 bg-gray-100"
      edges={["left", "right", "bottom"]}
    >
      <ScrollView className="flex-1">
        {/* Profile Header */}
        <View className="bg-white px-6 py-8 items-center">
          <Image
            source={{ uri: "https://picsum.photos/100/100" }}
            className="w-20 h-20 rounded-full mb-3"
          />
          <Text className="text-xl font-medium text-gray-900 mb-1">
            {user?.name ?? "Guest"}
          </Text>
          <Text className="text-sm text-gray-500">{user?.email ?? ""}</Text>
        </View>

        {/* Organized Events */}
        <View className="px-6 pt-6">
          {/* Create Event CTA */}
          <View className="bg-white rounded-xl p-4 mb-4 border border-gray-100">
            <Text className="text-gray-900 font-semibold mb-2">
              Create a new cupping event
            </Text>
            <Text className="text-gray-500 text-sm mb-3">
              Set up date, samples, and invite participants.
            </Text>
            <TouchableOpacity
              onPress={handleCreateEvent}
              className="bg-amber-600 rounded-lg py-3 items-center"
              activeOpacity={0.9}
            >
              <Text className="text-white font-medium">Create Event</Text>
            </TouchableOpacity>
          </View>

          <Text className="text-gray-800 text-xl font-semibold mb-4">
            Organized Events
          </Text>
          {organizedEvents.length === 0 ? (
            <Text className="text-gray-500">No events organized yet.</Text>
          ) : (
            <View className="space-y-3">
              {organizedEvents.map((ev) => (
                <EventPostCard
                  key={ev.id}
                  event={ev}
                  onPress={handlePressEvent}
                />
              ))}
            </View>
          )}
        </View>

        {/* Minimalist Menu */}
        {/* <View className="px-6 pt-6">
          <TouchableOpacity
            className="bg-white rounded-lg p-4 mb-3 flex-row items-center"
            onPress={handleAccountSecurity}
          >
            <Text className="flex-1 text-base text-gray-900">
              Account & Security
            </Text>
            <Text className="text-gray-400">›</Text>
          </TouchableOpacity>

          <TouchableOpacity
            className="bg-white rounded-lg p-4 mb-6 flex-row items-center"
            onPress={handleCuppingStyle}
          >
            <Text className="flex-1 text-base text-gray-900">
              Cupping Style
            </Text>
            <Text className="text-gray-400">›</Text>
          </TouchableOpacity>

          <TouchableOpacity
            className="bg-white rounded-lg p-4 mb-8 items-center border border-gray-200"
            onPress={handleLogout}
          >
            <Text className="text-base text-red-600 font-medium">Logout</Text>
          </TouchableOpacity>
        </View> */}
      </ScrollView>
    </SafeAreaView>
  );
}
