import React from 'react';
import {
  View,
  Text,
  FlatList,
  ActivityIndicator,
  TouchableOpacity,
  Alert,
} from 'react-native';
import { useUsers, useCreateUser } from '../services/queries/userQueries';
import { User } from '../types';

const UsersScreen: React.FC = () => {
  const { data: users, isLoading, error, refetch } = useUsers();
  const createUserMutation = useCreateUser();

  const handleCreateUser = () => {
    const newUser = {
      name: 'New User',
      username: 'newuser',
      email: 'newuser@example.com',
      phone: '123-456-7890',
      website: 'newuser.com',
    };

    createUserMutation.mutate(newUser, {
      onSuccess: () => {
        Alert.alert('Success', 'User created successfully!');
      },
      onError: (error) => {
        Alert.alert('Error', `Failed to create user: ${error.message}`);
      },
    });
  };

  const renderUser = ({ item }: { item: User }) => (
    <View className="bg-white mx-4 mb-3 p-4 rounded-lg shadow-sm border border-gray-100">
      <View className="flex-row justify-between items-start mb-2">
        <View className="flex-1">
          <Text className="text-lg font-bold text-gray-900 mb-1">
            {item.name}
          </Text>
          <Text className="text-sm text-gray-500">@{item.username}</Text>
        </View>
        <View className="bg-primary-100 px-2 py-1 rounded">
          <Text className="text-xs text-primary-700 font-medium">
            ID: {item.id}
          </Text>
        </View>
      </View>
      
      <View className="border-t border-gray-100 pt-3 mt-3">
        <View className="flex-row items-center mb-2">
          <View className="w-2 h-2 bg-green-500 rounded-full mr-2" />
          <Text className="text-sm text-gray-700 flex-1">{item.email}</Text>
        </View>
        <View className="flex-row items-center mb-2">
          <View className="w-2 h-2 bg-blue-500 rounded-full mr-2" />
          <Text className="text-sm text-gray-700 flex-1">{item.phone}</Text>
        </View>
        <View className="flex-row items-center">
          <View className="w-2 h-2 bg-purple-500 rounded-full mr-2" />
          <Text className="text-sm text-gray-700 flex-1">{item.website}</Text>
        </View>
      </View>
    </View>
  );

  if (isLoading) {
    return (
      <View className="flex-1 justify-center items-center bg-gray-50">
        <View className="bg-white p-6 rounded-lg shadow-sm">
          <ActivityIndicator size="large" color="#3b82f6" />
          <Text className="mt-4 text-gray-600 text-base text-center">
            Loading users...
          </Text>
        </View>
      </View>
    );
  }

  if (error) {
    return (
      <View className="flex-1 justify-center items-center bg-gray-50 px-6">
        <View className="bg-white p-6 rounded-lg shadow-sm border border-red-200">
          <View className="items-center mb-4">
            <View className="w-12 h-12 bg-red-100 rounded-full items-center justify-center mb-3">
              <Text className="text-red-600 text-xl">!</Text>
            </View>
            <Text className="text-red-800 text-center text-base font-medium">
              Oops! Something went wrong
            </Text>
            <Text className="text-red-600 text-center text-sm mt-1">
              {error.message}
            </Text>
          </View>
          <TouchableOpacity
            className="bg-red-600 py-3 px-6 rounded-lg active:bg-red-700"
            onPress={() => refetch()}
          >
            <Text className="text-white font-semibold text-center">
              Try Again
            </Text>
          </TouchableOpacity>
        </View>
      </View>
    );
  }

  return (
    <View className="flex-1 bg-gray-50">
      {/* Header */}
      <View className="bg-white shadow-sm">
        <View className="px-4 pt-12 pb-4">
          <Text className="text-2xl font-bold text-gray-900 mb-1">
            Users Directory
          </Text>
          <Text className="text-gray-600 text-sm mb-4">
            Manage and view all users
          </Text>
          
          <TouchableOpacity
            className={`py-3 px-6 rounded-lg flex-row items-center justify-center ${
              createUserMutation.isPending 
                ? 'bg-gray-400' 
                : 'bg-primary-600 active:bg-primary-700'
            }`}
            onPress={handleCreateUser}
            disabled={createUserMutation.isPending}
          >
            {createUserMutation.isPending && (
              <ActivityIndicator size="small" color="white" className="mr-2" />
            )}
            <Text className="text-white font-semibold">
              {createUserMutation.isPending ? 'Creating User...' : '+ Add New User'}
            </Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* Users List */}
      <FlatList
        data={users}
        keyExtractor={(item) => item.id.toString()}
        renderItem={renderUser}
        onRefresh={refetch}
        refreshing={isLoading}
        className="flex-1 pt-4"
        showsVerticalScrollIndicator={false}
        ListEmptyComponent={
          <View className="flex-1 justify-center items-center py-20">
            <View className="items-center">
              <View className="w-16 h-16 bg-gray-200 rounded-full items-center justify-center mb-4">
                <Text className="text-gray-400 text-2xl">👥</Text>
              </View>
              <Text className="text-gray-500 text-base font-medium">
                No users found
              </Text>
              <Text className="text-gray-400 text-sm text-center mt-1">
                Add your first user to get started
              </Text>
            </View>
          </View>
        }
      />
    </View>
  );
};

export default UsersScreen;