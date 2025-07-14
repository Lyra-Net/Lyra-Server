export const getDeviceId = (): string => {
  const key = 'device_id';
  let deviceId = localStorage.getItem(key);

  if (!deviceId) {
    deviceId = crypto.randomUUID(); // Hoặc dùng uuid()
    localStorage.setItem(key, deviceId);
  }

  return deviceId;
};
