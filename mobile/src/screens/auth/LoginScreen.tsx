import React from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  Alert,
  ActivityIndicator,
} from "react-native";
import { useForm, Controller } from "react-hook-form";
import { AuthScreenProps } from "../../types/navigation";
import { useAuth } from "../../contexts/AuthContext";

type LoginForm = {
  email: string;
  password: string;
};

export default function LoginScreen({ navigation }: AuthScreenProps<"Login">) {
  const { login, isLoading } = useAuth();
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    defaultValues: {
      email: "admin@gmail.com",
      password: "12345",
    },
  });

  const onLogin = async (data: LoginForm) => {
    const success = await login(data);
  };

  return (
    <View className="flex-1 bg-white px-6 justify-center">
      <View className="mb-8">
        <Text className="text-4xl font-extrabold text-amber-900 text-center tracking-widest mb-1">
          SPOONER
        </Text>
        <Text className="text-xl text-amber-700 text-center tracking-wide">
          Coffee Social Network
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
                } focus:border-amber-600`}
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
                value: 5,
                message: "Password must be at least 5 characters",
              },
            }}
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className={`border-2 rounded-lg px-4 py-3 ${
                  errors.password ? "border-red-500" : "border-amber-900/40"
                } focus:border-amber-600`}
                placeholder="Enter your password"
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

        <TouchableOpacity
          onPress={() => navigation.navigate("ForgotPassword")}
          className="self-end"
        >
          <Text className="text-amber-900 font-medium">Forgot Password?</Text>
        </TouchableOpacity>

        <TouchableOpacity
          onPress={handleSubmit(onLogin)}
          disabled={isLoading}
          className={`rounded-lg py-4 mt-6 ${
            isLoading ? "bg-amber-600" : "bg-amber-800"
          }`}
        >
          {isLoading ? (
            <ActivityIndicator color="white" />
          ) : (
            <Text className="text-white text-center font-semibold text-lg">
              Sign In
            </Text>
          )}
        </TouchableOpacity>

        <View className="flex-row justify-center mt-6">
          <Text className="text-amber-700">Don't have an account? </Text>
          <TouchableOpacity onPress={() => navigation.navigate("Register")}>
            <Text className="text-amber-800 font-semibold">Sign Up</Text>
          </TouchableOpacity>
        </View>
      </View>
    </View>
  );
}
