<script>
/* Hero motion (hero-motion.md): hex grid on a canvas with a travelling wave, plus the pause button. */
(function () {
  var hero = document.querySelector('.hero-motion'); if (!hero) return;
  var canvas = hero.querySelector('.hero-motion__grid'), btn = hero.querySelector('.hero-motion__pause');
  var ctx = canvas && canvas.getContext && canvas.getContext('2d'); if (!ctx) return;
  var HEX_R = 32, HEX_W = 56, ROW_H = 48, WAVE_LENGTH = 560, WAVE_PERIOD = 9, BASE_ALPHA = 0.09, MIN_WAVE_WIDTH = 768;
  var root = document.documentElement, reduced = matchMedia('(prefers-reduced-motion: reduce)');
  var st = { paused: false, visible: true, w: 0, h: 0, t0: performance.now(), frozen: 0, raf: 0 };
  var col = { grid: '100 116 139', wave: '2 122 137', fill: 0.07, line: 0.42 };
  function stored() { try { return localStorage.getItem('ivm-motion'); } catch (e) { return null; } }
  function store(v) { try { localStorage.setItem('ivm-motion', v); } catch (e) {} }
  function rgb(hex, fallback) {
    var m = /^#([0-9a-f]{6})$/i.exec(hex.trim()); if (!m) return fallback;
    var n = parseInt(m[1], 16); return (n >> 16) + ' ' + ((n >> 8) & 255) + ' ' + (n & 255);
  }
  function readColors() {
    var cs = getComputedStyle(hero);
    col.grid = rgb(cs.getPropertyValue('--grid-line'), col.grid);
    col.wave = rgb(cs.getPropertyValue('--accent'), col.wave);
    col.fill = parseFloat(cs.getPropertyValue('--wave-fill')) || col.fill;
    col.line = parseFloat(cs.getPropertyValue('--wave-line')) || col.line;
  }
  function now() { return st.paused ? st.frozen : (performance.now() - st.t0) / 1000; }
  function waveOn() { return !reduced.matches && st.w >= MIN_WAVE_WIDTH; }
  function running() { return waveOn() && !st.paused && st.visible && !document.hidden; }
  function draw(t) {
    var w = st.w, h = st.h, rows = Math.ceil(h / ROW_H) + 2, cols = Math.ceil(w / HEX_W) + 2, r, c, cx, cy;
    ctx.clearRect(0, 0, w, h); ctx.lineWidth = 1; ctx.beginPath();
    for (r = -1; r < rows; r++) for (c = -1; c < cols; c++) {   // three edges per cell: shared edges are drawn once
      cx = c * HEX_W + (r & 1 ? HEX_W / 2 : 0); cy = r * ROW_H;
      ctx.moveTo(cx - HEX_W / 2, cy - HEX_R / 2); ctx.lineTo(cx, cy - HEX_R);
      ctx.lineTo(cx + HEX_W / 2, cy - HEX_R / 2); ctx.lineTo(cx + HEX_W / 2, cy + HEX_R / 2);
    }
    ctx.strokeStyle = 'rgb(' + col.grid + ' / ' + BASE_ALPHA + ')'; ctx.stroke();
    if (!waveOn()) return;
    for (r = -1; r < rows; r++) for (c = -1; c < cols; c++) {
      cx = c * HEX_W + (r & 1 ? HEX_W / 2 : 0); cy = r * ROW_H;
      var crest = Math.max(0, Math.sin(((cx * 0.8 + cy * 0.6) / WAVE_LENGTH - t / WAVE_PERIOD) * Math.PI * 2));
      var i = crest * crest * crest * (0.7 + 0.3 * Math.sin(cx / 97 - cy / 61 + t * 0.5));
      if (i < 0.02) continue;
      ctx.beginPath(); ctx.moveTo(cx, cy - HEX_R); ctx.lineTo(cx + HEX_W / 2, cy - HEX_R / 2); ctx.lineTo(cx + HEX_W / 2, cy + HEX_R / 2);
      ctx.lineTo(cx, cy + HEX_R); ctx.lineTo(cx - HEX_W / 2, cy + HEX_R / 2); ctx.lineTo(cx - HEX_W / 2, cy - HEX_R / 2); ctx.closePath();
      ctx.fillStyle = 'rgb(' + col.wave + ' / ' + (i * col.fill).toFixed(3) + ')'; ctx.fill();
      ctx.strokeStyle = 'rgb(' + col.wave + ' / ' + (i * col.line).toFixed(3) + ')'; ctx.stroke();
    }
  }
  function loop() { draw(now()); st.raf = running() ? requestAnimationFrame(loop) : 0; }
  function restart() { cancelAnimationFrame(st.raf); st.raf = 0; loop(); }
  function resize() {
    var r = hero.getBoundingClientRect(), dpr = Math.min(window.devicePixelRatio || 1, 2);
    st.w = r.width; st.h = r.height;
    canvas.width = Math.round(r.width * dpr); canvas.height = Math.round(r.height * dpr);
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0); restart();
  }
  function setPaused(p) {
    if (p && !st.paused) st.frozen = now();
    if (!p && st.paused) st.t0 = performance.now() - st.frozen * 1000;
    st.paused = p; hero.classList.toggle('is-paused', p);
    if (btn) {
      btn.setAttribute('aria-pressed', String(p));
      btn.setAttribute('aria-label', p ? 'Play background motion' : 'Pause background motion');
      btn.querySelector('.i-pause').toggleAttribute('hidden', p); btn.querySelector('.i-play').toggleAttribute('hidden', !p);
    }
    restart();
  }
  if (btn) btn.addEventListener('click', function () { setPaused(!st.paused); store(st.paused ? 'off' : 'on'); });
  new MutationObserver(function () { readColors(); draw(now()); }).observe(root, { attributes: true, attributeFilter: ['data-theme'] });
  matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function () { readColors(); draw(now()); });
  reduced.addEventListener('change', restart);
  document.addEventListener('visibilitychange', restart);
  if ('IntersectionObserver' in window) new IntersectionObserver(function (e) { st.visible = e[0].isIntersecting; restart(); }).observe(hero);
  if ('ResizeObserver' in window) new ResizeObserver(resize).observe(hero); else addEventListener('resize', resize);
  hero.classList.add('has-canvas');
  readColors(); resize();
  if (stored() === 'off') setPaused(true);
})();
</script>
