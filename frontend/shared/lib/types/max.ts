export type InitData = {
  query_id: string;
  ip?: string;
  auth_date: number;
  hash: string;
  user: {
    id: number;
    first_name: string;
    last_name: string;
    username: string;
    language_code: string;
    photo_url: string;
  };
  chat: {
    id: number;
    type: 'DIALOG' | 'CHAT' | 'CHANNEL';
  };
  start_param: string;
};

export type MaxUser = {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string | null;
  language_code?: string | null;
  photo_url?: string | null;
};
