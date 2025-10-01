export interface Event {
  date_of_event: string;
  email_contact: string;
  end_time: string;
  event_image: string;
  event_address: EventAddress[];
  host_by: HostBy;
  id: string;
  is_public: boolean;
  limit: number;
  name: string;
  number_samples: number;
  phone_contact: string;
  register_date: string;
  register_status: string;
  start_time: string;
}

export interface EventAddress {
  district: string;
  latitude: string;
  longitude: string;
  phone: string;
  province: string;
  street: string;
  ward: string;
}

export interface HostBy {
  email: string;
  first_name: string;
  id: string;
  last_name: string;
  phone: string;
}

export interface Sample {
  altitude_grow: string;
  breed_name: string;
  grow_address: string;
  grow_nation: string;
  id: string;
  name: string;
  pre_processing: string;
  price: string;
  rating: number;
  roast_level: string;
  roastery_address: string;
  roastery_name: string;
  roasting_date: string;
}

export interface EventDetail {
  event: Event;
  samples: Sample[];
}
