import React from 'react';
import { View, Text, TouchableOpacity, ScrollView } from 'react-native';

const TestNativeWind: React.FC = () => {
  return (
    <ScrollView className="flex-1 bg-gray-50">
      {/* Header */}
      <View className="bg-primary-600 p-6 pt-12">
        <Text className="text-white text-2xl font-bold text-center">
          NativeWind Test
        </Text>
        <Text className="text-primary-100 text-center mt-2">
          Testing Tailwind CSS classes
        </Text>
      </View>

      {/* Content */}
      <View className="p-4 space-y-4">
        {/* Colors Test */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <Text className="text-lg font-semibold mb-3 text-gray-900">
            Colors Test
          </Text>
          <View className="flex-row space-x-2 mb-2">
            <View className="w-8 h-8 bg-primary-500 rounded" />
            <View className="w-8 h-8 bg-secondary-500 rounded" />
            <View className="w-8 h-8 bg-success-500 rounded" />
            <View className="w-8 h-8 bg-warning-500 rounded" />
            <View className="w-8 h-8 bg-error-500 rounded" />
          </View>
        </View>

        {/* Typography Test */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <Text className="text-lg font-semibold mb-3 text-gray-900">
            Typography Test
          </Text>
          <Text className="text-xs text-gray-500 mb-1">Extra Small Text</Text>
          <Text className="text-sm text-gray-600 mb-1">Small Text</Text>
          <Text className="text-base text-gray-700 mb-1">Base Text</Text>
          <Text className="text-lg text-gray-800 mb-1">Large Text</Text>
          <Text className="text-xl font-bold text-gray-900">Extra Large Bold</Text>
        </View>

        {/* Buttons Test */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <Text className="text-lg font-semibold mb-3 text-gray-900">
            Buttons Test
          </Text>
          <View className="space-y-3">
            <TouchableOpacity className="bg-primary-600 py-3 px-6 rounded-lg active:bg-primary-700">
              <Text className="text-white font-semibold text-center">
                Primary Button
              </Text>
            </TouchableOpacity>
            
            <TouchableOpacity className="bg-secondary-500 py-3 px-6 rounded-lg active:bg-secondary-600">
              <Text className="text-white font-semibold text-center">
                Secondary Button
              </Text>
            </TouchableOpacity>
            
            <TouchableOpacity className="border-2 border-primary-600 py-3 px-6 rounded-lg active:bg-primary-50">
              <Text className="text-primary-600 font-semibold text-center">
                Outline Button
              </Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* Layout Test */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <Text className="text-lg font-semibold mb-3 text-gray-900">
            Layout Test
          </Text>
          <View className="flex-row justify-between items-center mb-3">
            <View className="bg-blue-100 p-3 rounded flex-1 mr-2">
              <Text className="text-blue-800 text-center">Flex 1</Text>
            </View>
            <View className="bg-green-100 p-3 rounded flex-1 ml-2">
              <Text className="text-green-800 text-center">Flex 1</Text>
            </View>
          </View>
          
          <View className="flex-row justify-around">
            <View className="bg-red-100 p-2 rounded">
              <Text className="text-red-800">Around</Text>
            </View>
            <View className="bg-yellow-100 p-2 rounded">
              <Text className="text-yellow-800">Space</Text>
            </View>
            <View className="bg-purple-100 p-2 rounded">
              <Text className="text-purple-800">Items</Text>
            </View>
          </View>
        </View>

        {/* Spacing Test */}
        <View className="bg-white rounded-lg shadow-sm">
          <View className="p-4 border-b border-gray-200">
            <Text className="text-lg font-semibold text-gray-900">
              Spacing Test
            </Text>
          </View>
          <View className="p-2">
            <View className="bg-blue-50 p-1 m-1">
              <Text className="text-blue-800 text-sm">p-1 m-1</Text>
            </View>
            <View className="bg-green-50 p-2 m-2">
              <Text className="text-green-800 text-sm">p-2 m-2</Text>
            </View>
            <View className="bg-red-50 p-4 m-4">
              <Text className="text-red-800 text-sm">p-4 m-4</Text>
            </View>
          </View>
        </View>

        {/* Border & Shadow Test */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <Text className="text-lg font-semibold mb-3 text-gray-900">
            Border & Shadow Test
          </Text>
          <View className="space-y-3">
            <View className="border border-gray-300 p-3 rounded">
              <Text className="text-gray-700">Border Gray</Text>
            </View>
            <View className="border-2 border-primary-500 p-3 rounded-lg">
              <Text className="text-primary-700">Border Primary</Text>
            </View>
            <View className="bg-white p-3 rounded-lg shadow-lg border border-gray-100">
              <Text className="text-gray-700">Shadow Large</Text>
            </View>
          </View>
        </View>
      </View>
    </ScrollView>
  );
};

export default TestNativeWind;