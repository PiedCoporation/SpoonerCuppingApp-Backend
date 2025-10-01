import React from "react";
import { createNativeStackNavigator } from "@react-navigation/native-stack";
import UsersScreen from "../screens/UsersScreen";
import CuppingEventDetailScreen from "../screens/CuppingEventDetailScreen";
import CuppingRegistrationMinimalist from "@/screens/CuppingRegistrationMinimalist";
import CuppingRegistrationOverview from "@/screens/CuppingRegistrationOverview";
import CreateCuppingEvent from "@/screens/CreateCuppingEvent";
import ProtectedRoute from "../components/auth/ProtectedRoute";
import MainTabNavigator from "./MainTabNavigator";

const Stack = createNativeStackNavigator();

const MainNavigator: React.FC = () => {
  return (
    <Stack.Navigator initialRouteName="MainTabs">
      <Stack.Screen
        name="MainTabs"
        component={MainTabNavigator}
        options={{ headerShown: false }}
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
        options={({ route }) => ({
          title: (route.params as any)?.event?.name ?? "Event Details",
          headerBackTitle: "Home",
          headerStyle: { backgroundColor: "#552507" },
          headerTintColor: "#fff",
        })}
      />
      <Stack.Screen
        name="CreateCuppingEvent"
        component={CreateCuppingEvent}
        options={{
          title: "Create Event",
          headerBackTitle: "Profile",
          headerStyle: { backgroundColor: "#552507" },
          headerTintColor: "#fff",
        }}
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
