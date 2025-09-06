import React from "react";
import {
  View,
  Text,
  TouchableOpacity,
  Image,
  SafeAreaView,
} from "react-native";
import { useNavigation } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";
import { RootStackParamList } from "../../types/navigation";
import { CommonActions } from "@react-navigation/native";
import { useAuth } from "../../contexts/AuthContext";

type ProfileNavigationProp = NativeStackNavigationProp<RootStackParamList>;

export default function ProfileScreen() {
  const navigation = useNavigation<ProfileNavigationProp>();
  const { logout } = useAuth();

  const handleAccountSecurity = () => {
    navigation.navigate("AccountSecurity");
  };

  const handleCuppingStyle = () => {
    navigation.navigate("CuppingStyle");
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

  return (
    <SafeAreaView className="flex-1 bg-gray-100">
      {/* Minimalist Header */}
      <View className="bg-white px-6 py-8 items-center">
        <Image
          source={{
            uri: "https://picsum.photos/100/100",
          }}
          className="w-20 h-20 rounded-full mb-3"
        />
        <Text className="text-xl font-medium text-gray-900 mb-1">John Doe</Text>
        <Text className="text-sm text-gray-500">Coffee Enthusiast</Text>
      </View>

      {/* Minimalist Menu */}
      <View className="flex-1 px-6 pt-6">
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
          <Text className="flex-1 text-base text-gray-900">Cupping Style</Text>
          <Text className="text-gray-400">›</Text>
        </TouchableOpacity>

        {/* Logout Button */}
        <TouchableOpacity
          className="bg-white rounded-lg p-4 mt-auto mb-8 items-center border border-gray-200"
          onPress={handleLogout}
        >
          <Text className="text-base text-red-600 font-medium">Logout</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}
