import React from "react";
import { View, Text, ActivityIndicator, StyleSheet } from "react-native";
import { useAuth } from "../../contexts/AuthContext";
import { useRouteAccess } from "../../hooks/useAuthGuard";
import { ROUTE_ACCESS } from "../../types/auth";

interface ProtectedRouteProps {
  children: React.ReactNode;
  routeName: string;
  fallback?: React.ReactNode;
  loadingComponent?: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  children,
  routeName,
  fallback,
  loadingComponent,
}) => {
  const { isLoading } = useAuth();
  const { canAccess, reason } = useRouteAccess(routeName);

  // Show loading state
  if (isLoading) {
    return (
      loadingComponent || (
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color="#8B4513" />
          <Text style={styles.loadingText}>Loading...</Text>
        </View>
      )
    );
  }

  // Check access
  if (!canAccess) {
    return (
      fallback || (
        <View style={styles.errorContainer}>
          <Text style={styles.errorTitle}>Access Denied</Text>
          <Text style={styles.errorMessage}>{reason}</Text>
        </View>
      )
    );
  }

  return <>{children}</>;
};

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: "#FFFEF7",
  },
  loadingText: {
    marginTop: 16,
    fontSize: 16,
    color: "#8B7355",
  },
  errorContainer: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: "#FFFEF7",
    padding: 20,
  },
  errorTitle: {
    fontSize: 24,
    fontWeight: "bold",
    color: "#8B4513",
    marginBottom: 12,
  },
  errorMessage: {
    fontSize: 16,
    color: "#8B7355",
    textAlign: "center",
    lineHeight: 24,
  },
});

export default ProtectedRoute;
