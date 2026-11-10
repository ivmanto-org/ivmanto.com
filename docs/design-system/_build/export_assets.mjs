// Raster exports of the logo SVGs: PNG (transparent), JPG (solid background), favicon.ico, OG image.
// Run after build.py:  node _build/export_assets.mjs
// Needs puppeteer-core and a Chromium build:
//   PUPPETEER_CORE=/path/to/node_modules/puppeteer-core CHROME=/path/to/chrome node _build/export_assets.mjs
import { readFileSync, writeFileSync, mkdirSync, readdirSync } from "node:fs";
import { join, dirname } from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";

const DS = join(dirname(fileURLToPath(import.meta.url)), "..");
const A = join(DS, "assets");
const { default: puppeteer } = await import(pathToFileURL(join(process.env.PUPPETEER_CORE ?? "puppeteer-core", "lib/puppeteer/puppeteer-core.js")).href);
const BG = { "on-light": "#ffffff", "on-dark": "#0a1220", "mono-black": "#ffffff", "mono-white": "#0a1220" };

const browser = await puppeteer.launch({ executablePath: process.env.CHROME, args: ["--no-sandbox", "--disable-gpu"] });
const page = await browser.newPage();

// Render an SVG file into a w x h box; the logo is scaled to fit inside (w, h) minus pad on every side.
async function render(svgPath, { w, h, pad = 0, bg = null, type = "png", out }) {
  const svg = readFileSync(svgPath, "utf8").replace(/ width="[\d.]+" height="[\d.]+"/, "");
  const html = `<!doctype html><html><body style="margin:0;width:${w}px;height:${h}px;display:grid;place-items:center;background:${bg ?? "transparent"}">
    <div style="width:${w - 2 * pad}px;height:${h - 2 * pad}px;display:grid;place-items:center">${svg.replace("<svg ", '<svg style="width:100%;height:100%" ')}</div></body></html>`;
  await page.setViewport({ width: w, height: h, deviceScaleFactor: 1 });
  await page.setContent(html);
  const buf = await page.screenshot({ type, omitBackground: !bg, quality: type === "jpeg" ? 92 : undefined, clip: { x: 0, y: 0, width: w, height: h } });
  if (out) writeFileSync(out, buf);
  return buf;
}

const ratio = (p) => { const [, , vw, vh] = readFileSync(p, "utf8").match(/viewBox="([^"]+)"/)[1].split(" ").map(Number); return vw / vh; };

// Logo parts, every colourway
const parts = { wordmark: { png: [600, 1200, 2400], jpg: 2400, wide: true }, lockup: { png: [1200, 2400], jpg: 2400, wide: true },
                symbol: { png: [256, 512, 1024], jpg: 1024 }, short: { png: [512, 1024], jpg: 1024 } };
for (const [part, spec] of Object.entries(parts)) {
  const dir = join(A, "logo", part);
  for (const f of readdirSync(dir).filter((n) => n.endsWith(".svg"))) {
    const src = join(dir, f), base = f.replace(".svg", ""), colour = base.split(`${part}-`)[1], r = ratio(src);
    for (const s of spec.png) {
      const [w, h] = spec.wide ? [s, Math.round(s / r)] : [s, s];
      await render(src, { w, h, out: join(dir, `${base}-${s}.png`) });
    }
    const s = spec.jpg, [w, h] = spec.wide ? [s, Math.round(s / r)] : [s, s], pad = Math.round((spec.wide ? h : w) * 0.12);
    await render(src, { w: w + 2 * pad, h: h + 2 * pad, pad, bg: BG[colour], type: "jpeg", out: join(dir, `${base}-${s}.jpg`) });
  }
}

// Favicon set from the short form ("iv" + tick)
const F = join(A, "favicon"), short = join(A, "logo/short/ivmanto-short-on-light.svg");
const ico = [];
for (const s of [16, 32, 48]) ico.push([s, await render(short, { w: s, h: s, pad: s >= 32 ? 1 : 0, out: join(F, `favicon-${s}.png`) })]);
await render(short, { w: 180, h: 180, pad: 22, bg: "#ffffff", out: join(F, "apple-touch-icon.png") });
for (const s of [192, 512]) await render(short, { w: s, h: s, pad: Math.round(s * 0.12), bg: "#ffffff", out: join(F, `icon-${s}.png`) });
await render(short, { w: 512, h: 512, pad: 112, bg: "#ffffff", out: join(F, "icon-maskable-512.png") }); // content inside the 80% safe zone
await render(short, { w: 512, h: 512, pad: 60, bg: "#ffffff", type: "jpeg", out: join(F, "favicon-512.jpg") });

// favicon.ico: PNG-compressed entries (supported by every current browser)
const head = Buffer.alloc(6 + 16 * ico.length);
head.writeUInt16LE(0, 0); head.writeUInt16LE(1, 2); head.writeUInt16LE(ico.length, 4);
let offset = head.length;
ico.forEach(([s, png], i) => {
  const o = 6 + 16 * i;
  head.writeUInt8(s, o); head.writeUInt8(s, o + 1); head.writeUInt16LE(1, o + 4); head.writeUInt16LE(32, o + 6);
  head.writeUInt32LE(png.length, o + 8); head.writeUInt32LE(offset, o + 12); offset += png.length;
});
writeFileSync(join(F, "favicon.ico"), Buffer.concat([head, ...ico.map(([, b]) => b)]));

// Link preview: lockup on navy, 1200 x 630
mkdirSync(join(A, "social"), { recursive: true });
const lock = join(A, "logo/lockup/ivmanto-lockup-on-dark.svg");
await render(lock, { w: 1200, h: 630, pad: 230, bg: "#0a1220", out: join(A, "social/og-image.png") });
await render(lock, { w: 1200, h: 630, pad: 230, bg: "#0a1220", type: "jpeg", out: join(A, "social/og-image.jpg") });

await browser.close();
console.log("exported rasters into", A);
