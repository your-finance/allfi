/**
 * AllFi 主题系统 — 专业金融产品风格
 * 4套主题，统一低饱和度、高对比度配色
 * 统一字体: DM Sans + IBM Plex Sans + IBM Plex Mono
 *
 * 1. Nexus Pro - Bloomberg 风格深色（默认）
 * 2. Vestia   - GitHub 风格深色
 * 3. XChange  - 交易所风格深色
 * 4. Aurora   - 专业浅色
 */

// 统一字体配置
const sharedFonts = {
  heading: "'DM Sans', system-ui, -apple-system, sans-serif",
  body: "'IBM Plex Sans', system-ui, -apple-system, sans-serif",
  mono: "'IBM Plex Mono', 'SF Mono', 'Consolas', monospace",
  googleImport: "https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;600;700&family=IBM+Plex+Mono:wght@400;500;600&family=IBM+Plex+Sans:wght@300;400;500;600&display=swap"
}

export const themes = {
  // 主题1: Nexus Pro — Bloomberg 风格
  nexus: {
    id: 'nexus',
    name: 'Nexus Pro',
    description: '专业金融蓝，Bloomberg 风格',
    icon: '⬡',
    mode: 'dark',
    preview: '#0C0E12',
    fonts: sharedFonts,

    colors: {
      bgPrimary: '#0C0E12',
      bgSecondary: '#13161C',
      bgTertiary: '#1A1E26',
      bgElevated: '#1F242E',

      accentPrimary: '#4B83F0',
      accentSecondary: '#3B6FD4',
      accentTertiary: '#2EBD85',

      success: '#2EBD85',
      warning: '#D4A029',
      error: '#E25C5C',
      info: '#4B83F0',

      textPrimary: '#E1E4EA',
      textSecondary: '#8A919E',
      textMuted: '#555B66',
      textInverse: '#0C0E12',

      border: '#1E2330',
      borderHover: '#2A3040',
      borderActive: '#4B83F0',

      chartLine: '#4B83F0',
      chartGradientStart: 'rgba(75, 131, 240, 0.15)',
      chartGradientEnd: 'rgba(75, 131, 240, 0)',
    }
  },

  // 主题2: Vestia — GitHub 风格
  vestia: {
    id: 'vestia',
    name: 'Vestia',
    description: '投资平台风格，GitHub 深色',
    icon: '◆',
    mode: 'dark',
    preview: '#0D1117',
    fonts: sharedFonts,

    colors: {
      bgPrimary: '#0D1117',
      bgSecondary: '#161B22',
      bgTertiary: '#1D2028',
      bgElevated: '#252A34',

      accentPrimary: '#3B82D9',
      accentSecondary: '#2D6BC0',
      accentTertiary: '#2EA043',

      success: '#2EA043',
      warning: '#C89B2A',
      error: '#D14D4D',
      info: '#3B82D9',

      textPrimary: '#C9D1D9',
      textSecondary: '#848D97',
      textMuted: '#545D68',
      textInverse: '#0D1117',

      border: '#21262D',
      borderHover: '#2D333B',
      borderActive: '#3B82D9',

      chartLine: '#3B82D9',
      chartGradientStart: 'rgba(59, 130, 217, 0.15)',
      chartGradientEnd: 'rgba(59, 130, 217, 0)',
    }
  },

  // 主题3: XChange — 交易所风格
  xchange: {
    id: 'xchange',
    name: 'XChange',
    description: '交易所风格，沉稳绿色',
    icon: '✕',
    mode: 'dark',
    preview: '#0B0D10',
    fonts: sharedFonts,

    colors: {
      bgPrimary: '#0B0D10',
      bgSecondary: '#12151A',
      bgTertiary: '#191D23',
      bgElevated: '#1F242B',

      accentPrimary: '#3EA87A',
      accentSecondary: '#2F8F65',
      accentTertiary: '#5697CC',

      success: '#3EA87A',
      warning: '#C4952B',
      error: '#D55B5B',
      info: '#5697CC',

      textPrimary: '#D8DBE0',
      textSecondary: '#8A8F97',
      textMuted: '#585D65',
      textInverse: '#0B0D10',

      border: '#1C2028',
      borderHover: '#262D37',
      borderActive: '#3EA87A',

      chartLine: '#3EA87A',
      chartGradientStart: 'rgba(62, 168, 122, 0.15)',
      chartGradientEnd: 'rgba(62, 168, 122, 0)',
    }
  },

  // 主题4: Aurora — 专业浅色
  aurora: {
    id: 'aurora',
    name: 'Aurora',
    description: '专业浅色主题',
    icon: '◐',
    mode: 'light',
    preview: '#FAFBFC',
    fonts: sharedFonts,

    colors: {
      bgPrimary: '#FAFBFC',
      bgSecondary: '#F0F2F5',
      bgTertiary: '#E4E7EC',
      bgElevated: '#D5D9E0',

      accentPrimary: '#3574D4',
      accentSecondary: '#2A5FB8',
      accentTertiary: '#1A9B5A',

      success: '#1A9B5A',
      warning: '#B88A1F',
      error: '#D14545',
      info: '#3574D4',

      textPrimary: '#1A1D23',
      textSecondary: '#5F6670',
      textMuted: '#9CA3AE',
      textInverse: '#FFFFFF',

      border: '#D5D9E0',
      borderHover: '#B8BFC8',
      borderActive: '#3574D4',

      chartLine: '#3574D4',
      chartGradientStart: 'rgba(53, 116, 212, 0.12)',
      chartGradientEnd: 'rgba(53, 116, 212, 0)',
    }
  }
}

// 默认主题
export const defaultTheme = 'nexus'

// 主题列表
export const themeList = Object.values(themes)

/**
 * 将主题配置转换为CSS变量
 */
export function themeToCssVars(theme) {
  const { colors, fonts } = theme

  return {
    // 字体
    '--font-heading': fonts.heading,
    '--font-body': fonts.body,
    '--font-mono': fonts.mono,

    // 背景色
    '--color-bg-primary': colors.bgPrimary,
    '--color-bg-secondary': colors.bgSecondary,
    '--color-bg-tertiary': colors.bgTertiary,
    '--color-bg-elevated': colors.bgElevated,

    // 强调色
    '--color-accent-primary': colors.accentPrimary,
    '--color-accent-secondary': colors.accentSecondary,
    '--color-accent-tertiary': colors.accentTertiary,

    // 功能色
    '--color-success': colors.success,
    '--color-warning': colors.warning,
    '--color-error': colors.error,
    '--color-info': colors.info,

    // 文字色
    '--color-text-primary': colors.textPrimary,
    '--color-text-secondary': colors.textSecondary,
    '--color-text-muted': colors.textMuted,
    '--color-text-inverse': colors.textInverse,

    // 边框
    '--color-border': colors.border,
    '--color-border-hover': colors.borderHover,
    '--color-border-active': colors.borderActive,

    // 图表色
    '--chart-line': colors.chartLine,
    '--chart-gradient-start': colors.chartGradientStart,
    '--chart-gradient-end': colors.chartGradientEnd,
  }
}

/**
 * 动态加载主题字体
 */
export function loadThemeFonts(theme) {
  if (!theme.fonts?.googleImport) return

  const fontUrl = theme.fonts.googleImport
  const linkId = 'theme-font-shared'

  let link = document.getElementById(linkId)
  if (!link) {
    link = document.createElement('link')
    link.id = linkId
    link.rel = 'stylesheet'
    link.href = fontUrl
    document.head.appendChild(link)
  }
}

/**
 * 预加载字体（所有主题共用同一字体）
 */
export function preloadAllFonts() {
  const fontUrl = sharedFonts.googleImport
  const link = document.createElement('link')
  link.rel = 'preload'
  link.as = 'style'
  link.href = fontUrl
  document.head.appendChild(link)
}

/**
 * 应用主题到DOM
 */
export function applyTheme(themeId) {
  const theme = themes[themeId]
  if (!theme) {
    console.warn(`主题 "${themeId}" 不存在，使用默认主题`)
    return applyTheme(defaultTheme)
  }

  loadThemeFonts(theme)

  const cssVars = themeToCssVars(theme)
  const root = document.documentElement

  Object.entries(cssVars).forEach(([key, value]) => {
    root.style.setProperty(key, value)
  })

  root.setAttribute('data-theme', themeId)
  root.setAttribute('data-theme-mode', theme.mode)

  let metaTheme = document.querySelector('meta[name="theme-color"]')
  if (!metaTheme) {
    metaTheme = document.createElement('meta')
    metaTheme.name = 'theme-color'
    document.head.appendChild(metaTheme)
  }
  metaTheme.setAttribute('content', theme.colors.bgPrimary)

  return theme
}
