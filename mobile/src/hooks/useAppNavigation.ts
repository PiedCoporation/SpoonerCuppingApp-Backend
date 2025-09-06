import { useNavigation } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";
import { BottomTabNavigationProp } from "@react-navigation/bottom-tabs";
import {
  RootStackParamList,
  MainTabParamList,
  AuthStackParamList,
} from "../types/navigation";
import { AUTH_ROUTES, MAIN_ROUTES, APP_ROUTES } from "../navigation/routes";
import { useAuth } from "../contexts/AuthContext";
import { ROUTE_ACCESS } from "../types/auth";

// Type-safe navigation hooks
export const useRootNavigation = () => {
  return useNavigation<NativeStackNavigationProp<RootStackParamList>>();
};

export const useMainTabNavigation = () => {
  return useNavigation<BottomTabNavigationProp<MainTabParamList>>();
};

export const useAuthNavigation = () => {
  return useNavigation<NativeStackNavigationProp<AuthStackParamList>>();
};

// Navigation helper functions
export const useAppNavigation = () => {
  const rootNavigation = useRootNavigation();
  const mainTabNavigation = useMainTabNavigation();
  const authNavigation = useAuthNavigation();
  const { isAuthenticated, user } = useAuth();

  // Secure navigation helper
  const secureNavigate = (routeName: string, params?: any) => {
    const routeConfig = ROUTE_ACCESS[routeName];

    if (!routeConfig) {
      // No access control, allow navigation
      return rootNavigation.navigate(routeName as any, params);
    }

    const { access, roles } = routeConfig;

    switch (access) {
      case "protected":
        if (!isAuthenticated) {
          console.log(
            `Redirecting to Login - Authentication required for ${routeName}`
          );
          return authNavigation.navigate(AUTH_ROUTES.LOGIN);
        }

        if (roles && user && !roles.includes(user.role)) {
          console.log(
            `Access denied - Insufficient permissions for ${routeName}`
          );
          return mainTabNavigation.navigate(MAIN_ROUTES.HOME);
        }
        return rootNavigation.navigate(routeName as any, params);

      case "rejected":
        if (isAuthenticated) {
          console.log(
            `Redirecting to Home - Already authenticated, cannot access ${routeName}`
          );
          return mainTabNavigation.navigate(MAIN_ROUTES.HOME);
        }
        return rootNavigation.navigate(routeName as any, params);

      case "public":
      default:
        return rootNavigation.navigate(routeName as any, params);
    }
  };

  return {
    // Auth navigation
    navigateToLogin: () => authNavigation.navigate(AUTH_ROUTES.LOGIN),
    navigateToRegister: () => authNavigation.navigate(AUTH_ROUTES.REGISTER),
    navigateToForgotPassword: () =>
      authNavigation.navigate(AUTH_ROUTES.FORGOT_PASSWORD),

    // Main tab navigation
    navigateToHome: () => mainTabNavigation.navigate(MAIN_ROUTES.HOME),
    navigateToSessions: () => mainTabNavigation.navigate(MAIN_ROUTES.SESSIONS),
    navigateToExplore: () => mainTabNavigation.navigate(MAIN_ROUTES.EXPLORE),
    navigateToProfile: () => mainTabNavigation.navigate(MAIN_ROUTES.PROFILE),

    // App navigation
    navigateToSessionDetail: (sessionId: string) =>
      rootNavigation.navigate(APP_ROUTES.SESSION_DETAIL, { sessionId }),
    navigateToCreateSession: () =>
      rootNavigation.navigate(APP_ROUTES.CREATE_SESSION),
    navigateToSettings: () => rootNavigation.navigate(APP_ROUTES.SETTINGS),
    navigateToAccountSecurity: () =>
      rootNavigation.navigate(APP_ROUTES.ACCOUNT_SECURITY),
    navigateToCuppingStyle: () =>
      rootNavigation.navigate(APP_ROUTES.CUPPING_STYLE),

    // Cupping screens
    navigateToCuppingEventDetail: (eventData: any) =>
      rootNavigation.navigate("CuppingEventDetail", eventData),
    navigateToCuppingRegistrationMinimalist: () =>
      rootNavigation.navigate("CuppingRegistrationMinimalist"),
    navigateToCuppingRegistrationOverview: () =>
      rootNavigation.navigate("CuppingRegistrationOverview"),

    // Utility functions
    goBack: () => rootNavigation.goBack(),
    canGoBack: () => rootNavigation.canGoBack(),

    // Secure navigation
    secureNavigate,
  };
};
