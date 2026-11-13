<script setup lang="ts">
// AI service advisor (epic #112, spec docs/design-system/components/advisor.md).
// POST /api/advisor -> { reply, recommendations: [{ id, why }], next_action }.
// Names, prices and links come from the catalogue by id, never from model text.
import { trackEvent } from '~/services/analytics'

type Role = 'user' | 'assistant'
type NextAction = 'buy' | 'call' | 'schedule' | 'contact' | 'early_access'
interface ChatMessage { role: Role; content: string }
interface Recommendation { id: string; why: string }
interface AdvisorReply { reply: string; recommendations: Recommendation[]; next_action: NextAction }
interface Turn { role: Role; content: string; recommendations?: Recommendation[]; next_action?: NextAction }

const MAX_CHARS = 600
const MAX_TURNS = 4
const TIMEOUT_MS = 20000

const { pkg } = useCatalogue()

const turns = ref<Turn[]>([])
const draft = ref('')
const state = ref<'idle' | 'thinking' | 'answer' | 'error' | 'limit'>('idle')
const showLead = ref(false)
const textarea = ref<HTMLTextAreaElement | null>(null)
const sessionId = ref('')

const chips = ['We’re rolling out an AI assistant', 'Does the EU AI Act apply to us?', 'ISO 42001 readiness']
const userTurns = computed(() => turns.value.filter((t) => t.role === 'user').length)
const lastAnswer = computed(() => [...turns.value].reverse().find((t) => t.role === 'assistant'))

const summary = computed(() =>
  turns.value.map((t) => `${t.role === 'user' ? 'Me' : 'Advisor'}: ${t.content}`).join('\n\n').slice(0, 1800),
)

function newSessionId(): string {
  return typeof crypto !== 'undefined' && 'randomUUID' in crypto ? crypto.randomUUID() : String(Date.now())
}

function useChip(text: string) {
  draft.value = text
  textarea.value?.focus()
}

function ctaTarget(action: NextAction, rec?: Recommendation): string {
  if (action === 'buy') return rec ? `/packages/${rec.id}` : '/packages'
  if (action === 'call') return '/booking'
  if (action === 'schedule') return '/booking'
  if (action === 'early_access') return '/ivgil#early-access'
  return ''
}

async function send() {
  const text = draft.value.trim()
  if (!text || state.value === 'thinking') return
  if (userTurns.value >= MAX_TURNS) {
    state.value = 'limit'
    return
  }
  if (!sessionId.value) sessionId.value = newSessionId()

  turns.value = [...turns.value, { role: 'user', content: text }]
  state.value = 'thinking'
  showLead.value = false
  trackEvent('advisor_send', { turn: userTurns.value })

  const messages: ChatMessage[] = turns.value.map(({ role, content }) => ({ role, content }))
  try {
    const res = await $fetch<AdvisorReply>('/api/advisor', {
      method: 'POST',
      body: { session_id: sessionId.value, messages },
      timeout: TIMEOUT_MS,
    })
    const recs = (res.recommendations ?? []).filter((r) => pkg(r.id)).slice(0, 3)
    turns.value = [...turns.value, { role: 'assistant', content: res.reply, recommendations: recs, next_action: res.next_action ?? 'call' }]
    draft.value = ''
    state.value = userTurns.value >= MAX_TURNS ? 'limit' : 'answer'
    trackEvent('advisor_answer', { recommendations: recs.map((r) => r.id).join(','), next_action: res.next_action })
  } catch (err: unknown) {
    const status = (err as { status?: number })?.status
    // Keep the user's text in the box; drop the unanswered turn from the log.
    turns.value = turns.value.slice(0, -1)
    state.value = status === 429 ? 'limit' : 'error'
    trackEvent('advisor_error', { status: status ?? 0 })
  } finally {
    await nextTick()
    textarea.value?.focus()
  }
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    send()
  }
}

const ctas = computed(() => {
  const answer = lastAnswer.value
  const primary: NextAction = answer?.next_action ?? 'call'
  const firstRec = answer?.recommendations?.[0]
  // The ivGIL early-access CTA appears only when the advisor proposes it (IA proposal §6).
  const platform = primary === 'early_access'
    ? [{ key: 'early_access' as const, label: 'Request early access', sub: 'ivGIL, in development', to: ctaTarget('early_access') }]
    : []
  return [
    ...platform,
    { key: 'buy' as const, label: 'Order package', sub: 'Fixed scope and price', to: ctaTarget('buy', firstRec) },
    { key: 'call' as const, label: 'Book a call', sub: '[30 min, free]', to: ctaTarget('call') },
    { key: 'schedule' as const, label: 'Schedule later', sub: 'Pick a time', to: ctaTarget('schedule') },
    { key: 'contact' as const, label: 'Leave contact', sub: 'Reply within 1 working day', to: '' },
  ].slice(0, 4).map((c) => ({ ...c, primary: c.key === primary }))
})

function ctaClick(key: NextAction) {
  trackEvent('advisor_cta', { cta: key })
  if (key === 'contact') showLead.value = true
}
</script>

<template>
  <section class="advisor" aria-labelledby="advisor-title">
    <div class="advisor__head">
      <SiteAiAvatar />
      <strong id="advisor-title">Which service do you need?</strong>
      <span class="ai-label">AI</span>
    </div>

    <div class="advisor__log" aria-live="polite">
      <div class="msg msg--ai">
        <SiteAiAvatar />
        <div>
          <p style="margin: 0">
            Describe in a sentence or two what you’re working on or worried about. I’ll point you to the service that fits.
          </p>
        </div>
      </div>

      <template v-for="(t, i) in turns" :key="i">
        <div v-if="t.role === 'user'" class="msg msg--user">{{ t.content }}</div>
        <div v-else class="msg msg--ai">
          <SiteAiAvatar />
          <div>
            <p v-for="(para, j) in t.content.split(/\n{2,}/)" :key="j">{{ para }}</p>
            <template v-if="t.recommendations?.length">
              <p class="small muted" style="margin-bottom: var(--space-2)">Services that fit</p>
              <NuxtLink
                v-for="r in t.recommendations"
                :key="r.id"
                :to="`/packages/${r.id}`"
                class="rec"
                style="text-decoration: none; color: inherit"
                @click="trackEvent('advisor_recommendation_click', { id: r.id })"
              >
                <SitePillarIcon :pillar="pkg(r.id)!.pillar" :size="40" />
                <span class="rec__title">{{ pkg(r.id)!.name }}</span>
                <span class="rec__price">{{ formatPrice(pkg(r.id)!) }}</span>
                <p class="rec__why">{{ r.why }}</p>
              </NuxtLink>
            </template>
            <div v-if="t === lastAnswer" class="cta-set">
              <template v-for="c in ctas" :key="c.key">
                <NuxtLink v-if="c.to" :to="c.to" class="cta" :class="{ 'cta--primary': c.primary }" @click="ctaClick(c.key)">
                  <strong>{{ c.label }}</strong><span>{{ c.sub }}</span>
                </NuxtLink>
                <button v-else type="button" class="cta" :class="{ 'cta--primary': c.primary }" style="text-align: left; font: inherit; cursor: pointer" @click="ctaClick(c.key)">
                  <strong>{{ c.label }}</strong><span>{{ c.sub }}</span>
                </button>
              </template>
            </div>
          </div>
        </div>
      </template>

      <div v-if="state === 'thinking'" class="msg msg--ai">
        <SiteAiAvatar />
        <div class="typing" role="status" aria-label="Advisor is answering"><i /><i /><i /></div>
      </div>

      <div v-if="showLead" class="card">
        <h4>Leave your contact</h4>
        <AdvisorLeadForm source="advisor" :prefill="summary" />
      </div>

      <div v-if="state === 'error'" class="notice notice--error" role="alert">
        <p><strong>The advisor can’t answer right now.</strong> Your question wasn’t lost.</p>
        <div class="row">
          <NuxtLink to="/booking" class="btn btn--primary btn--sm">Book a call</NuxtLink>
          <NuxtLink to="/packages" class="btn btn--secondary btn--sm">Browse packages</NuxtLink>
        </div>
      </div>

      <div v-if="state === 'limit'" class="notice" role="status">
        <p>That’s as far as I can help here. A short call with Nick is the quickest next step.</p>
        <div class="row">
          <NuxtLink to="/booking" class="btn btn--primary btn--sm">Book a call</NuxtLink>
          <button type="button" class="btn btn--secondary btn--sm" @click="showLead = true">Leave contact</button>
        </div>
      </div>
    </div>

    <div v-if="turns.length === 0" class="prompt-chips">
      <button v-for="c in chips" :key="c" type="button" class="chip" @click="useChip(c)">{{ c }}</button>
    </div>

    <form class="advisor__form" @submit.prevent="send">
      <details v-if="turns.length === 0" class="notice ai-notice" style="margin-bottom: var(--space-3)">
        <summary>
          <span><strong>AI assistant, not a person.</strong> Please don’t include personal or confidential information.</span>
          <span class="ai-notice__more">Details</span>
        </summary>
        <p>
          It suggests which of our services fit. It is guidance, not legal advice. Your message is sent to our AI
          provider (Mistral AI, EU) to generate the answer. We don’t store the conversation.
          <NuxtLink to="/ai-use">How we handle your message</NuxtLink>
        </p>
      </details>
      <div class="advisor__input">
        <label class="sr-only" for="advisor-q">Describe your situation</label>
        <textarea
          id="advisor-q"
          ref="textarea"
          v-model="draft"
          rows="2"
          :maxlength="MAX_CHARS"
          placeholder="e.g. Our sales team uses ChatGPT. What do we have to document?"
          :disabled="state === 'limit'"
          @keydown="onKeydown"
        />
        <button type="submit" class="btn btn--primary" aria-label="Send" :disabled="state === 'thinking' || state === 'limit' || !draft.trim()">
          <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M5 12h14M13 6l6 6-6 6" />
          </svg>
          <span class="send-label">Send</span>
        </button>
      </div>
      <div class="advisor__meta">
        <span>AI assistant · guidance, not legal advice · <NuxtLink to="/ai-use">how we handle your message</NuxtLink></span>
        <span>{{ draft.length }} / {{ MAX_CHARS }}</span>
      </div>
    </form>
  </section>
</template>

<style scoped>
.ai-notice summary { display: flex; gap: var(--space-3); justify-content: space-between; align-items: baseline; cursor: pointer; list-style: none; }
.ai-notice summary::-webkit-details-marker { display: none; }
.ai-notice__more { flex: none; color: var(--accent); text-decoration: underline; text-underline-offset: 3px; white-space: nowrap; }
.ai-notice__more::after { content: ' ▾'; text-decoration: none; display: inline-block; margin-left: 2px; }
.ai-notice[open] .ai-notice__more::after { transform: rotate(180deg); }
.ai-notice[open] summary { margin-bottom: var(--space-2); }
.ai-notice p { margin: 0; }
/* Make the input the obvious place to type. */
.advisor__input textarea { border-width: 2px; border-color: var(--data); min-height: 3.5rem; }
.advisor__input textarea:focus { outline: 2px solid var(--focus); outline-offset: 1px; }
.send-label { display: none; }
/* Small screens: full-width input, send button on its own row below. */
@media (max-width: 560px) {
  .advisor__input { flex-direction: column; align-items: stretch; }
  .advisor__input textarea { width: 100%; }
  .advisor__input .btn { width: 100%; justify-content: center; gap: var(--space-2); }
  .send-label { display: inline; }
}
</style>
