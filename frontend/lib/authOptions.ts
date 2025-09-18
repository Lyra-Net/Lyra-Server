import CredentialsProvider from "next-auth/providers/credentials";
import { jwtDecode } from "jwt-decode";
import type { NextAuthOptions } from "next-auth";

interface JwtPayload {
  user_id: string;
  exp: number;
  iat: number;
  [key: string]: any;
}

export const authOptions: NextAuthOptions = {
  providers: [
    CredentialsProvider({
      name: "Credentials",
      credentials: {
        username: { label: "Username", type: "text" },
        password: { label: "Password", type: "password" },
        device_id: { label: "Device ID", type: "text" },
      },
      async authorize(credentials) {
        const res = await fetch(`${process.env.NEXT_BASE_URL}/api/auth/login`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            username: credentials?.username,
            password: credentials?.password,
            device_id: credentials?.device_id,
          }),
        });
        if (!res.ok) return null;

        const data = await res.json();
        const { access_token } = data;

        const decoded: JwtPayload = jwtDecode(access_token);

        return {
          id: decoded.user_id,
          username: decoded.username ?? credentials?.username,
          accessToken: access_token,
        };
      },
    }),
  ],
  session: { strategy: "jwt" },
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.accessToken = (user as any).accessToken;
        token.username = (user as any).username;
        token.id = (user as any).id;
      }
      return token;
    },
    async session({ session, token }) {
      session.user = { name: token.username as string };
      (session.user as any).id = token.id;
      (session as any).accessToken = token.accessToken;
      return session;
    },
  },
};
