export interface User {
  id: string;
  email: string;
  name: string;
  role: "user" | "admin" | "moderator";
  isEmailVerified: boolean;
  createdAt: string;
  lastLoginAt?: string;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
}

export interface AuthContextType extends AuthState {
  login: (credentials: LoginCredentials) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => Promise<void>;
  refreshToken: () => Promise<void>;
  clearError: () => void;
  updateUser: (user: Partial<User>) => void;
}

// Route protection types
export type RouteAccess = "public" | "protected" | "rejected";

export interface RouteConfig {
  access: RouteAccess;
  roles?: string[];
  redirectTo?: string;
}

// Route definitions with access control
export const ROUTE_ACCESS: Record<string, RouteConfig> = {
  // Public routes (accessible to everyone)
  Login: { access: "rejected", redirectTo: "Main" }, // Redirect if already authenticated
  Register: { access: "rejected", redirectTo: "Main" },
  ForgotPassword: { access: "rejected", redirectTo: "Main" },

  // Protected routes (require authentication)
  Main: { access: "protected" },
  Home: { access: "protected" },
  Sessions: { access: "protected" },
  Explore: { access: "protected" },
  Profile: { access: "protected" },
  Friend: { access: "protected" },
  Market: { access: "protected" },
  Notification: { access: "protected" },
  CuppingEventDetail: { access: "protected" },
  CuppingRegistrationMinimalist: { access: "protected" },
  CuppingRegistrationOverview: { access: "protected" },
  Users: { access: "protected", roles: ["admin", "moderator"] },
  SessionDetail: { access: "protected" },
  CreateSession: { access: "protected" },
  Settings: { access: "protected" },
  AccountSecurity: { access: "protected" },
  CuppingStyle: { access: "protected" },
} as const;
