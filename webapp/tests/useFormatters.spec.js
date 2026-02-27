import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useFormatters } from '../src/composables/useFormatters'

describe('useFormatters', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('formats percentages correctly', () => {
        const { formatPercent } = useFormatters()

        // Test base cases
        expect(formatPercent(5.1234)).toBe('+5.12%')
        expect(formatPercent(-3.456)).toBe('-3.46%')
        expect(formatPercent(0)).toBe('0.00%')

        // Test without sign
        expect(formatPercent(5.1234, false)).toBe('5.12%')

        // Test edge cases
        expect(formatPercent(null)).toBe('0.00%')
        expect(formatPercent(undefined)).toBe('0.00%')
        expect(formatPercent(NaN)).toBe('0.00%')
    })

    it('shortens addresses correctly', () => {
        const { shortenAddress } = useFormatters()

        const address = '0x1234567890abcdef1234567890abcdef12345678'
        expect(shortenAddress(address)).toBe('0x1234...5678')
        expect(shortenAddress(address, 6)).toBe('0x123456...345678')

        expect(shortenAddress('')).toBe('')
        expect(shortenAddress(null)).toBe('')
    })
})
