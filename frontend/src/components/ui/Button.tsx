import React from 'react';
import { TouchableOpacity, Text, ActivityIndicator, View } from 'react-native';

interface ButtonProps {
  title: string;
  onPress: () => void;
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  disabled?: boolean;
  className?: string;
  icon?: React.ReactNode;
}

const Button: React.FC<ButtonProps> = ({
  title,
  onPress,
  variant = 'primary',
  size = 'md',
  loading = false,
  disabled = false,
  className = '',
  icon,
}) => {
  const getVariantClasses = () => {
    if (disabled || loading) {
      return 'bg-gray-300';
    }
    
    switch (variant) {
      case 'primary':
        return 'bg-primary-600 active:bg-primary-700';
      case 'secondary':
        return 'bg-secondary-500 active:bg-secondary-600';
      case 'outline':
        return 'border-2 border-primary-600 bg-transparent active:bg-primary-50';
      case 'ghost':
        return 'bg-transparent active:bg-gray-100';
      case 'danger':
        return 'bg-error-600 active:bg-error-700';
      default:
        return 'bg-primary-600 active:bg-primary-700';
    }
  };

  const getSizeClasses = () => {
    switch (size) {
      case 'sm':
        return 'py-2 px-4';
      case 'md':
        return 'py-3 px-6';
      case 'lg':
        return 'py-4 px-8';
      default:
        return 'py-3 px-6';
    }
  };

  const getTextClasses = () => {
    let baseClasses = 'font-semibold';
    let colorClasses = '';
    let sizeClasses = '';

    // Color classes
    if (disabled || loading) {
      colorClasses = 'text-gray-500';
    } else {
      switch (variant) {
        case 'outline':
          colorClasses = 'text-primary-600';
          break;
        case 'ghost':
          colorClasses = 'text-gray-700';
          break;
        default:
          colorClasses = 'text-white';
      }
    }

    // Size classes
    switch (size) {
      case 'sm':
        sizeClasses = 'text-sm';
        break;
      case 'lg':
        sizeClasses = 'text-lg';
        break;
      default:
        sizeClasses = 'text-base';
    }

    return `${baseClasses} ${colorClasses} ${sizeClasses}`;
  };

  return (
    <TouchableOpacity
      className={`
        ${getVariantClasses()}
        ${getSizeClasses()}
        rounded-lg
        flex-row justify-center items-center
        ${className}
      `}
      onPress={onPress}
      disabled={disabled || loading}
    >
      {loading && (
        <ActivityIndicator 
          size="small" 
          color={variant === 'outline' || variant === 'ghost' ? '#3b82f6' : 'white'}
          className="mr-2"
        />
      )}
      {icon && !loading && (
        <View className="mr-2">
          {icon}
        </View>
      )}
      <Text className={getTextClasses()}>
        {title}
      </Text>
    </TouchableOpacity>
  );
};

export default Button;