import { NextRequest, NextResponse } from "next/server";
import decodeJwt from "@/utils/jwtdecode";
import axios, { AxiosError } from "axios";

interface JWTPayload {
  sub: string;
  iat: number;
  exp: number;
  display_name: string;
}


export async function POST(req: NextRequest) {
  const { username, password, device_id } = await req.json();

  const userAgent = req.headers.get("user-agent") || "unknown-UA";
  const clientIp = req.headers.get("x-forwarded-for") || req.headers.get("x-real-ip") || "0.0.0.0";
  // console.log("User-Agent:", userAgent);
  // console.log("Client IP:", clientIp);
  // console.log("req headers:", req.headers);
  try {
    const res = await axios.post(
      `${process.env.NEXT_PUBLIC_API_URL}/auth/login`,
      { username, password, device_id },
      {
        headers: {
          "user-agent": userAgent,
          "x-forwarded-for": clientIp,
        },
      }
    );

    const { access_token, refresh_token } = res.data;
    const user: JWTPayload = decodeJwt(access_token);
    console.log("Logged in user:", user);
    const response = NextResponse.json({
      access_token,
      userId: user.sub,
      displayName: user.display_name,
      exp: user.exp,
    });

    response.cookies.set({
      name: "refresh_token",
      value: refresh_token,
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "strict",
      path: "/",
      maxAge: 60 * 60 * 24 * 7, // 7 days
    });

    response.cookies.set({
      name: "access_token",
      value: access_token,
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "strict",
      path: "/",
      maxAge: 60 * 60, // 1h
    });

    return response;
  } catch (err) {
    const axiosErr = err as AxiosError;

    if (axiosErr.response?.status === 403) {
      const grpcMessage = axiosErr.response.headers["grpc-message"];

      let sessionId: string | null = null;
      if (grpcMessage && grpcMessage.startsWith("2FA_REQUIRED:SESSION:")) {
        sessionId = grpcMessage.split("2FA_REQUIRED:SESSION:")[1];
      }

      console.log("Extracted 2FA sessionId:", sessionId);

      return NextResponse.json(
        { message: "2FA required", session_id: sessionId },
        { status: 403 }
      );
    }

    console.error("Login error:", axiosErr);
    return NextResponse.json(
      { message: "Login failed" },
      { status: 401 }
    );
  }
}

