import { addEmail } from "@/lib/authApi";
import { use2FASessionStore } from "@/stores/2fa_session";
import { getDeviceId } from "@/utils/device";
import validateEmail from "@/utils/validate_email";
import { useState } from "react";

export default function EmailModal({ currentEmail, onOtpRequest, onClose }: any) {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const setSessionId = use2FASessionStore(state => state.setSessionId);

  const handleSendOtp = async (actionType: "add" | "remove") => {
    setLoading(true);
    console.log("Mock send OTP for", actionType, { email });
    if (actionType === "add" && !validateEmail(email)) {
      setLoading(false);
      return;
    }
    await addEmail({ email, device_id: getDeviceId() }).then(data => {
      console.log("Add email response:", data);
      if (data.is_success && data.session_id) {
        setSessionId(data.session_id);
      }
    }).catch(err => {
      console.error("Failed to add email:", err);
    }).finally(() => {
      setLoading(false);
    });
    
    onOtpRequest({
      action: actionType,
      email,
    });
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/70 z-50">
      <div className="bg-neutral-900 border border-gray-700 rounded-xl p-6 w-[420px] shadow-lg relative">
        {loading && (
          <div className="absolute inset-0 flex items-center justify-center bg-black/40 rounded-xl">
            <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-amber-500"></div>
          </div>
        )}

        <h2 className="text-xl font-semibold mb-4 text-gray-100">
          {currentEmail ? "Remove Your Email" : "Add Email Address"}
        </h2>

        {!currentEmail && (
          <input
            type="email"
            placeholder="Enter your email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full bg-neutral-800 border border-gray-600 rounded-md px-3 py-2 text-gray-300 focus:outline-none focus:ring-2 focus:ring-amber-500"
          />
        )}

        <div className="flex justify-end space-x-3 mt-6">
          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-md"
          >
            Cancel
          </button>

          {currentEmail ? (
            <button
              onClick={() => handleSendOtp("remove")}
              className="px-4 py-2 bg-red-600 hover:bg-red-700 rounded-md"
            >
              Remove Email
            </button>
          ) : (
            <button
              onClick={() => handleSendOtp("add")}
              disabled={!validateEmail(email) || loading}
              className={`px-4 py-2 rounded-md ${
                (validateEmail(email) && !loading)
                  ? "bg-amber-600 hover:bg-amber-700"
                  : "bg-gray-700 cursor-not-allowed"
              }`}
            >
              Send OTP
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
