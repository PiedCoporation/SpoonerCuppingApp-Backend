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
    <SafeAreaView className="flex-1 ">
      <View className="bg-amber-800 py-10 px-5 items-center border-b border-amber-900">
        <View className="items-center">
          <Image
            source={{
              uri: "https://picsum.photos/100/100",
            }}
            className="w-24 h-24 rounded-full mb-4 border-4 border-amber-200"
          />
          <Text className="text-2xl font-semibold text-white mb-1">
            John Doe
          </Text>
          <Text className="text-base text-amber-100">Coffee Enthusiast</Text>
        </View>
      </View>

      <View className="flex-1 pt-2">
        <TouchableOpacity
          className="bg-white rounded-xl p-5 flex-row items-center shadow-sm"
          onPress={handleAccountSecurity}
        >
          <View className="flex-1">
            <Text className="text-lg font-semibold text-amber-900 mb-1">
              Account and Security
            </Text>
          </View>
          <Text className="text-2xl text-orange-600 font-light">›</Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="bg-white rounded-xl p-5 mb-4 flex-row items-center shadow-sm"
          onPress={handleCuppingStyle}
        >
          <View className="flex-1">
            <Text className="text-lg font-semibold text-amber-900 mb-1">
              Cupping Style
            </Text>
          </View>
          <Text className="text-2xl text-orange-600 font-light">›</Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="bg-red-600 mx-5 rounded-xl p-4 mt-2 items-center shadow-sm"
          onPress={handleLogout}
        >
          <Text className="text-base font-semibold text-white">Logout</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}
