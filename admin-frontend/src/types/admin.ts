export interface GuestbookEntry {
  id: number;
  nickname: string;
  message: string;
  secret: boolean;
  ip: string;
  hidden: boolean;
  evaluated: boolean;
  created_at: string;
}

export interface GameScore {
  id: number;
  nickname: string;
  time_ms: number;
  ip: string;
  created_at: string;
}

export interface PhotoUpload {
  id: number;
  name: string;
  upload_date: string;
  ip_address: string;
  hashname: string;
  original_filename: string;
  hidden: boolean;
  evaluated: boolean;
  hearts: number;
  url?: string;
  original_url?: string;
  thumbnail?: string;
}

export interface BusInfoEntry {
  stop: string;
  routes: string;
}

export interface CharterBusEntry {
  location: string;
  company: string;
  bus_number: string;
  departure: string;
}

export interface MapLinks {
  google: string;
  kakao: string;
  naver: string;
  tmap: string;
}

export interface MapProviders {
  embed_provider: string;
  latitude: number;
  longitude: number;
  api_key: string;
  links: MapLinks;
}

export interface WeddingConfig {
  groom_eng_name: string;
  groom_kor_name: string;
  bride_eng_name: string;
  bride_kor_name: string;
  groom_father_kor_name: string;
  groom_mother_kor_name: string;
  bride_father_kor_name: string;
  bride_mother_kor_name: string;
  groom_bank_account: string;
  bride_bank_account: string;
  groom_father_bank_account: string;
  groom_mother_bank_account: string;
  bride_father_bank_account: string;
  bride_mother_bank_account: string;
  wedding_datetime: string;
  venue_name: string;
  venue_address: string;
  venue_floor: string;
  venue_hall: string;
  venue_phone: string;
  map_providers: MapProviders;
  subway_info: string[];
  bus_info: BusInfoEntry[];
  charter_bus: CharterBusEntry[];
  charter_bus_notice: string;
  car_info: string;
  groom_birth_order: string;
  bride_birth_order: string;
  card_game_timer: number;
  game_npc_message: string;
  avatar_colors: string;
  short_greeting: string;
  main_greet_text: string;
  photo_upload_enabled: boolean;
  photo_upload_hours_before: number;
  hearts_flush_interval_ms: number;
  hearts_flush_batch_size: number;
  moderation_thresholds: Record<string, number>;
}

export interface AssetPhoto {
  id: number;
  label: string;
  hashname: string;
  thumb_hashname: string;
  original_filename: string;
  use_for_game: boolean;
  is_main_photo: boolean;
  sort_order: number;
  created_at: string;
  url?: string;
  thumbnail_url?: string;
}

export interface ModerationCategoryStats {
  total: number;
  pending: number;
  approved: number;
  flagged: number;
}

export interface ModerationStatus {
  enabled: boolean;
  queue_depth: number;
  guestbook: ModerationCategoryStats;
  photos: ModerationCategoryStats;
  ws_connections: number;
  total_hearts: number;
}

export interface SystemSettings {
  bcrypt_cost: number;
  game_timer_ms: number;
  rate_limit_enabled: boolean;
  s3_bucket: string;
  s3_region: string;
  s3_endpoint: string;
  s3_access_key: string;
  s3_access_key_set: boolean;
  s3_secret_key: string;
  s3_secret_key_set: boolean;
  use_moderation: boolean;
  openai_api_key: string;
  openai_api_key_set: boolean;
}

export interface SystemSettingsUpdateResponse {
  status: string;
  restart_required: boolean;
  s3_reinitialized: boolean;
}

export interface S3TestResult {
  success: boolean;
  error?: string;
}

export interface HallOfFameEntry {
  id: number;
  nickname: string;
  ip: string;
  created_at: string;
}
