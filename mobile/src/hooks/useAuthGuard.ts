import { useEffect } from "react";
import { useAuth } from "../contexts/AuthContext";
import { useAppNavigation } from "./useAppNavigation";
import { ROUTE_ACCESS } from "../types/auth";

interface UseAuthGuardProps {
  routeName: string;
  onUnauthorized?: () => void;
  onAuthorized?: () => void;
}

export const useAuthGuard = ({
  routeName,
  onUnauthorized,
  onAuthorized,
}: UseAuthGuardProps) => {
  const { isAuthenticated, user, isLoading } = useAuth();
  const { navigateToLogin, navigateToHome } = useAppNavigation();

  useEffect(() => {
    if (isLoading) return; // Wait for auth state to load

    const routeConfig = ROUTE_ACCESS[routeName];
    if (!routeConfig) return; // No access control for this route

    const { access, roles, redirectTo } = routeConfig;

    switch (access) {
      case "protected":
        if (!isAuthenticated) {
          onUnauthorized?.();
          navigateToLogin();
          return;
        }

        // Check role-based access
        if (roles && user && !roles.includes(user.role)) {
          onUnauthorized?.();
          navigateToHome(); // Redirect to home instead of login
          return;
        }

        onAuthorized?.();
        break;

      case "rejected":
        if (isAuthenticated) {
          // User is authenticated but trying to access auth screens
          const redirectRoute = redirectTo || "Home";
          if (redirectRoute === "Main" || redirectRoute === "Home") {
            navigateToHome();
          }
          return;
        }
        onAuthorized?.();
        break;

      case "public":
      default:
        onAuthorized?.();
        break;
    }
  }, [
    isAuthenticated,
    user,
    isLoading,
    routeName,
    onUnauthorized,
    onAuthorized,
    navigateToLogin,
    navigateToHome,
  ]);

  return {
    isAuthorized: isAuthenticated,
    user,
    isLoading,
    hasAccess: (requiredRoles?: string[]) => {
      if (!isAuthenticated) return false;
      if (!requiredRoles) return true;
      return user ? requiredRoles.includes(user.role) : false;
    },
  };
};

// Hook for checking if user can access a specific route
export const useRouteAccess = (routeName: string) => {
  const { isAuthenticated, user } = useAuth();
  const routeConfig = ROUTE_ACCESS[routeName];

  if (!routeConfig) return { canAccess: true, reason: "No access control" };

  const { access, roles } = routeConfig;

  switch (access) {
    case "protected":
      if (!isAuthenticated) {
        return { canAccess: false, reason: "Authentication required" };
      }
      if (roles && user && !roles.includes(user.role)) {
        return { canAccess: false, reason: "Insufficient permissions" };
      }
      return { canAccess: true, reason: "Access granted" };

    case "rejected":
      if (isAuthenticated) {
        return { canAccess: false, reason: "Already authenticated" };
      }
      return { canAccess: true, reason: "Access granted" };

    case "public":
    default:
      return { canAccess: true, reason: "Public route" };
  }
};
