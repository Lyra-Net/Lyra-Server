import { useEffect, useRef, useState } from "react";

export default function OtpModal({
  actionLabel = "Verify Action",
  onSubmit,
  onClose,
  expectedLength = 6,
}: {
  actionLabel?: string;
  onSubmit: (otp: string) => Promise<void> | void;
  onClose: () => void;
  expectedLength?: number;
}) {
  const [otp, setOtp] = useState("");
  const [loading, setLoading] = useState(false);
  const otpContainerRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    otpContainerRef.current?.focus();
  }, []);

  const handleKeyDown = async (e: React.KeyboardEvent<HTMLDivElement>) => {
    if (loading) return;

    if (e.key === "Backspace") {
      e.preventDefault();
      setOtp((prev) => prev.slice(0, -1));
    } else if (/^[0-9]$/.test(e.key)) {
      e.preventDefault();
      setOtp((prev) => {
        if (prev.length >= expectedLength) return prev;
        return prev + e.key;
      });
    } else if (e.key === "Enter") {
      e.preventDefault();
      if (otp.length >= expectedLength) {
        await handleSubmit();
      }
    }
  };

  const handleSubmit = async () => {
    if (loading) return;
    setLoading(true);
    try {
      await onSubmit(otp);
    } finally {
      setLoading(false);
    }
  };

  const renderOtpBoxes = () => {
    const boxes = [];
    const chars = otp.split("");
    for (let i = 0; i < expectedLength; i++) {
      boxes.push(
        <div
          key={i}
          className={`w-10 h-12 flex items-center justify-center rounded-md text-xl font-semibold transition-colors
            ${i === chars.length
              ? "border-2 border-amber-500"
              : i < chars.length
              ? "bg-amber-600 text-white"
              : "bg-neutral-800 border border-gray-600 text-gray-400"
            }`}
        >
          {chars[i] || ""}
        </div>
      );
    }
    return (
      <div
        ref={otpContainerRef}
        tabIndex={0}
        onKeyDown={handleKeyDown}
        className="flex justify-center gap-2 mb-4 outline-none select-none cursor-text"
      >
        {boxes}
      </div>
    );
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/70 z-50">
      <div
        className="bg-neutral-900 border border-gray-700 rounded-xl p-6 w-[420px] shadow-lg relative"
        onClick={(e) => {
          const target = e.target as HTMLElement;
          if (target.tagName !== "BUTTON") {
            otpContainerRef.current?.focus();
          }
        }}
      >
        {loading && (
          <div className="absolute inset-0 flex items-center justify-center bg-black/40 rounded-xl">
            <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-amber-500"></div>
          </div>
        )}

        <h2 className="text-xl font-semibold mb-4 text-gray-100">
          {actionLabel}
        </h2>

        {renderOtpBoxes()}

        <div className="flex justify-end space-x-3 mt-6">
          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-md"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            disabled={otp.length < expectedLength}
            className={`px-4 py-2 rounded-md ${
              otp.length >= expectedLength
                ? "bg-amber-600 hover:bg-amber-700"
                : "bg-gray-700 cursor-not-allowed"
            }`}
          >
            Confirm
          </button>
        </div>
      </div>
    </div>
  );
}
