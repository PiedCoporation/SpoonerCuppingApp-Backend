import { createNavigationContainerRef } from "@react-navigation/native";
import { RootStackParamList } from "../types/navigation";

export const navigationRef = createNavigationContainerRef<RootStackParamList>();

export function navigate(name: keyof RootStackParamList, params?: any) {
  if (navigationRef.isReady()) {
    navigationRef.navigate({
      name,
      params,
    } as any);
  }
}

export function goBack() {
  if (navigationRef.isReady()) {
    navigationRef.goBack();
  }
}
