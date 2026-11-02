import fs from 'fs';
import path from 'path';
import { createRequire } from 'module';

const nodeRuntime = 'C:/Users/12495/.cache/codex-runtimes/codex-primary-runtime/dependencies/node';
const require = createRequire(path.join(nodeRuntime, 'index.js'));
const sharp = require('sharp');

const seedDir = 'backend/static/seed';
const localDir = 'frontend/assets/images';

const puppet = (color, accent, opacity = 1) => `
  <g opacity="${opacity}">
    <path d="M316 150c-42 0-74 30-78 70-3 27 6 48 6 72 0 34-20 54-48 70l70 26c30 14 62 16 92 6l64-22c-30-14-48-38-48-72 0-26 10-46 8-74-4-42-40-76-80-76z" fill="${color}"/>
    <path d="M236 208c22-40 70-58 108-42 8 4 14 8 20 16-30-14-70-12-96 8-10 8-18 18-24 30z" fill="${accent}" opacity=".92"/>
    <path d="M252 190c-14-34 4-70 40-84 26-10 54-4 74 14-30-8-62 0-82 20-14 14-20 32-18 50z" fill="${color}"/>
    <path d="M300 118c-10-26-2-52 20-64 4 24 16 40 36 48-14 16-34 24-56 16z" fill="${accent}"/>
    <path d="M286 236c16-6 34-4 46 6" stroke="${accent}" stroke-width="9" fill="none" stroke-linecap="round"/>
    <circle cx="300" cy="238" r="7" fill="${accent}"/>
    <path d="M262 268c26 12 54 12 78 0" stroke="${accent}" stroke-width="7" fill="none" stroke-linecap="round" opacity=".7"/>
    <path d="M236 300l-96 62" stroke="${color}" stroke-width="16" stroke-linecap="round"/>
    <path d="M372 306l96 40" stroke="${color}" stroke-width="16" stroke-linecap="round"/>
    <circle cx="140" cy="362" r="18" fill="${accent}"/>
    <circle cx="468" cy="346" r="16" fill="${accent}"/>
  </g>`;

const lantern = (x, y, r, glow) => `
  <g>
    <circle cx="${x}" cy="${y}" r="${r * 3.4}" fill="${glow}" opacity=".22"/>
    <circle cx="${x}" cy="${y}" r="${r * 1.9}" fill="${glow}" opacity=".26"/>
    <ellipse cx="${x}" cy="${y}" rx="${r}" ry="${r * 1.18}" fill="#D8452F"/>
    <ellipse cx="${x}" cy="${y}" rx="${r * 0.38}" ry="${r * 1.18}" fill="#F0A24A" opacity=".55"/>
    <rect x="${x - r * 0.42}" y="${y - r * 1.4}" width="${r * 0.84}" height="${r * 0.22}" fill="#8E2C1E"/>
    <rect x="${x - r * 0.42}" y="${y + r * 1.18}" width="${r * 0.84}" height="${r * 0.2}" fill="#8E2C1E"/>
    <path d="M${x} ${y + r * 1.4}v${r * 0.9}" stroke="#C39A4B" stroke-width="3"/>
  </g>`;

const dust = (count, w, h, seedStart = 1) => {
  const parts = [];
  let seed = seedStart;
  const rand = () => {
    seed = (seed * 9301 + 49297) % 233280;
    return seed / 233280;
  };
  for (let index = 0; index < count; index += 1) {
    const x = Math.round(rand() * w);
    const y = Math.round(rand() * h);
    const r = (rand() * 2.6 + 0.8).toFixed(1);
    const o = (rand() * 0.5 + 0.18).toFixed(2);
    parts.push(`<circle cx="${x}" cy="${y}" r="${r}" fill="#FFE3A8" opacity="${o}"/>`);
  }
  return `<g>${parts.join('')}</g>`;
};

const curtains = (w, h) => `
  <g>
    <path d="M0 0h150c-30 120-46 260-40 400 4 100 18 200 40 300H0z" fill="#5B1D18"/>
    <path d="M0 0h96c-20 140-26 300-14 460 6 82 14 162 26 240H0z" fill="#74251D" opacity=".85"/>
    <path d="M${w} 0h-150c30 120 46 260 40 400-4 100-18 200-40 300h150z" fill="#5B1D18"/>
    <path d="M${w} 0h-96c20 140 26 300 14 460-6 82-14 162-26 240h108z" fill="#74251D" opacity=".85"/>
  </g>`;

const stageScene = () => `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 700" width="1200" height="700">
  <defs>
    <linearGradient id="bg" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#31180F"/>
      <stop offset="0.55" stop-color="#1D100B"/>
      <stop offset="1" stop-color="#0F0806"/>
    </linearGradient>
    <radialGradient id="spot" cx="0.5" cy="0.44" r="0.62">
      <stop offset="0" stop-color="#FFC46B" stop-opacity=".46"/>
      <stop offset="0.55" stop-color="#E08A2E" stop-opacity=".18"/>
      <stop offset="1" stop-color="#E08A2E" stop-opacity="0"/>
    </radialGradient>
    <linearGradient id="screen" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#F7E7C0"/>
      <stop offset="1" stop-color="#E4C88E"/>
    </linearGradient>
    <radialGradient id="screenGlow" cx="0.5" cy="0.45" r="0.6">
      <stop offset="0" stop-color="#FFF3D4" stop-opacity=".95"/>
      <stop offset="1" stop-color="#D8B478" stop-opacity=".2"/>
    </radialGradient>
  </defs>
  <rect width="1200" height="700" fill="url(#bg)"/>
  <rect width="1200" height="700" fill="url(#spot)"/>
  <rect x="228" y="72" width="744" height="556" rx="14" fill="url(#screen)"/>
  <rect x="228" y="72" width="744" height="556" rx="14" fill="url(#screenGlow)"/>
  ${puppet('#26150F', '#8E2C1E', 0.94).replace('<g opacity="0.94">', '<g transform="translate(268 96) scale(1.16)" opacity="0.94">')}
  <rect x="228" y="72" width="744" height="556" rx="14" fill="none" stroke="#B8893F" stroke-width="6"/>
  <rect x="248" y="92" width="704" height="516" rx="6" fill="none" stroke="#8E2C1E" stroke-width="3" opacity=".55"/>
  <path d="M228 628h744" stroke="#3A1E14" stroke-width="26"/>
  <path d="M300 92c60 30 90 74 96 132" stroke="#FFF3D4" stroke-width="10" opacity=".16" fill="none"/>
  ${lantern(118, 132, 34, '#E8A33D')}
  ${lantern(1082, 132, 34, '#E8A33D')}
  ${curtains(1200, 700)}
  ${dust(90, 1200, 700, 7)}
</svg>`;

const riverScene = () => {
  const reflections = [];
  for (let index = 0; index < 26; index += 1) {
    const y = 508 + index * 7;
    const w = 60 + ((index * 37) % 190);
    const x = 430 + ((index * 83) % 420) - w / 2;
    reflections.push(`<rect x="${x.toFixed(0)}" y="${y}" width="${w}" height="3" rx="1.5" fill="#F2C177" opacity="${(0.34 - index * 0.01).toFixed(2)}"/>`);
  }
  return `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 700" width="1200" height="700">
  <defs>
    <linearGradient id="sky" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#141C2E"/>
      <stop offset="0.62" stop-color="#2A2131"/>
      <stop offset="1" stop-color="#3A2622"/>
    </linearGradient>
    <linearGradient id="water" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#2C2A33"/>
      <stop offset="1" stop-color="#161A26"/>
    </linearGradient>
    <radialGradient id="moonGlow" cx="0.5" cy="0.5" r="0.5">
      <stop offset="0" stop-color="#FFF0C8" stop-opacity=".85"/>
      <stop offset="1" stop-color="#FFF0C8" stop-opacity="0"/>
    </radialGradient>
    <linearGradient id="pavilion" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#F6E2B4"/>
      <stop offset="1" stop-color="#D9AE6B"/>
    </linearGradient>
  </defs>
  <rect width="1200" height="700" fill="url(#sky)"/>
  <circle cx="930" cy="132" r="150" fill="url(#moonGlow)"/>
  <circle cx="930" cy="132" r="52" fill="#FDF3D8" opacity=".92"/>
  <path d="M0 470h1200v230H0z" fill="url(#water)"/>
  <g fill="#20242F">
    <path d="M60 470V392h54v78zM124 470V362h40v108zM174 470V410h66v60zM250 470V378h48v92zM308 470V404h58v66z"/>
    <path d="M846 470V396h52v74zM908 470V366h44v104zM962 470V406h60v64zM1032 470V384h54v86zM1096 470V416h56v54z"/>
  </g>
  <g>
    <rect x="470" y="300" width="260" height="170" fill="url(#pavilion)"/>
    <rect x="492" y="322" width="216" height="126" fill="#FFF0CC" opacity=".92"/>
    <g transform="translate(520 336) scale(0.62)" opacity=".92">
      ${puppet('#241712', '#8E2C1E', 0.9)}
    </g>
    <path d="M452 300h296l-34-42H486z" fill="#7E2419"/>
    <path d="M470 470V300M730 470V300" stroke="#3B2119" stroke-width="14"/>
    <path d="M486 258c40-22 188-22 228 0" stroke="#C39A4B" stroke-width="6" fill="none"/>
  </g>
  ${reflections.join('')}
  <path d="M0 470h1200" stroke="#4A3A33" stroke-width="4"/>
  <g>
    <path d="M0 118c300 34 600 34 900 0s300-22 300-22" stroke="#4C3A2C" stroke-width="4" fill="none"/>
    ${lantern(200, 152, 24, '#F0B85C')}
    ${lantern(400, 170, 24, '#F0B85C')}
    ${lantern(600, 176, 24, '#F0B85C')}
    ${lantern(800, 170, 24, '#F0B85C')}
    ${lantern(1000, 152, 24, '#F0B85C')}
  </g>
  ${dust(60, 1200, 480, 23)}
</svg>`;
};

const workshopScene = () => `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 700" width="1200" height="700">
  <defs>
    <linearGradient id="wall" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#F7E9C8"/>
      <stop offset="1" stop-color="#E2C795"/>
    </linearGradient>
    <linearGradient id="table" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#B98A52"/>
      <stop offset="1" stop-color="#8E6337"/>
    </linearGradient>
    <radialGradient id="lamp" cx="0.5" cy="0.1" r="0.72">
      <stop offset="0" stop-color="#FFD98A" stop-opacity=".62"/>
      <stop offset="1" stop-color="#FFD98A" stop-opacity="0"/>
    </radialGradient>
  </defs>
  <rect width="1200" height="700" fill="url(#wall)"/>
  <rect width="1200" height="700" fill="url(#lamp)"/>
  <rect y="470" width="1200" height="230" fill="url(#table)"/>
  <path d="M0 470h1200" stroke="#6F4B27" stroke-width="8"/>
  <g>
    <rect x="150" y="316" width="330" height="176" rx="12" fill="#7C5330"/>
    <rect x="166" y="332" width="298" height="144" rx="8" fill="#5E3C21"/>
    <rect x="180" y="346" width="126" height="116" rx="6" fill="#43301C"/>
    <rect x="320" y="346" width="130" height="116" rx="6" fill="#43301C"/>
    <rect x="196" y="362" width="94" height="84" rx="4" fill="#C9A46A" opacity=".55"/>
    <path d="M330 360h110M330 386h110M330 412h110" stroke="#C39A4B" stroke-width="5" opacity=".6"/>
    <path d="M186 316c40-30 258-30 300 0" stroke="#A97C42" stroke-width="10" fill="none"/>
  </g>
  <g>
    <path d="M232 300c0-28 20-44 44-44s44 16 44 44z" fill="#5B2F1C" opacity=".9"/>
    <circle cx="234" cy="272" r="12" fill="#C23A2A"/>
    <path d="M596 300c0-28 20-44 44-44s44 16 44 44z" fill="#5B2F1C" opacity=".9"/>
    <circle cx="598" cy="272" r="12" fill="#C23A2A"/>
  </g>
  <g transform="translate(520 300) rotate(-8) scale(1.02)">
    ${puppet('#8E2C1E', '#C39A4B', 0.95)}
  </g>
  <g fill="#E8D5AE" stroke="#A9843F" stroke-width="4">
    <path d="M320 372l26-96h96l30 96z"/>
    <path d="M448 276h34v96h-34z"/>
  </g>
  <path d="M470 276c30-8 52-2 66 18" stroke="#A9843F" stroke-width="5" fill="none"/>
  <g stroke="#6E5A43" stroke-width="6" fill="#F6EEDC">
    <path d="M120 400l84 26-16 34-84-26z"/>
    <circle cx="122" cy="412" r="14" fill="none"/>
    <circle cx="150" cy="404" r="6" fill="#6E5A43" stroke="none"/>
  </g>
  <g>
    <path d="M880 330l150 0 0 108-150 0z" fill="#3B2A1E"/>
    <path d="M880 330l150 0-24 20-126 0z" fill="#54402C"/>
    <circle cx="955" cy="392" r="26" fill="#1F1712"/>
    <circle cx="955" cy="392" r="15" fill="#4A3A2A"/>
  </g>
  <g stroke="#7E2419" stroke-width="10" fill="none" stroke-linecap="round">
    <path d="M1030 300l40 118"/>
    <path d="M1052 300l-24 118"/>
  </g>
  <g>
    <rect x="768" y="250" width="72" height="72" rx="6" fill="#C0392B"/>
    <rect x="784" y="266" width="40" height="40" rx="4" fill="none" stroke="#F7E9C8" stroke-width="5"/>
    <path d="M804 250v-24" stroke="#C39A4B" stroke-width="6"/>
  </g>
  ${dust(50, 1200, 460, 41)}
</svg>`;

const puppetCover = () => `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 1000" width="800" height="1000">
  <defs>
    <linearGradient id="ab" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#8E2C1E"/>
      <stop offset="1" stop-color="#40190F"/>
    </linearGradient>
    <radialGradient id="ag" cx="0.5" cy="0.42" r="0.55">
      <stop offset="0" stop-color="#FFD98A" stop-opacity=".5"/>
      <stop offset="1" stop-color="#FFD98A" stop-opacity="0"/>
    </radialGradient>
  </defs>
  <rect width="800" height="1000" fill="url(#ab)"/>
  <rect width="800" height="1000" fill="url(#ag)"/>
  <g transform="translate(140 190) scale(1.36)">
    ${puppet('#F3DCB0', '#C0392B', 0.98)}
  </g>
  <g stroke="#F3DCB0" stroke-width="6" fill="none" opacity=".5">
    <rect x="40" y="40" width="720" height="920" rx="14"/>
  </g>
  ${dust(70, 800, 1000, 61)}
</svg>`;

const townCover = () => `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 1000" width="800" height="1000">
  <defs>
    <linearGradient id="tb" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#22303F"/>
      <stop offset="1" stop-color="#3B2830"/>
    </linearGradient>
    <radialGradient id="tg" cx="0.68" cy="0.22" r="0.5">
      <stop offset="0" stop-color="#FFF0C8" stop-opacity=".7"/>
      <stop offset="1" stop-color="#FFF0C8" stop-opacity="0"/>
    </radialGradient>
  </defs>
  <rect width="800" height="1000" fill="url(#tb)"/>
  <rect width="800" height="1000" fill="url(#tg)"/>
  <circle cx="560" cy="238" r="60" fill="#FDF3D8" opacity=".9"/>
  <g fill="#1B1F2A">
    <path d="M0 640h800v360H0z"/>
    <path d="M40 640V520h90v120zM150 640V470h70v170zM240 640V540h110v100zM370 640V500h80v140zM470 640V548h120v92zM610 640V478h90v162zM720 640V540h80v100z"/>
    <path d="M180 400h60l20 70h-100zM600 386h60l20 70h-100z"/>
  </g>
  <g>
    <rect x="260" y="560" width="280" height="200" fill="#F6E2B4"/>
    <g transform="translate(300 574) scale(0.66)" opacity=".92">
      ${puppet('#241712', '#8E2C1E', 0.9)}
    </g>
    <path d="M240 560h320l-40-46H280z" fill="#7E2419"/>
  </g>
  ${lantern(160, 300, 26, '#F0B85C')}
  ${lantern(660, 320, 26, '#F0B85C')}
  ${dust(60, 800, 1000, 83)}
</svg>`;

const classCover = () => `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 1000" width="800" height="1000">
  <defs>
    <linearGradient id="cb" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#F6E7C6"/>
      <stop offset="1" stop-color="#DCBE8B"/>
    </linearGradient>
    <radialGradient id="cg" cx="0.5" cy="0.34" r="0.6">
      <stop offset="0" stop-color="#FFF6DC" stop-opacity=".9"/>
      <stop offset="1" stop-color="#FFF6DC" stop-opacity="0"/>
    </radialGradient>
  </defs>
  <rect width="800" height="1000" fill="url(#cb)"/>
  <rect width="800" height="1000" fill="url(#cg)"/>
  <rect x="120" y="150" width="560" height="420" rx="10" fill="#FBF2DA"/>
  <rect x="120" y="150" width="560" height="420" rx="10" fill="none" stroke="#B8893F" stroke-width="6"/>
  <g transform="translate(210 190) scale(0.86)" opacity=".92">
    ${puppet('#26150F', '#8E2C1E', 0.94)}
  </g>
  <path d="M120 570h560" stroke="#6F4B27" stroke-width="12"/>
  <g fill="#5B3A24">
    <path d="M90 640c0-42 26-68 62-68s62 26 62 68z"/>
    <path d="M300 664c0-42 26-68 62-68s62 26 62 68z"/>
    <path d="M520 640c0-42 26-68 62-68s62 26 62 68z"/>
  </g>
  <g stroke="#8E2C1E" stroke-width="7" fill="none" stroke-linecap="round">
    <path d="M152 640v-96"/>
    <path d="M362 664v-96"/>
    <path d="M582 640v-96"/>
  </g>
  ${dust(45, 800, 1000, 107)}
</svg>`;

const wide = (body, id) => `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 600" width="800" height="600">
  ${body}
</svg>`;

const activityWorkshop = () => wide(`
  <defs>
    <linearGradient id="aw" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#F7E9C8"/>
      <stop offset="1" stop-color="#D9B77F"/>
    </linearGradient>
  </defs>
  <rect width="800" height="600" fill="url(#aw)"/>
  <rect y="380" width="800" height="220" fill="#A87843"/>
  <g transform="translate(250 130) scale(0.94)">${puppet('#8E2C1E', '#C39A4B', 0.95)}</g>
  <g fill="#E8D5AE" stroke="#A9843F" stroke-width="4">
    <path d="M110 250l26-96h96l30 96z"/>
  </g>
  <g>
    <rect x="580" y="220" width="120" height="160" rx="10" fill="#5E3C21"/>
    <path d="M596 250h88M596 280h88M596 310h88" stroke="#C39A4B" stroke-width="6"/>
  </g>
  ${dust(40, 800, 520, 131)}`);

const activityCampus = () => wide(`
  <defs>
    <linearGradient id="ac" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#E9D3A6"/>
      <stop offset="1" stop-color="#C9A86F"/>
    </linearGradient>
  </defs>
  <rect width="800" height="600" fill="url(#ac)"/>
  <rect x="90" y="90" width="620" height="330" rx="10" fill="#FBF2DA"/>
  <rect x="90" y="90" width="620" height="330" rx="10" fill="none" stroke="#8E2C1E" stroke-width="6"/>
  <g transform="translate(250 120) scale(0.72)">${puppet('#2A1A13', '#8E2C1E', 0.92)}</g>
  <path d="M90 420h620" stroke="#6F4B27" stroke-width="12"/>
  <g fill="#5B3A24">
    <path d="M120 500c0-34 22-56 50-56s50 22 50 56z"/>
    <path d="M330 520c0-34 22-56 50-56s50 22 50 56z"/>
    <path d="M540 500c0-34 22-56 50-56s50 22 50 56z"/>
  </g>
  ${dust(35, 800, 560, 149)}`);

const activityNight = () => wide(`
  <defs>
    <linearGradient id="an" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#1C2636"/>
      <stop offset="1" stop-color="#3A2622"/>
    </linearGradient>
  </defs>
  <rect width="800" height="600" fill="url(#an)"/>
  <circle cx="640" cy="120" r="44" fill="#FDF3D8" opacity=".9"/>
  <rect y="420" width="800" height="180" fill="#191C26"/>
  <g>
    <rect x="270" y="230" width="260" height="190" fill="#F6E2B4"/>
    <g transform="translate(310 244) scale(0.62)">${puppet('#241712', '#8E2C1E', 0.9)}</g>
    <path d="M250 230h300l-38-44H288z" fill="#7E2419"/>
    <path d="M270 420V230M530 420V230" stroke="#3B2119" stroke-width="12"/>
  </g>
  ${lantern(120, 200, 26, '#F0B85C')}
  ${lantern(680, 220, 26, '#F0B85C')}
  ${dust(50, 800, 500, 167)}`);

const puppetMark = () => `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 600 760" width="600" height="760">
  <g transform="translate(60 70) scale(1.1)" fill="#A93628" stroke="none" opacity=".9">
    ${puppet('#A93628', '#7E2419', 0.95)}
  </g>
</svg>`;

const png = (svg, file) => sharp(Buffer.from(svg)).png({ compressionLevel: 9 }).toFile(file);

await fs.promises.mkdir(seedDir, { recursive: true });
await fs.promises.mkdir(localDir, { recursive: true });

await Promise.all([
  png(stageScene(), `${seedDir}/shadow-stage.png`),
  png(riverScene(), `${seedDir}/yaowan-night.png`),
  png(workshopScene(), `${seedDir}/shadow-workshop.png`),
  png(puppetCover(), `${seedDir}/art-puppet.png`),
  png(townCover(), `${seedDir}/art-town.png`),
  png(classCover(), `${seedDir}/art-class.png`),
  png(activityWorkshop(), `${seedDir}/activity-workshop.png`),
  png(activityCampus(), `${seedDir}/activity-campus.png`),
  png(activityNight(), `${seedDir}/activity-night.png`),
  png(stageScene(), `${seedDir}/placeholder.png`),
  png(puppetMark(), `${localDir}/puppet-mark.png`)
]);

console.log('assets generated');
