export interface IUser {
  id: string;
  name: string;
  email: string;
  otp_enabled: string;
  telegram_chat_id: string;
}

export interface GenericResponse {
  status: string;
  message: string;
}

export interface ILoginResponse {
  status: string;
  user: IUser;
}
