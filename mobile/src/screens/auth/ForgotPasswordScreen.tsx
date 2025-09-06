import React, { useState } from "react";
import { View, Text, TextInput, TouchableOpacity, Alert } from "react-native";
import { useForm, Controller } from "react-hook-form";
import { AuthScreenProps } from "../../types/navigation";

type ForgotPasswordForm = {
  email: string;
};

export default function ForgotPasswordScreen({
  navigation,
}: AuthScreenProps<"ForgotPassword">) {
  const [isSubmitted, setIsSubmitted] = useState(false);
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<ForgotPasswordForm>();

  const onSubmit = (data: ForgotPasswordForm) => {
    // TODO: Implement forgot password logic
    console.log("Forgot Password:", data);
    setIsSubmitted(true);
  };

  if (isSubmitted) {
    return (
      <View className="flex-1 bg-white px-6 justify-center">
        <View className="items-center">
          <Text className="text-2xl font-bold text-amber-900 text-center mb-4">
            Check Your Email
          </Text>
          <Text className="text-amber-700 text-center mb-8 leading-6">
            We've sent password reset instructions to your email address.
          </Text>
          <TouchableOpacity
            onPress={() => navigation.navigate("Login")}
            className="bg-amber-800 rounded-lg py-4 px-8"
          >
            <Text className="text-white font-semibold text-lg">
              Back to Login
            </Text>
          </TouchableOpacity>
        </View>
      </View>
    );
  }

  return (
    <View className="flex-1 bg-white px-6 justify-center">
      <View className="mb-8">
        <Text className="text-3xl font-bold text-amber-900 text-center mb-2">
          Forgot Password
        </Text>
        <Text className="text-lg text-amber-700 text-center leading-6">
          Enter your email address and we'll send you instructions to reset your
          password.
        </Text>
      </View>

      <View className="space-y-4">
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

        <TouchableOpacity
          onPress={handleSubmit(onSubmit)}
          className="bg-amber-800 rounded-lg py-4 mt-6"
        >
          <Text className="text-white text-center font-semibold text-lg">
            Send Instructions
          </Text>
        </TouchableOpacity>

        <TouchableOpacity
          onPress={() => navigation.navigate("Login")}
          className="mt-6"
        >
          <Text className="text-amber-800 text-center font-semibold">
            Back to Login
          </Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}
