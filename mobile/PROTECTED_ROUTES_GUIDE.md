# Protected Routes & Authentication System

## 🛡️ **Overview**

This system provides comprehensive route protection with three types of access control:

1. **Protected Routes** - Require authentication
2. **Rejected Routes** - Block authenticated users (like login/register)
3. **Public Routes** - Accessible to everyone

## 🔧 **Implementation**

### **1. Route Access Configuration**

Routes are configured in `src/types/auth.ts`:

```typescript
export const ROUTE_ACCESS: Record<string, RouteConfig> = {
  // Rejected routes (redirect if already authenticated)
  Login: { access: "rejected", redirectTo: "Main" },
  Register: { access: "rejected", redirectTo: "Main" },

  // Protected routes (require authentication)
  Main: { access: "protected" },
  Users: { access: "protected", roles: ["admin", "moderator"] },
  Profile: { access: "protected" },

  // Public routes (no access control)
  PublicScreen: { access: "public" },
};
```

### **2. Authentication Context**

Enhanced `AuthContext` with:

- ✅ Persistent authentication state
- ✅ AsyncStorage integration
- ✅ Loading states
- ✅ Error handling
- ✅ User role management

```typescript
const { isAuthenticated, user, login, logout, isLoading } = useAuth();
```

### **3. Route Protection Components**

#### **ProtectedRoute Component**

```typescript
<ProtectedRoute routeName="Users">
  <UsersScreen />
</ProtectedRoute>
```

#### **useAuthGuard Hook**

```typescript
useAuthGuard({
  routeName: "Users",
  onUnauthorized: () => console.log("Access denied"),
  onAuthorized: () => console.log("Access granted"),
});
```

### **4. Secure Navigation**

#### **useAppNavigation Hook**

```typescript
const { secureNavigate, navigateToUsers } = useAppNavigation();

// Automatic permission checking
secureNavigate("Users"); // Will redirect to login if not authenticated
navigateToUsers(); // Type-safe navigation with built-in protection
```

## 🚀 **Usage Examples**

### **1. Basic Protected Screen**

```typescript
import { useAuth } from "../contexts/AuthContext";
import ProtectedRoute from "../components/auth/ProtectedRoute";

const MyScreen = () => {
  const { user } = useAuth();

  return (
    <View>
      <Text>Welcome, {user?.name}!</Text>
    </View>
  );
};

// Wrap with protection
export default () => (
  <ProtectedRoute routeName="MyScreen">
    <MyScreen />
  </ProtectedRoute>
);
```

### **2. Role-Based Access**

```typescript
const AdminScreen = () => {
  const { user } = useAuth();

  // Check if user has admin role
  if (user?.role !== 'admin') {
    return <Text>Access Denied</Text>;
  }

  return <Text>Admin Panel</Text>;
};

// Configure in ROUTE_ACCESS
'AdminScreen': { access: 'protected', roles: ['admin'] }
```

### **3. Navigation with Protection**

```typescript
const MyComponent = () => {
  const { secureNavigate } = useAppNavigation();

  const handleNavigate = () => {
    // This will automatically check permissions
    secureNavigate("Users"); // Redirects to login if not authenticated
  };

  return (
    <TouchableOpacity onPress={handleNavigate}>
      <Text>Go to Users</Text>
    </TouchableOpacity>
  );
};
```

## 🔄 **Authentication Flow**

### **Login Flow**

1. User enters credentials
2. `login()` function validates credentials
3. On success: Store token & user data in AsyncStorage
4. Update auth state: `isAuthenticated = true`
5. Navigate to main app

### **Logout Flow**

1. Call `logout()` function
2. Clear AsyncStorage
3. Update auth state: `isAuthenticated = false`
4. Navigate to login screen

### **Route Protection Flow**

1. User navigates to protected route
2. `ProtectedRoute` component checks access
3. If not authenticated: Redirect to login
4. If wrong role: Show access denied or redirect
5. If authorized: Render component

## 🛠️ **Configuration**

### **Adding New Protected Routes**

1. **Add to ROUTE_ACCESS:**

```typescript
'NewScreen': { access: 'protected', roles: ['admin'] }
```

2. **Wrap component:**

```typescript
<ProtectedRoute routeName="NewScreen">
  <NewScreen />
</ProtectedRoute>
```

3. **Add navigation function:**

```typescript
navigateToNewScreen: () => secureNavigate("NewScreen");
```

### **Custom Access Control**

```typescript
const CustomScreen = () => {
  const { user } = useAuth();
  const { canAccess } = useRouteAccess("CustomScreen");

  if (!canAccess) {
    return <Text>Access Denied</Text>;
  }

  return <Text>Custom Screen</Text>;
};
```

## 🔒 **Security Features**

### **1. Automatic Redirects**

- Unauthenticated users → Login screen
- Authenticated users → Home screen (when accessing auth screens)
- Wrong role → Access denied or redirect

### **2. Persistent Sessions**

- Authentication state survives app restarts
- Automatic token validation
- Secure storage with AsyncStorage

### **3. Role-Based Access**

- Admin-only screens
- Moderator permissions
- User-level access

### **4. Loading States**

- Prevents flash of wrong content
- Smooth authentication flow
- Proper error handling

## 📱 **Production Considerations**

### **1. Token Management**

- Implement JWT token refresh
- Add token expiration handling
- Secure token storage

### **2. Error Handling**

- Network error recovery
- Invalid token handling
- Graceful fallbacks

### **3. Analytics**

- Track authentication events
- Monitor route access patterns
- User behavior analytics

### **4. Testing**

- Unit tests for auth logic
- Integration tests for navigation
- E2E tests for user flows

## 🎯 **Best Practices**

1. **Always wrap protected screens** with `ProtectedRoute`
2. **Use `secureNavigate`** for programmatic navigation
3. **Define route access** in `ROUTE_ACCESS` configuration
4. **Handle loading states** properly
5. **Provide clear error messages** for access denied
6. **Test authentication flows** thoroughly
7. **Implement proper logout** functionality
8. **Use role-based access** for sensitive features

This system provides a robust foundation for secure navigation and authentication in your Coffee Cupping app!
