<script setup lang="ts">
import { trackEvent } from '~/services/analytics'

useSeoMeta({
  title: 'Data Maturity Assessment | ivmanto.com',
  description:
    'Take our free Data Maturity Assessment to evaluate your organization\'s data readiness across 5 key dimensions. Get personalized recommendations for your data strategy.',
  ogTitle: 'Data Maturity Assessment | ivmanto.com',
  ogDescription:
    'Evaluate your data readiness across Strategy, Architecture, Governance, AI Readiness, and Culture. Get a personalized maturity score and recommendations.',
  twitterTitle: 'Data Maturity Assessment | ivmanto.com',
  twitterDescription:
    'Free interactive assessment to evaluate your organization\'s data maturity and get actionable recommendations.',
})

useHead({
  script: [
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'WebPage',
        '@id': 'https://ivmanto.com/assessment',
        name: 'Data Maturity Assessment',
        url: 'https://ivmanto.com/assessment',
        description:
          'Free interactive assessment to evaluate your organization\'s data maturity across 5 key dimensions.',
        isPartOf: { '@id': 'https://ivmanto.com/#website' },
        about: {
          '@type': 'Quiz',
          name: 'Data Maturity Assessment',
          description:
            'Evaluate your organization\'s data readiness across Strategy, Architecture, Governance, AI Readiness, and Culture. Get a personalized maturity score and recommendations.',
          educationalLevel: 'Professional',
          numberOfQuestions: 10,
          timeRequired: 'PT3M',
        },
        provider: { '@id': 'https://ivmanto.com/#organization' },
      }),
    },
  ],
})

// ─── Assessment Data ────────────────────────────────────────────────
interface Question {
  id: number
  dimension: string
  dimensionIndex: number
  text: string
  options: { label: string; value: number }[]
}

const dimensions = ['Strategy', 'Architecture', 'Governance', 'AI Readiness', 'Culture']

const questions: Question[] = [
  // Strategy (0)
  {
    id: 1,
    dimension: 'Strategy',
    dimensionIndex: 0,
    text: 'How well is your data strategy aligned with your overall business objectives?',
    options: [
      { label: 'No formal data strategy exists', value: 1 },
      { label: 'Data strategy exists but is loosely connected to business goals', value: 2 },
      { label: 'Data strategy is documented and aligned with key business initiatives', value: 3 },
      { label: 'Data strategy drives business decisions and is reviewed regularly', value: 4 },
    ],
  },
  {
    id: 2,
    dimension: 'Strategy',
    dimensionIndex: 0,
    text: 'How does your organization measure the value of its data assets?',
    options: [
      { label: 'Data value is not measured or discussed', value: 1 },
      { label: 'Anecdotal evidence of data value in some departments', value: 2 },
      { label: 'KPIs and metrics track data-driven outcomes across the organization', value: 3 },
      { label: 'Data ROI is quantified and reported to leadership regularly', value: 4 },
    ],
  },
  // Architecture (1)
  {
    id: 3,
    dimension: 'Architecture',
    dimensionIndex: 1,
    text: 'How would you describe your current data platform infrastructure?',
    options: [
      { label: 'Fragmented — data lives in spreadsheets, local databases, and email', value: 1 },
      { label: 'Partially centralized — a data warehouse or lake exists but with major gaps', value: 2 },
      { label: 'Modern platform — cloud-based with automated pipelines and a catalog', value: 3 },
      { label: 'Advanced ecosystem — lakehouse/mesh architecture with real-time capabilities', value: 4 },
    ],
  },
  {
    id: 4,
    dimension: 'Architecture',
    dimensionIndex: 1,
    text: 'How easily can teams access the data they need for analysis or machine learning?',
    options: [
      { label: 'Data access requires manual requests and takes days or weeks', value: 1 },
      { label: 'Some self-service access exists, but many datasets are siloed', value: 2 },
      { label: 'Most teams have self-service access with proper permissioning', value: 3 },
      { label: 'Data is discoverable, well-documented, and accessible within minutes', value: 4 },
    ],
  },
  // Governance (2)
  {
    id: 5,
    dimension: 'Governance',
    dimensionIndex: 2,
    text: 'How does your organization manage data quality?',
    options: [
      { label: 'No formal data quality processes — issues are fixed ad-hoc', value: 1 },
      { label: 'Basic validation exists at ingestion but no ongoing quality monitoring', value: 2 },
      { label: 'Data quality rules are defined and monitored with automated alerts', value: 3 },
      { label: 'Comprehensive data quality framework with SLAs and ownership', value: 4 },
    ],
  },
  {
    id: 6,
    dimension: 'Governance',
    dimensionIndex: 2,
    text: 'How do you handle data security, privacy, and compliance requirements?',
    options: [
      { label: 'Minimal — we react to issues and audits as they come', value: 1 },
      { label: 'Basic controls are in place but not consistently enforced', value: 2 },
      { label: 'Policies are documented and enforced with role-based access controls', value: 3 },
      { label: 'Security is "baked in" — automated compliance checks, encryption, and audit trails', value: 4 },
    ],
  },
  // AI Readiness (3)
  {
    id: 7,
    dimension: 'AI Readiness',
    dimensionIndex: 3,
    text: 'What is the current state of AI/ML adoption in your organization?',
    options: [
      { label: 'No AI/ML initiatives — it\'s not on our roadmap yet', value: 1 },
      { label: 'Exploring or piloting AI in isolated use cases', value: 2 },
      { label: 'Multiple AI models in production driving real business value', value: 3 },
      { label: 'AI is deeply embedded in operations with MLOps and continuous improvement', value: 4 },
    ],
  },
  {
    id: 8,
    dimension: 'AI Readiness',
    dimensionIndex: 3,
    text: 'How prepared is your data infrastructure to support AI/ML workloads?',
    options: [
      { label: 'Data is too messy / siloed for any meaningful AI work', value: 1 },
      { label: 'Some clean datasets exist but significant prep is needed each time', value: 2 },
      { label: 'Feature stores and clean training data are available for priority use cases', value: 3 },
      { label: 'Production-grade ML pipelines with automated retraining and monitoring', value: 4 },
    ],
  },
  // Culture (4)
  {
    id: 9,
    dimension: 'Culture',
    dimensionIndex: 4,
    text: 'How data-literate is your organization overall?',
    options: [
      { label: 'Most decisions are made on gut feeling — data is rarely consulted', value: 1 },
      { label: 'Some teams use data, but there\'s no organization-wide data literacy effort', value: 2 },
      { label: 'Data literacy training is available and most managers use dashboards', value: 3 },
      { label: 'Data-driven culture — all levels use data for decisions with strong literacy programs', value: 4 },
    ],
  },
  {
    id: 10,
    dimension: 'Culture',
    dimensionIndex: 4,
    text: 'Who owns data in your organization?',
    options: [
      { label: 'Nobody — data ownership is unclear and unassigned', value: 1 },
      { label: 'IT owns most data — business units have minimal involvement', value: 2 },
      { label: 'Shared ownership — domain teams manage their data with IT support', value: 3 },
      { label: 'Clear data ownership model with accountable data stewards per domain', value: 4 },
    ],
  },
]

// ─── State ──────────────────────────────────────────────────────────
const phase = ref<'intro' | 'questions' | 'results'>('intro')
const currentQuestion = ref(0)
const answers = ref<number[]>(new Array(questions.length).fill(0))
const animating = ref(false)
const slideDirection = ref<'left' | 'right'>('left')

const progress = computed(() => ((currentQuestion.value + 1) / questions.length) * 100)
const isLastQuestion = computed(() => currentQuestion.value === questions.length - 1)
const canProceed = computed(() => answers.value[currentQuestion.value] > 0)

// ─── Scoring ────────────────────────────────────────────────────────
const dimensionScores = computed(() => {
  const scores: number[] = new Array(dimensions.length).fill(0)
  const counts: number[] = new Array(dimensions.length).fill(0)
  questions.forEach((q, i) => {
    if (answers.value[i] > 0) {
      scores[q.dimensionIndex] += answers.value[i]
      counts[q.dimensionIndex]++
    }
  })
  return scores.map((s, i) => (counts[i] > 0 ? s / counts[i] : 0))
})

const overallScore = computed(() => {
  const filled = dimensionScores.value.filter((s) => s > 0)
  if (filled.length === 0) return 0
  return filled.reduce((a, b) => a + b, 0) / filled.length
})

const maturityLevel = computed(() => {
  const s = overallScore.value
  if (s <= 1.5) return { label: 'Foundational', color: '#ef4444', description: 'Your data practice is in its early stages. There are significant opportunities to establish formal processes and begin building a strategic data foundation.' }
  if (s <= 2.5) return { label: 'Managed', color: '#f59e0b', description: 'You have some data practices in place, but they are inconsistent across the organization. Formalizing and standardizing will unlock major value.' }
  if (s <= 3.5) return { label: 'Integrated', color: '#00a896', description: 'Your organization has a solid data foundation with consistent practices. Focus on optimization and advanced capabilities like AI/ML to accelerate further.' }
  return { label: 'Optimized', color: '#059669', description: 'You are among the most data-mature organizations. Continue innovating and ensure your data practices scale with business growth.' }
})

const dimensionRecommendations: Record<string, string[]> = {
  Strategy: [
    'Develop a formal data strategy document aligned with business OKRs',
    'Establish a Data Council with representation from each business unit',
    'Define measurable KPIs to track data-driven business impact',
  ],
  Architecture: [
    'Evaluate modern cloud data platforms like GCP BigQuery or a Lakehouse architecture',
    'Implement a data catalog for discoverability and self-service access',
    'Design automated data pipelines to eliminate manual ETL processes',
  ],
  Governance: [
    'Define data quality rules and implement automated validation frameworks',
    'Establish clear data ownership with accountable stewards per domain',
    'Implement role-based access controls and audit trail capabilities',
  ],
  'AI Readiness': [
    'Start with well-scoped AI pilot projects that address specific business problems',
    'Build feature stores and clean datasets for priority ML use cases',
    'Invest in MLOps tooling for model monitoring, retraining, and governance',
  ],
  Culture: [
    'Launch a data literacy program with training tailored to different roles',
    'Create data champions in each department to promote data-driven decisions',
    'Make key metrics visible through accessible dashboards and self-service tools',
  ],
}

// ─── Radar Chart SVG ────────────────────────────────────────────────
const radarSize = 280
const radarCenter = radarSize / 2
const radarRadius = 110

function polarToCartesian(angle: number, radius: number): { x: number; y: number } {
  const radian = ((angle - 90) * Math.PI) / 180
  return {
    x: radarCenter + radius * Math.cos(radian),
    y: radarCenter + radius * Math.sin(radian),
  }
}

const radarGridLevels = [0.25, 0.5, 0.75, 1.0]

function radarPolygonPoints(level: number): string {
  return dimensions
    .map((_, i) => {
      const angle = (360 / dimensions.length) * i
      const p = polarToCartesian(angle, radarRadius * level)
      return `${p.x},${p.y}`
    })
    .join(' ')
}

const radarDataPoints = computed(() => {
  return dimensions
    .map((_, i) => {
      const score = dimensionScores.value[i] / 4 // normalize 0–1
      const angle = (360 / dimensions.length) * i
      const p = polarToCartesian(angle, radarRadius * score)
      return `${p.x},${p.y}`
    })
    .join(' ')
})

const radarLabelPositions = computed(() => {
  return dimensions.map((label, i) => {
    const angle = (360 / dimensions.length) * i
    const p = polarToCartesian(angle, radarRadius + 24)
    return { label, x: p.x, y: p.y }
  })
})

// ─── Navigation ─────────────────────────────────────────────────────
function startAssessment() {
  phase.value = 'questions'
  trackEvent('assessment_start', {})
}

function selectAnswer(value: number) {
  answers.value[currentQuestion.value] = value
}

function nextQuestion() {
  if (!canProceed.value) return
  if (isLastQuestion.value) {
    phase.value = 'results'
    trackEvent('assessment_complete', {
      overall_score: overallScore.value.toFixed(2),
      maturity_level: maturityLevel.value.label,
    })
    return
  }
  slideDirection.value = 'left'
  animating.value = true
  setTimeout(() => {
    currentQuestion.value++
    animating.value = false
  }, 200)
}

function prevQuestion() {
  if (currentQuestion.value === 0) return
  slideDirection.value = 'right'
  animating.value = true
  setTimeout(() => {
    currentQuestion.value--
    animating.value = false
  }, 200)
}

function restartAssessment() {
  answers.value = new Array(questions.length).fill(0)
  currentQuestion.value = 0
  phase.value = 'intro'
}

async function shareResults() {
  const shareText = `My Data Maturity Score: ${overallScore.value.toFixed(1)}/4.0 (${maturityLevel.value.label}). Take the free assessment:`
  const shareUrl = 'https://ivmanto.com/assessment'
  if (navigator.share) {
    try {
      await navigator.share({ title: 'Data Maturity Assessment', text: shareText, url: shareUrl })
    } catch (e) {
      console.error('Share failed:', e)
    }
  } else {
    await navigator.clipboard.writeText(`${shareText} ${shareUrl}`)
  }
}
</script>

<template>
  <div class="min-h-screen bg-light-gray">
    <!-- ━━━ INTRO PHASE ━━━ -->
    <div v-if="phase === 'intro'" class="py-16 md:py-24">
      <div class="container mx-auto px-6 max-w-4xl">
        <!-- Hero -->
        <div class="text-center">
          <span
            class="inline-block bg-primary bg-opacity-10 text-primary font-semibold text-sm px-4 py-1.5 rounded-full mb-6"
          >
            Free Interactive Tool
          </span>
          <h1 class="text-4xl md:text-5xl font-bold text-dark-slate leading-tight">
            Data Maturity<br />
            <span class="text-primary">Assessment</span>
          </h1>
          <p class="text-lg md:text-xl text-gray-600 mt-6 max-w-2xl mx-auto">
            Evaluate your organization's data readiness across 5 critical dimensions and get
            personalized recommendations to accelerate your data strategy.
          </p>
        </div>

        <!-- Dimensions Preview -->
        <div class="mt-16 grid grid-cols-2 md:grid-cols-5 gap-4">
          <div
            v-for="dim in dimensions"
            :key="dim"
            class="bg-white p-4 rounded-xl text-center shadow-sm border border-gray-100"
          >
            <div
              class="w-10 h-10 mx-auto mb-3 bg-primary bg-opacity-10 rounded-lg flex items-center justify-center"
            >
              <svg
                class="w-5 h-5 text-primary"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <p class="text-sm font-semibold text-dark-slate">{{ dim }}</p>
          </div>
        </div>

        <!-- Info Cards -->
        <div class="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
            <p class="text-3xl font-bold text-primary">10</p>
            <p class="text-sm text-gray-600 mt-1">Questions</p>
          </div>
          <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
            <p class="text-3xl font-bold text-primary">3 min</p>
            <p class="text-sm text-gray-600 mt-1">To complete</p>
          </div>
          <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
            <p class="text-3xl font-bold text-primary">Free</p>
            <p class="text-sm text-gray-600 mt-1">No signup required</p>
          </div>
        </div>

        <!-- CTA -->
        <div class="mt-12 text-center">
          <button
            id="start-assessment"
            @click="startAssessment"
            class="bg-primary text-white font-bold py-4 px-10 rounded-xl text-lg hover:bg-opacity-90 transition-all shadow-lg shadow-primary/20 hover:shadow-xl hover:shadow-primary/30 hover:-translate-y-0.5"
          >
            Start Assessment
          </button>
          <p class="text-sm text-gray-500 mt-4">Your answers are processed locally — nothing is stored.</p>
        </div>
      </div>
    </div>

    <!-- ━━━ QUESTIONS PHASE ━━━ -->
    <div v-else-if="phase === 'questions'" class="py-12 md:py-20">
      <div class="container mx-auto px-6 max-w-2xl">
        <!-- Progress Bar -->
        <div class="mb-10">
          <div class="flex justify-between items-center mb-2">
            <span class="text-sm font-medium text-gray-600">Question {{ currentQuestion + 1 }} of {{ questions.length }}</span>
            <span class="text-sm font-medium text-primary">{{ questions[currentQuestion].dimension }}</span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
            <div
              class="h-2 rounded-full transition-all duration-500 ease-out"
              :style="{ width: `${progress}%`, background: 'linear-gradient(90deg, #00a896, #02c39a)' }"
            />
          </div>
        </div>

        <!-- Question Card -->
        <div
          class="bg-white rounded-2xl shadow-lg border border-gray-100 p-8 md:p-10 transition-all duration-200"
          :class="{
            'opacity-0 translate-x-8': animating && slideDirection === 'left',
            'opacity-0 -translate-x-8': animating && slideDirection === 'right',
          }"
        >
          <h2 class="text-xl md:text-2xl font-bold text-dark-slate leading-snug">
            {{ questions[currentQuestion].text }}
          </h2>

          <!-- Options -->
          <div class="mt-8 space-y-3">
            <button
              v-for="option in questions[currentQuestion].options"
              :key="option.value"
              @click="selectAnswer(option.value)"
              class="w-full text-left p-4 rounded-xl border-2 transition-all duration-200"
              :class="
                answers[currentQuestion] === option.value
                  ? 'border-primary bg-primary bg-opacity-5 text-dark-slate shadow-sm'
                  : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50 text-gray-700'
              "
            >
              <div class="flex items-start gap-3">
                <div
                  class="w-5 h-5 mt-0.5 rounded-full border-2 flex-shrink-0 flex items-center justify-center transition-all"
                  :class="
                    answers[currentQuestion] === option.value
                      ? 'border-primary bg-primary'
                      : 'border-gray-300'
                  "
                >
                  <svg
                    v-if="answers[currentQuestion] === option.value"
                    class="w-3 h-3 text-white"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </div>
                <span class="text-sm md:text-base">{{ option.label }}</span>
              </div>
            </button>
          </div>
        </div>

        <!-- Navigation -->
        <div class="mt-8 flex justify-between items-center">
          <button
            @click="prevQuestion"
            :disabled="currentQuestion === 0"
            class="px-6 py-2.5 rounded-lg font-medium transition-all text-sm"
            :class="
              currentQuestion === 0
                ? 'text-gray-400 cursor-not-allowed'
                : 'text-gray-600 hover:text-dark-slate hover:bg-white'
            "
          >
            &larr; Previous
          </button>
          <button
            @click="nextQuestion"
            :disabled="!canProceed"
            class="px-8 py-2.5 rounded-lg font-bold transition-all text-sm"
            :class="
              canProceed
                ? 'bg-primary text-white hover:bg-opacity-90 shadow-md hover:shadow-lg'
                : 'bg-gray-200 text-gray-400 cursor-not-allowed'
            "
          >
            {{ isLastQuestion ? 'See Results' : 'Next &rarr;' }}
          </button>
        </div>
      </div>
    </div>

    <!-- ━━━ RESULTS PHASE ━━━ -->
    <div v-else class="py-12 md:py-20">
      <div class="container mx-auto px-6 max-w-5xl">
        <!-- Overall Score Header -->
        <div class="bg-white rounded-2xl shadow-lg overflow-hidden">
          <div
            class="p-8 md:p-12 text-center text-white"
            :style="{ background: `linear-gradient(135deg, ${maturityLevel.color}, ${maturityLevel.color}dd)` }"
          >
            <p class="text-sm font-medium uppercase tracking-widest opacity-80">Your Data Maturity Level</p>
            <h2 class="text-4xl md:text-5xl font-bold mt-2">{{ maturityLevel.label }}</h2>
            <p class="text-6xl md:text-7xl font-bold mt-4 opacity-90">
              {{ overallScore.toFixed(1) }}<span class="text-2xl opacity-60">/4.0</span>
            </p>
          </div>
          <div class="p-8 text-center">
            <p class="text-lg text-gray-600 max-w-2xl mx-auto">{{ maturityLevel.description }}</p>
          </div>
        </div>

        <!-- Radar Chart + Dimension Breakdown -->
        <div class="mt-12 grid grid-cols-1 lg:grid-cols-2 gap-8">
          <!-- Radar Chart -->
          <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-8">
            <h3 class="text-xl font-bold text-dark-slate mb-6 text-center">Dimension Scores</h3>
            <div class="flex justify-center">
              <svg :width="radarSize" :height="radarSize" class="overflow-visible">
                <!-- Grid levels -->
                <polygon
                  v-for="level in radarGridLevels"
                  :key="level"
                  :points="radarPolygonPoints(level)"
                  fill="none"
                  stroke="#e5e7eb"
                  stroke-width="1"
                />
                <!-- Axes -->
                <line
                  v-for="(_, i) in dimensions"
                  :key="'axis-' + i"
                  :x1="radarCenter"
                  :y1="radarCenter"
                  :x2="polarToCartesian((360 / dimensions.length) * i, radarRadius).x"
                  :y2="polarToCartesian((360 / dimensions.length) * i, radarRadius).y"
                  stroke="#e5e7eb"
                  stroke-width="1"
                />
                <!-- Data polygon -->
                <polygon
                  :points="radarDataPoints"
                  fill="rgba(0, 168, 150, 0.2)"
                  stroke="#00a896"
                  stroke-width="2.5"
                />
                <!-- Data points -->
                <circle
                  v-for="(_, i) in dimensions"
                  :key="'point-' + i"
                  :cx="polarToCartesian((360 / dimensions.length) * i, radarRadius * (dimensionScores[i] / 4)).x"
                  :cy="polarToCartesian((360 / dimensions.length) * i, radarRadius * (dimensionScores[i] / 4)).y"
                  r="5"
                  fill="#00a896"
                  stroke="white"
                  stroke-width="2"
                />
                <!-- Labels -->
                <text
                  v-for="lp in radarLabelPositions"
                  :key="'label-' + lp.label"
                  :x="lp.x"
                  :y="lp.y"
                  text-anchor="middle"
                  dominant-baseline="middle"
                  class="text-xs font-semibold fill-gray-600"
                >
                  {{ lp.label }}
                </text>
              </svg>
            </div>
          </div>

          <!-- Score Bars -->
          <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-8">
            <h3 class="text-xl font-bold text-dark-slate mb-6">Score Breakdown</h3>
            <div class="space-y-5">
              <div v-for="(dim, i) in dimensions" :key="dim">
                <div class="flex justify-between items-center mb-1.5">
                  <span class="font-medium text-dark-slate text-sm">{{ dim }}</span>
                  <span class="font-bold text-primary text-sm">{{ dimensionScores[i].toFixed(1) }}/4.0</span>
                </div>
                <div class="w-full bg-gray-100 rounded-full h-3 overflow-hidden">
                  <div
                    class="h-3 rounded-full transition-all duration-700 ease-out"
                    :style="{
                      width: `${(dimensionScores[i] / 4) * 100}%`,
                      background: 'linear-gradient(90deg, #00a896, #02c39a)',
                    }"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Recommendations -->
        <div class="mt-12">
          <h3 class="text-2xl font-bold text-dark-slate mb-6">Personalized Recommendations</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div
              v-for="(dim, i) in dimensions"
              :key="'rec-' + dim"
              class="bg-white rounded-xl shadow-sm border border-gray-100 p-6"
              :class="{ 'ring-2 ring-primary ring-opacity-30': dimensionScores[i] < 2.5 }"
            >
              <div class="flex items-center gap-3 mb-4">
                <div
                  class="w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold"
                  :class="
                    dimensionScores[i] < 2
                      ? 'bg-red-50 text-red-600'
                      : dimensionScores[i] < 3
                        ? 'bg-amber-50 text-amber-600'
                        : 'bg-green-50 text-green-600'
                  "
                >
                  {{ dimensionScores[i].toFixed(1) }}
                </div>
                <h4 class="font-bold text-dark-slate">{{ dim }}</h4>
                <span
                  v-if="dimensionScores[i] < 2.5"
                  class="ml-auto text-xs font-semibold text-red-500 bg-red-50 px-2 py-0.5 rounded-full"
                >
                  Priority
                </span>
              </div>
              <ul class="space-y-2">
                <li
                  v-for="rec in dimensionRecommendations[dim]"
                  :key="rec"
                  class="text-sm text-gray-600 flex items-start gap-2"
                >
                  <svg class="w-4 h-4 text-primary flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                  </svg>
                  <span>{{ rec }}</span>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- CTAs -->
        <div class="mt-12 bg-white rounded-2xl shadow-lg border border-gray-100 p-8 md:p-12">
          <div class="text-center max-w-2xl mx-auto">
            <h3 class="text-2xl md:text-3xl font-bold text-dark-slate">Ready to Improve Your Score?</h3>
            <p class="text-gray-600 mt-4">
              Whether you scored Foundational or Integrated, we can help you design and implement a roadmap
              to take your data practice to the next level.
            </p>
            <div class="mt-8 flex flex-col sm:flex-row justify-center gap-4">
              <NuxtLink
                to="/booking"
                class="bg-primary text-white font-bold py-3 px-8 rounded-xl hover:bg-opacity-90 transition-all shadow-md hover:shadow-lg"
              >
                Book a Free Consultation
              </NuxtLink>
              <NuxtLink
                to="/services"
                class="bg-white text-primary font-bold py-3 px-8 rounded-xl border-2 border-primary hover:bg-primary hover:bg-opacity-5 transition-all"
              >
                Explore Our Services
              </NuxtLink>
            </div>
            <div class="mt-6 flex justify-center gap-4">
              <button
                @click="shareResults"
                class="text-sm text-gray-500 hover:text-primary transition-colors flex items-center gap-1"
              >
                <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M15 8a3 3 0 10-2.977-2.63l-4.94 2.47a3 3 0 100 4.319l4.94 2.47a3 3 0 10.895-1.789l-4.94-2.47a3.027 3.027 0 000-.74l4.94-2.47C13.456 7.68 14.19 8 15 8z" />
                </svg>
                Share Results
              </button>
              <button
                @click="restartAssessment"
                class="text-sm text-gray-500 hover:text-primary transition-colors flex items-center gap-1"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                Retake Assessment
              </button>
            </div>
          </div>
        </div>

        <!-- Related Services -->
        <div class="mt-12">
          <h3 class="text-xl font-bold text-dark-slate mb-6">Services That Can Help</h3>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <NuxtLink
              to="/services"
              class="block p-6 bg-white rounded-xl shadow-sm border border-gray-100 hover:shadow-lg hover:-translate-y-0.5 transition-all"
            >
              <h4 class="font-bold text-dark-slate">Data Governance</h4>
              <p class="text-sm text-gray-600 mt-2">DAMA-aligned framework for data quality and security</p>
              <span class="mt-3 inline-block text-primary font-semibold text-sm">Learn More &rarr;</span>
            </NuxtLink>
            <NuxtLink
              to="/services"
              class="block p-6 bg-white rounded-xl shadow-sm border border-gray-100 hover:shadow-lg hover:-translate-y-0.5 transition-all"
            >
              <h4 class="font-bold text-dark-slate">Data Architecture</h4>
              <p class="text-sm text-gray-600 mt-2">Scalable, secure data platforms on GCP</p>
              <span class="mt-3 inline-block text-primary font-semibold text-sm">Learn More &rarr;</span>
            </NuxtLink>
            <NuxtLink
              to="/services"
              class="block p-6 bg-white rounded-xl shadow-sm border border-gray-100 hover:shadow-lg hover:-translate-y-0.5 transition-all"
            >
              <h4 class="font-bold text-dark-slate">AI & ML Solutions</h4>
              <p class="text-sm text-gray-600 mt-2">From prototype to production AI systems</p>
              <span class="mt-3 inline-block text-primary font-semibold text-sm">Learn More &rarr;</span>
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
