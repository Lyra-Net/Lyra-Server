'use client'
import { useEffect, useState } from "react";
import { useAuthStore } from "@/stores/useAuth";
import DashboardLayout from "@/app/ui/dashboardLayout";
import EmailModal from "../ui/modals/email";
import OtpModal from "../ui/modals/otp";
import ConfirmModal from "../ui/modals/confirm";
import EditNameModal from "../ui/modals/editName";
import AvatarModal from "../ui/modals/avatar";
import { getProfile, toggle2Fa, verifyCode, updateProfile, addEmail } from "@/lib/authApi";
import { Toggle2FAPayload, VerifyCodePayload } from "@/declarations/auth";
import { getDeviceId } from "@/utils/device";
import { use2FASessionStore } from "@/stores/2fa_session";

export default function ProfilePage() {
  const { userId, username, displayName: storeName, email: storeEmail } = useAuthStore();

  // 📦 State
  const [displayName, setDisplayName] = useState(storeName || "");
  const [email, setEmail] = useState(storeEmail || "");
  const [avatar, setAvatar] = useState<string>("https://picsum.photos/seed/picsum/300/300");
  const [is2FA, setIs2FA] = useState(false);
  const [createdAt, setCreatedAt] = useState<number | null>(null);
  const [updatedAt, setUpdatedAt] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const {sessionId, setSessionId} = use2FASessionStore();

  const [showConfirm, setShowConfirm] = useState(false);
  const [modalField, setModalField] = useState<null | "email" | "display_name" | "avatar">(null);
  const [otpConfig, setOtpConfig] = useState<null | {
    action: "enable" | "disable" | "email" | "add" | "remove";
    email: string;
    expectedLength: number;
  }>(null);

  // 🧠 Fetch profile khi mount
  useEffect(() => {
    (async () => {
      try {
        const data = await getProfile();
        setDisplayName(data.display_name);
        setEmail(data.email);
        setIs2FA(data.is_2fa);
        setAvatar(data.avatar_url || "https://picsum.photos/seed/picsum/300/300");
        setCreatedAt(data.created_at);
        setUpdatedAt(data.updated_at);
      } catch (err) {
        console.error("Failed to fetch profile:", err);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  // 📅 Format timestamp
  const formatDate = (timestamp: number | null) =>
    timestamp ? new Date(timestamp * 1000).toLocaleString() : "-";

  // 🟡 Confirm bật/tắt 2FA
  const handleToggle2FA = () => setShowConfirm(true);

  // ✅ Xác nhận bật/tắt 2FA
  const confirm2FA = async () => {
    setShowConfirm(false);
    const action = is2FA ? "disable" : "enable";
    setOtpConfig({
        action,
        email: email || "",
        expectedLength: 8,
      });
    try {
      const payload: Toggle2FAPayload = {
        device_id: getDeviceId(),
        enable: !is2FA,
      };
      console.log("Payload for toggling 2FA: ", payload);
      const res = await toggle2Fa(payload);
      console.log("Toggle 2FA response: ", res);
      setSessionId(res.session_id);
    } catch (err) {
      console.error("Error toggling 2FA:", err);
      alert("Failed to initiate 2FA request");
    }
  };

  // 🔐 Xác thực OTP (bật/tắt 2FA hoặc verify email)
  const handleVerifyOtp = async (otp: string) => {
    if (!otpConfig) return;
    try {
      console.log("sessionId in handleVerifyOtp: ", sessionId);
      console.log("action: ", otpConfig.action);
      const payload: VerifyCodePayload = {
        session_id: sessionId || "",
        otp,
        device_id: getDeviceId(),
      };

      const data = await verifyCode(payload);
      console.log("Email verification response: ", data);

      if (otpConfig.action === "enable" || otpConfig.action === "disable") {
        if (data.is_success) {
          const payload: Toggle2FAPayload = {
            device_id: getDeviceId(),
            enable: !is2FA,
            verified_2fa_session_id: sessionId || "",
          };
          const res = await toggle2Fa(payload);
          console.log("Final toggle 2FA response: ", res);
          setSessionId("");
        }
        setIs2FA(otpConfig.action === "enable");
        alert(`2FA ${otpConfig.action === "enable" ? "enabled" : "disabled"} successfully`);
      }

      if (otpConfig.action === "add") {
        if (data.is_success) {
          const res = await addEmail({ email: otpConfig.email, device_id: getDeviceId(), verified_2fa_session_id: sessionId || "" });
          console.log("Add email response: ", res);
          setSessionId("");
          setEmail(otpConfig.email);
          alert("Email added and verified successfully");
        } 

      }
    } catch (err) {
      console.log("Error verifying OTP:", err);
    } finally {
      setOtpConfig(null);
    }
  };

  // 🧾 Update profile field (display_name or avatar)
  const handleProfileUpdate = async (field: "display_name" | "avatar", value: string) => {
    try {
      await updateProfile({ [field]: value });
      if (field === "display_name") setDisplayName(value);
      if (field === "avatar") setAvatar(value);
      alert(`${field === "display_name" ? "Display name" : "Avatar"} updated successfully`);
    } catch (err) {
      console.error("Failed to update profile:", err);
      alert("Update failed");
    }
  };

  if (loading) return <div className="text-center py-10 text-gray-400">Loading...</div>;

  return (
    <DashboardLayout>
      <div className="mt-8 px-6 max-w-2xl mx-auto text-gray-200">
        {/* Avatar */}
        <div className="flex flex-col items-center mb-10 relative">
          <div
            className="rounded-full overflow-hidden w-32 h-32 cursor-pointer border border-gray-700"
            onClick={() => setModalField("avatar")}
          >
            <img src={avatar} alt="profile" className="object-cover w-full h-full" />
          </div>
          <button
            className="mt-2 text-sm text-amber-500 hover:underline"
            onClick={() => setModalField("avatar")}
          >
            Change avatar
          </button>
        </div>

        <h2 className="text-lg font-semibold mb-3">Account Info</h2>

        <div className="space-y-4">
          <InfoRow label="User ID" value={userId} />
          <InfoRow label="Username" value={username} />
          <InfoRow label="Display Name" value={displayName} editable onEdit={() => setModalField("display_name")} />
          <InfoRow label="Email" value={email || "Not set"} editable onEdit={() => setModalField("email")} />
          <InfoRow label="Created At" value={formatDate(createdAt)} />
          <InfoRow label="Updated At" value={formatDate(updatedAt)} />
        </div>

        <div className="mt-10 border border-gray-600 rounded-md p-4 bg-neutral-900">
          <div className="flex justify-between items-center">
            <h2 className="text-lg font-semibold">Two Factor Authentication</h2>
            <button
              disabled={!email || email === ""}
              onClick={handleToggle2FA}
              className={`px-4 py-2 rounded-md font-medium ${
                is2FA ? "bg-red-600 hover:bg-red-700" : "bg-amber-600 hover:bg-amber-700"
              } disabled:bg-gray-600`}
            >
              {is2FA ? "Disable 2FA" : "Enable 2FA"}
            </button>
          </div>
          {!email && <p className="text-xs mt-2 text-gray-400">⚠️ You must set an email before enabling 2FA.</p>}
        </div>
      </div>

      {/* Modals */}
      {showConfirm && (
        <ConfirmModal
          title={is2FA ? "Disable 2FA" : "Enable 2FA"}
          message={`Are you sure you want to ${is2FA ? "disable" : "enable"} two-factor authentication?`}
          onCancel={() => setShowConfirm(false)}
          onConfirm={confirm2FA}
        />
      )}

      {modalField === "email" && (
        <EmailModal
          currentEmail={email}
          onOtpRequest={(config: any) => setOtpConfig({ action: "email", ...config })}
          onClose={() => setModalField(null)}
        />
      )}

      {modalField === "display_name" && (
        <EditNameModal
          currentName={displayName}
          onSave={(name) => handleProfileUpdate("display_name", name)}
          onClose={() => setModalField(null)}
        />
      )}

      {modalField === "avatar" && (
        <AvatarModal
          currentAvatar={avatar}
          onSave={(url) => handleProfileUpdate("avatar", url)}
          onClose={() => setModalField(null)}
        />
      )}

      {otpConfig && (
        <OtpModal
          actionLabel={
            otpConfig.action === "enable"
              ? "Verify to Enable 2FA"
              : otpConfig.action === "disable"
              ? "Verify to Disable 2FA"
              : "Verify Email"
          }
          expectedLength={otpConfig.expectedLength}
          onSubmit={handleVerifyOtp}
          onClose={() => setOtpConfig(null)}
        />
      )}
    </DashboardLayout>
  );
}

function InfoRow({ label, value, editable = false, onEdit }: {
  label: string;
  value: string | null;
  editable?: boolean;
  onEdit?: () => void;
}) {
  return (
    <div
      className={`flex justify-between border-b border-gray-700 py-2 text-sm ${editable ? "cursor-pointer hover:bg-neutral-800" : ""}`}
      onClick={editable ? onEdit : undefined}
    >
      <span className="text-gray-400">{label}</span>
      <div className="flex items-center gap-2">
        <span className="text-gray-200">{value || "-"}</span>
        {editable && <span className="text-xs text-amber-500">Edit</span>}
      </div>
    </div>
  );
}
