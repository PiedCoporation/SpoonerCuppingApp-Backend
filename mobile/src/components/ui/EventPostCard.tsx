import React, { useState } from "react";
import {
  View,
  Text,
  TouchableOpacity,
  Image,
  ImageErrorEventData,
  NativeSyntheticEvent,
} from "react-native";
import { Event } from "@/types/event";

// Import Expo Vector Icons with proper typing
import { MaterialIcons } from "@expo/vector-icons";

// Define valid icon names for type safety using Expo Vector Icons
type IconName =
  | "people"
  | "coffee"
  | "heart"
  | "heart-outline"
  | "chatbubble-outline"
  | "share-outline";

// Create a type-safe Icon component
interface SafeIconProps {
  name: IconName;
  size?: number;
  color?: string;
}

const Icon = (props: SafeIconProps) => {
  const { name, size = 16, color = "#6b7280" } = props;

  // Map icon names to appropriate icon sets with verified names
  switch (name) {
    case "people":
      return <MaterialIcons name="people" size={size} color={color} />;
    case "coffee":
      return <MaterialIcons name="local-cafe" size={size} color={color} />;
    case "heart":
      return <MaterialIcons name="favorite" size={size} color={color} />;
    case "heart-outline":
      return <MaterialIcons name="favorite-border" size={size} color={color} />;
    case "chatbubble-outline":
      return (
        <MaterialIcons name="chat-bubble-outline" size={size} color={color} />
      );
    case "share-outline":
      return <MaterialIcons name="share" size={size} color={color} />;
    default:
      return <MaterialIcons name="help-outline" size={size} color={color} />;
  }
};

type Props = {
  event: Event;
  onPress?: (event: Event) => void;
  onLikePress?: (event: Event) => void;
  onCommentPress?: (event: Event) => void;
  onSharePress?: (event: Event) => void;
  initialLiked?: boolean;
  likeCount?: number;
  commentCount?: number;
};

const formatAddress = (event: Event) => {
  const addr = event.event_address?.[0];
  if (!addr) return "";
  const parts = [addr.street, addr.ward, addr.district, addr.province].filter(
    Boolean
  );
  return parts.join(", ");
};

export const EventPostCard: React.FC<Props> = React.memo(
  ({
    event,
    onPress,
    onLikePress,
    onCommentPress,
    onSharePress,
    initialLiked = false,
    likeCount = 0,
    commentCount = 0,
  }) => {
    const [isLiked, setIsLiked] = useState(initialLiked);
    const [currentLikeCount, setCurrentLikeCount] = useState(likeCount);
    const [imageError, setImageError] = useState(false);

    // Validate required props
    if (!event) {
      console.error("EventPostCard: event prop is required");
      return null;
    }

    const handleImageError = (
      error: NativeSyntheticEvent<ImageErrorEventData>
    ) => {
      console.warn(
        "EventPostCard: Failed to load image",
        error.nativeEvent.error
      );
      setImageError(true);
    };

    const handleLikePress = () => {
      setIsLiked((prev) => {
        const next = !prev;
        setCurrentLikeCount((c) => (next ? c + 1 : Math.max(0, c - 1)));
        return next;
      });
      onLikePress?.(event);
    };

    const handleCommentPress = () => {
      onCommentPress?.(event);
    };

    const handleSharePress = () => {
      onSharePress?.(event);
    };
    return (
      <TouchableOpacity
        className="bg-white rounded-xl p-3 shadow-sm border border-gray-100 mb-6"
        onPress={() => onPress?.(event)}
        activeOpacity={0.9}
      >
        {/* Header: User image + Event title */}
        <View className="flex-row items-center mb-3">
          {/* User Avatar */}
          <View className="w-10 h-10 rounded-full bg-gray-200 items-center justify-center mr-2">
            <Text className="text-gray-600 font-semibold text-lg">
              {event.host_by?.first_name?.[0]}
              {event.host_by?.last_name?.[0]}
            </Text>
          </View>

          {/* Event Title and Host Info */}
          <View className="flex-1">
            <Text
              className="text-gray-900 font-medium text-sm"
              numberOfLines={1}
            >
              {event.host_by?.first_name} {event.host_by?.last_name}
            </Text>
            <Text className="text-gray-500 text-xs" numberOfLines={1}>
              {/* Should be created_at */}
              15 hours
            </Text>
          </View>
        </View>
        <View className="flex-row items-center justify-between mb-0">
          <Text
            className="text-gray-800 flex-1 font-semibold text-sm"
            numberOfLines={1}
          >
            {event.name}
          </Text>
          <Text
            className="text-gray-800 font-semibold text-sm ml-2"
            numberOfLines={1}
          >
            {event.date_of_event}
          </Text>
        </View>
        <View className="mb-3">
          <View className="flex-row items-center justify-between">
            {!!formatAddress(event) && (
              <Text className="text-gray-700 text-xs" numberOfLines={1}>
                {formatAddress(event)}
              </Text>
            )}
            <Text className="text-gray-700 text-xs ml-1">
              {event.limit} people
            </Text>
          </View>
        </View>

        {Boolean(event.event_image) && !imageError && (
          <View className="mb-3 overflow-hidden rounded-lg">
            <Image
              source={{ uri: event.event_image as string }}
              className="w-full h-48"
              resizeMode="cover"
              onError={handleImageError}
              accessibilityLabel={`Event image for ${event.name}`}
            />
          </View>
        )}

        {Boolean(event.event_image) && imageError && (
          <View className="mb-3 h-48 bg-gray-100 rounded-lg items-center justify-center">
            <Icon name="coffee" size={32} color="#9ca3af" />
            <Text className="text-gray-500 text-xs mt-1">
              Image unavailable
            </Text>
          </View>
        )}

        {/* Like and Comment Counts */}
        {(currentLikeCount > 0 || commentCount > 0) && (
          <View className="mb-2">
            <View className="flex flex-row items-center justify-between">
              {currentLikeCount > 0 && (
                <View className="flex-row items-center mr-4">
                  <Icon name="heart" size={14} color="#e11d48" />
                  <Text className="text-gray-700 text-[13px] ml-1">
                    {currentLikeCount}{" "}
                    {currentLikeCount === 1 ? "like" : "likes"}
                  </Text>
                </View>
              )}
              {commentCount > 0 && (
                <View className="flex-row items-center">
                  <Text className="text-gray-700 text-[13px] ml-1">
                    {commentCount} {commentCount === 1 ? "comment" : "comments"}
                  </Text>
                </View>
              )}
            </View>
          </View>
        )}
        {/* Actions */}
        <View className="border-t border-gray-100 pt-2">
          <View className="flex-row items-center">
            <TouchableOpacity
              className="flex-1 py-2 items-center justify-center"
              activeOpacity={0.7}
              onPress={handleLikePress}
            >
              <Icon
                name={isLiked ? "heart" : "heart-outline"}
                size={20}
                color={isLiked ? "#e11d48" : "#6b7280"}
              />
            </TouchableOpacity>
            <View className="w-px h-4 bg-gray-200" />
            <TouchableOpacity
              className="flex-1 py-2 items-center justify-center"
              activeOpacity={0.7}
              onPress={handleCommentPress}
            >
              <Icon name="chatbubble-outline" size={20} color="#6b7280" />
            </TouchableOpacity>
            <View className="w-px h-4 bg-gray-200" />
            <TouchableOpacity
              className="flex-1 py-2 items-center justify-center"
              activeOpacity={0.7}
              onPress={handleSharePress}
            >
              <Icon name="share-outline" size={20} color="#6b7280" />
            </TouchableOpacity>
          </View>
        </View>
      </TouchableOpacity>
    );
  }
);

EventPostCard.displayName = "EventPostCard";

export default EventPostCard;
