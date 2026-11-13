declare global {
  interface Window {
    dataLayer: Record<string, any>[]
    gtmInitialized?: boolean
  }
}

/**
 * Pushes a custom event to the Google Tag Manager data layer.
 * @param {string} eventName The name of the event.
 * @param {Record<string, any>} params Additional data to send with the event.
 */
export const trackEvent = (eventName: string, params: Record<string, any> = {}) => {
  // Ensure this only runs in the browser where window.dataLayer exists.
  if (typeof window === 'undefined') {
    return
  }
  window.dataLayer = window.dataLayer || []
  window.dataLayer.push({
    event: eventName,
    ...params,
  })
}

const PRODUCTION_HOSTS = ['ivmanto.com', 'www.ivmanto.com']

/**
 * Injects the Google Tag Manager script into the page.
 * This should be called once, after cookie consent is given.
 */
export const initGtm = () => {
  // IMPORTANT: Replace with your actual GTM ID, which starts with "GTM-".
  const gtmId = 'GTM-PXRF8FPQ'

  // Ensure this only runs in the browser and hasn't been run before.
  if (typeof window === 'undefined' || window.gtmInitialized) {
    return
  }
  // Analytics runs on the production domain only: staging, previews and local
  // builds must not send visits into the production data (Nick, 2026-11-12).
  if (!PRODUCTION_HOSTS.includes(window.location.hostname)) {
    return
  }
  window.gtmInitialized = true

  // This is the standard GTM script, made TypeScript-friendly.
  // The original uses an IIFE with implicit 'any' types, which can cause compiler errors.
  ;(function (w: Window, d: Document, s: string, l: string, i: string) {
    w.dataLayer = w.dataLayer || []
    w.dataLayer.push({ 'gtm.start': new Date().getTime(), event: 'gtm.js' })
    const f = d.getElementsByTagName(s)[0]
    const j = d.createElement(s) as HTMLScriptElement
    const dl = l !== 'dataLayer' ? '&l=' + l : ''
    j.async = true
    j.src = 'https://www.googletagmanager.com/gtm.js?id=' + i + dl
    f.parentNode?.insertBefore(j, f)
  })(window, document, 'script', 'dataLayer', gtmId)
}

type GaSessionInfo = {
  clientId: string | null
  sessionId: string | null
}

/**
 * Retrieves the Google Analytics Client ID and Session ID from browser cookies.
 * This is crucial for stitching frontend user sessions with backend (server-to-server) events.
 * @param {string} measurementId The GA4 Measurement ID (e.g., 'G-XXXXXXXXXX').
 * @returns {GaSessionInfo} An object containing the clientId and sessionId.
 */
export const getGaSessionInfo = (measurementId: string): GaSessionInfo => {
  // Ensure this only runs in the browser where document.cookie exists.
  if (typeof document === 'undefined' || typeof document.cookie === 'undefined') {
    return { clientId: null, sessionId: null }
  }

  let clientId: string | null = null
  let sessionId: string | null = null
  const measurementCookieName = `_ga_${measurementId.replace('G-', '')}`

  const cookies = document.cookie.split(';')
  for (const cookie of cookies) {
    const cookieStr = cookie.trim()

    if (cookieStr.startsWith('_ga=')) {
      const value = cookieStr.substring(4) // length of '_ga='
      clientId = value.split('.').slice(2).join('.')
    }

    if (cookieStr.startsWith(measurementCookieName + '=')) {
      const value = cookieStr.substring(measurementCookieName.length + 1)
      const parts = value.split('.')
      if (parts.length >= 3) sessionId = parts[2]
    }
  }

  return { clientId, sessionId }
}
