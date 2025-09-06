import { NativeStackScreenProps } from "@react-navigation/native-stack";
import { BottomTabScreenProps } from "@react-navigation/bottom-tabs";

export type AuthStackParamList = {
  Login: undefined;
  Register: undefined;
  ForgotPassword: undefined;
};

export type MainTabParamList = {
  Home: undefined;
  Sessions: undefined;
  Explore: undefined;
  Profile: undefined;
  Friend: undefined;
  Market: undefined;
  Notification: undefined;
};

export type AppStackParamList = {
  SessionDetail: { sessionId: string };
  CreateSession: undefined;
  CoffeeDetail: { coffeeId: string };
  Settings: undefined;
  AccountSecurity: undefined;
  CuppingStyle: undefined;
};

export type RootStackParamList = {
  Auth: undefined;
  Main: undefined;
} & AppStackParamList;

export type AuthScreenProps<T extends keyof AuthStackParamList> =
  NativeStackScreenProps<AuthStackParamList, T>;

export type MainTabProps<T extends keyof MainTabParamList> =
  BottomTabScreenProps<MainTabParamList, T>;

export type AppScreenProps<T extends keyof AppStackParamList> =
  NativeStackScreenProps<AppStackParamList, T>;
