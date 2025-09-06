import React from "react";
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import { Ionicons } from "@expo/vector-icons";
import { MainTabParamList } from "../types/navigation";
import SessionsScreen from "../screens/main/SessionsScreen";
import ExploreScreen from "../screens/main/ExploreScreen";
import ProfileScreen from "../screens/main/ProfileScreen";
import FriendScreen from "../screens/main/FriendScreen";
import MarketScreen from "../screens/main/MarketScreen";
import NotificationScreen from "../screens/main/NotificationScreen";
import HomeScreen from "@/screens/HomeScreen";

const Tab = createBottomTabNavigator<MainTabParamList>();

export default function MainTabNavigator() {
  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        headerShown: false,
        tabBarActiveTintColor: "#8B4513",
        tabBarInactiveTintColor: "#8B7355",
        tabBarStyle: {
          backgroundColor: "#FFFEF7",
          borderTopColor: "#E6D7C0",
        },
        tabBarIcon: ({ focused, color, size }) => {
          let iconName: keyof typeof Ionicons.glyphMap;

          switch (route.name) {
            case "Home":
              iconName = focused ? "home" : "home-outline";
              break;
            case "Friend":
              iconName = focused ? "people" : "people-outline";
              break;
            case "Market":
              iconName = focused ? "storefront" : "storefront-outline";
              break;
            case "Notification":
              iconName = focused ? "notifications" : "notifications-outline";
              break;
            case "Profile":
              iconName = focused ? "person" : "person-outline";
              break;
            default:
              iconName = "home-outline";
          }

          return <Ionicons name={iconName} size={size} color={color} />;
        },
      })}
    >
      <Tab.Screen
        name="Home"
        component={HomeScreen}
        options={{
          tabBarLabel: "Home",
        }}
      />
      <Tab.Screen
        name="Friend"
        component={FriendScreen}
        options={{
          tabBarLabel: "Friends",
        }}
      />
      <Tab.Screen
        name="Market"
        component={MarketScreen}
        options={{
          tabBarLabel: "Market",
        }}
      />
      <Tab.Screen
        name="Notification"
        component={NotificationScreen}
        options={{
          tabBarLabel: "Notifications",
        }}
      />
      <Tab.Screen
        name="Profile"
        component={ProfileScreen}
        options={{
          tabBarLabel: "Profile",
        }}
      />
    </Tab.Navigator>
  );
}
