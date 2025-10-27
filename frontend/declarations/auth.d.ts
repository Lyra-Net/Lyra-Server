export interface LoginPayload {
  username: string;
  password: string;
  device_id: string;
  remember_device?: boolean;
  session_id?: string;
}

export interface RegisterPayload {
  username: string;
  display_name: string;
  password: string;
  email: string;
  device_id: string;
}

export interface RefreshPayload {}

export interface LogoutPayload {}

export interface ChangePasswordPayload {}

export interface ForgotPasswordPayload {}

export interface AddEmailPayload {
  email: string;
  device_id: string;
  verified_2fa_session_id?: string;
}

export interface RemoveEmailPayload {}

export interface VerifyEmailPayload {}

export interface Toggle2FAPayload {
  device_id: string;
  enable: boolean;
  verified_2fa_session_id?: string;
}

export interface VerifyCodePayload {
  session_id: string;
  otp: string;
  device_id: string;
  remember_device?: boolean;
}