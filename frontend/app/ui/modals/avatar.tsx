import { useState } from "react";

export default function AvatarModal({
  currentAvatar,
  onSave,
  onClose,
}: {
  currentAvatar: string;
  onSave: (url: string) => void;
  onClose: () => void;
}) {
  const [preview, setPreview] = useState(currentAvatar);
  const [file, setFile] = useState<File | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0];
    if (f) {
      setFile(f);
      setPreview(URL.createObjectURL(f));
    }
  };

  const handleUpload = async () => {
    if (!file) return;
    // 🔧 TODO: gọi API upload avatar, ví dụ: uploadAvatar(file)
    // const url = await uploadAvatar(file);
    const url = preview; // demo
    onSave(url);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
      <div className="bg-neutral-900 p-6 rounded-xl w-96 flex flex-col items-center">
        <h2 className="text-lg font-semibold mb-4 text-gray-200">Change Avatar</h2>
        <div className="w-32 h-32 rounded-full overflow-hidden mb-4">
          <img src={preview} className="object-cover w-full h-full" alt="preview" />
        </div>
        <input
          type="file"
          accept="image/*"
          onChange={handleFileChange}
          className="text-sm text-gray-300 mb-4"
        />
        <div className="flex justify-end w-full gap-2">
          <button onClick={onClose} className="px-4 py-2 text-sm text-gray-300 hover:text-white">
            Cancel
          </button>
          <button
            onClick={handleUpload}
            disabled={!file}
            className="px-4 py-2 text-sm bg-amber-600 hover:bg-amber-700 rounded disabled:bg-gray-600"
          >
            Save
          </button>
        </div>
      </div>
    </div>
  );
}
