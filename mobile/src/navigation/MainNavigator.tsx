import React from "react";
import { NavigationContainer } from "@react-navigation/native";
import { createNativeStackNavigator } from "@react-navigation/native-stack";
import HomeScreen from "../screens/HomeScreen";
import UsersScreen from "../screens/UsersScreen";
import CuppingEventDetailScreen from "../screens/CuppingEventDetailScreen";
import CuppingRegistration from "@/screens/CuppingRegistration";

const Stack = createNativeStackNavigator();

const MainNavigator: React.FC = () => {
  return (
    <NavigationContainer>
      <Stack.Navigator initialRouteName="Home">
        <Stack.Screen
          name="Home"
          component={HomeScreen}
          options={{ title: "Coffee Cupping Mall" }}
        />
        <Stack.Screen
          name="Users"
          component={UsersScreen}
          options={{ title: "Users List" }}
        />
        <Stack.Screen
          name="CuppingEventDetail"
          component={CuppingEventDetailScreen}
          options={{ title: "Event Details" }}
        />
        <Stack.Screen
          name="CuppingRegistration"
          component={CuppingRegistration}
          options={{
            title: "Register for Event",
            headerBackTitle: "Events",
            headerStyle: { backgroundColor: "#552507" },
            headerTintColor: "#fff",
          }}
        />
      </Stack.Navigator>
    </NavigationContainer>
  );
};

export default MainNavigator;
