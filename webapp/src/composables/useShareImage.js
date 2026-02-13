/**
 * 分享图生成工具
 * 使用 html2canvas 将 DOM 元素渲染为 PNG 图片并下载
 */
import html2canvas from 'html2canvas'

/**
 * 生成分享图片并下载
 * @param {HTMLElement} element - 要截图的 DOM 元素
 * @param {string} filename - 文件名（不含扩展名）
 */
export async function generateShareImage(element, filename = 'allfi-portfolio') {
  const canvas = await html2canvas(element, {
    width: 1200,
    height: 630,
    scale: 2,
    backgroundColor: null,
    useCORS: true,
    logging: false
  })

  // 转换为 blob 并下载
  canvas.toBlob((blob) => {
    if (!blob) return
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${filename}.png`
    link.click()
    URL.revokeObjectURL(url)
  }, 'image/png')
}
