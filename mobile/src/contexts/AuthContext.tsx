import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import AsyncStorage from "@react-native-async-storage/async-storage";
import {
  AuthContextType,
  AuthState,
  LoginCredentials,
  RegisterData,
  User,
} from "../types/auth";

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Mock user data for development
const MOCK_USER: User = {
  id: "1",
  email: "admin@gmail.com",
  name: "John Doe",
  role: "admin",
  isEmailVerified: true,
  createdAt: new Date().toISOString(),
  lastLoginAt: new Date().toISOString(),
};

export function AuthProvider({ children }: { children: ReactNode }) {
  const [authState, setAuthState] = useState<AuthState>({
    user: null,
    isAuthenticated: false,
    isLoading: true,
    error: null,
  });

  // Check for existing session on app start
  useEffect(() => {
    checkAuthState();
  }, []);

  const checkAuthState = async () => {
    try {
      setAuthState((prev) => ({ ...prev, isLoading: true }));

      // Check for stored auth token
      const token = await AsyncStorage.getItem("auth_token");
      const userData = await AsyncStorage.getItem("user_data");

      if (token && userData) {
        const user = JSON.parse(userData);
        setAuthState({
          user,
          isAuthenticated: true,
          isLoading: false,
          error: null,
        });
      } else {
        setAuthState({
          user: null,
          isAuthenticated: false,
          isLoading: false,
          error: null,
        });
      }
    } catch (error) {
      console.error("Error checking auth state:", error);
      setAuthState({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: "Failed to check authentication state",
      });
    }
  };

  const login = async (credentials: LoginCredentials) => {
    try {
      setAuthState((prev) => ({ ...prev, isLoading: true, error: null }));

      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // Mock authentication logic
      if (
        credentials.email === "admin@gmail.com" &&
        credentials.password === "12345"
      ) {
        const user = { ...MOCK_USER, lastLoginAt: new Date().toISOString() };
        const token = "mock_jwt_token_" + Date.now();

        // Store auth data
        await AsyncStorage.setItem("auth_token", token);
        await AsyncStorage.setItem("user_data", JSON.stringify(user));

        setAuthState({
          user,
          isAuthenticated: true,
          isLoading: false,
          error: null,
        });
      } else {
        throw new Error("Invalid email or password");
      }
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Login failed";
      setAuthState((prev) => ({
        ...prev,
        isLoading: false,
        error: errorMessage,
      }));
      throw error;
    }
  };

  const register = async (data: RegisterData) => {
    try {
      setAuthState((prev) => ({ ...prev, isLoading: true, error: null }));

      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // Mock registration logic
      if (data.password !== data.confirmPassword) {
        throw new Error("Passwords do not match");
      }

      const newUser: User = {
        id: Date.now().toString(),
        email: data.email,
        name: data.name,
        role: "user",
        isEmailVerified: false,
        createdAt: new Date().toISOString(),
      };

      const token = "mock_jwt_token_" + Date.now();

      // Store auth data
      await AsyncStorage.setItem("auth_token", token);
      await AsyncStorage.setItem("user_data", JSON.stringify(newUser));

      setAuthState({
        user: newUser,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      });
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Registration failed";
      setAuthState((prev) => ({
        ...prev,
        isLoading: false,
        error: errorMessage,
      }));
      throw error;
    }
  };

  const logout = async () => {
    try {
      setAuthState((prev) => ({ ...prev, isLoading: true }));

      // Clear stored auth data
      await AsyncStorage.removeItem("auth_token");
      await AsyncStorage.removeItem("user_data");

      setAuthState({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: null,
      });
    } catch (error) {
      console.error("Error during logout:", error);
      setAuthState((prev) => ({
        ...prev,
        isLoading: false,
        error: "Logout failed",
      }));
    }
  };

  const refreshToken = async () => {
    try {
      const token = await AsyncStorage.getItem("auth_token");
      if (!token) {
        throw new Error("No token found");
      }

      // Simulate token refresh
      await new Promise((resolve) => setTimeout(resolve, 500));

      // In a real app, you would validate the token with your backend
      console.log("Token refreshed successfully");
    } catch (error) {
      console.error("Token refresh failed:", error);
      await logout();
    }
  };

  const clearError = () => {
    setAuthState((prev) => ({ ...prev, error: null }));
  };

  const updateUser = (userUpdates: Partial<User>) => {
    if (authState.user) {
      const updatedUser = { ...authState.user, ...userUpdates };
      setAuthState((prev) => ({ ...prev, user: updatedUser }));

      // Update stored user data
      AsyncStorage.setItem("user_data", JSON.stringify(updatedUser));
    }
  };

  const value: AuthContextType = {
    ...authState,
    login,
    register,
    logout,
    refreshToken,
    clearError,
    updateUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
