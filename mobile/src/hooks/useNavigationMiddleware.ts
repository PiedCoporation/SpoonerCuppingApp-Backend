import { useEffect } from "react";
import { useNavigation, useRoute } from "@react-navigation/native";
import { useAuth } from "../contexts/AuthContext";
import { useAppNavigation } from "./useAppNavigation";
import { ROUTE_ACCESS } from "../types/auth";

export const useNavigationMiddleware = () => {
  const navigation = useNavigation();
  const route = useRoute();
  const { isAuthenticated, user, isLoading } = useAuth();
  const { navigateToLogin, navigateToHome } = useAppNavigation();

  useEffect(() => {
    if (isLoading) return; // Wait for auth state to load

    const routeName = route.name;
    const routeConfig = ROUTE_ACCESS[routeName];

    if (!routeConfig) return; // No access control for this route

    const { access, roles, redirectTo } = routeConfig;

    switch (access) {
      case "protected":
        if (!isAuthenticated) {
          console.log(
            `Redirecting from ${routeName} to Login - Authentication required`
          );
          navigateToLogin();
          return;
        }

        // Check role-based access
        if (roles && user && !roles.includes(user.role)) {
          console.log(
            `Redirecting from ${routeName} to Home - Insufficient permissions`
          );
          navigateToHome();
          return;
        }
        break;

      case "rejected":
        if (isAuthenticated) {
          console.log(
            `Redirecting from ${routeName} to ${
              redirectTo || "Home"
            } - Already authenticated`
          );
          const redirectRoute = redirectTo || "Home";
          if (redirectRoute === "Main" || redirectRoute === "Home") {
            navigateToHome();
          }
          return;
        }
        break;

      case "public":
      default:
        // No redirect needed for public routes
        break;
    }
  }, [
    isAuthenticated,
    user,
    isLoading,
    route.name,
    navigateToLogin,
    navigateToHome,
  ]);

  return {
    isAuthorized: isAuthenticated,
    user,
    isLoading,
  };
};

// Hook for programmatic navigation with access control
export const useSecureNavigation = () => {
  const { isAuthenticated, user } = useAuth();
  const { navigateToLogin } = useAppNavigation();

  const navigateWithAuth = (routeName: string, params?: any) => {
    const routeConfig = ROUTE_ACCESS[routeName];

    if (!routeConfig) {
      // No access control, allow navigation
      return true;
    }

    const { access, roles } = routeConfig;

    switch (access) {
      case "protected":
        if (!isAuthenticated) {
          console.log(
            `Blocking navigation to ${routeName} - Authentication required`
          );
          navigateToLogin();
          return false;
        }

        if (roles && user && !roles.includes(user.role)) {
          console.log(
            `Blocking navigation to ${routeName} - Insufficient permissions`
          );
          return false;
        }
        return true;

      case "rejected":
        if (isAuthenticated) {
          console.log(
            `Blocking navigation to ${routeName} - Already authenticated`
          );
          return false;
        }
        return true;

      case "public":
      default:
        return true;
    }
  };

  return { navigateWithAuth };
};
