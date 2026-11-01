import { ImageResponse } from "next/og";

export const size = { width: 180, height: 180 };
export const contentType = "image/png";

export default function AppleIcon() {
  return new ImageResponse(
    <div style={{ width: "100%", height: "100%", display: "flex", alignItems: "center", justifyContent: "center", background: "transparent" }}>
      <div style={{ width: 156, height: 156, display: "flex", alignItems: "center", justifyContent: "center", border: "3px solid #c1d5e8", borderRadius: 47, background: "linear-gradient(150deg, #d3e1ed, #b9cfe3)", transform: "rotate(-6deg)" }}>
        <svg width="105" height="105" viewBox="0 0 32 32" fill="none" style={{ transform: "rotate(6deg)" }}>
          <path d="M12 3h8a9 9 0 0 1 9 9v8a9 9 0 0 1-9 9h-8a9 9 0 0 1-9-9v-8a9 9 0 0 1 9-9Z" stroke="#172536" strokeWidth="2.5" />
          <path d="m9 11 5 5-5 5m9 0h5" stroke="#172536" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </div>
    </div>,
    size,
  );
}
