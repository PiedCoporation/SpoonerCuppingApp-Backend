import React from "react";
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import { MaterialIcons } from "@expo/vector-icons";
import { MainTabParamList } from "../types/navigation";
import ProfileScreen from "../screens/main/ProfileScreen";
import FriendScreen from "../screens/main/FriendScreen";
import MarketScreen from "../screens/main/MarketScreen";
import NotificationScreen from "../screens/main/NotificationScreen";
import HomeScreen from "@/screens/HomeScreen";
import EventScreen from "@/screens/main/EventScreen";

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
          height: 80, // Increased height to accommodate text and bottom spacing
          paddingBottom: 20, // Extra bottom padding for safe area
          paddingTop: 10,
          paddingHorizontal: 10, // Add horizontal padding
        },
        tabBarLabelStyle: {
          fontSize: 10, // Smaller font size
          fontWeight: "500",
          marginBottom: 4,
          marginTop: 2,
        },
        tabBarItemStyle: {
          paddingHorizontal: 2, // Reduce horizontal padding
        },
        tabBarIcon: ({ focused, color, size }) => {
          let iconName: keyof typeof MaterialIcons.glyphMap;

          switch (route.name) {
            case "Home":
              iconName = "home";
              break;
            case "Event":
              iconName = "local-cafe";
              break;
            case "Friend":
              iconName = "people";
              break;
            case "Market":
              iconName = "storefront";
              break;
            case "Notification":
              iconName = "notifications";
              break;
            case "Profile":
              iconName = "person";
              break;
            default:
              iconName = "home";
          }

          return <MaterialIcons name={iconName} size={size} color={color} />;
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
        name="Event"
        component={EventScreen}
        options={{
          tabBarLabel: "Samples",
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
          tabBarLabel: "Alerts",
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
