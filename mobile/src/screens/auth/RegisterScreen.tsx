import React from "react";
import { View, Text, TextInput, TouchableOpacity, Alert } from "react-native";
import { useForm, Controller } from "react-hook-form";
import { AuthScreenProps } from "../../types/navigation";

type RegisterForm = {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
};

export default function RegisterScreen({
  navigation,
}: AuthScreenProps<"Register">) {
  const {
    control,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<RegisterForm>();
  const password = watch("password");

  const onRegister = (data: RegisterForm) => {
    // TODO: Implement registration logic
    console.log("Register:", data);
    Alert.alert("Success", "Account created successfully!", [
      {
        text: "Continue",
        onPress: () => navigation.navigate("MainApp" as any),
      },
    ]);
  };

  return (
    <View className="flex-1 bg-white justify-center">
      <View className="px-6">
        <View className="mb-8">
          <Text className="text-3xl font-bold text-amber-900 text-center mb-2">
            Join Spooner
          </Text>
          <Text className="text-lg text-amber-700 text-center">
            Create your account
          </Text>
        </View>

        <View className="space-y-4">
          <View>
            <Text className="text-amber-900 font-medium mb-2">Full Name</Text>
            <Controller
              control={control}
              name="name"
              rules={{
                required: "Name is required",
                minLength: {
                  value: 2,
                  message: "Name must be at least 2 characters",
                },
              }}
              render={({ field: { onChange, onBlur, value } }) => (
                <TextInput
                  className={`border-2 rounded-lg px-4 py-3 ${
                    errors.name ? "border-red-500" : "border-amber-900/40"
                  } focus:border-amber-500`}
                  placeholder="Enter your full name"
                  value={value}
                  onChangeText={onChange}
                  onBlur={onBlur}
                />
              )}
            />
            {errors.name && (
              <Text className="text-red-500 text-sm mt-1">
                {errors.name.message}
              </Text>
            )}
          </View>

          <View>
            <Text className="text-amber-900 font-medium mb-2">Email</Text>
            <Controller
              control={control}
              name="email"
              rules={{
                required: "Email is required",
                pattern: {
                  value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                  message: "Invalid email address",
                },
              }}
              render={({ field: { onChange, onBlur, value } }) => (
                <TextInput
                  className={`border-2 rounded-lg px-4 py-3 ${
                    errors.email ? "border-red-500" : "border-amber-900/40"
                  } focus:border-amber-500`}
                  placeholder="Enter your email"
                  value={value}
                  onChangeText={onChange}
                  onBlur={onBlur}
                  keyboardType="email-address"
                  autoCapitalize="none"
                />
              )}
            />
            {errors.email && (
              <Text className="text-red-500 text-sm mt-1">
                {errors.email.message}
              </Text>
            )}
          </View>

          <View>
            <Text className="text-amber-900 font-medium mb-2">Password</Text>
            <Controller
              control={control}
              name="password"
              rules={{
                required: "Password is required",
                minLength: {
                  value: 8,
                  message: "Password must be at least 8 characters",
                },
              }}
              render={({ field: { onChange, onBlur, value } }) => (
                <TextInput
                  className={`border-2 rounded-lg px-4 py-3 ${
                    errors.password ? "border-red-500" : "border-amber-900/40"
                  } focus:border-amber-500`}
                  placeholder="Create a password"
                  value={value}
                  onChangeText={onChange}
                  onBlur={onBlur}
                  secureTextEntry
                />
              )}
            />
            {errors.password && (
              <Text className="text-red-500 text-sm mt-1">
                {errors.password.message}
              </Text>
            )}
          </View>

          <View>
            <Text className="text-amber-900 font-medium mb-2">
              Confirm Password
            </Text>
            <Controller
              control={control}
              name="confirmPassword"
              rules={{
                required: "Please confirm your password",
                validate: (value) =>
                  value === password || "Passwords do not match",
              }}
              render={({ field: { onChange, onBlur, value } }) => (
                <TextInput
                  className={`border-2 rounded-lg px-4 py-3 ${
                    errors.confirmPassword
                      ? "border-red-500"
                      : "border-amber-900/40"
                  } focus:border-amber-500`}
                  placeholder="Confirm your password"
                  value={value}
                  onChangeText={onChange}
                  onBlur={onBlur}
                  secureTextEntry
                />
              )}
            />
            {errors.confirmPassword && (
              <Text className="text-red-500 text-sm mt-1">
                {errors.confirmPassword.message}
              </Text>
            )}
          </View>

          <TouchableOpacity
            onPress={handleSubmit(onRegister)}
            className="bg-amber-800 rounded-lg py-4 mt-6"
          >
            <Text className="text-white text-center font-semibold text-lg">
              Create Account
            </Text>
          </TouchableOpacity>

          <View className="flex-row justify-center mt-6">
            <Text className="text-amber-700">Already have an account? </Text>
            <TouchableOpacity onPress={() => navigation.navigate("Login")}>
              <Text className="text-amber-800 font-semibold">Sign In</Text>
            </TouchableOpacity>
          </View>
        </View>
      </View>
    </View>
  );
}
