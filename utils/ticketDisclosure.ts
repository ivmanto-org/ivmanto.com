// "What's included" on the compact tickets (DS v1.4 cards.md): from 900px the tickets of
// one row open and close together so the row stays level; below that each opens alone.
export const TICKET_ROW_MIN_WIDTH = '(min-width: 900px)'

/** Returns the ids open after `id` was toggled. The clicked ticket decides the direction. */
export function toggleTickets(open: readonly string[], id: string, row: readonly string[], isWide: boolean): string[] {
  const opening = !open.includes(id)
  if (isWide) return opening ? [...row] : []
  return opening ? [...open, id] : open.filter((o) => o !== id)
}
