import { ImageResponse } from "next/og";

export const alt = "Podlet — Your server. Your deployments.";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default function OpenGraphImage() {
  return new ImageResponse(
    <div style={{ position: "relative", width: "100%", height: "100%", display: "flex", flexDirection: "column", justifyContent: "space-between", overflow: "hidden", padding: "64px 72px", background: "linear-gradient(115deg, #aebfce 0%, #b9cee0 38%, #86bbfa 76%, #acbdef 100%)", color: "#172536", fontFamily: "Arial, sans-serif" }}>
      <div style={{ position: "absolute", top: -190, right: -110, width: 750, height: 750, display: "flex", borderRadius: 999, background: "radial-gradient(circle, rgba(185,241,255,.92) 0%, rgba(138,200,255,.4) 34%, rgba(138,200,255,0) 70%)" }} />
      <div style={{ display: "flex", alignItems: "center", gap: 18, fontSize: 42, fontWeight: 700, letterSpacing: -2 }}>
        <div style={{ width: 62, height: 62, display: "flex", alignItems: "center", justifyContent: "center", borderRadius: 16, background: "#25374c", color: "#dceaff", fontFamily: "monospace", fontSize: 30, fontWeight: 700, letterSpacing: -4 }}>&gt;_</div>
        <span>podlet</span>
      </div>
      <div style={{ display: "flex", flexDirection: "column", gap: 24 }}>
        <div style={{ maxWidth: 900, display: "flex", flexDirection: "column", fontSize: 84, fontWeight: 500, letterSpacing: -5, lineHeight: 1.02 }}>
          <span>Your server.</span>
          <span>Your deployments.</span>
        </div>
        <div style={{ fontSize: 28, color: "#35485f" }}>Deploy from your terminal. Keep control of your code.</div>
      </div>
    </div>,
    size,
  );
}
