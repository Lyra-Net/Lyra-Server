import { NextRequest, NextResponse } from "next/server";
import axios from "axios";

export async function POST(req: NextRequest) {
  const { username, password, device_id } = await req.json();

  try {
    const res = await axios.post(`${process.env.NEXT_PUBLIC_API_URL}/auth/login`, {
      username,
      password,
      device_id,
    });

    const { access_token, refresh_token } = res.data;

    const response = NextResponse.json({ access_token });

    response.cookies.set({
      name: "refresh_token",
      value: refresh_token,
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "strict",
      path: "/",
      maxAge: 60 * 60 * 24 * 7, // 7 days
    });
    console.log({response})
    return response;
  } catch (err) {
    console.error(err);
    return NextResponse.json({ message: "Login failed" }, { status: 401 });
  }
}
