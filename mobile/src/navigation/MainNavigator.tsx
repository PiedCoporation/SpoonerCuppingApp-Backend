import React from "react";
import { createNativeStackNavigator } from "@react-navigation/native-stack";
import HomeScreen from "../screens/HomeScreen";
import UsersScreen from "../screens/UsersScreen";
import CuppingEventDetailScreen from "../screens/CuppingEventDetailScreen";
import CuppingRegistrationMinimalist from "@/screens/CuppingRegistrationMinimalist";
import CuppingRegistrationOverview from "@/screens/CuppingRegistrationOverview";
import ProtectedRoute from "../components/auth/ProtectedRoute";

const Stack = createNativeStackNavigator();

const MainNavigator: React.FC = () => {
  return (
    <Stack.Navigator initialRouteName="Home">
      <Stack.Screen
        name="Home"
        component={HomeScreen}
        options={{ title: "Coffee Cupping Mall" }}
      />
      <Stack.Screen name="Users" options={{ title: "Users List" }}>
        {() => (
          <ProtectedRoute routeName="Users">
            <UsersScreen />
          </ProtectedRoute>
        )}
      </Stack.Screen>
      <Stack.Screen
        name="CuppingEventDetail"
        component={CuppingEventDetailScreen}
        options={{ title: "Event Details" }}
      />
      <Stack.Screen
        name="CuppingRegistrationMinimalist"
        component={CuppingRegistrationMinimalist}
        options={{
          title: "Register for Event",
          headerBackTitle: "Events",
          headerStyle: { backgroundColor: "#552507" },
          headerTintColor: "#fff",
        }}
      />
      <Stack.Screen
        name="CuppingRegistrationOverview"
        component={CuppingRegistrationOverview}
        options={{
          title: "Register for Event1",
          headerBackTitle: "Events1",
          headerStyle: { backgroundColor: "#552507" },
          headerTintColor: "#fff",
        }}
      />
    </Stack.Navigator>
  );
};

export default MainNavigator;
