import AsyncStorage from "@react-native-async-storage/async-storage";

// Keys for AsyncStorage
const ACCESS_TOKEN_KEY = "access_token";
const PROFILE_KEY = "user_profile";

/**
 * Get access token from AsyncStorage
 * @returns Promise<string | null> - The access token or null if not found
 */
export const getAccessTokenToAS = async (): Promise<string | null> => {
  try {
    const token = await AsyncStorage.getItem(ACCESS_TOKEN_KEY);
    return token;
  } catch (error) {
    console.error("Error getting access token from AsyncStorage:", error);
    return null;
  }
};

/**
 * Save access token to AsyncStorage
 * @param token - The access token to save
 * @returns Promise<boolean> - True if successful, false otherwise
 */
export const saveAccessTokenToAS = async (token: string): Promise<boolean> => {
  try {
    await AsyncStorage.setItem(ACCESS_TOKEN_KEY, token);
    return true;
  } catch (error) {
    console.error("Error saving access token to AsyncStorage:", error);
    return false;
  }
};

/**
 * Set user profile to AsyncStorage
 * @param profile - The user profile object to save
 * @returns Promise<boolean> - True if successful, false otherwise
 */
export const setProfileToAS = async (profile: any): Promise<boolean> => {
  try {
    const profileString = JSON.stringify(profile);
    await AsyncStorage.setItem(PROFILE_KEY, profileString);
    return true;
  } catch (error) {
    console.error("Error saving profile to AsyncStorage:", error);
    return false;
  }
};

/**
 * Get user profile from AsyncStorage
 * @returns Promise<any | null> - The user profile or null if not found
 */
export const getProfileFromAS = async (): Promise<any | null> => {
  try {
    const profileString = await AsyncStorage.getItem(PROFILE_KEY);
    if (profileString) {
      return JSON.parse(profileString);
    }
    return null;
  } catch (error) {
    console.error("Error getting profile from AsyncStorage:", error);
    return null;
  }
};

/**
 * Clear all auth-related data from AsyncStorage
 * @returns Promise<boolean> - True if successful, false otherwise
 */
export const clearAuthDataFromAS = async (): Promise<boolean> => {
  try {
    await AsyncStorage.multiRemove([ACCESS_TOKEN_KEY, PROFILE_KEY]);
    return true;
  } catch (error) {
    console.error("Error clearing auth data from AsyncStorage:", error);
    return false;
  }
};

/**
 * Remove access token from AsyncStorage
 * @returns Promise<boolean> - True if successful, false otherwise
 */
export const removeAccessTokenFromAS = async (): Promise<boolean> => {
  try {
    await AsyncStorage.removeItem(ACCESS_TOKEN_KEY);
    return true;
  } catch (error) {
    console.error("Error removing access token from AsyncStorage:", error);
    return false;
  }
};

/**
 * Remove user profile from AsyncStorage
 * @returns Promise<boolean> - True if successful, false otherwise
 */
export const removeProfileFromAS = async (): Promise<boolean> => {
  try {
    await AsyncStorage.removeItem(PROFILE_KEY);
    return true;
  } catch (error) {
    console.error("Error removing profile from AsyncStorage:", error);
    return false;
  }
};
