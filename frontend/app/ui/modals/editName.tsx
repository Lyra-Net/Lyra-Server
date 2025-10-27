import { useState } from "react";

export default function EditNameModal({
  currentName,
  onSave,
  onClose,
}: {
  currentName: string;
  onSave: (name: string) => void;
  onClose: () => void;
}) {
  const [name, setName] = useState(currentName);

  return (
    <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
      <div className="bg-neutral-900 p-6 rounded-xl w-96">
        <h2 className="text-lg font-semibold mb-4 text-gray-200">Edit Display Name</h2>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="w-full p-2 rounded bg-neutral-800 border border-gray-700 text-gray-100 mb-4"
          placeholder="Enter new display name"
        />
        <div className="flex justify-end gap-2">
          <button onClick={onClose} className="px-4 py-2 text-sm text-gray-300 hover:text-white">
            Cancel
          </button>
          <button
            onClick={() => {
              onSave(name);
              onClose();
            }}
            className="px-4 py-2 text-sm bg-amber-600 hover:bg-amber-700 rounded"
          >
            Save
          </button>
        </div>
      </div>
    </div>
  );
}
