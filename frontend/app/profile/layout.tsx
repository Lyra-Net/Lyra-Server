import { Metadata } from "next";

export const metadata: Metadata = {
  title: 'Profile | BypassBeats',
  description: 'Edit Profile',
};

export default function ProfileLayout({ children }: { children: React.ReactNode }) {
  return (
    <>{children}</>
  );
}   