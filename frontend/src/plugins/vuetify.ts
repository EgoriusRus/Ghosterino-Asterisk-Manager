import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi-svg'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { ru } from 'vuetify/locale'

export default createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: { mdi }
  },
  locale: {
    locale: 'ru',
    messages: { ru }
  },
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          primary: '#4f46e5',
          secondary: '#6b7280',
          error: '#ef4444',
          success: '#22c55e',
          warning: '#f59e0b',
          info: '#3b82f6'
        }
      }
    }
  },
  defaults: {
    VBtn: {
      variant: 'flat'
    },
    VTextField: {
      variant: 'outlined',
      density: 'comfortable'
    },
    VSelect: {
      variant: 'outlined',
      density: 'comfortable'
    },
    VDataTableServer: {
      density: 'comfortable'
    }
  }
})
