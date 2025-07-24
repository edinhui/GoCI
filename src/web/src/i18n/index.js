import { createI18n } from 'vue-i18n'
import en from '../locales/en'
import zh from '../locales/zh'

const messages = {
  en,
  zh
}

// 从localStorage获取语言设置，如果没有则使用中文作为默认语言
const locale = localStorage.getItem('language') || 'zh'

const i18n = createI18n({
  legacy: false, // 使用Vue 3 Composition API
  locale: locale,
  fallbackLocale: 'en',
  globalInjection: true, // 全局注入$t函数
  messages
})

export default i18n
