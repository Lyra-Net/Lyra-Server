
export default function decodeJwt(token: string) {
  try {
    const payload = token.split(".")[1];
    return JSON.parse(Buffer.from(payload, "base64").toString("utf8"));
  } catch (e) {
    console.error("Failed to decode JWT", e);
    return null;
  }
}