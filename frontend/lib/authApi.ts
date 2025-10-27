import { AddEmailPayload, Toggle2FAPayload, VerifyCodePayload } from './../declarations/auth.d';
import api from "@/lib/api";


export async function getProfile() {
  const response = await api.get("/auth/profile");
  return response.data;
}

export async function toggle2Fa(payload: Toggle2FAPayload) {
  const response = await api.post("/auth/toggle-2fa", { ...payload
   });
  return response.data;
}

export async function verifyCode(payload: VerifyCodePayload) {
  const response = await api.post("/auth/verify-code", { ...payload });
  return response.data;
}

export async function updateProfile(arg: any) {
  const response = await api.put("/auth/profile", { arg });
  return response.data;
}

export async function addEmail(addEmailPayload: AddEmailPayload) {
  const response = await api.post("/auth/add-email", { ...addEmailPayload });
  return response.data;
}