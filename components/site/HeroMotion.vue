<script setup lang="ts">
// Home hero motion (DS v1.1 components/hero-motion.md, variant C). Place inside a
// `.hero-motion` section: draws the glow, the hex grid with a travelling wave, and the
// pause button (WCAG 2.2.2, stored as ivm-motion). Ported from the DS reference script
// _build/templates/_hero_motion.js. Without JS the still .hexgrid texture stays.
const HEX_R = 32
const HEX_W = 56
const ROW_H = 48
const WAVE_LENGTH = 560
const WAVE_PERIOD_S = 9
const BASE_ALPHA = 0.09
const MIN_WAVE_WIDTH = 768
const MAX_DPR = 2
const STORAGE_KEY = 'ivm-motion'

const canvas = ref<HTMLCanvasElement | null>(null)
const paused = ref(false)

let cleanup: (() => void) | null = null

function readStored(): string | null {
  try { return localStorage.getItem(STORAGE_KEY) } catch { return null }
}
function writeStored(value: 'on' | 'off') {
  try { localStorage.setItem(STORAGE_KEY, value) } catch { /* private mode: preference not kept */ }
}

function rgb(hex: string, fallback: string): string {
  const m = /^#([0-9a-f]{6})$/i.exec(hex.trim())
  if (!m) return fallback
  const n = parseInt(m[1], 16)
  return `${n >> 16} ${(n >> 8) & 255} ${n & 255}`
}

onMounted(() => {
  const el = canvas.value
  const hero = el?.parentElement
  const ctx = el?.getContext('2d')
  if (!el || !hero || !ctx) return

  const root = document.documentElement
  const reduced = matchMedia('(prefers-reduced-motion: reduce)')
  const darkScheme = matchMedia('(prefers-color-scheme: dark)')
  const st = { visible: true, w: 0, h: 0, t0: performance.now(), frozen: 0, raf: 0 }
  const col = { grid: '100 116 139', wave: '2 122 137', fill: 0.07, line: 0.42 }

  function readColors() {
    const cs = getComputedStyle(hero!)
    col.grid = rgb(cs.getPropertyValue('--grid-line'), col.grid)
    col.wave = rgb(cs.getPropertyValue('--accent'), col.wave)
    col.fill = parseFloat(cs.getPropertyValue('--wave-fill')) || col.fill
    col.line = parseFloat(cs.getPropertyValue('--wave-line')) || col.line
  }
  const now = () => (paused.value ? st.frozen : (performance.now() - st.t0) / 1000)
  const waveOn = () => !reduced.matches && st.w >= MIN_WAVE_WIDTH
  const running = () => waveOn() && !paused.value && st.visible && !document.hidden

  function cellCentre(r: number, c: number): [number, number] {
    return [c * HEX_W + (r & 1 ? HEX_W / 2 : 0), r * ROW_H]
  }

  function draw(t: number) {
    const { w, h } = st
    const rows = Math.ceil(h / ROW_H) + 2
    const cols = Math.ceil(w / HEX_W) + 2
    ctx!.clearRect(0, 0, w, h)
    ctx!.lineWidth = 1
    ctx!.beginPath()
    for (let r = -1; r < rows; r++) for (let c = -1; c < cols; c++) {
      // Three edges per cell: shared edges are drawn once.
      const [cx, cy] = cellCentre(r, c)
      ctx!.moveTo(cx - HEX_W / 2, cy - HEX_R / 2)
      ctx!.lineTo(cx, cy - HEX_R)
      ctx!.lineTo(cx + HEX_W / 2, cy - HEX_R / 2)
      ctx!.lineTo(cx + HEX_W / 2, cy + HEX_R / 2)
    }
    ctx!.strokeStyle = `rgb(${col.grid} / ${BASE_ALPHA})`
    ctx!.stroke()
    if (!waveOn()) return
    for (let r = -1; r < rows; r++) for (let c = -1; c < cols; c++) {
      const [cx, cy] = cellCentre(r, c)
      const crest = Math.max(0, Math.sin(((cx * 0.8 + cy * 0.6) / WAVE_LENGTH - t / WAVE_PERIOD_S) * Math.PI * 2))
      const i = crest ** 3 * (0.7 + 0.3 * Math.sin(cx / 97 - cy / 61 + t * 0.5))
      if (i < 0.02) continue
      ctx!.beginPath()
      ctx!.moveTo(cx, cy - HEX_R)
      ctx!.lineTo(cx + HEX_W / 2, cy - HEX_R / 2)
      ctx!.lineTo(cx + HEX_W / 2, cy + HEX_R / 2)
      ctx!.lineTo(cx, cy + HEX_R)
      ctx!.lineTo(cx - HEX_W / 2, cy + HEX_R / 2)
      ctx!.lineTo(cx - HEX_W / 2, cy - HEX_R / 2)
      ctx!.closePath()
      ctx!.fillStyle = `rgb(${col.wave} / ${(i * col.fill).toFixed(3)})`
      ctx!.fill()
      ctx!.strokeStyle = `rgb(${col.wave} / ${(i * col.line).toFixed(3)})`
      ctx!.stroke()
    }
  }

  function loop() {
    draw(now())
    st.raf = running() ? requestAnimationFrame(loop) : 0
  }
  function restart() {
    cancelAnimationFrame(st.raf)
    st.raf = 0
    loop()
  }
  function resize() {
    const rect = hero!.getBoundingClientRect()
    const dpr = Math.min(window.devicePixelRatio || 1, MAX_DPR)
    st.w = rect.width
    st.h = rect.height
    el!.width = Math.round(rect.width * dpr)
    el!.height = Math.round(rect.height * dpr)
    ctx!.setTransform(dpr, 0, 0, dpr, 0, 0)
    restart()
  }
  const recolor = () => { readColors(); draw(now()) }
  const onVisibility = () => restart()

  const stopPausedWatch = watch(paused, (p, was) => {
    if (p && !was) st.frozen = (performance.now() - st.t0) / 1000
    if (!p && was) st.t0 = performance.now() - st.frozen * 1000
    hero.classList.toggle('is-paused', p)
    restart()
  })

  const themeObserver = new MutationObserver(recolor)
  themeObserver.observe(root, { attributes: true, attributeFilter: ['data-theme'] })
  darkScheme.addEventListener('change', recolor)
  reduced.addEventListener('change', restart)
  document.addEventListener('visibilitychange', onVisibility)
  const io = 'IntersectionObserver' in window
    ? new IntersectionObserver((e) => { st.visible = e[0].isIntersecting; restart() })
    : null
  io?.observe(hero)
  const ro = 'ResizeObserver' in window ? new ResizeObserver(resize) : null
  if (ro) ro.observe(hero)
  else window.addEventListener('resize', resize)

  hero.classList.add('has-canvas')
  readColors()
  resize()
  if (readStored() === 'off') paused.value = true

  cleanup = () => {
    cancelAnimationFrame(st.raf)
    stopPausedWatch()
    themeObserver.disconnect()
    darkScheme.removeEventListener('change', recolor)
    reduced.removeEventListener('change', restart)
    document.removeEventListener('visibilitychange', onVisibility)
    io?.disconnect()
    if (ro) ro.disconnect()
    else window.removeEventListener('resize', resize)
    hero.classList.remove('has-canvas', 'is-paused')
  }
})

onBeforeUnmount(() => cleanup?.())

function toggle() {
  paused.value = !paused.value
  writeStored(paused.value ? 'off' : 'on')
}
</script>

<template>
  <div class="hero-motion__glow" aria-hidden="true" />
  <canvas ref="canvas" class="hero-motion__grid" aria-hidden="true" />
  <button
    type="button"
    class="icon-btn hero-motion__pause"
    :aria-pressed="paused"
    :aria-label="paused ? 'Play background motion' : 'Pause background motion'"
    @click="toggle"
  >
    <svg v-if="!paused" viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true"><path d="M9 6v12M15 6v12" /></svg>
    <svg v-else viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linejoin="round" aria-hidden="true"><path d="M8 5.5l10 6.5-10 6.5z" /></svg>
  </button>
</template>
