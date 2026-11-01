import { ImageResponse } from "next/og";

export const alt = "Podlet — Your server. Your deployments.";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default function OpenGraphImage() {
  return new ImageResponse(
    <div style={{ position: "relative", width: "100%", height: "100%", display: "flex", flexDirection: "column", justifyContent: "space-between", overflow: "hidden", padding: "64px 72px", background: "linear-gradient(115deg, #aebfce 0%, #b9cee0 38%, #86bbfa 76%, #acbdef 100%)", color: "#172536", fontFamily: "Arial, sans-serif" }}>
      <div style={{ position: "absolute", top: -190, right: -110, width: 750, height: 750, display: "flex", borderRadius: 999, background: "radial-gradient(circle, rgba(185,241,255,.92) 0%, rgba(138,200,255,.4) 34%, rgba(138,200,255,0) 70%)" }} />
      <div style={{ display: "flex", alignItems: "center", gap: 22, fontSize: 58, fontWeight: 700, letterSpacing: -2.5 }}>
        <div style={{ width: 76, height: 76, display: "flex", alignItems: "center", justifyContent: "center", border: "1px solid #c1d5e8", borderRadius: 23, background: "linear-gradient(150deg, #d3e1ed, #b9cfe3)", boxShadow: "0 5px 8px rgba(36, 59, 91, .08), inset 0 -2px 2px #d9e3f3", transform: "rotate(-6deg)" }}>
          <svg width="46" height="46" viewBox="0 0 32 32" fill="none" style={{ transform: "rotate(6deg)" }}>
            <path d="M12 3h8a9 9 0 0 1 9 9v8a9 9 0 0 1-9 9h-8a9 9 0 0 1-9-9v-8a9 9 0 0 1 9-9Z" stroke="#172536" strokeWidth="2.5" />
            <path d="m9 11 5 5-5 5m9 0h5" stroke="#172536" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
        </div>
        <span>podlet</span>
      </div>
      <div style={{ display: "flex", flexDirection: "column", gap: 24 }}>
        <div style={{ maxWidth: 900, display: "flex", flexDirection: "column", fontSize: 90, fontWeight: 500, letterSpacing: -5, lineHeight: 1.02 }}>
          <span>Your server.</span>
          <span>Your deployments.</span>
        </div>
        <div style={{ fontSize: 32, color: "#35485f" }}>Deploy from your terminal. Keep control of your code.</div>
      </div>
    </div>,
    size,
  );
}
