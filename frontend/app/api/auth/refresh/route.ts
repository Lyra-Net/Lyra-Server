import { NextRequest, NextResponse } from "next/server";
import axios from "axios";

export async function POST(req: NextRequest) {
  try {

    const { accessToken, device_id } = await req.json();
    const refresh_token = req.cookies.get("refresh_token")?.value;

    if (!refresh_token) {
      return NextResponse.json({ message: "No refresh token" }, { status: 401 });
    }

    const res = await axios.post(`${process.env.API_URL}/auth/refresh`, {
      accessToken,
      refresh_token,
      device_id,
    });

    const { access_token, refresh_token: new_refresh } = res.data;

    const response = NextResponse.json({ access_token });

    if (new_refresh) {
      response.cookies.set("refresh_token", new_refresh, {
        httpOnly: true,
        secure: process.env.NODE_ENV === "production",
        sameSite: "strict",
        path: "/",
        maxAge: 60 * 60 * 24 * 7,
      });
    }

    return response;
  } catch (err) {
    console.error("Refresh error", err);
    return NextResponse.json({ message: "Refresh failed" }, { status: 401 });
  }
}
