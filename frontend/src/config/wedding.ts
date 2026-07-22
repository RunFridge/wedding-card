// ─── Names ───
export let GROOM_ENG_NAME = 'Groom';
export let GROOM_KOR_NAME = '김철수';
export let BRIDE_ENG_NAME = 'Bride';
export let BRIDE_KOR_NAME = '이영희';
export let GROOM_FATHER_KOR_NAME = '김아버지';
export let GROOM_MOTHER_KOR_NAME = '박어머니';
export let BRIDE_FATHER_KOR_NAME = '이아버지';
export let BRIDE_MOTHER_KOR_NAME = '최어머니';

// ─── Bank Accounts ───
export let GROOM_BANK_ACCOUNT = '카카오뱅크 0000-00-0000000';
export let BRIDE_BANK_ACCOUNT = '은행 000-000000-00000';
export let GROOM_FATHER_BANK_ACCOUNT = '은행 000-000000-00000';
export let GROOM_MOTHER_BANK_ACCOUNT = '은행 000-000000-00000';
export let BRIDE_FATHER_BANK_ACCOUNT = '은행 000-000000-00000';
export let BRIDE_MOTHER_BANK_ACCOUNT = '은행 000-000000-00000';

// ─── Wedding Details ───
export let WEDDING_DATETIME = new Date('2030-01-01T11:00:00+09:00');
export let VENUE_NAME = 'OO웨딩홀';
export let VENUE_ADDRESS = '서울특별시 중구 세종대로 110';
export let VENUE_FLOOR = '3층';
export let VENUE_HALL = '그랜드홀';
export let VENUE_PHONE = '02-000-0000';

// ─── Map Providers ───
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

export let MAP_PROVIDERS: MapProviders = {
  embed_provider: '',
  latitude: 0,
  longitude: 0,
  api_key: '',
  links: {
    google: '',
    kakao: '',
    naver: '',
    tmap: '',
  },
};

// ─── Transport Info ───
export let SUBWAY_INFO = ['1호선', '시청역 1번 출구 도보 5분 거리'];
export let BUS_INFO: { stop: string; routes: string }[] = [
  {
    stop: '시청역 1번 출구 정류장',
    routes: '100, 200, 300',
  },
];
export let CAR_INFO = '웨딩홀 전용 주차장 이용 가능';
export let CHARTER_BUS: {
  location: string;
  company: string;
  bus_number: string;
  departure: string;
}[] = [];
export let CHARTER_BUS_NOTICE = '';

// ─── Birth Order ───
export let GROOM_BIRTH_ORDER = '장남';
export let BRIDE_BIRTH_ORDER = '장녀';

// ─── Game ───
export let CARD_GAME_TIMER = 30000;
export let GAME_NPC_MESSAGE =
  '안녕! 게임을 클리어하면 신랑 신부의 특별한 사진들을 볼 수 있어요... 도전해 보세요!';

// ─── Dynamic Assets ───
export let MAIN_PHOTO_URL = '';
export let GAME_PHOTO_COUNT = 0;

// ─── Avatar ───
export let AVATAR_COLORS = '8B6914,A0722A,C4943A,5C3A0E,F5E6C8';

// ─── Photo Upload ───
export let PHOTO_UPLOAD_ENABLED = false;
export let PHOTO_UPLOAD_HOURS_BEFORE = 1;

// ─── Greeting ───
export let SHORT_GREETING = '저희 결혼합니다';
export let MAIN_GREET_TEXT =
  '소중한 분들을 초대합니다.\n함께 축복해 주시면 더없는 기쁨으로 간직하겠습니다.';

// ─── Demo Mode ───
export let DEMO_MODE = false;

// ─── Runtime Config Loader ───
export async function loadConfig(): Promise<void> {
  try {
    const res = await fetch('/api/config');
    if (!res.ok) return;
    const c = await res.json();

    if (c.demo === true) DEMO_MODE = true;

    if (c.groom_eng_name) GROOM_ENG_NAME = c.groom_eng_name;
    if (c.groom_kor_name) GROOM_KOR_NAME = c.groom_kor_name;
    if (c.bride_eng_name) BRIDE_ENG_NAME = c.bride_eng_name;
    if (c.bride_kor_name) BRIDE_KOR_NAME = c.bride_kor_name;
    if (c.groom_father_kor_name)
      GROOM_FATHER_KOR_NAME = c.groom_father_kor_name;
    if (c.groom_mother_kor_name)
      GROOM_MOTHER_KOR_NAME = c.groom_mother_kor_name;
    if (c.bride_father_kor_name)
      BRIDE_FATHER_KOR_NAME = c.bride_father_kor_name;
    if (c.bride_mother_kor_name)
      BRIDE_MOTHER_KOR_NAME = c.bride_mother_kor_name;

    if (c.groom_bank_account) GROOM_BANK_ACCOUNT = c.groom_bank_account;
    if (c.bride_bank_account) BRIDE_BANK_ACCOUNT = c.bride_bank_account;
    if (c.groom_father_bank_account)
      GROOM_FATHER_BANK_ACCOUNT = c.groom_father_bank_account;
    if (c.groom_mother_bank_account)
      GROOM_MOTHER_BANK_ACCOUNT = c.groom_mother_bank_account;
    if (c.bride_father_bank_account)
      BRIDE_FATHER_BANK_ACCOUNT = c.bride_father_bank_account;
    if (c.bride_mother_bank_account)
      BRIDE_MOTHER_BANK_ACCOUNT = c.bride_mother_bank_account;

    if (c.wedding_datetime) {
      WEDDING_DATETIME = new Date(c.wedding_datetime);
    }

    if (c.venue_name) VENUE_NAME = c.venue_name;
    if (c.venue_address) VENUE_ADDRESS = c.venue_address;
    if (c.venue_floor) VENUE_FLOOR = c.venue_floor;
    if (c.venue_hall) VENUE_HALL = c.venue_hall;
    if (c.venue_phone) VENUE_PHONE = c.venue_phone;

    if (c.map_providers) MAP_PROVIDERS = c.map_providers;

    if (c.subway_info) SUBWAY_INFO = c.subway_info;
    if (c.bus_info) BUS_INFO = c.bus_info;
    if (c.car_info) CAR_INFO = c.car_info;
    if (c.charter_bus) CHARTER_BUS = c.charter_bus;
    if (c.charter_bus_notice) CHARTER_BUS_NOTICE = c.charter_bus_notice;

    if (c.groom_birth_order) GROOM_BIRTH_ORDER = c.groom_birth_order;
    if (c.bride_birth_order) BRIDE_BIRTH_ORDER = c.bride_birth_order;

    if (c.card_game_timer) CARD_GAME_TIMER = c.card_game_timer;
    if (c.game_npc_message) GAME_NPC_MESSAGE = c.game_npc_message;

    if (c.avatar_colors) AVATAR_COLORS = c.avatar_colors;
    if (c.short_greeting) SHORT_GREETING = c.short_greeting;
    if (c.main_greet_text) MAIN_GREET_TEXT = c.main_greet_text;

    if (c.photo_upload_enabled === true) PHOTO_UPLOAD_ENABLED = true;
    if (
      typeof c.photo_upload_hours_before === 'number' &&
      c.photo_upload_hours_before >= 0
    )
      PHOTO_UPLOAD_HOURS_BEFORE = c.photo_upload_hours_before;

    if (c.main_photo_url) MAIN_PHOTO_URL = c.main_photo_url;
    if (c.game_photo_count != null) GAME_PHOTO_COUNT = c.game_photo_count;
  } catch {
    // silently use defaults
  }
}
