import { NextRequest, NextResponse } from 'next/server';
import axios from 'axios';
import { getServerSession } from "next-auth/next";
import { authOptions } from '@/lib/authOptions';

export async function POST(req: NextRequest) {
  try {
    const session = await getServerSession(authOptions);
    if (!session) {
      return NextResponse.json({ message: 'No session' }, { status: 401 });
    }

    const accessToken = (session as any).accessToken;
    const refreshToken = req.cookies.get('refresh_token')?.value;

    if (!refreshToken) {
      return NextResponse.json({ message: 'No refresh token' }, { status: 401 });
    }

    const backendRes = await axios.post(`${process.env.NEXT_PUBLIC_API_URL}/auth/refreshToken`, {
      access_token: accessToken,
      refresh_token: refreshToken,
    });

    const { access_token: newAccessToken, refresh_token: newRefreshToken } = backendRes.data;

    const res = NextResponse.json({ accessToken: newAccessToken });
    res.cookies.set({
      name: 'refresh_token',
      value: newRefreshToken,
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict',
      path: '/',
      maxAge: 60 * 60 * 24 * 7, // 7 days
    });

    return res;
  } catch (err) {
    console.error(err);
    return NextResponse.json({ message: 'Failed to refresh token' }, { status: 500 });
  }
}
