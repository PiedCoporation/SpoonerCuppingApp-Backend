import React, { useMemo, useRef, useState, useCallback } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  TouchableWithoutFeedback,
  ScrollView,
  Dimensions,
  KeyboardAvoidingView,
  Platform,
  Keyboard,
} from "react-native";
import { useHeaderHeight } from "@react-navigation/elements";
import { SafeAreaView } from "react-native-safe-area-context";
import { useForm, Controller } from "react-hook-form";

const { width: SCREEN_WIDTH } = Dimensions.get("window");

type EventInfo = {
  name: string;
  dateOfEvent: string;
  registerDate: string;
  startTime: string;
  endTime: string;
  limit: string;
  numberSamples: string;
  phoneContact: string;
  emailContact: string;
  isPublic: boolean;
};

type Address = {
  line1: string;
  line2?: string;
  city?: string;
  state?: string;
  postalCode?: string;
  country?: string;
};

type Sample = {
  name: string;
  roastingDate: string;
  roastLevel: string;
  altitudeGrow: string;
  roasteryName: string;
  roasteryAddress: string;
  breedName: string;
  preProcessing: string;
  growNation: string;
  growAddress: string;
  price: string;
};

type FormData = {
  eventInfo: EventInfo;
  addresses: Address[];
  samples: Sample[];
};

export default function CreateCuppingEvent() {
  const pagerRef = useRef<ScrollView>(null);
  const [currentIndex, setCurrentIndex] = useState(0);
  const headerHeight = useHeaderHeight();

  // Initialize react-hook-form
  const { control, watch, setValue, getValues, handleSubmit } =
    useForm<FormData>({
      defaultValues: {
        eventInfo: {
          name: "",
          dateOfEvent: "",
          registerDate: "",
          startTime: "",
          endTime: "",
          limit: "",
          numberSamples: "",
          phoneContact: "",
          emailContact: "",
          isPublic: true,
        },
        addresses: [{ line1: "" }],
        samples: [
          {
            name: "",
            roastingDate: "",
            roastLevel: "",
            altitudeGrow: "",
            roasteryName: "",
            roasteryAddress: "",
            breedName: "",
            preProcessing: "",
            growNation: "",
            growAddress: "",
            price: "",
          },
        ],
      },
    });

  const addresses = watch("addresses");
  const samples = watch("samples");
  const totalSlides = useMemo(() => 2 + samples.length, [samples.length]);

  // No keyboard listeners; avoid layout thrash that can blur inputs

  const goToIndex = useCallback(
    (index: number) => {
      const clamped = Math.max(0, Math.min(index, totalSlides - 1));
      setCurrentIndex(clamped);
      pagerRef.current?.scrollTo({ x: clamped * SCREEN_WIDTH, animated: true });
    },
    [totalSlides]
  );

  const handleAddAddress = useCallback(() => {
    const currentAddresses = getValues("addresses");
    setValue("addresses", [...currentAddresses, { line1: "" }]);
  }, [setValue, getValues]);

  const handleRemoveAddress = useCallback(
    (idx: number) => {
      const currentAddresses = getValues("addresses");
      if (currentAddresses.length > 1) {
        setValue(
          "addresses",
          currentAddresses.filter((_, i) => i !== idx)
        );
      }
    },
    [setValue, getValues]
  );

  const handleAddSample = useCallback(() => {
    const currentSamples = getValues("samples");
    setValue("samples", [
      ...currentSamples,
      {
        name: "",
        roastingDate: "",
        roastLevel: "",
        altitudeGrow: "",
        roasteryName: "",
        roasteryAddress: "",
        breedName: "",
        preProcessing: "",
        growNation: "",
        growAddress: "",
        price: "",
      },
    ]);
  }, [setValue, getValues]);

  const handleRemoveSample = useCallback(
    (idx: number) => {
      const currentSamples = getValues("samples");
      if (currentSamples.length > 1) {
        setValue(
          "samples",
          currentSamples.filter((_, i) => i !== idx)
        );
      }
    },
    [setValue, getValues]
  );

  const onSubmit = useCallback((data: FormData) => {
    const payload = {
      ...data.eventInfo,
      limit: Number(data.eventInfo.limit || 0),
      numberSamples: Number(data.eventInfo.numberSamples || 0),
      addresses: data.addresses,
      samples: data.samples.map((s) => ({ ...s, price: Number(s.price || 0) })),
    };
    console.log("Create Event Payload", payload);
  }, []);

  const renderHeader = () => (
    <View className="px-6 py-4 border-b border-gray-200 bg-white">
      <View className="mt-2 flex-row items-center">
        {Array.from({ length: totalSlides }).map((_, i) => (
          <View
            key={`dot-${i}`}
            className={`h-2 rounded-full mr-2 ${
              i === currentIndex ? "bg-amber-700 w-6" : "bg-gray-300 w-2"
            }`}
          />
        ))}
      </View>
    </View>
  );

  const NavControls = () => (
    <View className="flex-row justify-between items-center px-6 py-3 border-t border-gray-200 bg-white">
      <TouchableOpacity
        onPress={() => goToIndex(currentIndex - 1)}
        disabled={currentIndex === 0}
      >
        <Text
          className={`text-base ${
            currentIndex === 0 ? "text-gray-400" : "text-amber-700"
          }`}
        >
          Back
        </Text>
      </TouchableOpacity>
      <View className="flex-row items-center">
        <TouchableOpacity
          onPress={() => goToIndex(currentIndex + 1)}
          disabled={currentIndex >= totalSlides - 1}
        >
          <Text
            className={`text-base ${
              currentIndex >= totalSlides - 1
                ? "text-gray-400"
                : "text-amber-700"
            }`}
          >
            Next
          </Text>
        </TouchableOpacity>
      </View>
    </View>
  );

  const ControlledInput = ({
    name,
    label,
    placeholder,
    keyboardType = "default",
    multiline = false,
    rules = {},
  }: {
    name: string;
    label: string;
    placeholder?: string;
    keyboardType?: "default" | "email-address" | "numeric" | "phone-pad";
    multiline?: boolean;
    rules?: any;
  }) => (
    <View className="mb-4">
      <Text className="mb-1 text-sm text-gray-600">{label}</Text>
      <Controller
        control={control}
        name={name as any}
        rules={rules}
        render={({ field: { onChange, onBlur, value } }) => (
          <TextInput
            className={`border border-gray-300 rounded-md px-3 py-2 ${
              multiline ? "min-h-[80px]" : ""
            }`}
            value={value || ""}
            onChangeText={onChange}
            onBlur={onBlur}
            placeholder={placeholder}
            keyboardType={keyboardType}
            multiline={multiline}
            textAlignVertical={multiline ? "top" : "center"}
            // iOS optimizations
            autoCorrect={false}
            spellCheck={false}
            autoCapitalize="none"
            returnKeyType={multiline ? "default" : "done"}
            enablesReturnKeyAutomatically={true}
            clearButtonMode="while-editing"
            blurOnSubmit={!multiline}
            // Prevent keyboard dismissal issues
            onSubmitEditing={() => {}}
          />
        )}
      />
    </View>
  );

  const SlideEventInfo = () => (
    <View style={{ width: SCREEN_WIDTH, flex: 1 }}>
      <ScrollView
        className="px-6 py-4"
        keyboardShouldPersistTaps="always"
        contentContainerStyle={{
          paddingBottom: 120,
        }}
        showsVerticalScrollIndicator={false}
        automaticallyAdjustContentInsets={false}
        contentInsetAdjustmentBehavior="never"
      >
        <Text className="text-lg font-semibold mb-3">Event Information</Text>
        <ControlledInput name="eventInfo.name" label="Name" multiline={true} />
        <ControlledInput
          name="eventInfo.dateOfEvent"
          label="Date Of Event"
          placeholder="YYYY-MM-DD"
        />
        <ControlledInput
          name="eventInfo.registerDate"
          label="Register Date"
          placeholder="YYYY-MM-DD"
        />
        <View className="flex-row gap-3">
          <View className="flex-1">
            <ControlledInput
              name="eventInfo.startTime"
              label="Start Time"
              placeholder="HH:mm"
            />
          </View>
          <View className="flex-1">
            <ControlledInput
              name="eventInfo.endTime"
              label="End Time"
              placeholder="HH:mm"
            />
          </View>
        </View>
        <View className="flex-row gap-3">
          <View className="flex-1">
            <ControlledInput
              name="eventInfo.limit"
              label="Limit"
              keyboardType="numeric"
            />
          </View>
          <View className="flex-1">
            <ControlledInput
              name="eventInfo.numberSamples"
              label="Number of Samples"
              keyboardType="numeric"
            />
          </View>
        </View>
        <ControlledInput
          name="eventInfo.phoneContact"
          label="Phone Contact"
          keyboardType="phone-pad"
        />
        <ControlledInput
          name="eventInfo.emailContact"
          label="Email Contact"
          keyboardType="email-address"
        />
        <View className="mt-2" />
        <Controller
          control={control}
          name="eventInfo.isPublic"
          render={({ field: { onChange, value } }) => (
            <View className="flex-row items-center justify-between mt-2">
              <Text className="text-sm text-gray-700">Public Event</Text>
              <TouchableOpacity onPress={() => onChange(!value)}>
                <Text className="text-amber-700">{value ? "Yes" : "No"}</Text>
              </TouchableOpacity>
            </View>
          )}
        />
      </ScrollView>
    </View>
  );

  const SlideAddresses = () => (
    <View style={{ width: SCREEN_WIDTH, flex: 1 }}>
      <ScrollView
        className="px-6 py-4"
        keyboardShouldPersistTaps="always"
        contentContainerStyle={{
          paddingBottom: 120,
        }}
        showsVerticalScrollIndicator={false}
        automaticallyAdjustContentInsets={false}
        contentInsetAdjustmentBehavior="never"
      >
        <Text className="text-lg font-semibold mb-3">Addresses</Text>
        {addresses.map((_, idx) => (
          <View
            key={`addr-${idx}`}
            className="mb-4 border border-gray-200 rounded-md p-3"
          >
            <View className="flex-row justify-between items-center mb-2">
              <Text className="font-medium">Address {idx + 1}</Text>
              <TouchableOpacity
                onPress={() => handleRemoveAddress(idx)}
                disabled={addresses.length === 1}
              >
                <Text
                  className={`text-sm ${
                    addresses.length === 1 ? "text-gray-400" : "text-red-600"
                  }`}
                >
                  Remove
                </Text>
              </TouchableOpacity>
            </View>
            <ControlledInput
              name={`addresses.${idx}.line1`}
              label="Line 1"
              multiline={true}
            />
            <ControlledInput
              name={`addresses.${idx}.line2`}
              label="Line 2"
              multiline={true}
            />
            <View className="flex-row gap-3">
              <View className="flex-1">
                <ControlledInput name={`addresses.${idx}.city`} label="City" />
              </View>
              <View className="flex-1">
                <ControlledInput
                  name={`addresses.${idx}.state`}
                  label="State"
                />
              </View>
            </View>
            <View className="flex-row gap-3">
              <View className="flex-1">
                <ControlledInput
                  name={`addresses.${idx}.postalCode`}
                  label="Postal Code"
                />
              </View>
              <View className="flex-1">
                <ControlledInput
                  name={`addresses.${idx}.country`}
                  label="Country"
                />
              </View>
            </View>
          </View>
        ))}
        <TouchableOpacity className="mt-2" onPress={handleAddAddress}>
          <Text className="text-amber-700">+ Add another address</Text>
        </TouchableOpacity>
      </ScrollView>
    </View>
  );

  const SampleCard = ({ idx }: { idx: number }) => (
    <View className="mb-4 border border-gray-200 rounded-md p-3">
      <View className="flex-row justify-between items-center mb-2">
        <Text className="font-medium">Sample {idx + 1}</Text>
        <TouchableOpacity
          onPress={() => handleRemoveSample(idx)}
          disabled={samples.length === 1}
        >
          <Text
            className={`text-sm ${
              samples.length === 1 ? "text-gray-400" : "text-red-600"
            }`}
          >
            Remove
          </Text>
        </TouchableOpacity>
      </View>
      <ControlledInput
        name={`samples.${idx}.name`}
        label="Name"
        multiline={true}
      />
      <ControlledInput
        name={`samples.${idx}.roastingDate`}
        label="Roasting Date"
        placeholder="YYYY-MM-DD"
      />
      <View className="flex-row gap-3">
        <View className="flex-1">
          <ControlledInput
            name={`samples.${idx}.roastLevel`}
            label="Roast Level"
          />
        </View>
        <View className="flex-1">
          <ControlledInput
            name={`samples.${idx}.altitudeGrow`}
            label="Altitude Grow"
          />
        </View>
      </View>
      <ControlledInput
        name={`samples.${idx}.roasteryName`}
        label="Roastery Name"
        multiline={true}
      />
      <ControlledInput
        name={`samples.${idx}.roasteryAddress`}
        label="Roastery Address"
        multiline={true}
      />
      <View className="flex-row gap-3">
        <View className="flex-1">
          <ControlledInput
            name={`samples.${idx}.breedName`}
            label="Breed Name"
          />
        </View>
        <View className="flex-1">
          <ControlledInput
            name={`samples.${idx}.preProcessing`}
            label="Pre Processing"
          />
        </View>
      </View>
      <View className="flex-row gap-3">
        <View className="flex-1">
          <ControlledInput
            name={`samples.${idx}.growNation`}
            label="Grow Nation"
          />
        </View>
        <View className="flex-1">
          <ControlledInput
            name={`samples.${idx}.growAddress`}
            label="Grow Address"
            multiline={true}
          />
        </View>
      </View>
      <ControlledInput
        name={`samples.${idx}.price`}
        label="Price"
        keyboardType="numeric"
      />
    </View>
  );

  const SlideSample = ({ sampleIndex }: { sampleIndex: number }) => (
    <View style={{ width: SCREEN_WIDTH, flex: 1 }}>
      <ScrollView
        className="px-6 py-4"
        keyboardShouldPersistTaps="always"
        contentContainerStyle={{
          paddingBottom: 120,
        }}
        showsVerticalScrollIndicator={false}
        automaticallyAdjustContentInsets={false}
        contentInsetAdjustmentBehavior="never"
      >
        <Text className="text-lg font-semibold mb-3">
          Sample {sampleIndex + 1}
        </Text>
        <SampleCard idx={sampleIndex} />
        <View className="flex-row justify-between items-center">
          <TouchableOpacity onPress={handleAddSample}>
            <Text className="text-amber-700">+ Add another sample</Text>
          </TouchableOpacity>
          <TouchableOpacity
            className="bg-amber-700 px-4 py-2 rounded-md"
            onPress={handleSubmit(onSubmit)}
          >
            <Text className="text-white font-semibold">Create Event</Text>
          </TouchableOpacity>
        </View>
      </ScrollView>
    </View>
  );

  const onMomentumScrollEnd = useCallback(
    (e: any) => {
      const x = e.nativeEvent.contentOffset.x as number;
      const idx = Math.round(x / SCREEN_WIDTH);
      if (idx !== currentIndex) setCurrentIndex(idx);
    },
    [currentIndex]
  );

  return (
    <SafeAreaView
      className="flex-1 bg-white"
      edges={["left", "right", "bottom"]}
    >
      <KeyboardAvoidingView
        style={{ flex: 1 }}
        behavior={Platform.OS === "ios" ? "padding" : "height"}
        keyboardVerticalOffset={Platform.OS === "ios" ? headerHeight : 20}
      >
        <TouchableWithoutFeedback onPress={Keyboard.dismiss} accessible={false}>
          <View style={{ flex: 1 }}>
            {renderHeader()}
            <ScrollView
              ref={pagerRef}
              horizontal
              pagingEnabled
              showsHorizontalScrollIndicator={false}
              onMomentumScrollEnd={onMomentumScrollEnd}
              style={{ flex: 1 }}
              keyboardShouldPersistTaps="always"
              scrollEnabled={false}
              keyboardDismissMode="none"
              automaticallyAdjustContentInsets={false}
              contentInsetAdjustmentBehavior="never"
            >
              <SlideEventInfo />
              <SlideAddresses />
              {samples.map((_, i) => (
                <SlideSample key={`sample-slide-${i}`} sampleIndex={i} />
              ))}
            </ScrollView>
            <NavControls />
          </View>
        </TouchableWithoutFeedback>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}
