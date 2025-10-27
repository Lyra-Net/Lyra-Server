"use client";

export default function ConfirmModal({ title, message, onCancel, onConfirm }: any) {
  return (
    <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50">
      <div className="bg-neutral-900 border border-gray-700 p-6 rounded-lg w-[380px]">
        <h2 className="text-lg font-semibold mb-3">{title}</h2>
        <p className="text-gray-400 text-sm mb-6">{message}</p>
        <div className="flex justify-end space-x-3">
          <button onClick={onCancel} className="px-4 py-2 bg-gray-700 rounded-md hover:bg-gray-600">
            Cancel
          </button>
          <button onClick={onConfirm} className="px-4 py-2 bg-amber-600 rounded-md hover:bg-amber-700">
            Confirm
          </button>
        </div>
      </div>
    </div>
  );
}
