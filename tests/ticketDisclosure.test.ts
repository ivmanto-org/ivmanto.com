import { test } from 'node:test'
import assert from 'node:assert/strict'
import { toggleTickets } from '../utils/ticketDisclosure.ts'

const row = ['gov', 'risk', 'comp']

test('opens every ticket of the row together on a wide screen', () => {
  assert.deepEqual(toggleTickets([], 'risk', row, true), row)
})

test('closes every ticket of the row together on a wide screen', () => {
  assert.deepEqual(toggleTickets(row, 'gov', row, true), [])
})

test('levels the row on a wide screen when only some tickets were opened on a phone', () => {
  assert.deepEqual(toggleTickets(['gov'], 'risk', row, true), row)
  assert.deepEqual(toggleTickets(['gov'], 'gov', row, true), [])
})

test('opens and closes one ticket alone on a narrow screen', () => {
  assert.deepEqual(toggleTickets([], 'risk', row, false), ['risk'])
  assert.deepEqual(toggleTickets(['gov', 'risk'], 'gov', row, false), ['risk'])
})

test('does not change the list it was given', () => {
  const open = ['gov']
  toggleTickets(open, 'risk', row, false)
  assert.deepEqual(open, ['gov'])
})
