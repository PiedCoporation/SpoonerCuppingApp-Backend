import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import UsersScreen from '../screens/UsersScreen';

const Stack = createNativeStackNavigator();

const MainNavigator: React.FC = () => {
  return (
    <NavigationContainer>
      <Stack.Navigator initialRouteName="Users">
        <Stack.Screen 
          name="Users" 
          component={UsersScreen}
          options={{ title: 'Users List' }}
        />
      </Stack.Navigator>
    </NavigationContainer>
  );
};


export default MainNavigator;